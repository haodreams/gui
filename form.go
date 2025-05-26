package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/x/component"
)

// 表单
type Form struct {
	widget.Clickable
	color color.NRGBA
	show  bool
	Widget[Form]
	layout.Direction
	content Contenter
}

func (m *Form) Hide() *Form {
	m.show = false
	return m
}
func (m *Form) Show() *Form {
	m.show = true
	return m
}

func NewForm(win *Window) *Form {
	m := &Form{}
	m.Init(win, m)
	m.color = win.theme.Bg
	return m
}

func (m *Form) SetDirection(direction layout.Direction) *Form {
	m.Direction = direction
	return m
}

func (m *Form) SetContent(c Contenter) *Form {
	m.content = c
	return m
}

func (m *Form) SetColor(c color.NRGBA) *Form {
	m.color = c
	return m
}

// Layout draws the scrim using the provided animation. If the animation indicates
// that the scrim is not visible, this is a no-op.
func (m *Form) Layout(gtx layout.Context) D {
	if !m.show {
		return D{}
	}
	if m.content == nil {
		return D{}
	}
	m.CheckDimensions(&gtx)

	m.Clickable.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min = gtx.Constraints.Max
		component.Rect{Color: m.color, Size: gtx.Constraints.Max}.Layout(gtx)
		return D{Size: gtx.Constraints.Max}
	})

	return m.content.Layout(gtx)
}
