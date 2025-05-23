/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2025-05-23 20:18:06
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-23 22:58:15
 * @FilePath: \gui\example\alarm\page_me.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"log"

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

	root.AddWidget(gui.NewTextField(w, "首页地址").SetBlur(func(text string) {
		log.Println("首页地址()", text)
	}))

	root.AddWidget(gui.NewTextField(w, "代理服务器").SetBlur(func(text string) {
		log.Println("代理服务器()", text)
	}))

	root.AddWidget(gui.NewSplit())

	//消息通知配置
	AddMsgWidgets(w, root)

	//设置主界面
	m.SetContent(root)
}

// func NewEdit(w *gui.Window, editHint, btnTitle string, f func(string)) (*gui.Container, *gui.TextField) {
// 	edit := gui.NewTextField(w, editHint)
// 	btn := gui.NewButton(w, btnTitle, func() {
// 		f(edit.Text())
// 	}).UniformInset(boss.offset)
// 	content := gui.NewContainer(w).SetAlignment(layout.Baseline) //从底部对齐
// 	content.AddWidget(edit, 1)
// 	content.AddWidget(btn)
// 	return content, edit
// }

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
	parent.AddWidget(gui.NewTextField(w, "服务器地址").SetBlur(func(text string) {
		log.Println("服务器地址()", text)
	}))
	parent.AddWidget(gui.NewTextField(w, "服务器账号").SetBlur(func(text string) {
		log.Println("服务器账号()", text)
	}))
	parent.AddWidget(gui.NewTextField(w, "服务器密码").SetMask('*').SetBlur(func(text string) {
		log.Println("服务器密码()", text)
	}))
	parent.AddWidget(gui.NewTextField(w, "订阅主题"))
	parent.AddWidget(gui.NewTextField(w, "订阅消费组").SetBlur(func(text string) {
		log.Println("订阅消费组()", text)
	}))
	parent.AddWidget(gui.NewTextField(w, "订阅消费者").SetBlur(func(text string) {
		log.Println("订阅消费者()", text)
	}))

	// content := gui.NewContainer(w).SetAlignment(layout.Baseline) //从底部对齐
	// content.AddWidget(gui.NewTextField(w, "消息通知订阅组"), 1)
	// content.AddWidget(gui.NewButton(w, "应用", func() {
	// 	w.Info("测试", "订阅消息通知成功")
	// 	w.Log("应用")
	// }).UniformInset(boss.offset))
	// parent.AddWidget(content)
}
