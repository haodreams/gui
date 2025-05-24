package gui

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Config struct {
	theme          *material.Theme
	ButtonInset    layout.Inset
	CheckboxInset  layout.Inset
	EditInset      layout.Inset
	LabelInset     layout.Inset
	TextFieldInset layout.Inset
	SelectInset    layout.Inset
	TabsInset      layout.Inset
	insetInited    bool // 初始化过Inset
	showStatusBar  bool
	plugin         func(*app.Window) event.Event
	logf           func(v ...any)
}

type Option func(*Config)

// NewConfig returns a new Config with default values.
func Inset(v unit.Dp) Option {
	return func(m *Config) {
		m.ButtonInset = layout.UniformInset(v)
		m.CheckboxInset = layout.UniformInset(v)
		m.EditInset = layout.UniformInset(v)
		m.LabelInset = layout.UniformInset(v)
		m.TextFieldInset = layout.UniformInset(v)
		m.SelectInset = layout.UniformInset(v)
		m.TabsInset = layout.UniformInset(v)
		m.TabsInset.Right += 10
		m.TabsInset.Left += 10
		m.insetInited = true
	}
}

// 设置主题
func ThemeShaper(ts *text.Shaper) Option {
	return func(m *Config) {
		m.theme.Shaper = ts
	}
}

// 隐藏状态栏
func HideStatusBar() Option {
	return func(m *Config) {
		m.showStatusBar = false
	}
}

func WithLog(f func(v ...any)) Option {
	return func(m *Config) {
		m.logf = f
	}
}

func WithPlugin(plugin func(*app.Window) event.Event) Option {
	return func(m *Config) {
		m.plugin = plugin
	}
}
