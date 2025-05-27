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
	content Contenter
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

func (m *Form)Setup(win *Window)*Form{
	m.Init(win, m)
	m.color = win.theme.Bg
	return m
}

func (m *Form)SetBottomOffset(offset unit.Dp) *Form{
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

	maxWidth := unit.Dp(m.MaxWidth)
	width:= gtx.Metric.PxToDp(m.win.width)
	if maxWidth > 0 {
		offset.Left = (unit.Dp(width-maxWidth) / 2)
		if offset.Left < 0 {
			offset.Left = 0
		} else {
			x := offset.Left + maxWidth
			if x > width{
				x = width
			}
			if gtx.Constraints.Max.X > int(x) {
				gtx.Constraints.Max.X = int(x)
			}
		}
	}

	//高度不需要进行转换， 高度是已经计算后的数据
	maxHeight := gtx.Dp(unit.Dp(m.MaxHeight))
	if maxHeight > 0 {
		offset.Top = unit.Dp(m.win.height-maxHeight) / unit.Dp(gtx.Metric.PxPerDp)
		if offset.Top < 0 {
			offset.Top = 0
		} else {
			y := int(gtx.Dp(offset.Top)) + maxHeight
			if y > m.win.height{
				y = m.win.height
			}
			y -= int(gtx.Dp(m.bottomOffset))
			gtx.Constraints.Max.Y = y 
			// if gtx.Constraints.Max.Y > y {
			// }
		}
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
