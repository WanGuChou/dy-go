package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Config 应用程序配置
type Config struct {
	ProxyPort      int  `json:"proxy_port"`
	HTTPSDecrypt   bool `json:"https_decrypt"`
	AutoStart      bool `json:"auto_start"`
	AutoSetProxy   bool `json:"auto_set_proxy"`
	mu             sync.RWMutex
	configFilePath string
}

// LoadOrCreate 加载或创建配置文件
func LoadOrCreate(path string) *Config {
	cfg := &Config{
		ProxyPort:      8888,
		HTTPSDecrypt:   true,
		AutoStart:      false,
		AutoSetProxy:   false,
		configFilePath: path,
	}

	// 尝试从文件加载配置
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			log.Printf("Failed to parse config file: %v", err)
		}
	}

	cfg.configFilePath = path
	return cfg
}

// Save 保存配置到文件
func (c *Config) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.configFilePath, data, 0644)
}

// GetProxyPort 获取代理端口
func (c *Config) GetProxyPort() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ProxyPort
}

// SetProxyPort 设置代理端口
func (c *Config) SetProxyPort(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ProxyPort = port
}

// GetHTTPSDecrypt 获取HTTPS解密开关
func (c *Config) GetHTTPSDecrypt() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.HTTPSDecrypt
}

// SetHTTPSDecrypt 设置HTTPS解密开关
func (c *Config) SetHTTPSDecrypt(enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.HTTPSDecrypt = enabled
}
