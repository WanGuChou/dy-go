package main

import (
	"log"
	"os"
	"path/filepath"

	"packet-capture-tool/internal/app"
	"packet-capture-tool/internal/cert"
	"packet-capture-tool/internal/config"
)

func main() {
	// 初始化日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 获取程序数据目录
	dataDir, err := getDataDir()
	if err != nil {
		log.Fatalf("Failed to get data directory: %v", err)
	}

	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// 加载配置
	cfg := config.LoadOrCreate(filepath.Join(dataDir, "config.json"))

	// 初始化证书管理器
	certManager, err := cert.NewManager(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize certificate manager: %v", err)
	}

	// 检查并安装根证书
	if !certManager.IsRootCertInstalled() {
		log.Println("Root certificate not installed, installing...")
		if err := certManager.InstallRootCert(); err != nil {
			log.Printf("Warning: Failed to install root certificate: %v", err)
			log.Println("You may need to run as administrator to install the certificate.")
		} else {
			log.Println("Root certificate installed successfully!")
		}
	}

	// 启动应用程序
	application := app.NewApp(cfg, certManager, dataDir)
	application.Run()
}

// getDataDir 获取程序数据存储目录
func getDataDir() (string, error) {
	userDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userDataDir, "PacketCaptureTool"), nil
}
