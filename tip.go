package gui

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tip struct {
	Message string
	Visible bool

	Widget[Tip]
	lastWidth unit.Dp
}

func NewTip(win *Window, msg string) *Tip {
	m := &Tip{
		Message: msg,
	}
	m.Init(win, m)
	m.Width = 400
	m.lastWidth = m.Width
	m.Height = 35
	m.win = win

	return m
}
func (m *Tip) RestoreWidth() *Tip {
	m.Width = m.lastWidth
	return m
}

func (m *Tip) Show() {
	m.Visible = true
}

func (m *Tip) Hide() {
	m.Visible = false
}

func (m *Tip) IsVisible() bool {
	return m.Visible
}

func (m *Tip) Tip(msg string, msDelay ...int64) {
	m.Message = msg
	delay := int64(2000)
	if len(msDelay) > 0 {
		delay = msDelay[0]
	}
	m.Show()
	go func() {
		time.Sleep(time.Millisecond * time.Duration(delay))
		m.Hide()
	}()
}

func (m *Tip) Layout(gtx layout.Context) D {
	if !m.Visible {
		return D{}
	}

	offset := layout.Inset{}

	maxWidth := gtx.Dp(unit.Dp(m.Width))
	if maxWidth > 0 {
		offset.Left = (unit.Dp(gtx.Constraints.Max.X-maxWidth) / 2) / unit.Dp(gtx.Metric.PxPerDp)
		if offset.Left < 0 {
			offset.Left = 0
		} else {
			x := int(gtx.Dp(offset.Left)) + maxWidth
			if gtx.Constraints.Max.X > x {
				gtx.Constraints.Max.X = x
			}
		}
	}

	// 高度不需要进行转换， 高度是已经计算后的数据
	maxHeight := gtx.Dp(unit.Dp(m.Height))
	if maxHeight > 0 {
		offset.Top = unit.Dp(45)
		if offset.Top < 0 {
			offset.Top = 0
		} else {
			y := int(gtx.Dp(offset.Top)) + maxHeight
			if gtx.Constraints.Max.Y > y {
				gtx.Constraints.Max.Y = y
			}
		}
	} else {
		h := gtx.Dp(unit.Dp(30))
		if m.win.height-gtx.Constraints.Max.Y < h {
			offset.Top = 0
			gtx.Constraints.Max.Y = m.win.height - h
		}
	}
	border := widget.Border{
		Color:        Border,
		Width:        2,
		CornerRadius: 2,
	}

	return offset.Layout(gtx, func(gtx C) D {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) D {
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 8).Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, dlgBgs[0])
				return D{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) D {
				return border.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(8).Layout(gtx, func(gtx C) D {
						return material.Body1(m.Theme(), m.Message).Layout(gtx)
					})
				})
			},
		)
	})
	// return border.Layout(gtx, func(gtx C) D {
	// 		return layout.UniformInset(5).Layout(gtx, func(gtx C) D {
	// 			return material.Body1(m.Theme(), m.Message).Layout(gtx)
	// 		})
	// 	})
	// return offset.Layout(gtx, func(gtx C) D {
	// 	return border.Layout(gtx, func(gtx C) D {
	// 		return layout.UniformInset(5).Layout(gtx, func(gtx C) D {
	// 			return layout.Background{}.Layout(gtx,
	// 				func(gtx layout.Context) D {
	// 					defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 8).Push(gtx.Ops).Pop()
	// 					paint.Fill(gtx.Ops, dlgBgs[0])
	// 					return D{Size: gtx.Constraints.Min}
	// 				},
	// 				func(gtx layout.Context) D {
	// 					return material.Body1(m.Theme(), m.Message).Layout(gtx)
	// 				},
	// 			)
	// 		})
	// 	})
	// })

}
