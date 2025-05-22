package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var NavigationExpandMore = []byte{
	0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50, 0x50, 0xb0, 0xb0, 0xc0, 0x2d, 0x89, 0x2d, 0x79,
	0x00, 0x80, 0x59, 0x82, 0x20, 0xd5, 0x76, 0xd5, 0x76, 0x00, 0x68, 0x78, 0x21, 0x98, 0x98, 0x98,
	0x68, 0xe1,
}

type Select struct {
	Widget[Select]
	menuContextArea ContextArea
	menu            component.MenuState
	list            *widget.List

	// MinWidth unit.Dp
	// MaxWidth unit.Dp
	menuInit bool

	isOpen              bool
	selectedOptionIndex int
	lastSelectedIndex   int
	options             []*SelectOption

	size image.Point

	borderWidth  unit.Dp
	cornerRadius unit.Dp

	onValueChange func(value string)

	icon *widget.Icon
}

func NewSelect(win *Window, options ...*SelectOption) *Select {
	c := &Select{
		menuContextArea: ContextArea{
			AbsolutePosition: true,
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		options:      options,
		borderWidth:  unit.Dp(1),
		cornerRadius: unit.Dp(4),
		menuInit:     true,
	}
	c.win = win
	c.icon = ExpandIcon
	c.Init(win, c)
	c.Inset = win.SelectInset
	return c
}

func NewSelectWithoutBorder(win *Window, options ...*SelectOption) *Select {
	c := &Select{
		menuContextArea: ContextArea{
			AbsolutePosition: true,
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		options:  options,
		menuInit: true,
	}
	c.win = win
	c.Init(win, c)
	c.icon, _ = widget.NewIcon(icons.NavigationExpandMore)
	c.Inset = win.SelectInset

	return c
}

type SelectOption struct {
	Text      string
	Value     string
	clickable widget.Clickable

	Icon      *widget.Icon
	IconColor color.NRGBA

	isDivider bool
	isDefault bool
}

func NewSelectOption(text string) *SelectOption {
	return &SelectOption{
		Text:      text,
		isDivider: false,
	}
}

func NewSelectDivider() *SelectOption {
	return &SelectOption{
		isDivider: true,
	}
}

func (o *SelectOption) WithValue(value string) *SelectOption {
	o.Value = value
	return o
}

func (o *SelectOption) WithIcon(icon *widget.Icon, color color.NRGBA) *SelectOption {
	o.Icon = icon
	o.IconColor = color
	return o
}

func (o *SelectOption) DefaultSelected() *SelectOption {
	o.isDefault = true
	return o
}

func (o *SelectOption) GetText() string {
	if o == nil {
		return ""
	}

	return o.Text
}

func (o *SelectOption) GetValue() string {
	if o == nil {
		return ""
	}

	return o.Value
}

func (c *Select) SetSelected(index int) *Select {
	c.selectedOptionIndex = index
	c.lastSelectedIndex = index
	return c
}

func (c *Select) SetOnChanged(f func(value string)) *Select {
	c.onValueChange = f
	return c
}

func (c *Select) SetSelectedByTitle(title string) {
	if len(c.options) == 0 {
		return
	}

	for i, opt := range c.options {
		if opt.Text == title {
			c.selectedOptionIndex = i
			c.lastSelectedIndex = i
			break
		}
	}
}

func (c *Select) SetSelectedByValue(value string) {
	for i, opt := range c.options {
		if opt.Value == value {
			c.selectedOptionIndex = i
			c.lastSelectedIndex = i
			break
		}
	}
}

func (c *Select) SelectedIndex() int {
	return c.selectedOptionIndex
}

func (c *Select) SetOptions(options ...*SelectOption) {
	c.selectedOptionIndex = 0
	c.options = options
	if len(c.options) > 0 {
		c.menuInit = true
	}
}

func (c *Select) GetSelected() *SelectOption {
	if len(c.options) == 0 {
		return nil
	}

	return c.options[c.selectedOptionIndex]
}

func (c *Select) box(gtx layout.Context, theme *material.Theme, text string) D {
	borderColor := Border
	if c.isOpen {
		borderColor = BorderFocused
	}

	border := widget.Border{
		Color:        borderColor,
		Width:        c.borderWidth,
		CornerRadius: c.cornerRadius,
	}

	// if maxWidth == 0 {
	// 	maxWidth = unit.Dp(gtx.Constraints.Max.X)
	// }

	// c.size.X = int(maxWidth)

	c.CheckDimensions(&gtx)

	return border.Layout(gtx, func(gtx layout.Context) D {
		return c.Inset.Layout(gtx, func(gtx layout.Context) D {
			// calculate the minimum width of the box, considering icon and padding
			gtx.Constraints.Min.X -= gtx.Dp(8)

			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) D {
					return material.Label(theme, theme.TextSize, text).Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) D {
					gtx.Constraints.Min.X = gtx.Dp(16)
					return c.icon.Layout(gtx, theme.Palette.Fg)
				}),
			)
		})
	})

}

