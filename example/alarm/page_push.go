package main

import (
	"fmt"
	"strings"
	"time"

	"gitee.com/haodreams/gui"
)

type PagePush struct {
	gui.Page
}

func NewPagePush(w *gui.Window, name, title string) *PagePush {
	m := &PagePush{}
	m.Setup(w, name, title)
	return m
}

func (m *PagePush) Setup(w *gui.Window, name, title string) {
	m.Page.Setup(w, name, title)
	root := gui.NewContainer(m.Parent()).SetVertical()
	root.AddWidget(gui.NewLabel(w, "消息模拟发送测试:"))
	editTitle := gui.NewTextField(w, "消息标题")
	var edit *gui.TextField
	var conta *gui.Container
	conta, edit = NewEdit(w, "消息内容", "测试", func(s string) {
		title := editTitle.Text()
		title = strings.TrimSpace(title)
		if title == "" {
			w.Error("输入错误", "标题不能为空")
			return
		}
		msg := edit.Text()
		msg = strings.TrimSpace(msg)
		if msg == "" {
			w.Error("输入错误", "输入内容不能为空")
			return
		}
		now := time.Now().Format("15:04:05")

		boss.msg.AddMsg(title, now, msg, boss.info)
		boss.win.Info("提示", "添加完成")
	})

	root.AddWidget(editTitle)
	root.AddWidget(conta)

	root.AddWidget(gui.NewButton(w, "获取数据目录", func() {
		w.Info("数据目录", w.DataDir())
	}))
	root.AddWidget(gui.NewButton(w, "获取时区", func() {
		now := time.Now()
		name, offset := now.Zone()
		w.Info("time zone", fmt.Sprintf("%s %d", name, offset))
	}))

	root.AddWidget(gui.NewButton(w, "设置时区", func() {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			w.Error("错误", err.Error())
			return
		}
		time.Local = loc
		w.Info("设置时区", "设置成功")
	}))

	m.SetContent(root)
}
