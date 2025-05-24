/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2025-05-23 20:18:06
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-23 23:24:15
 * @FilePath: \gui\example\alarm\page_push.go
 * @Description:
 */
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
	editTitle := gui.NewEdit(w, "消息标题")
	editMsg := gui.NewEdit(w, "消息内容")

	root.AddWidget(editTitle)
	root.AddWidget(editMsg)
	level := gui.NewRadioGroup(w, []string{"info", "warn", "error"},
		[]string{"提示 ", "告警 ", "错误 "}).SetValue("info")

	btns := gui.NewContainer(w)
	btns.AddWidget(level)
	btns.Add(gui.NewSpace(10, 10))
	btns.AddWidget(gui.NewButton(w, "提交", func() {
		title := editTitle.Text()
		title = strings.TrimSpace(title)
		if title == "" {
			w.Error("输入错误", "标题不能为空")
			return
		}
		msg := editMsg.Text()
		msg = strings.TrimSpace(msg)
		if msg == "" {
			w.Error("输入错误", "输入内容不能为空")
			return
		}
		now := time.Now().Unix()
		boss.PushMsg("测试", now, level.Value, title, msg)
		boss.win.Info("提示", "添加完成")
	}))
	root.AddWidget(btns)

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
