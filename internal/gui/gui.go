package gui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"packet-capture-tool/internal/config"
	"packet-capture-tool/internal/storage"
)

// GUI 图形用户界面
type GUI struct {
	app         fyne.App
	mainWindow  fyne.Window
	config      *config.Config
	storage     *storage.Storage
	requestList *widget.Table
	searchEntry *widget.Entry
	statusLabel *widget.Label

	// 回调函数
	onProxyPortChange    func(int)
	onHTTPSDecryptChange func(bool)
	onReinstallCert      func()
	onClearRequests      func()

	// 数据
	filteredRequests []*storage.CapturedRequest
}

// NewGUI 创建GUI实例
func NewGUI(cfg *config.Config, store *storage.Storage) *GUI {
	myApp := app.New()

	gui := &GUI{
		app:              myApp,
		mainWindow:       myApp.NewWindow("Packet Capture Tool - 本地抓包工具"),
		config:           cfg,
		storage:          store,
		filteredRequests: make([]*storage.CapturedRequest, 0),
	}

	gui.setupUI()
	return gui
}

// setupUI 设置用户界面
func (g *GUI) setupUI() {
	// 顶部配置区域
	configArea := g.createConfigArea()

	// 中间搜索区域
	searchArea := g.createSearchArea()

	// 请求列表区域
	requestArea := g.createRequestArea()

	// 底部状态栏
	g.statusLabel = widget.NewLabel("就绪")
	statusBar := container.NewBorder(nil, nil, nil, nil, g.statusLabel)

	// 组合布局
	content := container.NewBorder(
		container.NewVBox(configArea, searchArea),
		statusBar,
		nil,
		nil,
		requestArea,
	)

	g.mainWindow.SetContent(content)
	g.mainWindow.Resize(fyne.NewSize(1200, 700))

	// 启动定期刷新
	go g.refreshLoop()
}

// createConfigArea 创建配置区域
func (g *GUI) createConfigArea() fyne.CanvasObject {
	// 代理端口设置
	portLabel := widget.NewLabel("代理端口:")
	portEntry := widget.NewEntry()
	portEntry.SetText(strconv.Itoa(g.config.GetProxyPort()))
	portEntry.OnChanged = func(s string) {
		if port, err := strconv.Atoi(s); err == nil && port > 0 && port < 65536 {
			if g.onProxyPortChange != nil {
				g.onProxyPortChange(port)
			}
		}
	}

	// HTTPS解密开关
	httpsCheck := widget.NewCheck("启用HTTPS解密", func(checked bool) {
		if g.onHTTPSDecryptChange != nil {
			g.onHTTPSDecryptChange(checked)
		}
	})
	httpsCheck.SetChecked(g.config.GetHTTPSDecrypt())

	// 重新安装证书按钮
	reinstallCertBtn := widget.NewButton("重新安装证书", func() {
		if g.onReinstallCert != nil {
			g.onReinstallCert()
		}
	})

	// 清空请求按钮
	clearBtn := widget.NewButton("清空请求", func() {
		if g.onClearRequests != nil {
			g.onClearRequests()
		}
		g.filteredRequests = make([]*storage.CapturedRequest, 0)
		g.requestList.Refresh()
	})

	return container.NewVBox(
		container.NewHBox(
			portLabel, portEntry,
			layout.NewSpacer(),
			httpsCheck,
			reinstallCertBtn,
			clearBtn,
		),
	)
}

// createSearchArea 创建搜索区域
func (g *GUI) createSearchArea() fyne.CanvasObject {
	g.searchEntry = widget.NewEntry()
	g.searchEntry.SetPlaceHolder("搜索: URL、Host、Method、Headers、Body...")
	g.searchEntry.OnChanged = func(keyword string) {
		g.filterRequests(keyword)
	}

	return container.NewBorder(nil, nil, widget.NewLabel("🔍"), nil, g.searchEntry)
}

