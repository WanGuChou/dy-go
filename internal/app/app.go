package app

import (
	"log"

	"packet-capture-tool/internal/cert"
	"packet-capture-tool/internal/config"
	"packet-capture-tool/internal/gui"
	"packet-capture-tool/internal/proxy"
	"packet-capture-tool/internal/storage"
	"packet-capture-tool/internal/systray_manager"
)

// App 应用程序主控制器
type App struct {
	config         *config.Config
	certManager    *cert.Manager
	storage        *storage.Storage
	proxyServer    *proxy.Server
	gui            *gui.GUI
	systrayManager *systray_manager.Manager
	dataDir        string
}

// NewApp 创建应用程序实例
func NewApp(cfg *config.Config, certMgr *cert.Manager, dataDir string) *App {
	// 创建存储
	store := storage.NewStorage()

	// 创建代理服务器
	proxyServer := proxy.NewServer(cfg.GetProxyPort(), certMgr, store)

	// 创建GUI
	guiInstance := gui.NewGUI(cfg, store)

	app := &App{
		config:      cfg,
		certManager: certMgr,
		storage:     store,
		proxyServer: proxyServer,
		gui:         guiInstance,
		dataDir:     dataDir,
	}

	// 设置GUI回调
	app.setupGUICallbacks()

	// 创建系统托盘管理器
	app.systrayManager = systray_manager.NewManager(
		cfg.GetProxyPort(),
		app.onProxyEnabled,
		app.onProxyDisabled,
		app.onExit,
	)

	return app
}

// setupGUICallbacks 设置GUI回调函数
func (a *App) setupGUICallbacks() {
	// 代理端口变化
	a.gui.SetOnProxyPortChange(func(port int) {
		log.Printf("Proxy port changed to: %d", port)
		a.config.SetProxyPort(port)
		a.config.Save()

		// 重启代理服务器
		a.proxyServer.Stop()
		a.proxyServer = proxy.NewServer(port, a.certManager, a.storage)
		if err := a.proxyServer.Start(); err != nil {
			log.Printf("Failed to restart proxy server: %v", err)
		}
	})

	// HTTPS解密开关变化
	a.gui.SetOnHTTPSDecryptChange(func(enabled bool) {
		log.Printf("HTTPS decrypt changed to: %v", enabled)
		a.config.SetHTTPSDecrypt(enabled)
		a.config.Save()
		a.proxyServer.SetHTTPSEnabled(enabled)
	})

	// 重新安装证书
	a.gui.SetOnReinstallCert(func() {
		log.Println("Reinstalling certificate...")
		if err := a.certManager.InstallRootCert(); err != nil {
			log.Printf("Failed to reinstall certificate: %v", err)
		} else {
			log.Println("Certificate reinstalled successfully!")
		}
	})

	// 清空请求
	a.gui.SetOnClearRequests(func() {
		log.Println("Clearing all requests...")
		a.storage.Clear()
	})
}

// onProxyEnabled 代理启用回调
func (a *App) onProxyEnabled() {
	log.Println("System proxy enabled")
}

// onProxyDisabled 代理禁用回调
func (a *App) onProxyDisabled() {
	log.Println("System proxy disabled")
}

// onExit 退出回调
func (a *App) onExit() {
	log.Println("Application exiting...")

	// 停止代理服务器
	if a.proxyServer != nil {
		a.proxyServer.Stop()
	}

	// 保存配置
	a.config.Save()
}

// Run 运行应用程序
func (a *App) Run() {
	// 启动代理服务器
	if err := a.proxyServer.Start(); err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}

	log.Printf("Proxy server started on port %d", a.config.GetProxyPort())
	log.Printf("HTTPS decryption: %v", a.config.GetHTTPSDecrypt())
	log.Printf("Certificate fingerprint: %s", a.certManager.GetCertFingerprint())

	// 启动系统托盘
	a.systrayManager.Start()

	// 运行GUI（阻塞）
	a.gui.Run()
}
