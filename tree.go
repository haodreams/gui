/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2026-11-04
 * @FilePath: \gui\tree.go
 * @Description: 通用树控件, 由树形数据(TreeChild)构建展示
 */
package gui

import (
	"image"
	"image/color"
	"sync"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// TreeChild 树形数据, Label 为节点显示名, Children 为空表示叶子节点
type TreeChild[G any] struct {
	Label    string
	Data     *G
	Children []*TreeChild[G]
}

// TreeNode 树节点, Data 为 nil 表示目录节点, 否则为叶子(数据)节点
type TreeNode[G any] struct {
	Name     string         //节点名称
	Data     *G             //叶子节点携带的数据
	Parent   *TreeNode[G]   //父节点
	Children []*TreeNode[G] //子节点
	Depth    int            //层级深度, 根为0
	Expanded bool           //是否展开

	path      string //全路径, 用于重建树时恢复展开和选中状态
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

// TreeRightClickFunc 节点右键回调, 返回 true 时弹出右键菜单
type TreeRightClickFunc[G any] func(gtx layout.Context, node *TreeNode[G]) bool

type Tree[G any] struct {
	Widget[Tree[G]]
	list widget.List

	roots      []*TreeNode[G]
	flat       []*TreeNode[G] //当前可见节点(按展开状态扁平化)
	selected   *TreeNode[G]
	selectedId string
	clickNum   int
	dataSize   int  //叶子(数据)节点总数
	dirty      bool //展开状态已变化, 需要在下一帧重建可见节点

	lock      sync.RWMutex
	textFunc  TreeTextFunc[G]
	clickFunc TreeClickFunc[G]
	rightFunc TreeRightClickFunc[G]

	menu      *Menu
	menuOpen  bool
	menuPos   image.Point
	rowHeight int
}

// NewTree 创建树控件
// textFunc 返回节点显示文本, 可为nil(默认显示节点名)
func NewTree[G any](win *Window, textFunc TreeTextFunc[G]) *Tree[G] {
	m := new(Tree[G])
	m.Init(win, m)
	m.textFunc = textFunc
	m.list.Axis = layout.Vertical
	return m
}

// SetOnNodeClick 设置节点单击回调
func (m *Tree[G]) SetOnNodeClick(f TreeClickFunc[G]) *Tree[G] {
	m.clickFunc = f
	return m
}

// SetOnNodeRightClick 设置节点右键回调
func (m *Tree[G]) SetOnNodeRightClick(f TreeRightClickFunc[G]) *Tree[G] {
	m.rightFunc = f
	return m
}

// SetMenu 设置右键弹出菜单
func (m *Tree[G]) SetMenu(menu *Menu) *Tree[G] {
	m.menu = menu
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

// UpdateData 用树形数据重建树, 可传入一个或多个根节点
// 重建时保留已有节点的展开状态和选中状态
func (m *Tree[G]) UpdateData(children ...TreeChild[G]) {
	m.lock.Lock()
	defer m.lock.Unlock()

	expanded := map[string]bool{}
	var collect func(nodes []*TreeNode[G])
	collect = func(nodes []*TreeNode[G]) {
		for _, n := range nodes {
			expanded[n.path] = n.Expanded
			collect(n.Children)
		}
	}
	collect(m.roots)

	m.roots = m.roots[:0]
	m.selected = nil
	m.clickNum = 0
	m.dataSize = 0

	for i := range children {
		if node := m.buildChild(&children[i], nil, 0, expanded); node != nil {
			m.roots = append(m.roots, node)
		}
	}
	m.rebuildFlat()

	if m.selectedId != "" {
		m.selected = m.findByPath(m.roots, m.selectedId)
	}
}

// buildChild 由 TreeChild 递归构建 TreeNode
func (m *Tree[G]) buildChild(c *TreeChild[G], parent *TreeNode[G], depth int, expanded map[string]bool) *TreeNode[G] {
	if c == nil {
		return nil
	}
	node := &TreeNode[G]{
		Name:      c.Label,
		Data:      c.Data,
		Parent:    parent,
		Depth:     depth,
		clickable: new(widget.Clickable),
	}
	if parent == nil {
		node.path = c.Label
	} else {
		node.path = parent.path + "\x1f" + c.Label
	}
	if c.Data != nil {
		m.dataSize++
	}
	if len(c.Children) > 0 {
		if v, ok := expanded[node.path]; ok {
			node.Expanded = v
		}
		for _, child := range c.Children {
			if n := m.buildChild(child, node, depth+1, expanded); n != nil {
				node.Children = append(node.Children, n)
			}
		}
	}
	return node
}

// findByPath 按全路径查找节点
func (m *Tree[G]) findByPath(nodes []*TreeNode[G], path string) *TreeNode[G] {
	for _, n := range nodes {
		if n.path == path {
			return n
		}
		if found := m.findByPath(n.Children, path); found != nil {
			return found
		}
	}
	return nil
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

// ExpandRoots 展开所有根目录
func (m *Tree[G]) ExpandRoots() *Tree[G] {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, n := range m.roots {
		if len(n.Children) > 0 {
			n.Expanded = true
		}
	}
	m.rebuildFlat()
	return m.t
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

	// 展开状态在上一帧已变化, 这里重建可见节点
	m.lock.Lock()
	if m.dirty {
		m.rebuildFlat()
		m.dirty = false
	}
	count := len(m.flat)
	m.lock.Unlock()

	list := m.list.Layout(gtx, count, func(gtx C, index int) D {
		m.lock.RLock()
		//数据可能被异步重新查询替换, 防止索引越界
		if index >= len(m.flat) {
			m.lock.RUnlock()
			return D{}
		}
		node := m.flat[index]
		isSelected := m.selected == node
		m.lock.RUnlock()

		if node.clickable.Clicked(gtx) {
			m.lock.Lock()
			if m.selected == node {
				m.clickNum++
			} else {
				m.clickNum = 1
			}
			m.selected = node
			m.selectedId = node.path
			num := m.clickNum
			if !node.IsLeaf() {
				//不能在布局过程中改变可见节点数量, 留到下一帧重建
				node.Expanded = !node.Expanded
				m.dirty = true
			}
			m.menuOpen = false
			m.lock.Unlock()
			m.win.Invalidate()
			if m.clickFunc != nil {
				m.clickFunc(gtx, node, num)
			}
		}

		for {
			evt, ok := gtx.Source.Event(pointer.Filter{
				Target: node.clickable,
				Kinds:  pointer.Press,
			})
			if !ok {
				break
			}
			e, ok := evt.(pointer.Event)
			if !ok {
				continue
			}
			if e.Buttons != pointer.ButtonSecondary {
				continue
			}
			//右键: 按回调返回决定是否弹出菜单
			show := m.rightFunc != nil && m.rightFunc(gtx, node)
			m.menuOpen = show && m.menu != nil
			if m.menuOpen {
				y := index*m.rowHeight - m.list.Position.Offset + int(e.Position.Y)
				m.menuPos = image.Pt(int(e.Position.X), y)
				m.win.Invalidate()
			}
		}

		dims := material.Clickable(gtx, node.clickable, func(gtx C) D {
			return m.layoutNode(gtx, th, node, isSelected)
		})
		m.rowHeight = dims.Size.Y
		return dims
	})

	if !m.menuOpen || m.menu == nil {
		return list
	}
	if m.menu.Clicked(gtx) {
		m.menuOpen = false
		m.win.Invalidate()
	}
	//菜单弹出层, 定位到右键位置
	return layout.Stack{Alignment: layout.NW}.Layout(gtx,
		layout.Stacked(func(gtx C) D { return list }),
		layout.Stacked(func(gtx C) D {
			trans := op.Offset(m.menuPos).Push(gtx.Ops)
			defer trans.Pop()
			return m.menu.Layout(gtx)
		}),
	)
}

func (m *Tree[G]) layoutNode(gtx C, th *material.Theme, node *TreeNode[G], selected bool) D {
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

	if selected {
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
