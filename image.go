/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-23 11:14:35
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-23 09:52:31
 * @FilePath: \gui\image.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/iconvg"
)

// var img Image
// imgOp := paint.NewImageOp(img2)
type Image struct {
	Widget[Image]
	widget.Image

	cb func(gtx layout.Context, ie *Image)
}

func IconToImage(src []byte, size int, co color.Color, bgColor color.Color) (image.Image, error) {
	// 解码 IconVG 数据
	var ivo iconvg.Rasterizer
	if err := iconvg.Decode(&ivo, src, nil); err != nil {
		return nil, err
	}

	// 设置画布尺寸（例如 48x48）
	width, height := size, size
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色（例如白色）
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 创建调色板并设置图标颜色（索引0通常是默认填充色）
	var palette iconvg.Palette
	r, g, b, a := co.RGBA()
	palette[0] = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)} // 将指定的颜色设置为调色板的第一个条目

	ivo.SetDstImage(img, img.Bounds(), draw.Src)
	iconvg.Decode(&ivo, src, &iconvg.DecodeOptions{
		Palette: &palette,
	})

	return img, nil
}

func LoadJpeg(data []byte) (img image.Image, err error) {
	buf := bytes.NewBuffer(data)
	return jpeg.Decode(buf)
}

func LoadPng(data []byte) (img image.Image, err error) {
	buf := bytes.NewBuffer(data)
	return png.Decode(buf)
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
