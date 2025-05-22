/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:58:38
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 18:53:20
 * @FilePath: \goui\gui\container.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image"
	"image/color"
	"sync/atomic"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

var cardID uint32

type Card struct {
	Widget[Card]
	layout.Flex
	top   layout.Flex
	right layout.Flex
	img   *Image
	title *Label
	time  *Label
	msg   *Label
	split *Split
	f     func(card *Card)
	id    int //删除的识别字段
	widget.Clickable
}

func NewCard(win *Window, title, time, msg string) *Card {
	m := new(Card)
	m.Init(win, m)
	m.top.Spacing = layout.SpaceBetween
	m.right.Axis = layout.Vertical
	m.title = NewLabel(win, title).UniformInset(1).SetInsetLeft(5)
	m.time = NewLabel(win, time).UniformInset(1).SetInsetRight(5)
	m.msg = NewLabel(win, msg).UniformInset(1).SetInsetLeft(5)
	m.split = NewSplit(color.NRGBA{R: 200, G: 200, B: 200, A: 200})
	m.split.Left = 5
	m.time.TextSize -= 2
	m.msg.TextSize -= 2
	m.SetHeight(55)
	m.id = int(atomic.AddUint32(&cardID, 1))
	return m
}
func (m *Card) ID() int {
	return m.id
}

func (m *Card) OnClick(f func(card *Card)) *Card {
	m.f = f
	return m
}

func (m *Card) SetImage(img image.Image) *Card {
	if m.img != nil {
		m.img.Src = paint.NewImageOp(img)
		return m
	}
	m.img = NewImage(m.win, img).Resize(50, 50).SetInsetLeft(5)
	m.img.Fit = widget.Fill
	return m
}

func (m *Card) SetTitle(title string) *Card {
	return m
}

func (m *Card) Time() string {
	return m.time.text
}

func (m *Card) Title() string {
	return m.title.text
}

func (m *Card) Msg() string {
	return m.msg.Text
}

func (m *Card) Layout(gtx C) D {
	if m.f != nil {
		if m.Clicked(gtx) {
			m.win.Log("Card Clicked")
			m.f(m)
		}
	}

	m.CheckDimensions(&gtx)
	if m.img != nil {
		//有图标左右排列
		return m.Clickable.Layout(gtx, func(gtx C) D { //点击

			return m.Flex.Layout(gtx, //左右排列
				layout.Rigid(func(gtx C) D { return m.img.Layout(gtx) }), //图标
				layout.Flexed(1, func(gtx C) D { //右边上下排列
					return m.right.Layout(gtx, //上下排列
						layout.Rigid(func(gtx C) D {
							return m.top.Layout(gtx, layout.Rigid(m.title.Layout), layout.Rigid(m.time.Layout)) //左右排列 标题和时间
						}),
						layout.Flexed(1, m.msg.Layout), //消息
						layout.Rigid(m.split.Layout),   //底部分割线
						layout.Rigid(m.split.Layout),   //底部分割线
					)
				}),
			)
		})
	}
	//有图标左右排列
	return m.Clickable.Layout(gtx, func(gtx C) D { //点击
		return m.right.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return m.top.Layout(gtx, layout.Rigid(m.title.Layout), layout.Rigid(m.time.Layout)) //左右排列 标题和时间
			}),
			layout.Flexed(1, m.msg.Layout), //消息
			layout.Rigid(m.split.Layout),   //底部分割线
		)
	})
}
