/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2026-11-04
 * @FilePath: \gui\tree.go
 * @Description: 通用树控件, 按key(默认以"."分隔)将扁平数据构建为树形结构展示
 */
package gui

import (
	"image"
	"image/color"
	"strings"
	"sync"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// TreeNode 树节点, Data 为 nil 表示目录节点, 否则为叶子(数据)节点
type TreeNode[G any] struct {
	Name     string         //当前层级名称(路径中的一段)
	Data     *G             //叶子节点携带的数据
	Parent   *TreeNode[G]   //父节点
	Children []*TreeNode[G] //子节点
	Depth    int            //层级深度, 根为0
	Expanded bool           //是否展开

	childIdx  map[string]*TreeNode[G]
	clickable *widget.Clickable
}

// IsLeaf 是否为叶子节点
func (m *TreeNode[G]) IsLeaf() bool {
	return len(m.Children) == 0
}

// FullName 完整路径名
func (m *TreeNode[G]) FullName(sep string) string {
	if m.Parent == nil {
		return m.Name
	}
	return m.Parent.FullName(sep) + sep + m.Name
}

// TreeTextFunc 返回节点显示文本和颜色(颜色可为nil)
type TreeTextFunc[G any] func(node *TreeNode[G]) (text string, c *color.NRGBA)

// TreeClickFunc 节点单击回调, num 为连续单击次数(双击为2)
type TreeClickFunc[G any] func(gtx layout.Context, node *TreeNode[G], num int)

type Tree[G any] struct {
	Widget[Tree[G]]
	list widget.List

	roots    []*TreeNode[G]
	rootIdx  map[string]*TreeNode[G]
	flat     []*TreeNode[G] //当前可见节点(按展开状态扁平化)
	selected *TreeNode[G]
	clickNum int
	dataSize int //叶子(数据)节点总数

	lock      sync.RWMutex
	sep       string
	keyFunc   func(data *G) string
	textFunc  TreeTextFunc[G]
	clickFunc TreeClickFunc[G]
}

// NewTree 创建树控件
// keyFunc 从数据中取出层级key, 默认按"."分隔, 如 "W3.GPS"
// textFunc 返回节点显示文本, 可为nil(默认显示节点名)
func NewTree[G any](win *Window, keyFunc func(data *G) string, textFunc TreeTextFunc[G]) *Tree[G] {
	m := new(Tree[G])
	m.Init(win, m)
	m.keyFunc = keyFunc
	m.textFunc = textFunc
	m.sep = "."
	m.rootIdx = map[string]*TreeNode[G]{}
	m.list.Axis = layout.Vertical
	return m
}

// SetSeparator 设置key的分隔符, 默认为"."
func (m *Tree[G]) SetSeparator(sep string) *Tree[G] {
	m.sep = sep
	return m
}

// SetOnNodeClick 设置节点单击回调
func (m *Tree[G]) SetOnNodeClick(f TreeClickFunc[G]) *Tree[G] {
	m.clickFunc = f
	return m
}

// Size 叶子(数据)节点总数
func (m *Tree[G]) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.dataSize
}

// GetSelected 获取当前选中的节点
func (m *Tree[G]) GetSelected() *TreeNode[G] {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.selected
}

// UpdateList 用扁平数据重建树
func (m *Tree[G]) UpdateList(list []*G) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.roots = m.roots[:0]
	m.rootIdx = map[string]*TreeNode[G]{}
	m.selected = nil
	m.clickNum = 0
	m.dataSize = len(list)

	for _, item := range list {
		if item == nil {
			continue
		}
		key := ""
		if m.keyFunc != nil {
			key = m.keyFunc(item)
		}
		segs := strings.Split(key, m.sep)
		parent := (*TreeNode[G])(nil)
		idx := m.rootIdx
		siblings := &m.roots
		for i, seg := range segs {
			node, ok := idx[seg]
			if !ok {
				node = &TreeNode[G]{
					Name:      seg,
					Depth:     i,
					Parent:    parent,
					childIdx:  map[string]*TreeNode[G]{},
					clickable: new(widget.Clickable),
				}
				idx[seg] = node
				*siblings = append(*siblings, node)
			}
			if i == len(segs)-1 {
				node.Data = item
			}
			parent = node
			idx = node.childIdx
			siblings = &node.Children
		}
	}
	m.rebuildFlat()
}

