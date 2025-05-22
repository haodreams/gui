package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tabs struct {
	list     layout.List
	tabs     []*Tab
	selected int
	Widget[Tabs]

	cb func()
	// widget.Clickable
	// f    func()
	// text string
}

type Tab struct {
	btn   widget.Clickable
	Title string
	Widget[Tab]
	*Container
}

func NewTab(win *Window, title string) *Tab {
	m := new(Tab)
	m.Title = title
	m.Widget.Init(win, m)
	m.Inset = win.TabsInset
	m.Container = NewContainer(win)
	m.SetAxis(layout.Vertical)
	return m
}

func (m *Tabs) SetCallback(cb func()) *Tabs {
	m.cb = cb
	return m
}

func NewTabs(win *Window) *Tabs {
	m := new(Tabs)
	m.Init(win, m)
	m.Inset = win.TabsInset
	return m
}

func (m *Tabs) GetSelected() int {
	return m.selected
}

func (m *Tabs) GetSelectedTab() *Tab {
	return m.tabs[m.selected]
}

func (m *Tabs) SetSelectedTab(selected int) {
	m.selected = selected
}

func (m *Tabs) GetSelectedTitle() string {
	return m.tabs[m.selected].Title
}

func (m *Tabs) SetSelectedTabByTitle(title string) {
	list := m.tabs
	for i, t := range list {
		if t.Title == title {
			m.selected = i
			return
		}
	}
}

func (m *Tabs) Size() int {
	return len(m.tabs)
}

func (m *Tabs) Layout(gtx layout.Context) D {
	th := m.Theme()
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			// gtx.Constraints.Min.Y = 35
			// gtx.Constraints.Max.Y = 35
			return m.list.Layout(gtx, len(m.tabs), func(gtx C, tabIdx int) D {
				t := m.tabs[tabIdx]
				if t.btn.Clicked(gtx) {
					m.selected = tabIdx
					if m.cb != nil {
						m.cb()
					}
				}
				var tabWidth int
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						dims := material.Clickable(gtx, &t.btn, func(gtx C) D {
							return m.Inset.Layout(gtx,
								material.Body1(th, t.Title).Layout,
							)
						})
						tabWidth = dims.Size.X
						return dims
					}),
					layout.Stacked(func(gtx C) D {
						if m.selected != tabIdx {
							return D{}
						}
						tabHeight := gtx.Dp(unit.Dp(4))
						tabRect := image.Rect(0, 0, tabWidth, tabHeight)
						paint.FillShape(gtx.Ops, th.Palette.ContrastBg, clip.Rect(tabRect).Op())
						return D{
							Size: image.Point{X: tabWidth, Y: tabHeight},
						}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return m.tabs[m.selected].Layout(gtx)
			// return m.slider.Layout(gtx, func(gtx C) D {
			// 	//fill(gtx, dynamicColor(tabs.selected), dynamicColor(tabs.selected+1))
			// 	// return layout.Center.Layout(gtx,
			// 	// 	material.H1(th, fmt.Sprintf("Tab content #%d", m.selected+1)).Layout,
			// 	// )
			// })
		}),
	)
}

func (m *Tabs) AddTab(tab *Tab) {
	m.tabs = append(m.tabs, tab)
}
