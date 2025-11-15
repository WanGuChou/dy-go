package systray_manager

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/getlantern/systray"
)

// Manager 系统托盘管理器
type Manager struct {
	proxyEnabled bool
	proxyPort    int
	onEnable     func()
	onDisable    func()
	onExit       func()
}

// NewManager 创建系统托盘管理器
func NewManager(proxyPort int, onEnable, onDisable, onExit func()) *Manager {
	return &Manager{
		proxyPort: proxyPort,
		onEnable:  onEnable,
		onDisable: onDisable,
		onExit:    onExit,
	}
}

// Start 启动系统托盘
func (m *Manager) Start() {
	go systray.Run(m.onReady, m.onQuit)
}

// onReady 托盘准备就绪
func (m *Manager) onReady() {
	// 设置图标（使用默认图标，实际应该提供.ico文件）
	systray.SetIcon(getIcon())
	systray.SetTitle("Packet Capture")
	systray.SetTooltip("Packet Capture Tool")

	// 创建菜单项
	mEnable := systray.AddMenuItem("✅ 启用代理", "设置系统代理")
	mDisable := systray.AddMenuItem("❌ 禁用代理", "清除系统代理")
	mDisable.Disable() // 初始状态禁用此选项

	systray.AddSeparator()

	mShowWindow := systray.AddMenuItem("📊 显示主窗口", "打开主界面")

	systray.AddSeparator()

	mExit := systray.AddMenuItem("🚪 退出", "退出程序")

	// 处理菜单点击事件
	go func() {
		for {
			select {
			case <-mEnable.ClickedCh:
				m.enableProxy()
				mEnable.Disable()
				mDisable.Enable()
				systray.SetTooltip("Packet Capture Tool (代理已启用)")
			case <-mDisable.ClickedCh:
				m.disableProxy()
				mDisable.Disable()
				mEnable.Enable()
				systray.SetTooltip("Packet Capture Tool (代理已禁用)")
			case <-mShowWindow.ClickedCh:
				// 显示主窗口（需要从App传递过来）
				log.Println("Show main window requested")
			case <-mExit.ClickedCh:
				m.onExit()
				systray.Quit()
				return
			}
		}
	}()
}

// onQuit 托盘退出
func (m *Manager) onQuit() {
	// 清理工作
	if m.proxyEnabled {
		m.disableProxy()
	}
}

// enableProxy 启用系统代理
func (m *Manager) enableProxy() {
	proxyAddr := fmt.Sprintf("127.0.0.1:%d", m.proxyPort)

	// 使用netsh设置系统代理（Windows）
	cmd := exec.Command("reg", "add",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to enable proxy (step 1): %v", err)
		return
	}

	cmd = exec.Command("reg", "add",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v", "ProxyServer", "/t", "REG_SZ", "/d", proxyAddr, "/f")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to enable proxy (step 2): %v", err)
		return
	}

	m.proxyEnabled = true
	log.Printf("System proxy enabled: %s", proxyAddr)

	if m.onEnable != nil {
		m.onEnable()
	}
}

// disableProxy 禁用系统代理
func (m *Manager) disableProxy() {
	cmd := exec.Command("reg", "add",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to disable proxy: %v", err)
		return
	}

	m.proxyEnabled = false
	log.Println("System proxy disabled")

	if m.onDisable != nil {
		m.onDisable()
	}
}

// getIcon 获取托盘图标（简单的PNG数据）
func getIcon() []byte {
	// 这里返回一个简单的图标数据
	// 实际应该使用.ico文件或PNG数据
	// 为简化，返回空数据，systray会使用默认图标
	return []byte{}
}
