/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 23:59:39
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-23 18:22:38
 * @FilePath: \dataviewe:\go\src\gitee.com\haodreams\gui\windows.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
)

type Contenter interface {
	Layout(gtx C) D
}
type Config struct {
	theme          *material.Theme
	ButtonInset    layout.Inset
	CheckboxInset  layout.Inset
	EditInset      layout.Inset
	LabelInset     layout.Inset
	TextFieldInset layout.Inset
	SelectInset    layout.Inset
	TabsInset      layout.Inset
	insetInited    bool // 初始化过Inset
	showStatusBar  bool
	plugin         func(*app.Window) event.Event
	logf           func(v ...any)
}

type Option func(*Config)

// NewConfig returns a new Config with default values.
func Inset(v unit.Dp) Option {
	return func(m *Config) {
		m.ButtonInset = layout.UniformInset(v)
		m.CheckboxInset = layout.UniformInset(v)
		m.EditInset = layout.UniformInset(v)
		m.LabelInset = layout.UniformInset(v)
		m.TextFieldInset = layout.UniformInset(v)
		m.SelectInset = layout.UniformInset(v)
		m.TabsInset = layout.UniformInset(v)
		m.TabsInset.Right += 10
		m.TabsInset.Left += 10
		m.insetInited = true
	}
}

// 设置主题
func ThemeShaper(ts *text.Shaper) Option {
	return func(m *Config) {
		m.theme.Shaper = ts
	}
}

// 隐藏状态栏
func HideStatusBar() Option {
	return func(m *Config) {
		m.showStatusBar = false
	}
}

func WithLog(f func(v ...any)) Option {
	return func(m *Config) {
		m.logf = f
	}
}

func WithPlugin(plugin func(*app.Window) event.Event) Option {
	return func(m *Config) {
		m.plugin = plugin
	}
}

// Window holds all of the application state.
type Window struct {
	Config
	// Theme is used to hold the fonts used throughout the application.
	container Contenter // 内容容器
	win       *app.Window
	prompt    *MessageBox //消息对话框
	*explorer.Explorer

	width     int    // 窗口宽度
	height    int    // 窗口高度
	statusBar        //状态栏
	dataDir   string // 数据目录

	onClose func() // 关闭窗口事件
}

// NewUI creates a new UI using the Go Fonts.
func NewWindow(opts ...Option) *Window {
	m := &Window{}
	m.win = new(app.Window)
	// 初始化边距
	m.theme = material.NewTheme() // 创建 Material Theme

	for _, opt := range opts {
		opt(&m.Config)
	}

	// 初始化Inset
	if !m.insetInited {
		Inset(8)(&m.Config)
	}

	// Load Go Fonts to theme
	if m.theme.Shaper == nil {
		m.theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection())) // 加载Go Fonts字库
	}
	m.prompt = NewMessage(m, "", "", ModalInfo)       // 创建消息对话框
	m.statusBar.Init(m).show = m.Config.showStatusBar // 初始化状态栏, 显示状态栏
	m.Explorer = explorer.NewExplorer(m.win)          // 初始化文件浏览器

	var err error
	m.dataDir, err = app.DataDir() //
	if err != nil {
		m.Log("Init data dir error,", err)
		m.dataDir = "./"
	}
	return m
}

// SetContent sets the widget that will be displayed in the window.
func (m *Window) Option(opts ...app.Option) {
	m.win.Option(opts...)
}

func (m *Window) DataDir() string {
	return m.dataDir
}

// 设置关闭时的操作
func (m *Window) SetOnClose(cb func()) {
	m.onClose = cb
}

// 重新设置新的数据目录文件
func (m *Window) SetDataDir(dir string) {
	m.dataDir = strings.TrimSuffix(dir, "/")
}

func (m *Window) Log(v ...any) {
	if m.logf != nil {
		m.logf(v...)
	}
}

// Invalidate forces the application to repaint.  Note that this should only be used for performance debugging or very specific situations,
// as it can lead to unexpected behavior in a production environment.  If you're not sure what you're doing, it's probably best to use
func (m *Window) Invalidate() {
	m.win.Invalidate()
}

// 获取窗口句
func (m *Window) Owner() *app.Window {
	return m.win
}

// 获取窗口主题
func (m *Window) Theme() *material.Theme {
	return m.theme
}

// 设置窗口标题
func (m *Window) SetTitle(title string) {
	m.win.Option(app.Title(title))
}

func (m *Window) Info(title string, contant string) {
	m.MessageBox(title, contant, ModalInfo)
}
func (m *Window) Warn(title string, contant string) {
	m.MessageBox(title, contant, ModalWarn)
}
func (m *Window) Error(title string, contant string) {
	m.MessageBox(title, contant, ModalError)
}

