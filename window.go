/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 23:59:39
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 12:41:13
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
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
)

type Contenter interface {
	Layout(gtx C) D
}

// Window holds all of the application state.
type Window struct {
	Config
	// Theme is used to hold the fonts used throughout the application.
	container Contenter // 内容容器
	win       *app.Window
	prompt    *MessageBox //消息对话框
	*explorer.Explorer

	width      int                // 窗口宽度
	height     int                // 窗口高度
	statusBar                     //状态栏
	dataDir    string             // 数据目录
	Shield                        // 屏蔽层
	onClose    func()             // 关闭窗口事件
	onWinEvent func(focused bool) // 出现窗口事件
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
	m.Shield.SetColor(m.theme.Bg)                     // 设置屏蔽层颜色

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

func (m *Window) SetOnWinEvent(cb func(bool)) {
	m.onWinEvent = cb
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
		layout.Expanded(func(gtx C) D {
			return m.Shield.Layout(gtx)
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
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			m.layout(gtx)
			e.Frame(gtx.Ops)
		case app.ConfigEvent:
			//检查窗口是否有焦点
			if m.onWinEvent != nil {
				m.onWinEvent(e.Config.Focused)
			}
		case app.DestroyEvent:
			if m.onClose != nil {
				m.onClose()
			}
			return e.Err
		}
	}
}

// 运行窗体
func Run(init func() *Window) {
	go func() {
		win := init()
		err := win.Run()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}
