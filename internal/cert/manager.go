package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Manager 证书管理器
type Manager struct {
	dataDir     string
	rootCert    *x509.Certificate
	rootKey     *rsa.PrivateKey
	certCache   map[string]*tls.Certificate
	certCacheMu sync.RWMutex
}

// NewManager 创建证书管理器
func NewManager(dataDir string) (*Manager, error) {
	m := &Manager{
		dataDir:   dataDir,
		certCache: make(map[string]*tls.Certificate),
	}

	// 加载或生成根证书
	if err := m.loadOrGenerateRootCert(); err != nil {
		return nil, err
	}

	return m, nil
}

// loadOrGenerateRootCert 加载或生成根证书
func (m *Manager) loadOrGenerateRootCert() error {
	certPath := filepath.Join(m.dataDir, "root-ca.crt")
	keyPath := filepath.Join(m.dataDir, "root-ca.key")

	// 尝试加载现有证书
	if _, err := os.Stat(certPath); err == nil {
		if _, err := os.Stat(keyPath); err == nil {
			return m.loadRootCert(certPath, keyPath)
		}
	}

	// 生成新的根证书
	return m.generateRootCert(certPath, keyPath)
}

// loadRootCert 加载根证书
func (m *Manager) loadRootCert(certPath, keyPath string) error {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return err
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return fmt.Errorf("failed to decode key PEM")
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return err
	}

	m.rootCert = cert
	m.rootKey = key

	log.Println("Loaded existing root certificate")
	return nil
}

// generateRootCert 生成根证书
func (m *Manager) generateRootCert(certPath, keyPath string) error {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Packet Capture Tool"},
			CommonName:   "Packet Capture Tool Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10年有效期
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 自签名证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// 保存证书
	certOut, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	// 保存私钥
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return err
	}

	// 解析证书
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return err
	}

	m.rootCert = cert
	m.rootKey = privateKey

	log.Println("Generated new root certificate")
	return nil
}

// InstallRootCert 安装根证书到Windows系统信任存储
func (m *Manager) InstallRootCert() error {
	certPath := filepath.Join(m.dataDir, "root-ca.crt")

	// 使用certutil命令安装证书
	cmd := exec.Command("certutil", "-addstore", "-f", "Root", certPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install certificate: %v, output: %s", err, output)
	}

	log.Println("Root certificate installed to system trust store")
	return nil
}

// IsRootCertInstalled 检查根证书是否已安装
func (m *Manager) IsRootCertInstalled() bool {
	if m.rootCert == nil {
		return false
	}

	// 使用certutil检查证书是否在系统存储中
	cmd := exec.Command("certutil", "-verifystore", "Root", m.rootCert.Subject.CommonName)
	err := cmd.Run()
	return err == nil
}

// GetCertFingerprint 获取根证书指纹
func (m *Manager) GetCertFingerprint() string {
	if m.rootCert == nil {
		return ""
	}
	return fmt.Sprintf("%X", m.rootCert.SerialNumber)
}

// GenerateCertForHost 为指定主机动态生成证书（用于MITM）
func (m *Manager) GenerateCertForHost(host string) (*tls.Certificate, error) {
	// 检查缓存
	m.certCacheMu.RLock()
	if cert, ok := m.certCache[host]; ok {
		m.certCacheMu.RUnlock()
		return cert, nil
	}
	m.certCacheMu.RUnlock()

	// 生成新证书
	m.certCacheMu.Lock()
	defer m.certCacheMu.Unlock()

	// 再次检查（防止并发生成）
	if cert, ok := m.certCache[host]; ok {
		return cert, nil
	}

	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// 创建证书模板
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Packet Capture Tool"},
			CommonName:   host,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // 1年有效期
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{host},
	}

	// 使用根证书签名
	certDER, err := x509.CreateCertificate(rand.Reader, &template, m.rootCert, &privateKey.PublicKey, m.rootKey)
	if err != nil {
		return nil, err
	}

	// 构建tls.Certificate
	cert := &tls.Certificate{
		Certificate: [][]byte{certDER, m.rootCert.Raw},
		PrivateKey:  privateKey,
	}

	// 缓存证书
	m.certCache[host] = cert

	return cert, nil
}

// GetTLSConfig 获取TLS配置（用于代理服务器）
func (m *Manager) GetTLSConfig() *tls.Config {
	return &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return m.GenerateCertForHost(hello.ServerName)
		},
	}
}
