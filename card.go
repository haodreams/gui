/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:58:38
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 21:15:30
 * @FilePath: \goui\gui\container.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image"
	"sync/atomic"
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

var cardID uint32

type Card struct {
	Widget[Card]
	layout.Flex
	top       layout.Flex
	right     layout.Flex
	img       *Image
	title     *Label
	time      *Label
	msg       *Label
	split     *Split
	level     string             //级别
	timeStamp int64              //时间戳
	tsFormat  func(int64) string //时间戳格式化函数(默认是时间戳转时间字符串)
	f         func(card *Card)
	id        int //删除的识别字段
	widget.Clickable
}

func NewCard(win *Window, ts int64, title, msg string) *Card {
	m := new(Card)
	m.Init(win, m)
	m.top.Spacing = layout.SpaceBetween
	m.right.Axis = layout.Vertical
	m.title = NewLabel(win, title).UniformInset(1).SetInsetLeft(5)
	m.time = NewLabel(win, "").UniformInset(1).SetInsetRight(5)
	m.msg = NewLabel(win, msg).UniformInset(1).SetInsetLeft(5)
	m.split = NewSplit()
	m.split.Left = 5
	m.time.TextSize -= 2
	m.msg.TextSize -= 2
	m.SetHeight(55)
	m.id = int(atomic.AddUint32(&cardID, 1))
	m.timeStamp = ts
	m.tsFormat = TimeStampFormat
	return m
}

// 时间戳格式化为字符串
func TimeStampFormat(t int64) string {
	if t == 0 {
		return ""
	}
	//计算今天0点的时间，如果当前时间大于今天0点，则返回今天时分
	//如果时间大于昨天的0点，返回字符串昨天
	//如果小于7天返回周几
	//其他时间返回年月日
	now := time.Now()
	_, offset := now.Zone()
	//计算今天0点的时间戳
	secs := now.Unix()
	today := (secs + int64(offset)) //加上时区偏移
	today -= today % 86400
	today -= int64(offset) //减去时区偏移
	if t > today {
		return time.Unix(t, 0).Format("15:04")
	}
	if t > today-86400 {
		return "昨天"
	}

	if t > today-86400*7 {
		weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
		weekday := time.Unix(t, 0).Weekday()
		return weekdays[weekday]
	}

	return time.Unix(t, 0).Format("01月02日")
}

func (m *Card) ID() int {
	return m.id
}

func (m *Card) OnClick(f func(card *Card)) *Card {
	m.f = f
	return m
}

func (m *Card) SetLevel(level string) *Card {
	m.level = level
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

func (m *Card) SetTimeFormat(format func(int64) string) *Card {
	if format == nil { //如果为空，使用默认的时间戳格式化函数
		m.tsFormat = TimeStampFormat
	}
	m.tsFormat = format
	return m
}

func (m *Card) SetTime(time int64) *Card {
	m.timeStamp = time
	return m
}

func (m *Card) SetStringTime(time string) *Card {
	m.time.text = time
	return m
}

func (m *Card) SetMsg(msg string) *Card {
	m.msg.Text = msg
	return m
}

func (m *Card) Time() int64 {
	return m.timeStamp
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
			//m.win.Log("Card Clicked")
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
							if m.tsFormat != nil {
								m.time.text = m.tsFormat(m.timeStamp)
							}
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
				if m.tsFormat != nil {
					m.time.text = m.tsFormat(m.timeStamp)
				}
				return m.top.Layout(gtx, layout.Rigid(m.title.Layout), layout.Rigid(m.time.Layout)) //左右排列 标题和时间
			}),
			layout.Flexed(1, m.msg.Layout), //消息
			layout.Rigid(m.split.Layout),   //底部分割线
		)
	})
}
