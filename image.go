/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-23 11:14:35
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 00:06:43
 * @FilePath: \gui\image.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strconv"

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

// strconv.ParseIntColor 解析十六进制颜色代码为 RGBA 值
func ParseIntColor(hex string) (r, g, b, a byte, err error) {
	// 移除前缀 #
	if hex[0] == '#' {
		hex = hex[1:]
	}

	var r0, g0, b0 int64
	// 检查长度
	switch len(hex) {
	case 3: // 简写格式 #rgb
		r0, err = strconv.ParseInt(hex[0:1]+hex[0:1], 16, 32)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		g0, err = strconv.ParseInt(hex[1:2]+hex[1:2], 16, 32)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		b0, err = strconv.ParseInt(hex[2:3]+hex[2:3], 16, 32)
		if err != nil {
			return 0, 0, 0, 0, err
		}
	case 6: // 标准格式 #rrggbb
		r0, err = strconv.ParseInt(hex[0:2], 16, 32)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		g0, err = strconv.ParseInt(hex[2:4], 16, 32)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		b0, err = strconv.ParseInt(hex[4:6], 16, 32)
		if err != nil {
			return 0, 0, 0, 0, err
		}
	default:
		return 0, 0, 0, 0, errors.New("无效的颜色格式, 必须是3位或6位十六进制")
	}

	return byte(r0), byte(g0), byte(b0), 255, nil
}

// 解析字符串颜色，格式为#RRGGBB 或者#RGB
func RGB(src string) color.RGBA {
	r, g, b, a, err := ParseIntColor(src)
	if err != nil {
		return color.RGBA{R: 0, G: 0, B: 0, A: 0} // 默认黑色
	}
	return color.RGBA{R: r, G: g, B: b, A: a}
}

// 解析字符串颜色，格式为#RRGGBB 或者#RGB
func NRGB(src string) color.NRGBA {
	r, g, b, a, err := ParseIntColor(src)
	if err != nil {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0} // 默认黑色
	}
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// cs[0]=fg; cs[1]:bg
func IconToImage(src []byte, size int, cs ...color.Color) (image.Image, error) {
	// 设置画布尺寸（例如 48x48）
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	var opts *iconvg.DecodeOptions

	// 创建调色板并设置图标颜色（索引0通常是默认填充色）
	if len(cs) > 0 {
		r, g, b, a := cs[0].RGBA()
		palette := iconvg.Palette{}
		palette[0] = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)} // 将指定的颜色设置为调色板的第一个条目
		opts = new(iconvg.DecodeOptions)
		opts.Palette = &palette
	}

	var z iconvg.Rasterizer
	z.SetDstImage(img, img.Bounds(), draw.Src)
	err := iconvg.Decode(&z, src, opts)

	return img, err
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
