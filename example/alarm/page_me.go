package main

import (
	"gioui.org/layout"
	"gitee.com/haodreams/gui"
)

type PageMe struct {
	gui.Page
}

func NewPageMe(w *gui.Window, name, title string) *PageMe {
	m := &PageMe{}
	m.Setup(w, name, title)
	return m
}

func (m *PageMe) Setup(w *gui.Window, name, title string) {
	m.Page.Setup(w, name, title)
	root := gui.NewContainer(m.Parent()).SetVertical()
	conta, _ := NewEdit(w, "首页地址", "应用", func(s string) {
		w.Info("测试", s)
		w.Log(s)
	})
	root.AddWidget(conta)
	conta, _ = NewEdit(w, "代理服务器", "应用", func(s string) {
		w.Info("测试", s)
		w.Log(s)
	})

	root.AddWidget(conta)

	//消息通知配置
	AddMsgWidgets(w, root)

	//设置主界面
	m.SetContent(root)
}

func NewEdit(w *gui.Window, editHint, btnTitle string, f func(string)) (*gui.Container, *gui.TextField) {
	edit := gui.NewTextField(w, editHint)
	btn := gui.NewButton(w, btnTitle, func() {
		f(edit.Text())
	}).UniformInset(boss.offset)
	content := gui.NewContainer(w).SetAlignment(layout.Baseline) //从底部对齐
	content.AddWidget(edit, 1)
	content.AddWidget(btn)
	return content, edit
}

func NewButton(w *gui.Window, title string, f func()) *gui.Container {
	btn := gui.NewButton(w, title, f)
	content := gui.NewContainer(w)
	content.Add(gui.NewSpace())
	content.AddWidget(btn)
	return content
}

func AddMsgWidgets(w *gui.Window, parent *gui.Container) {
	parent.Add(gui.NewSpace(10, 30))
	parent.AddWidget(gui.NewLabel(w, "消息通知配置:"))
	parent.AddWidget(gui.NewTextField(w, "消息通知服务器地址"))
	parent.AddWidget(gui.NewTextField(w, "消息通知服务器账号"))
	parent.AddWidget(gui.NewTextField(w, "消息通知服务器密码").SetMask('*'))
	parent.AddWidget(gui.NewTextField(w, "消息通知订阅主题"))
	parent.AddWidget(gui.NewTextField(w, "消息通知订阅组"))

	content := gui.NewContainer(w).SetAlignment(layout.Baseline) //从底部对齐
	content.AddWidget(gui.NewTextField(w, "消息通知订阅组"), 1)
	content.AddWidget(gui.NewButton(w, "应用", func() {
		w.Info("测试", "订阅消息通知成功")
		w.Log("应用")
	}).UniformInset(boss.offset))
	parent.AddWidget(content)
}
