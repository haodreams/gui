package main

import "gitee.com/haodreams/gui"

type PageMain struct {
	gui.Page
}

func NewPageMain(w *gui.Window, name, title string) *PageMain {
	m := &PageMain{}
	m.Setup(w, name, title)
	return m
}

func (m *PageMain) Setup(w *gui.Window, name, title string) {
	m.Page.Setup(w, name, title)
	root := gui.NewContainer(m.Parent())
	root.AddWidget(gui.NewLabel(w, "首页内容"))
	m.SetContent(root)
}
