package gui

import (
	"time"

	"gioui.org/layout"
	"gioui.org/x/component"
)

type NaviSide struct {
	win   *Window
	ms    *component.ModalSheet
	va    *component.VisibilityAnimation
	width int
	Widget[NaviSide]
	content layout.Widget
	mode    int //0 扩展模式 1 附加模式
}

func (m *NaviSide) Init(win *Window) {
	m.win = win
	m.width = 300
	modal := component.NewModal()
	m.ms = component.NewModalSheet(modal)
	m.va = &component.VisibilityAnimation{
		State:    component.Invisible,
		Duration: time.Millisecond * 250,
	}
	m.Widget.Init(win, m)
}

func (m *NaviSide) SetNaviMode(mode int) {
	m.mode = mode
}

func (m *NaviSide) SetNaviWidth(width int) {
	m.width = width
	m.va.Appear(time.Now())
}

func (m *NaviSide) SetNaviContent(content layout.Widget) {
	m.content = content
	m.va.Appear(time.Now())
}

func (m *NaviSide) ShowNavi(b bool) {
	if b {
		m.va.Appear(time.Now())
	} else {
		m.va.Disappear(time.Now())
	}
}

func (m *NaviSide) Layout(gtx layout.Context) layout.Dimensions {
	if m.content == nil {
		return layout.Dimensions{}
	}
	if m.width > 0 {
		if gtx.Constraints.Max.X > m.width {
			gtx.Constraints.Max.X = m.width
		}
	}
	return m.ms.Layout(gtx, m.Theme(), m.va, func(gtx layout.Context) layout.Dimensions {
		return m.content(gtx)
	})
}
