package examples

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

// TestProxyBasicHTTP 测试基本的 HTTP 代理功能
// 注意: 此测试需要代理服务器运行在 localhost:8888
func TestProxyBasicHTTP(t *testing.T) {
	// 跳过自动化测试，仅用于手动测试
	t.Skip("Manual test only - requires proxy server running")

	// 创建使用代理的 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(mustParseURL("http://127.0.0.1:8888")),
		},
		Timeout: 10 * time.Second,
	}

	// 发送 HTTP 请求
	resp, err := client.Get("http://httpbin.org/get")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", body)
}

// TestProxyHTTPS 测试 HTTPS 代理功能
func TestProxyHTTPS(t *testing.T) {
	t.Skip("Manual test only - requires proxy server running")

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(mustParseURL("http://127.0.0.1:8888")),
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		t.Fatalf("Failed to send HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("HTTPS Response: %s\n", body)
}

// TestProxyPOST 测试 POST 请求
func TestProxyPOST(t *testing.T) {
	t.Skip("Manual test only - requires proxy server running")

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(mustParseURL("http://127.0.0.1:8888")),
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(
		"https://httpbin.org/post",
		"application/json",
		mustReadCloser(`{"test": "data"}`),
	)
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func mustReadCloser(s string) io.ReadCloser {
	return io.NopCloser(bytes.NewBufferString(s))
}
