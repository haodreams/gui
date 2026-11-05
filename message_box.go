package gui

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
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

// 对话框配色(参考 Windows 风格): 白色背景 + 深色文字
var (
	dlgBgs = []color.NRGBA{
		{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, //info  白
		{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, //warn  白
		{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, //error 白
	}
	dlgTextColor = []color.NRGBA{
		{R: 0x1F, G: 0x1F, B: 0x1F, A: 0xFF}, //info  深灰
		{R: 0x1F, G: 0x1F, B: 0x1F, A: 0xFF}, //warn  深灰
		{R: 0x1F, G: 0x1F, B: 0x1F, A: 0xFF}, //error 深灰
	}
	//对话框图标底色(参考 Windows): info蓝 / warn黄 / error红
	dlgIconBgs = []color.NRGBA{
		{R: 0x00, G: 0x78, B: 0xD7, A: 0xFF}, //info  #0078D7 蓝
		{R: 0xFF, G: 0xC7, B: 0x00, A: 0xFF}, //warn  #FFC700 黄
		{R: 0xE8, G: 0x11, B: 0x23, A: 0xFF}, //error #E81123 红
	}
)

// Windows风格标题栏配色
var (
	dlgTitleBarBg   = color.NRGBA{R: 0xF3, G: 0xF3, B: 0xF3, A: 0xFF} //标题栏背景
	dlgTitleFg      = color.NRGBA{R: 0x1F, G: 0x1F, B: 0x1F, A: 0xFF} //标题文字
	dlgCloseHoverBg = color.NRGBA{R: 0xE8, G: 0x11, B: 0x23, A: 0xFF} //关闭按钮悬停背景(红)
	dlgCloseHoverFg = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} //关闭按钮悬停图标(白)
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
	lastWidth    unit.Dp
	rememberBool *widget.Bool //是否记住选择
	options      []MBOption
	result       string
	refreshCount byte
	onSubmit     func(selectedOption string, remember bool) //确定后的回调函数
	content      Contenter
	closeBtn     widget.Clickable //标题栏X关闭按钮

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
	m.Width = 400
	m.lastWidth = m.Width
	m.Height = 200
	m.win = win

	m.Scrim.FinalAlpha = 82 //default
	return m
}
func (m *MessageBox) RestoreWidth() *MessageBox {
	m.Width = m.lastWidth
	return m
}

func (m *MessageBox) SetOptions(options ...MBOption) {
	m.options = options
}

func (m *MessageBox) Show() {
	m.Visible = true
	m.refreshCount = 0
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
func (m *MessageBox) Layout(gtx layout.Context) D {
	if !m.Visible {
		return D{}
	}
	if m.Scrim.Clicked(gtx) {
		//m.win.Log("click scrim")
	}
	d1 := m.Scrim.Layout(gtx, m.Theme())

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

	//高度不需要进行转换， 高度是已经计算后的数据
	maxHeight := int(m.Height)
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
		if m.win.height-gtx.Constraints.Max.Y < h {
			offset.Top = 0
			gtx.Constraints.Max.Y = m.win.height - h
		}
	}
	d2 := offset.Layout(gtx, func(gtx C) D {
		//window.go的Stack.Expanded会把约束Min撑满整个窗口, 这里重置Min:
		//高度按内容自适应, 宽度仍保持Width(窗口不足时铺满)
		gtx.Constraints.Min = image.Point{X: min(maxWidth, gtx.Constraints.Max.X)}
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
		Width:        0,
		CornerRadius: 3, //Windows风格圆角
	}
	return border.Layout(gtx, func(gtx layout.Context) D {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) D {
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 8).Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, dlgBgs[m.Type])
				return D{Size: gtx.Constraints.Min}
			}, func(gtx layout.Context) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(m.layoutTitleBar),
					layout.Rigid(m.layoutBody),
				)
			},
		)
	})
}

// layoutTitleBar 绘制Windows风格标题栏: 左侧标题, 右侧X关闭按钮(悬停变红)
func (m *MessageBox) layoutTitleBar(gtx layout.Context) D {
	if m.closeBtn.Clicked(gtx) {
		m.Hide()
		//设置页面需要刷新
		gtx.Execute(op.InvalidateCmd{})
	}
	return layout.Background{}.Layout(gtx,
		func(gtx C) D { //标题栏背景(顶部圆角与对话框一致)
			r := image.Rectangle{Max: gtx.Constraints.Min}
			paint.FillShape(gtx.Ops, dlgTitleBarBg, clip.RRect{Rect: r, NW: 8, NE: 8}.Op(gtx.Ops))
			return D{Size: gtx.Constraints.Min}
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Left: unit.Dp(10), Top: unit.Dp(6), Bottom: unit.Dp(6)}.Layout(gtx, func(gtx C) D {
						h := material.Body1(m.Theme(), m.Title)
						h.Color = dlgTitleFg
						h.MaxLines = 1
						return h.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D { //X关闭按钮
					size := gtx.Dp(32)
					gtx.Constraints.Min = image.Pt(size, size)
					gtx.Constraints.Max = gtx.Constraints.Min
					hovered := m.closeBtn.Hovered()
					if hovered { //悬停画红色背景(右上角与对话框圆角一致)
						r := image.Rectangle{Max: gtx.Constraints.Min}
						paint.FillShape(gtx.Ops, dlgCloseHoverBg, clip.RRect{Rect: r, NE: 8}.Op(gtx.Ops))
					}
					return m.closeBtn.Layout(gtx, func(gtx C) D {
						iconColor := dlgTitleFg
						if hovered {
							iconColor = dlgCloseHoverFg
						}
						//Center必须在完整按钮区内居中, 约束缩小只能作用于图标本身
						return layout.Center.Layout(gtx, func(gtx C) D {
							s := gtx.Dp(16)
							gtx.Constraints.Min = image.Pt(s, s)
							gtx.Constraints.Max = gtx.Constraints.Min
							return CloseIcon.Layout(gtx, iconColor)
						})
					})
				}),
			)
		},
	)
}

