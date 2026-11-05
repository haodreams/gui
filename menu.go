package gui

import (
	"gioui.org/widget"
	"gioui.org/x/component"
)

type MenuItem struct {
	Label  string
	Icon   *widget.Icon
	Action func()
	widget.Clickable
}
type Menu struct {
	win   *Window
	Items []*MenuItem
	component.MenuState
}

func NewMenu(win *Window) *Menu {
	m := new(Menu)
	m.win = win
	return m
}

func (m *Menu) AddItem(label string, icon *widget.Icon, action func()) {
	item := &MenuItem{Label: label, Icon: icon, Action: action}
	it := component.MenuItem(m.win.Theme(), &item.Clickable, label)
	it.Icon = icon
	m.Options = append(m.Options, func(gtx C) D {
		return it.Layout(gtx)
	})
	m.Items = append(m.Items, item)
}

// Clicked 处理菜单项点击, 返回是否有菜单项被点击
func (m *Menu) Clicked(gtx C) bool {
	clicked := false
	for _, item := range m.Items {
		if item.Clicked(gtx) {
			clicked = true
			if item.Action != nil {
				item.Action()
			}
		}
	}
	return clicked
}

// Layout 显示菜单
func (m *Menu) Layout(gtx C) D {
	return component.Menu(m.win.Theme(), &m.MenuState).Layout(gtx)
}
