/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-20 16:58:38
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-06-04 11:10:17
 * @FilePath: \goui\gui\container.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image"
	"image/color"
	"sync/atomic"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
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
	State     int                //设置当前卡片的状态
	level     string             //级别
	timeStamp int64              //时间戳
	tsFormat  func(int64) string //时间戳格式化函数(默认是时间戳转时间字符串)
	onClick   func(card *Card)
	id        int  //删除的识别字段
	hide      bool //是否不显示
	widget.Clickable
	UnRead bool //未读
}

func NewCard(win *Window, ts int64, title, msg string) *Card {
	m := new(Card)
	m.Init(win, m)
	m.top.Spacing = layout.SpaceBetween
	m.right.Axis = layout.Vertical
	m.title = NewLabel(win, title).UniformInset(1).SetInsetLeft(6).SetInsetTop(2)
	m.time = NewLabel(win, "").UniformInset(1).SetInsetRight(5)
	m.msg = NewLabel(win, msg).SetInsetRight(1).SetInsetLeft(5).SetInsetTop(4).SetInsetBottom(2)
	m.split = NewSplit()
	m.split.Left = 5
	m.time.TextSize -= 2
	m.msg.TextSize -= 2
	m.SetHeight(60)
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
		//超过上周写日期
		wd := int(weekday+7) % 8
		twd := int(now.Weekday()+7) % 8
		if wd >= twd {
			return time.Unix(t, 0).Format("01月02日")
		}
		return weekdays[weekday]
	}

	return time.Unix(t, 0).Format("01月02日")
}

// ID 获取卡片的识别字段
// 返回卡片的识别字段
func (m *Card) ID() int {
	return m.id
}

// Show 显示卡片
// 返回卡片
func (m *Card) Show() *Card {
	m.hide = false
	return m
}

// Hide 隐藏卡片
// 返回卡片
func (m *Card) Hide() *Card {
	m.hide = true
	return m
}

// IsHide 是否隐藏
// 返回是否隐藏
func (m *Card) IsHide() bool {
	return m.hide
}

// SetOnClick 设置点击回调函数
// onClick 点击回调函数，参数为当前卡片
func (m *Card) SetOnClick(onClick func(card *Card)) *Card {
	m.onClick = onClick
	return m
}

// SetLevel 设置卡片的级别
// level 卡片的级别
func (m *Card) SetLevel(level string) *Card {
	m.level = level
	return m
}

// SetImage 设置卡片的图标
// img 卡片的图标
func (m *Card) SetImage(img image.Image) *Card {
	if m.img != nil {
		m.img.SetImage(img)
		return m
	}
	m.img = NewImage(m.win, img).Resize(50, 50).UniformInset(5).SetInsetTop(10).SetInsetBottom(0)
	m.img.Fit = widget.Fill
	return m
}

// GetImage 获取卡片的图标
// 返回卡片的图标
func (m *Card) GetImage() *Image {
	return m.img
}

// SetTitle 设置卡片的标题
// title 卡片的标题
func (m *Card) SetTitle(title string) *Card {
	m.title.Text = title
	return m
}

// SetTimeFormat 设置时间戳格式化函数
// format 时间戳格式化函数，参数为时间戳，返回时间字符串
func (m *Card) SetTimeFormat(format func(int64) string) *Card {
	if format == nil { //如果为空，使用默认的时间戳格式化函数
		m.tsFormat = TimeStampFormat
	}
	m.tsFormat = format
	return m
}

// SetTime 设置卡片的时间戳
// time 卡片的时间戳
func (m *Card) SetTime(time int64) *Card {
	m.timeStamp = time
	return m
}

// SetStringTime 设置卡片的时间戳为字符串
// time 时间戳字符串
func (m *Card) SetStringTime(time string) *Card {
	m.time.Text = time
	return m
}

// SetMsg 设置卡片的消息
// msg 卡片的消息
func (m *Card) SetMsg(msg string) *Card {
	m.msg.Text = msg
	return m
}

// Time 获取卡片的时间戳
// 返回卡片的时间戳
func (m *Card) Time() int64 {
	return m.timeStamp
}

// Title 获取卡片的标题
// 返回卡片的标题
func (m *Card) Title() string {
	return m.title.Text
}

// Msg 获取卡片的消息
// 返回卡片的消息
func (m *Card) Msg() string {
	return m.msg.Text
}

// Layout 布局卡片
// 返回卡片的布局
func (m *Card) Layout(gtx C) D {
	if m.hide {
		return D{}
	}
	if m.onClick != nil {
		if m.Clicked(gtx) {
			//m.win.Log("Card Clicked")
			m.onClick(m)
		}
	}

	m.CheckDimensions(&gtx)
	if m.img != nil {
		//有图标左右排列
		return m.Clickable.Layout(gtx, func(gtx C) D { //点击
			return m.Flex.Layout(gtx, //左右排列
				layout.Rigid(func(gtx C) D {
					d := m.img.Layout(gtx)
					if m.UnRead {
						x := gtx.Dp(12)
						y := gtx.Dp(12)
						c := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
						paint.FillShape(gtx.Ops, c, clip.Ellipse{Min: image.Point{d.Size.X - x, 0}, Max: image.Point{d.Size.X, y}}.Op(gtx.Ops))
						//todo 输出文字
						// 使用内置字体，字号 16px，黑色
						// textStyle := text.Style{
						// 	Font:  gofont.Collection(),       // 字体（必填）
						// 	Size:  unit.Sp(16),               // 字号
						// 	Color: color.NRGBA{0, 0, 0, 255}, // 颜色
						// }
						// layout.Inset{Top: 10, Left: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// 	return textStyle.Layout(gtx, "Hello, World!") // 在左上角边距 10px 处绘制文字
						// })
					}
					return d
				}), //图标
				layout.Flexed(1, func(gtx C) D { //右边上下排列
					return m.right.Layout(gtx, //上下排列
						layout.Rigid(func(gtx C) D {
							if m.tsFormat != nil {
								m.time.Text = m.tsFormat(m.timeStamp)
							}
							return m.top.Layout(gtx, layout.Rigid(m.title.Layout), layout.Rigid(m.time.Layout)) //左右排列 标题和时间
						}),
						layout.Flexed(1, m.msg.Layout), //消息
						layout.Rigid(m.split.Layout),   //底部分割线
					)
				}),
			)
		})
	}
	//没有有图标左右排列
	return m.Clickable.Layout(gtx, func(gtx C) D { //点击
		return m.right.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if m.tsFormat != nil {
					m.time.SetText(m.tsFormat(m.timeStamp))
				}
				return m.top.Layout(gtx, layout.Rigid(m.title.Layout), layout.Rigid(m.time.Layout)) //左右排列 标题和时间
			}),
			layout.Flexed(1, m.msg.Layout), //消息
			layout.Rigid(m.split.Layout),   //底部分割线
		)
	})
}
