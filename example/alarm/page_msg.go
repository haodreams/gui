/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2025-05-23 20:18:06
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 10:06:12
 * @FilePath: \gui\example\alarm\page_msg.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"log"

	"gioui.org/layout"
	"gitee.com/haodreams/gui"
)

type PageMsg struct {
	gui.Page
	list *gui.FixList
}

func NewPageMsg(w *gui.Window, name, title string) *PageMsg {
	m := &PageMsg{}
	m.Setup(w, name, title)
	return m
}

func (m *PageMsg) Setup(w *gui.Window, name, title string) {
	m.Page.Setup(w, name, title)
	content := gui.NewContainer(w).SetVertical()

	btnSearch := gui.NewButton(w, "", func() {
		log.Println("search")
	}).SetIcon(gui.SearchIcon).SetColor(w.Theme().Fg).SetBackground(w.Theme().Bg)
	btnDelete := gui.NewButton(w, "", func() {
		log.Println("delete")
	}).SetIcon(gui.DeleteIcon).SetColor(w.Theme().Fg).SetBackground(w.Theme().Bg)

	line := gui.NewContainer(w)

	line.Add(gui.NewSpace())
	line.AddWidget(btnSearch)
	line.AddWidget(btnDelete)
	content.AddWidget(line)

	content.Add(gui.NewSpace(5, 5))
	content.AddWidget(gui.NewSplit().SetHeight(2)) //加上分割线

	m.list = gui.NewFixList(m.Parent(), 500).SetDesc(true)
	m.list.Axis = layout.Vertical

	content.AddWidget(m.list)
	m.SetContent(content)
}

