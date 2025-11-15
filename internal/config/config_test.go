package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_LoadOrCreate(t *testing.T) {
	// 创建临时文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// 测试创建新配置
	cfg := LoadOrCreate(configPath)

	if cfg.ProxyPort != 8888 {
		t.Errorf("Expected default port 8888, got %d", cfg.ProxyPort)
	}

	if !cfg.HTTPSDecrypt {
		t.Error("Expected HTTPS decrypt to be true by default")
	}
}

func TestConfig_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := LoadOrCreate(configPath)
	cfg.SetProxyPort(9999)
	cfg.SetHTTPSDecrypt(false)

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// 重新加载配置
	cfg2 := LoadOrCreate(configPath)

	if cfg2.GetProxyPort() != 9999 {
		t.Errorf("Expected port 9999, got %d", cfg2.GetProxyPort())
	}

	if cfg2.GetHTTPSDecrypt() {
		t.Error("Expected HTTPS decrypt to be false")
	}
}

func TestConfig_ThreadSafety(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := LoadOrCreate(configPath)

	// 并发读写测试
	done := make(chan bool)

	// Writer goroutines
	for i := 0; i < 10; i++ {
		go func(port int) {
			cfg.SetProxyPort(8000 + port)
			cfg.SetHTTPSDecrypt(port%2 == 0)
			done <- true
		}(i)
	}

	// Reader goroutines
	for i := 0; i < 10; i++ {
		go func() {
			_ = cfg.GetProxyPort()
			_ = cfg.GetHTTPSDecrypt()
			done <- true
		}()
	}

	// 等待所有 goroutines 完成
	for i := 0; i < 20; i++ {
		<-done
	}
}
