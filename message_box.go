package gui

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

// Prompt is a modal dialog that prompts the user for a response.
type ModelType byte

const (
	ModalInfo  ModelType = 0
	ModalWarn  ModelType = 1
	ModalError ModelType = 2
)

const (
	IconPositionStart = 0
	IconPositionEnd   = 1
)

var (
	dlgBgs = []color.NRGBA{
		{R: 224, G: 236, B: 255, A: 0xFF}, //info
		{R: 255, G: 255, B: 224, A: 0xFF}, //warn
		{R: 255, G: 224, B: 224, A: 0xFF}, //error
	}
	dlgTextColor = []color.NRGBA{
		{R: 0, G: 0, B: 255, A: 255},  //info
		{R: 0, G: 0, B: 0, A: 255},    //warn
		{R: 255, G: 0, B: 0, A: 0xFF}, //error
	}
)

// Scrim implments a clickable translucent overlay. It can animate appearing
// and disappearing as a fade-in, fade-out transition from zero opacity
// to a fixed maximum opacity.
type Scrim struct {
	// FinalAlpha is the final opacity of the scrim on a scale from 0 to 255.
	FinalAlpha uint8
	widget.Clickable
}

// Layout draws the scrim using the provided animation. If the animation indicates
// that the scrim is not visible, this is a no-op.
func (m *Scrim) Layout(gtx layout.Context, th *material.Theme) D {
	return m.Clickable.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min = gtx.Constraints.Max
		currentAlpha := m.FinalAlpha

		color := th.Fg
		color.A = currentAlpha
		component.Rect{Color: color, Size: gtx.Constraints.Max}.Layout(gtx)
		return D{Size: gtx.Constraints.Max}
	})
}

type MessageBox struct {
	Title   string
	Message string
	Type    ModelType
	Visible bool

	Widget[MessageBox]
	lastWidth    int
	rememberBool *widget.Bool //是否记住选择
	options      []MBOption
	result       string
	needRefresh  bool
	onSubmit     func(selectedOption string, remember bool) //确定后的回调函数
	content      Contenter

	Scrim //遮罩层，防止鼠标穿透
}

type MBOption struct {
	Text   string
	Button widget.Clickable
	Icon   *widget.Icon
}

func NewMessage(win *Window, title, content string, modalType ModelType, options ...MBOption) *MessageBox {
	m := &MessageBox{
		Title:   title,
		Message: content,
		Type:    modalType,
		options: options,
	}
	m.Init(win, m)
	m.MaxWidth = 400
	m.lastWidth = m.MaxWidth
	m.MaxHeight = 200
	m.win = win

	m.Scrim.FinalAlpha = 82 //default
	return m
}
func (m *MessageBox) RestoreWidth() *MessageBox {
	m.MaxWidth = m.lastWidth
	return m
}

func (m *MessageBox) SetOptions(options ...MBOption) {
	m.options = options
}

func (m *MessageBox) Show() {
	m.Visible = true
	m.needRefresh = true
}

func (m *MessageBox) Hide() {
	m.Visible = false
}

func (m *MessageBox) IsVisible() bool {
	return m.Visible
}

func (m *MessageBox) WithRememberBool() {
	m.rememberBool = &widget.Bool{Value: false}
}

func (m *MessageBox) WithoutRememberBool() {
	m.rememberBool = nil
}

func (m *MessageBox) SetOnSubmit(f func(selectedOption string, remember bool)) {
	m.onSubmit = f
}

func (m *MessageBox) submit() {
	if m.onSubmit == nil {
		return
	}

	if !m.Visible {
		return
	}

	if m.rememberBool == nil {
		m.onSubmit(m.result, false)
		return
	}

	m.onSubmit(m.result, m.rememberBool.Value)
}

