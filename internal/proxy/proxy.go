package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"packet-capture-tool/internal/cert"
	"packet-capture-tool/internal/storage"
)

// Server 代理服务器
type Server struct {
	port         int
	certManager  *cert.Manager
	storage      *storage.Storage
	server       *http.Server
	listener     net.Listener
	httpsEnabled bool
}

// NewServer 创建代理服务器
func NewServer(port int, certManager *cert.Manager, storage *storage.Storage) *Server {
	return &Server{
		port:         port,
		certManager:  certManager,
		storage:      storage,
		httpsEnabled: true,
	}
}

// Start 启动代理服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)

	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}

	s.server = &http.Server{
		Handler: http.HandlerFunc(s.handleRequest),
	}

	log.Printf("Proxy server started on %s", addr)

	go func() {
		if err := s.server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Proxy server error: %v", err)
		}
	}()

	return nil
}

// Stop 停止代理服务器
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// handleRequest 处理请求
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// 检查是否是CONNECT请求（HTTPS）
	if r.Method == http.MethodConnect {
		s.handleHTTPS(w, r, startTime)
		return
	}

	// 检查是否是WebSocket升级请求
	if isWebSocketRequest(r) {
		s.handleWebSocket(w, r, startTime)
		return
	}

	// 处理普通HTTP请求
	s.handleHTTP(w, r, startTime)
}

// handleHTTP 处理HTTP请求
func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	// 读取请求体
	var reqBody []byte
	if r.Body != nil {
		reqBody, _ = io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	// 创建新的请求
	targetURL := r.URL.String()
	if !strings.HasPrefix(targetURL, "http") {
		targetURL = "http://" + r.Host + r.RequestURI
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL, bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// 复制请求头
	copyHeaders(proxyReq.Header, r.Header)
	proxyReq.Header.Del("Proxy-Connection")

	// 发送请求
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不自动跟随重定向
		},
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to proxy request", http.StatusBadGateway)
		s.logRequest(r, nil, reqBody, nil, 0, startTime)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, _ := io.ReadAll(resp.Body)

	// 复制响应头
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	// 记录请求
	s.logRequest(r, resp, reqBody, respBody, resp.StatusCode, startTime)
}

// handleHTTPS 处理HTTPS CONNECT请求
func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	if !s.httpsEnabled {
		http.Error(w, "HTTPS interception is disabled", http.StatusForbidden)
		return
	}

	// 建立到目标服务器的连接
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, "Failed to connect to target", http.StatusBadGateway)
		return
	}
	defer destConn.Close()

	// 劫持客户端连接
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Failed to hijack connection", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// 发送200 Connection Established响应
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// 使用我们的证书进行TLS握手
	tlsConfig := s.certManager.GetTLSConfig()
	tlsClientConn := tls.Server(clientConn, tlsConfig)
	defer tlsClientConn.Close()

	if err := tlsClientConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}

	// 读取客户端的HTTP请求
	reader := bufio.NewReader(tlsClientConn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Printf("Failed to read HTTPS request: %v", err)
		return
	}

	// 构建完整的URL
	req.URL.Scheme = "https"
	req.URL.Host = r.Host

	// 处理HTTPS请求（类似HTTP）
	s.handleDecryptedHTTPS(tlsClientConn, req, r.Host, startTime)
}

// handleDecryptedHTTPS 处理解密后的HTTPS请求
func (s *Server) handleDecryptedHTTPS(clientConn net.Conn, r *http.Request, host string, startTime time.Time) {
	// 读取请求体
	var reqBody []byte
	if r.Body != nil {
		reqBody, _ = io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	// 创建到目标服务器的TLS连接
	targetConn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: true, // 在生产环境中应该验证证书
	})
	if err != nil {
		log.Printf("Failed to connect to target server: %v", err)
		return
	}
	defer targetConn.Close()

	// 发送请求到目标服务器
	if err := r.Write(targetConn); err != nil {
		log.Printf("Failed to write request to target: %v", err)
		return
	}

	// 读取响应
	resp, err := http.ReadResponse(bufio.NewReader(targetConn), r)
	if err != nil {
		log.Printf("Failed to read response from target: %v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, _ := io.ReadAll(resp.Body)

	// 将响应发送回客户端
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	if err := resp.Write(clientConn); err != nil {
		log.Printf("Failed to write response to client: %v", err)
		return
	}

	// 记录请求
	s.logRequest(r, resp, reqBody, respBody, resp.StatusCode, startTime)
}

// handleWebSocket 处理WebSocket请求
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	// WebSocket需要特殊处理，这里简化处理
	// 实际生产环境需要完整的WebSocket协议支持

	targetURL := r.URL.String()
	if !strings.HasPrefix(targetURL, "ws") {
		targetURL = "ws://" + r.Host + r.RequestURI
	}

	// 记录WebSocket连接
	s.storage.Add(&storage.CapturedRequest{
		Timestamp:      startTime,
		Method:         r.Method,
		URL:            targetURL,
		Host:           r.Host,
		Path:           r.URL.Path,
		Protocol:       "WebSocket",
		RequestHeaders: cloneHeaders(r.Header),
		Duration:       time.Since(startTime),
	})

	// 这里简化处理，直接返回
	http.Error(w, "WebSocket support is limited in this version", http.StatusNotImplemented)
}

// logRequest 记录请求到存储
func (s *Server) logRequest(r *http.Request, resp *http.Response, reqBody, respBody []byte, statusCode int, startTime time.Time) {
	protocol := "HTTP"
	if r.TLS != nil || r.URL.Scheme == "https" {
		protocol = "HTTPS"
	}

	urlStr := r.URL.String()
	if !strings.HasPrefix(urlStr, "http") {
		urlStr = protocol + "://" + r.Host + r.RequestURI
	}

	captured := &storage.CapturedRequest{
		Timestamp:      startTime,
		Method:         r.Method,
		URL:            urlStr,
		Host:           r.Host,
		Path:           r.URL.Path,
		Protocol:       protocol,
		StatusCode:     statusCode,
		RequestHeaders: cloneHeaders(r.Header),
		RequestBody:    reqBody,
		Duration:       time.Since(startTime),
	}

	if resp != nil {
		captured.ResponseHeaders = cloneHeaders(resp.Header)
		captured.ResponseBody = respBody
	}

	s.storage.Add(captured)
}

// isWebSocketRequest 检查是否是WebSocket升级请求
func isWebSocketRequest(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

// copyHeaders 复制HTTP头
func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

// cloneHeaders 克隆HTTP头
func cloneHeaders(src http.Header) map[string][]string {
	dst := make(map[string][]string)
	for key, values := range src {
		dst[key] = append([]string{}, values...)
	}
	return dst
}

// SetHTTPSEnabled 设置HTTPS拦截开关
func (s *Server) SetHTTPSEnabled(enabled bool) {
	s.httpsEnabled = enabled
}