// 消息对话框
func (m *Window) MessageBox(title string, contant string, mode ModelType) *MessageBox {
	m.prompt.content = nil
	m.prompt.Title = title
	m.prompt.Message = contant
	m.prompt.Type = mode
	m.prompt.SetOptions([]MBOption{{Text: "确定"}}...)
	m.prompt.SetOnSubmit(func(selectedOption string, remember bool) {
		if selectedOption == "确定" {
			m.prompt.Hide()
			return
		}
	})
	m.prompt.Show()
	return m.prompt
}

// 消息对话框
func (m *Window) MsgBox(title string, contant string, mode ModelType, cb func(selectedOption string, remember bool), options ...MBOption) *MessageBox {
	m.prompt.content = nil
	m.prompt.Title = title
	m.prompt.Message = contant
	m.prompt.Type = mode
	if len(options) > 0 {
		m.prompt.SetOptions(options...)
	} else {
		m.prompt.SetOptions([]MBOption{{Text: "确定"}}...)
	}
	if cb != nil {
		m.prompt.SetOnSubmit(cb)
	} else {
		m.prompt.SetOnSubmit(func(selectedOption string, remember bool) {
			if selectedOption == "确定" {
				m.prompt.Hide()
				return
			}
		})
	}

	m.prompt.Show()
	return m.prompt
}

// 输入对话框
func (m *Window) Prompt(title string, message string, f func(text string)) *MessageBox {
	m.prompt.Title = title
	edit := NewEdit(m, message).SetWidth(m.prompt.MaxWidth - 20)
	edit.SetOnSubmit(f)
	m.prompt.content = edit
	m.prompt.Type = ModalInfo

	m.prompt.SetOptions([]MBOption{{Text: "确定"}, {Text: "取消"}}...)
	m.prompt.SetOnSubmit(func(selectedOption string, remember bool) {
		m.prompt.Hide()
		if selectedOption == "确定" {
			if f != nil {
				f(edit.Text())
			}
			return
		}
	})
	m.prompt.Show()
	return m.prompt
}

// 自定义对话框ModalInfo
func (m *Window) ShowDialog(title string, content Contenter, f func(text string), options ...MBOption) *MessageBox {
	m.prompt.Title = title
	m.prompt.Type = ModalInfo
	m.prompt.content = content
	m.prompt.SetOnSubmit(func(selectedOption string, remember bool) {
		m.prompt.Hide()
		f(selectedOption)
	})

	if len(options) > 0 {
		m.prompt.SetOptions(options...)
	} else {
		m.prompt.SetOptions([]MBOption{{Text: "确定"}}...)
	}
	m.prompt.Show()
	return m.prompt
}

// 对话框恢复上一次的大小
func (m *Window) RestoreDialogWidth() {
	m.prompt.RestoreWidth()
}

// 设置窗体布局
func (m *Window) layout(gtx C) D {
	m.width = gtx.Constraints.Max.X
	m.height = gtx.Constraints.Max.Y

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, //垂直布局，才能显示状态栏
				layout.Flexed(1, func(gtx C) D { //同级别才会计算全部空间
					return layout.Flex{}.Layout(gtx, //水平布局才能扩展全部空间
						layout.Flexed(1, func(gtx C) D {
							if m.container == nil {
								return D{}
							}
							return layout.UniformInset(2).Layout(gtx, func(gtx C) D {
								return m.container.Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if m.statusBar.show {
						return m.statusBar.Layout(gtx)
					}
					return D{}
				}),
			)
		}),

		//显示对话框
		layout.Expanded(func(gtx C) D {
			return m.prompt.Layout(gtx, m.width, m.height)
		}),
	)
}

// 设置窗体内容
func (m *Window) SetContent(layout Contenter) {
	m.container = layout
}

// Run handles window events and renders the application.
func (m *Window) Run() error {
	var ops op.Ops
	// listen for events happening on the window.
	for {
		var et event.Event
		if m.plugin != nil {
			et = m.plugin(m.win)
		} else {
			et = m.win.Event()
		}
		m.Explorer.ListenEvents(et)
		// detect the type of the event.
		switch e := et.(type) {
		// this is sent when the application should re-render.
		case app.FrameEvent:
			// gtx is used to pass around rendering and event information.
			gtx := app.NewContext(&ops, e)
			m.layout(gtx)
			//m.container.Layout(gtx)
			// render and handle the operations from the UI.
			e.Frame(gtx.Ops)

		// this is sent when the application is closed.
		case app.DestroyEvent:
			return e.Err
		}
	}
}

// 运行窗体
func Run(init func() *Window) {
	go func() {
		win := init()
		err := win.Run()
		if win.onClose != nil {
			win.onClose()
		}
		if err != nil {
			win.Log(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}
