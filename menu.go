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

func (m *Menu) Clicked(gtx C) {
	for _, item := range m.Items {
		if item.Clicked(gtx) {
			if item.Action != nil {
				item.Action()
			}
		}
	}
}