// layoutDialogIcon 绘制Windows风格对话框图标(32dp): info蓝圈i / warn黄三角! / error红圈×
func layoutDialogIcon(gtx C, tp ModelType) D {
	size := gtx.Dp(32)
	s := float32(size)
	sz := image.Rectangle{Max: image.Pt(size, size)}
	white := color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	stroke := int(s * 0.125) //符号笔画宽度

	//fillRect 在offset处填充一个矩形
	fillRect := func(offset image.Point, wh image.Point, c color.NRGBA) {
		defer op.Offset(offset).Push(gtx.Ops).Pop()
		paint.FillShape(gtx.Ops, c, clip.Rect(image.Rectangle{Max: wh}).Op())
	}

	switch tp {
	case ModalWarn: //黄三角 + 黑!
		var p clip.Path
		p.Begin(gtx.Ops)
		p.MoveTo(f32.Pt(s/2, s*0.06))
		p.LineTo(f32.Pt(s*0.97, s*0.9))
		p.LineTo(f32.Pt(s*0.03, s*0.9))
		p.Close()
		paint.FillShape(gtx.Ops, dlgIconBgs[tp], clip.Outline{Path: p.End()}.Op())
		dark := color.NRGBA{A: 0xFF}
		fillRect(image.Pt(size/2-stroke/2, int(s*0.36)), image.Pt(stroke, int(s*0.32)), dark) //竖线
		fillRect(image.Pt(size/2-stroke/2, int(s*0.74)), image.Pt(stroke, stroke), dark)      //点
	case ModalError: //红圈 + 白×
		paint.FillShape(gtx.Ops, dlgIconBgs[tp], clip.Ellipse(sz).Op(gtx.Ops))
		var p clip.Path
		p.Begin(gtx.Ops)
		p.MoveTo(f32.Pt(s*0.3, s*0.3))
		p.LineTo(f32.Pt(s*0.7, s*0.7))
		p.MoveTo(f32.Pt(s*0.7, s*0.3))
		p.LineTo(f32.Pt(s*0.3, s*0.7))
		paint.FillShape(gtx.Ops, white, clip.Stroke{Path: p.End(), Width: float32(stroke)}.Op())
	default: //ModalInfo 蓝圈 + 白i
		paint.FillShape(gtx.Ops, dlgIconBgs[tp], clip.Ellipse(sz).Op(gtx.Ops))
		fillRect(image.Pt(size/2-stroke/2, int(s*0.22)), image.Pt(stroke, stroke), white)     //点
		fillRect(image.Pt(size/2-stroke/2, int(s*0.4)), image.Pt(stroke, int(s*0.38)), white) //竖线
	}
	return D{Size: image.Pt(size, size)}
}

// layoutBody 绘制消息正文与操作按钮
func (m *MessageBox) layoutBody(gtx layout.Context) D {
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) D {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) D {
				return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) D {
					var d D
					//预留高度 = 标题栏32 + 正文上下边距20 + 按钮区(上间距5+下间距10+按钮高)
					btnH := gtx.Metric.SpToDp(m.Theme().TextSize)*1.2 + 10 //与下方按钮高度公式一致
					dp := gtx.Dp(70 + btnH)
					if m.content != nil {
						//替换为自定义容器
						d = m.content.Layout(gtx)
						if m.win.height-d.Size.Y < dp {
							d.Size.Y = m.win.height - dp
						}
					} else {
						d = layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D { //类型图标
								return layoutDialogIcon(gtx, m.Type)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
							layout.Flexed(1, func(gtx C) D {
								b := material.Body1(m.Theme(), m.Message)
								b.Color = dlgTextColor[m.Type]
								return b.Layout(gtx)
							}),
						)
					}
					m.Height = unit.Dp(d.Size.Y + dp)
					if m.refreshCount < 3 { //动态刷新
						// m.win.Log(m.MaxHeight)
						m.refreshCount++
						gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Millisecond * 300)})
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
					if ops[i].Button.Clicked(gtx) {
						m.result = ops[i].Text
						m.submit()
					}

					items = append(
						items,
						layout.Rigid(func(gtx layout.Context) D {
							height := gtx.Metric.SpToDp(m.Theme().TextSize)*1.2 + 10
							gtx.Constraints.Max.Y = gtx.Dp(height)
							gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
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
				return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) D {
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
}
