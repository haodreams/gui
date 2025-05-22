/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-23 11:14:35
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2024-07-29 10:51:24
 * @FilePath: \gui\image.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image"

	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/widget"
)

// var img Image
// imgOp := paint.NewImageOp(img2)
type ImageEvent struct {
	Widget[ImageEvent]
	widget.Image
	history  [2]*pointer.Event
	es       [2]*pointer.Event
	LastSize image.Point
	Inited   bool

	cb func(gtx layout.Context, ie *ImageEvent)
}

func NewImageEvent(win *Window) (m *ImageEvent) {
	m = new(ImageEvent)
	m.win = win
	m.Init(win, m)
	return m
}

func (m *ImageEvent) SetCallback(cb func(gtx layout.Context, ie *ImageEvent)) *ImageEvent {
	m.cb = cb
	return m
}

func (m *ImageEvent) GetHistory() (e []*pointer.Event) {
	copy(m.es[:], m.history[:])
	m.history[0] = nil
	m.history[1] = nil
	return m.es[:]
}

// Add draws the widget to the given graphics context.
// 没有这段代码不会出现事件
func (m *ImageEvent) Add(gtx layout.Context, d D) {
	pr := clip.Rect(image.Rectangle{Max: d.Size}).Push(gtx.Ops)
	event.Op(gtx.Ops, m)
	m.update(gtx.Source)
	pr.Pop()
}

func (m *ImageEvent) Layout(gtx layout.Context) D {
	if m.cb != nil {
		m.cb(gtx, m)
	}
	if !m.Inited {
		return D{}
	}
	m.CheckDimensions(&gtx)
	d := m.Image.Layout(gtx)
	m.Add(gtx, d)
	return d
}

func (m *ImageEvent) push(e *pointer.Event) {
	if m.history[0] == nil {
		m.history[0] = e
		return
	}
	if m.history[1] == nil {
		m.history[1] = e
		return
	}
	m.history[0] = m.history[1]
	m.history[1] = e
}

// Update state and report the scroll distance along axis.
func (m *ImageEvent) update(q input.Source) {
	f := pointer.Filter{
		Target: m,
		Kinds:  pointer.Move | pointer.Press | pointer.Release | pointer.Drag | pointer.Scroll,
	}
	for {
		evt, ok := q.Event(f)
		if !ok {
			break
		}
		e, ok := evt.(pointer.Event)
		if !ok {
			continue
		}

		// if e.Kind == pointer.Release {
		// 	m.push(&e)
		// }
		m.push(&e)
	}
}
