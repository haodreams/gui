package main

import (
	"gioui.org/layout"
	"gitee.com/haodreams/gui"
)

func NewLoginDialog(w *gui.Window) gui.Contenter {
	content := gui.NewContainer(w).SetVertical().Resize(300, 300)
	content.AddWidget(gui.NewLabel(w, "用户名:"))
	content.AddWidget(gui.NewEdit(w, "请输入用户名"))
	content.AddWidget(gui.NewLabel(w, "密码:"))
	content.AddWidget(gui.NewEdit(w, "请输入密码").SetMask('*'))
	content.Add(gui.NewSpace(10, 10))
	content.AddWidget(gui.NewButton(w, "登录", func() {
		w.Shield.Hide()
	}))
	content.Add(gui.NewSpace(10, 10))

	root := gui.NewContainer(w)
	root.AddWidget(gui.NewDirection(layout.Center, content.Layout), 1) //垂直居中
	return root
}
