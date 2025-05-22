/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-23 11:14:35
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 17:20:24
 * @FilePath: \gui\image.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"bytes"
	"image"
	"image/jpeg"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

// var img Image
// imgOp := paint.NewImageOp(img2)
type Image struct {
	Widget[Image]
	widget.Image

	cb func(gtx layout.Context, ie *Image)
}

func LoadJpeg(data []byte) (img image.Image, err error) {
	buf := bytes.NewBuffer(data)
	return jpeg.Decode(buf)
}

func NewImage(win *Window, img image.Image) (m *Image) {
	m = new(Image)
	m.win = win
	m.Init(win, m)
	m.Src = paint.NewImageOp(img)
	return m
}

func (m *Image) Layout(gtx layout.Context) D {
	if m.cb != nil {
		m.cb(gtx, m)
	}

	m.CheckDimensions(&gtx)
	return m.Image.Layout(gtx)
}
