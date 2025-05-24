package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/x/component"
)

// 屏蔽层
type Shield struct {
	widget.Clickable
	color   color.NRGBA
	hide    bool
	content Contenter
}

func (m *Shield) Hide() *Shield {
	m.hide = true
	return m
}
func (m *Shield) Show() *Shield {
	m.hide = false
	return m
}

func NewShield(win *Window) *Shield {
	m := &Shield{}
	m.color = win.theme.Bg
	return m
}

func (m *Shield) SetContent(c Contenter) *Shield {
	m.content = c
	return m
}

func (m *Shield) SetColor(c color.NRGBA) *Shield {
	m.color = c
	return m
}

// Layout draws the scrim using the provided animation. If the animation indicates
// that the scrim is not visible, this is a no-op.
func (m *Shield) Layout(gtx layout.Context) D {
	if m.hide {
		return D{}
	}
	if m.content == nil {
		return D{}
	}
	m.Clickable.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min = gtx.Constraints.Max
		component.Rect{Color: m.color, Size: gtx.Constraints.Max}.Layout(gtx)
		return D{Size: gtx.Constraints.Max}
	})

	return m.content.Layout(gtx)
}