// rebuildFlat 按展开状态扁平化可见节点
func (m *Tree[G]) rebuildFlat() {
	m.flat = m.flat[:0]
	var walk func(nodes []*TreeNode[G])
	walk = func(nodes []*TreeNode[G]) {
		for _, n := range nodes {
			m.flat = append(m.flat, n)
			if n.Expanded && len(n.Children) > 0 {
				walk(n.Children)
			}
		}
	}
	walk(m.roots)
}

// ExpandAll 展开/收起所有目录节点
func (m *Tree[G]) ExpandAll(expanded bool) *Tree[G] {
	m.lock.Lock()
	defer m.lock.Unlock()
	var walk func(nodes []*TreeNode[G])
	walk = func(nodes []*TreeNode[G]) {
		for _, n := range nodes {
			if len(n.Children) > 0 {
				n.Expanded = expanded
				walk(n.Children)
			}
		}
	}
	walk(m.roots)
	m.rebuildFlat()
	return m.t
}

func (m *Tree[G]) text(node *TreeNode[G]) (string, *color.NRGBA) {
	if m.textFunc != nil {
		return m.textFunc(node)
	}
	return node.Name, nil
}

func (m *Tree[G]) Layout(gtx C) D {
	m.CheckDimensions(&gtx)
	th := m.Theme()
	return m.list.Layout(gtx, len(m.flat), func(gtx C, index int) D {
		m.lock.RLock()
		node := m.flat[index]
		m.lock.RUnlock()

		if node.clickable.Clicked(gtx) {
			m.lock.Lock()
			if m.selected == node {
				m.clickNum++
			} else {
				m.clickNum = 1
			}
			m.selected = node
			num := m.clickNum
			if !node.IsLeaf() {
				node.Expanded = !node.Expanded
				m.rebuildFlat()
			}
			m.lock.Unlock()
			if m.clickFunc != nil {
				m.clickFunc(gtx, node, num)
			}
		}

		return material.Clickable(gtx, node.clickable, func(gtx C) D {
			return m.layoutNode(gtx, th, node)
		})
	})
}

func (m *Tree[G]) layoutNode(gtx C, th *material.Theme, node *TreeNode[G]) D {
	row := func(gtx C) D {
		return layout.Inset{
			Left:   unit.Dp(float32(node.Depth)*18 + 4),
			Right:  unit.Dp(4),
			Top:    unit.Dp(4),
			Bottom: unit.Dp(4),
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					size := gtx.Dp(unit.Dp(16))
					if node.IsLeaf() {
						return D{Size: image.Pt(size, size)}
					}
					gtx.Constraints.Min.X = size
					gtx.Constraints.Max.X = size
					gtx.Constraints.Min.Y = size
					gtx.Constraints.Max.Y = size
					icon := NaviRightIcon
					if node.Expanded {
						icon = ArrowDropDown
					}
					return icon.Layout(gtx, th.Fg)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
				layout.Flexed(1, func(gtx C) D {
					text, c := m.text(node)
					lbl := material.Body1(th, text)
					if c != nil {
						lbl.Color = *c
					}
					lbl.MaxLines = 1
					lbl.Truncator = "."
					return lbl.Layout(gtx)
				}),
			)
		})
	}

	if m.selected == node {
		macro := op.Record(gtx.Ops)
		dims := row(gtx)
		call := macro.Stop()
		bg := th.Palette.ContrastBg
		bg.A = 40
		paint.FillShape(gtx.Ops, bg, clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Op())
		call.Add(gtx.Ops)
		return dims
	}
	return row(gtx)
}
