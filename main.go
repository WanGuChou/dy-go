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
	
	log.Println("========================================")
	log.Println("Packet Capture Tool Starting...")
	log.Println("========================================")

	// 获取程序数据目录
	log.Println("[1/6] Getting data directory...")
	dataDir, err := getDataDir()
	if err != nil {
		log.Fatalf("Failed to get data directory: %v", err)
	}
	log.Printf("Data directory: %s", dataDir)

	// 确保数据目录存在
	log.Println("[2/6] Creating data directory...")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
	log.Println("✓ Data directory ready")

	// 加载配置
	log.Println("[3/6] Loading configuration...")
	cfg := config.LoadOrCreate(filepath.Join(dataDir, "config.json"))
	log.Printf("✓ Configuration loaded (Port: %d, HTTPS: %v)", cfg.GetProxyPort(), cfg.GetHTTPSDecrypt())

	// 初始化证书管理器
	log.Println("[4/6] Initializing certificate manager...")
	certManager, err := cert.NewManager(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize certificate manager: %v", err)
	}
	log.Println("✓ Certificate manager initialized")

	// 检查并安装根证书
	log.Println("[5/6] Checking root certificate...")
	if !certManager.IsRootCertInstalled() {
		log.Println("⚠ Root certificate not installed, installing...")
		if err := certManager.InstallRootCert(); err != nil {
			log.Printf("⚠ Warning: Failed to install root certificate: %v", err)
			log.Println("⚠ You may need to run as administrator to install the certificate.")
			log.Println("⚠ You can manually install it later from GUI")
		} else {
			log.Println("✓ Root certificate installed successfully!")
		}
	} else {
		log.Println("✓ Root certificate already installed")
	}

	// 启动应用程序
	log.Println("[6/6] Starting application...")
	log.Println("✓ GUI will open shortly...")
	log.Println("✓ System tray icon will appear in taskbar")
	log.Println("========================================")
	log.Println("Application is running...")
	log.Println("Close this window or press Ctrl+C to exit")
	log.Println("========================================")
	
	application := app.NewApp(cfg, certManager, dataDir)
	application.Run()
	
	log.Println("Application stopped.")
}

// getDataDir 获取程序数据存储目录
func getDataDir() (string, error) {
	userDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userDataDir, "PacketCaptureTool"), nil
}