func (m *MessageBox) Result() (string, bool) {
	if m.result == "" {
		return "", false
	}

	if m.rememberBool != nil {
		return m.result, m.rememberBool.Value
	}

	return m.result, false
}
func (m *MessageBox) Layout(gtx layout.Context, width, height int) D {
	if !m.Visible {
		return D{}
	}
	if m.Scrim.Clicked(gtx) {
		m.win.Log("click scrim")
	}
	d1 := m.Scrim.Layout(gtx, m.Theme())

	offset := layout.Inset{}

	maxWidth := gtx.Dp(unit.Dp(m.MaxWidth))
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

	//高度不需要进行转换， 高度是已经计算后的数据
	maxHeight := m.MaxHeight
	if maxHeight > 0 {
		offset.Top = unit.Dp(gtx.Constraints.Max.Y-maxHeight) / 2 / unit.Dp(gtx.Metric.PxPerDp)
		if offset.Top < 0 {
			offset.Top = 0
		} else {
			y := int(gtx.Dp(offset.Top)) + maxHeight
			if gtx.Constraints.Max.Y > y {
				gtx.Constraints.Max.Y = y
			}
		}
	} else {
		h := gtx.Dp(unit.Dp(15))
		if height-gtx.Constraints.Max.Y < h {
			offset.Top = 0
			gtx.Constraints.Max.Y = height - h
		}
	}
	d2 := offset.Layout(gtx, func(gtx C) D {
		return m.layout(gtx)
	})
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) D {
			return d1
		}),
		layout.Stacked(func(gtx layout.Context) D {
			return d2
		}),
	)
}

func (m *MessageBox) layout(gtx layout.Context) D {
	border := widget.Border{
		Color:        Border,
		Width:        2,
		CornerRadius: 2,
	}
	return border.Layout(gtx, func(gtx layout.Context) D {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) D {
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 8).Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, dlgBgs[m.Type])
				return D{Size: gtx.Constraints.Min}
			}, func(gtx layout.Context) D {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) D {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) D {
							return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) D {
								h := material.H6(m.Theme(), m.Title)
								h.Color = dlgTextColor[m.Type]
								return h.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx layout.Context) D {
							return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) D {
								var d D
								dp := gtx.Dp(110)
								if m.content != nil {
									//替换为自定义容器
									d = m.content.Layout(gtx)
									if m.win.height-d.Size.Y < dp {
										d.Size.Y = m.win.height - dp
									}
								} else {
									b := material.Body1(m.Theme(), m.Message)
									b.Color = dlgTextColor[m.Type]
									d = b.Layout(gtx)
								}
								m.MaxHeight = d.Size.Y + dp
								if m.needRefresh {
									// m.win.Log(m.MaxHeight)
									m.needRefresh = false
									gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 100)})
								}
								return d
							})
						}),

						layout.Rigid(func(gtx layout.Context) D {
							ops := m.options
							count := len(ops)
							if m.rememberBool != nil {
								count++
							}

							items := make([]layout.FlexChild, 0, count)
							if m.rememberBool != nil {
								items = append(
									items,
									layout.Rigid(func(gtx layout.Context) D {
										return material.CheckBox(m.Theme(), m.rememberBool, "Don't ask again").Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
								)
							}

							for i := range ops {
								i := i

								if ops[i].Button.Clicked(gtx) {
									m.result = ops[i].Text
									m.submit()
								}

								items = append(
									items,
									layout.Rigid(func(gtx layout.Context) D {
										gtx.Constraints.Max.Y = gtx.Dp(40)
										btn := material.Button(m.Theme(), &ops[i].Button, ops[i].Text)
										btn.Inset = layout.Inset{
											Top: 5, Bottom: 5,
											Left: 12, Right: 12,
										}
										return btn.Layout(gtx)
										//return material.Button(p.Theme(), &p.options[i].Button, p.options[i].Text).Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
								)
							}
							return layout.Inset{Top: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) D {
								return layout.Flex{
									Axis:      layout.Horizontal,
									Alignment: layout.Middle,
									Spacing:   layout.SpaceStart,
								}.Layout(gtx,
									items...,
								)
							})
						}),
					)
				})

			},
		)
	})
}
