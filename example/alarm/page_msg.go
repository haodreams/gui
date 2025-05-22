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
	m.list = gui.NewList(m.Parent())
	m.list.Axis = layout.Vertical
	m.AddMsg("Info", "18:04:21", "这是一条提示信息", boss.info)
	m.AddMsg("Warn", "18:07:25", "这是一条警告信息", boss.warn)
	m.AddMsg("Error", "18:09:36", "这是一条错误信息", boss.erron)
	m.SetContent(m.list)

}

func (m *PageMsg) do(card *gui.Card) {
	log.Println("cart id:", card.ID(), card.Msg())
	m.Parent().Info(card.Title(), card.Msg())
}

func (m *PageMsg) AddMsg(title, time, msg string, icon image.Image) {
	item := gui.NewCard(m.Parent(), title, time, msg).OnClick(m.do).SetImage(icon)
	m.list.AddItem(item)
}
