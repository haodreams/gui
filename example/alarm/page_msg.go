/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2025-05-23 20:18:06
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 01:43:00
 * @FilePath: \gui\example\alarm\page_msg.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"image"
	"log"

	"gioui.org/layout"
	"gitee.com/haodreams/gui"
)

type PageMsg struct {
	gui.Page
	list *gui.List
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

	m.list = gui.NewList(m.Parent())
	m.list.Axis = layout.Vertical
	m.AddMsg("Info", "18:04:21", "这是一条提示信息", boss.info)
	m.AddMsg("Warn", "18:07:25", "这是一条警告信息", boss.warn)
	m.AddMsg("Error", "18:09:36", "这是一条错误信息", boss.erron)

	content.AddWidget(m.list)
	m.SetContent(content)
}

func (m *PageMsg) do(card *gui.Card) {
	log.Println("cart id:", card.ID(), card.Msg())
	m.Parent().Info(card.Title(), card.Msg())
}

func (m *PageMsg) AddMsg(title, time, msg string, icon image.Image) {
	item := gui.NewCard(m.Parent(), title, time, msg).OnClick(m.do).SetImage(icon)
	m.list.AddItem(item)
}
