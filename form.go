package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/component"
)

// 表单
type Form struct {
	widget.Clickable
	color color.NRGBA
	show  bool
	Widget[Form]
	content      Contenter
	bottomOffset unit.Dp
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
	return m.Setup(win)
}

func (m *Form) Setup(win *Window) *Form {
	m.Init(win, m)
	m.color = win.theme.Bg
	return m
}

func (m *Form) SetBottomOffset(offset unit.Dp) *Form {
	m.bottomOffset = offset
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

	offset := layout.Inset{}

	//gtx.Constraints.Max.X 存放的是像素
	if m.Width > 0 {
		maxWidth := gtx.Dp(unit.Dp(m.Width))         //转为像素
		offset.Left = unit.Dp(m.win.width-maxWidth) / 2 // (unit.Dp(gtx.Constraints.Max.X-maxWidth) / 2) / unit.Dp(gtx.Metric.PxPerDp)
		if offset.Left < 0 {
			offset.Left = 0
		} else {
			x := int(offset.Left) + maxWidth
			if x > m.win.width {
				x = m.win.width
			}
			gtx.Constraints.Max.X = x
		}
		offset.Left = gtx.Metric.PxToDp(int(offset.Left))
	}

	//gtx.Constraints.Max.Y 存放的是像素
	if m.Height > 0 {
		maxHeight := gtx.Dp(unit.Dp(m.Height))
		offset.Top = unit.Dp(m.win.height - maxHeight)
		if offset.Top < 0 {
			offset.Top = 0
		} else {
			y := int(offset.Top) + maxHeight
			if y > m.win.height {
				y = m.win.height
			}
			y -= int(gtx.Dp(m.bottomOffset))
			gtx.Constraints.Max.Y = y
		}
		offset.Top = gtx.Metric.PxToDp(int(offset.Top))
	}

	offset.Layout(gtx, func(gtx C) D {
		return m.Clickable.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min = gtx.Constraints.Max
			component.Rect{Color: m.color, Size: gtx.Constraints.Max}.Layout(gtx)
			return D{Size: gtx.Constraints.Max}
		})
	})

	return offset.Layout(gtx, m.content.Layout)
}