func (c *Select) SetSize(size image.Point) {
	c.size = size
}

// Layout the DropDown.
func (c *Select) Layout(gtx layout.Context) D {
	c.isOpen = c.menuContextArea.Active()

	for i, opt := range c.options {
		if opt.isDefault {
			c.selectedOptionIndex = i
		}

		for opt.clickable.Clicked(gtx) {
			c.isOpen = false
			c.selectedOptionIndex = i
		}
	}

	if c.selectedOptionIndex != c.lastSelectedIndex {
		if c.onValueChange != nil {
			c.onValueChange(c.options[c.selectedOptionIndex].Value)
		}
		c.lastSelectedIndex = c.selectedOptionIndex
	}

	// Update menu items only if options change
	if c.menuInit {
		c.menuInit = false
		c.updateMenuItems(c.Theme())
	}

	// if c.MinWidth == 0 {
	// 	c.MinWidth = gtx.Dp(100)
	// }

	c.CheckDimensions(&gtx)

	text := ""
	if c.selectedOptionIndex >= 0 && c.selectedOptionIndex < len(c.options) {
		text = c.options[c.selectedOptionIndex].Text
	}

	box := c.box(gtx, c.Theme(), text)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) D {
			return box
		}),
		layout.Expanded(func(gtx layout.Context) D {
			return c.menuContextArea.Layout(gtx, func(gtx layout.Context) D {
				offset := layout.Inset{
					Top:  unit.Dp(float32(box.Size.Y)/gtx.Metric.PxPerDp + 1),
					Left: unit.Dp(4),
				}
				return offset.Layout(gtx, func(gtx layout.Context) D {
					// gtx.Constraints.Min.X = gtx.Dp(c.MinWidth)
					// if c.MaxWidth != 0 {
					// 	gtx.Constraints.Max.X = gtx.Dp(c.MaxWidth)
					// }
					m := component.Menu(c.Theme(), &c.menu)
					m.SurfaceStyle.Fill = c.Theme().Palette.Bg
					return m.Layout(gtx)
				})
			})
		}),
	)
}

// updateMenuItems creates or updates menu items based on options and calculates minWidth.
func (c *Select) updateMenuItems(theme *material.Theme) {
	c.menu.Options = c.menu.Options[:0]
	for _, opt := range c.options {
		opt := opt
		c.menu.Options = append(c.menu.Options, func(gtx layout.Context) D {
			if opt.isDivider {
				dv := component.Divider(theme)
				dv.Fill = Border
				return dv.Layout(gtx)
			}

			itm := component.MenuItem(theme, &opt.clickable, opt.Text)
			if opt.Icon != nil {
				itm.Icon = opt.Icon
				itm.IconColor = opt.IconColor
			}

			itm.Label.Color = theme.Palette.Fg
			return itm.Layout(gtx)
		})
	}
}