// createRequestArea 创建请求列表区域
func (g *GUI) createRequestArea() fyne.CanvasObject {
	// 创建表格
	g.requestList = widget.NewTable(
		func() (int, int) {
			return len(g.filteredRequests), 7 // 行数, 列数
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row >= len(g.filteredRequests) {
				label.SetText("")
				return
			}

			req := g.filteredRequests[id.Row]

			switch id.Col {
			case 0: // ID
				label.SetText(fmt.Sprintf("%d", req.ID))
			case 1: // Method
				label.SetText(req.Method)
			case 2: // Host
				label.SetText(req.Host)
			case 3: // Path
				label.SetText(req.Path)
			case 4: // Status
				label.SetText(fmt.Sprintf("%d", req.StatusCode))
			case 5: // Protocol
				label.SetText(req.Protocol)
			case 6: // Time
				label.SetText(req.Timestamp.Format("15:04:05"))
			}
		},
	)

	// 设置列宽
	g.requestList.SetColumnWidth(0, 50)  // ID
	g.requestList.SetColumnWidth(1, 80)  // Method
	g.requestList.SetColumnWidth(2, 200) // Host
	g.requestList.SetColumnWidth(3, 300) // Path
	g.requestList.SetColumnWidth(4, 70)  // Status
	g.requestList.SetColumnWidth(5, 80)  // Protocol
	g.requestList.SetColumnWidth(6, 80)  // Time

	// 添加表头
	headerLabels := []string{"ID", "Method", "Host", "Path", "Status", "Protocol", "Time"}
	header := container.NewHBox()
	for i, text := range headerLabels {
		label := widget.NewLabel(text)
		label.TextStyle = fyne.TextStyle{Bold: true}
		header.Add(label)
		if i < len(headerLabels)-1 {
			header.Add(layout.NewSpacer())
		}
	}

	return container.NewBorder(header, nil, nil, nil, g.requestList)
}

// filterRequests 过滤请求（线程安全版本）
func (g *GUI) filterRequests(keyword string) {
	var requests []*storage.CapturedRequest
	if keyword == "" {
		requests = g.storage.GetAll()
	} else {
		requests = g.storage.Search(keyword)
	}
	
	// 在主线程更新 GUI
	if g.app != nil {
		g.app.Driver().DoEventually(func() {
			g.filteredRequests = requests
			g.requestList.Refresh()
			g.updateStatus()
		})
	} else {
		// 如果 app 还未初始化，直接更新
		g.filteredRequests = requests
		g.requestList.Refresh()
		g.updateStatus()
	}
}

// refreshLoop 定期刷新界面
func (g *GUI) refreshLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		// 使用 Fyne 的主线程来更新 GUI
		keyword := g.searchEntry.Text
		
		// 在后台获取数据
		if keyword == "" {
			requests := g.storage.GetAll()
			// 在主线程更新 GUI
			g.updateRequestsInMainThread(requests)
		} else {
			requests := g.storage.Search(keyword)
			// 在主线程更新 GUI
			g.updateRequestsInMainThread(requests)
		}
	}
}

// updateRequestsInMainThread 在主线程中更新请求列表
func (g *GUI) updateRequestsInMainThread(requests []*storage.CapturedRequest) {
	// 使用 Fyne 的 Do 确保在主线程中执行
	if g.app != nil {
		g.app.Driver().DoEventually(func() {
			g.filteredRequests = requests
			g.requestList.Refresh()
			g.updateStatus()
		})
	}
}

// updateStatus 更新状态栏
func (g *GUI) updateStatus() {
	total := g.storage.Count()
	filtered := len(g.filteredRequests)

	if filtered == total {
		g.statusLabel.SetText(fmt.Sprintf("总请求数: %d", total))
	} else {
		g.statusLabel.SetText(fmt.Sprintf("显示: %d / 总数: %d", filtered, total))
	}
}

// Show 显示窗口
func (g *GUI) Show() {
	g.mainWindow.Show()
}

// Run 运行GUI
func (g *GUI) Run() {
	g.mainWindow.ShowAndRun()
}

// SetOnProxyPortChange 设置端口变化回调
func (g *GUI) SetOnProxyPortChange(callback func(int)) {
	g.onProxyPortChange = callback
}

// SetOnHTTPSDecryptChange 设置HTTPS解密变化回调
func (g *GUI) SetOnHTTPSDecryptChange(callback func(bool)) {
	g.onHTTPSDecryptChange = callback
}

// SetOnReinstallCert 设置重新安装证书回调
func (g *GUI) SetOnReinstallCert(callback func()) {
	g.onReinstallCert = callback
}

// SetOnClearRequests 设置清空请求回调
func (g *GUI) SetOnClearRequests(callback func()) {
	g.onClearRequests = callback
}
