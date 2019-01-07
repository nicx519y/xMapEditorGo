package main

// ROOT 根
const ROOT = 0

//Box 控件
type Box struct {
	x          int
	y          int
	width      int
	height     int
	angle      float64
	styleClass string
	parent     *Box
	isSelected bool
	isCorrect  bool
	isUsed     bool
	canBubble  bool   //是否继续冒泡
	children   []*Box //子节点 按照z-index排序
}

// NewBox 构造函数
func NewBox(x, y, width, height int, styleClass string) (box *Box) {
	box = &Box{}
	box.x = x
	box.y = y
	box.width = width
	box.height = height
	box.angle = 0
	box.styleClass = styleClass

	box.parent = nil
	box.isCorrect = true
	box.isSelected = false
	box.isUsed = true
	box.canBubble = true
	box.children = make([]*Box, 0)

	return box
}

// IsSelected getter 选中状态
func (t *Box) IsSelected() bool {
	return t.isSelected
}

// SetIsSelected setter 选中状态
func (t *Box) SetIsSelected(status bool) {
	if status == t.IsSelected() {
		return
	}
	t.isSelected = status
}

// GetPosition 获取绝对位置坐标
func (t *Box) GetPosition() (int, int) {
	b := t
	x := b.x
	y := b.y
	var parent *Box
	for {
		parent = b.parent
		if parent != nil {
			x += parent.x
			y += parent.y
			b = parent
		} else {
			break
		}
	}
	return x, y
}

// GetBounds 获取bounds
func (t *Box) GetBounds() (bounds Bounds) {
	bounds = Bounds{}

	px, py := t.GetPosition()

	if t.angle == 0 {
		bounds.x = px
		bounds.y = py
		bounds.width = t.width
		bounds.height = t.height
	} else {
		cx, cy := t.x+t.width/2, t.y+t.height/2

		x1, y1 := px, py
		x2, y2 := px+t.width, py
		x3, y3 := px+t.width, py+t.height
		x4, y4 := px, py+t.height

		rx1, ry1 := srotate(t.angle, x1, y1, cx, cy)
		rx2, ry2 := srotate(t.angle, x2, y2, cx, cy)
		rx3, ry3 := srotate(t.angle, x3, y3, cx, cy)
		rx4, ry4 := srotate(t.angle, x4, y4, cx, cy)

		bounds.x = intMin(rx1, rx2, rx3, rx4) - 2
		bounds.y = intMin(ry1, ry2, ry3, ry4) - 2

		bounds.width = intMax(rx1, rx2, rx3, rx4) - bounds.x + 2
		bounds.height = intMax(ry1, ry2, ry3, ry4) - bounds.y + 2

	}

	return bounds
}

// GetCenterPoint 获取box原点
func (t *Box) GetCenterPoint() (x int, y int) {
	px, py := t.GetPosition()
	x = round(float64(t.x) + float64(t.width)/2)
	y = round(float64(t.y) + float64(t.height)/2)
	return x + px, y + py
}

//Bounds 元素外框
type Bounds struct {
	x      int
	y      int
	width  int
	height int
}

// Position 位置
type Position struct {
	x, y int
}

// BoxTree box树
type BoxTree struct {
	boxeslist            []*Box
	interactionBoxeslist []*Box
}

// NewBoxTree 构造函数
func NewBoxTree(width, height int) (tree *BoxTree) {
	tree = &BoxTree{}
	tree.boxeslist = make([]*Box, 0)
	tree.interactionBoxeslist = make([]*Box, 0)
	boxROOT := NewBox(0, 0, width, height, "")
	tree.AddBox(boxROOT, nil)
	interactionROOT := NewBox(0, 0, width, height, "")
	tree.AddInteractionBox(interactionROOT, nil)
	return tree
}

// GetBoxROOT 获取控件根节点
func (t *BoxTree) GetBoxROOT() *Box {
	return t.boxeslist[ROOT]
}

// GetBoxlist 回去控件列表
func (t *BoxTree) GetBoxlist() []*Box {
	return t.boxeslist
}

// GetInteractionROOT 获取交互层根节点
func (t *BoxTree) GetInteractionROOT() *Box {
	return t.interactionBoxeslist[ROOT]
}

// GetInteractionBoxeslist 获取交互层的控件列表
func (t *BoxTree) GetInteractionBoxeslist() []*Box {
	return t.interactionBoxeslist
}

// AddBox 添加一个box到boxtree
func (t *BoxTree) AddBox(box *Box, parent *Box) {
	boxlist := t.boxeslist
	if parent != nil {
		box.parent = parent
		children := box.parent.children
		box.parent.children = append(children, box)
	}
	t.boxeslist = append(boxlist, box)

}

// DisableBox 禁用控件
func (t *BoxTree) DisableBox(box *Box) {
	//设置为不可用
	box.isUsed = false
	boxeslist := t.GetBoxlist()
	for _, b := range boxeslist {
		//将以本控件为容器的控件设置为不可用
		if b.parent == box {
			t.DisableBox(b)
		}
	}
}

// EnableBox 恢复控件
func (t *BoxTree) EnableBox(box *Box) {
	box.isUsed = true
	boxeslist := t.GetBoxlist()
	for _, b := range boxeslist {
		//将以本控件为容器的控件设置为不可用
		if b.parent == box {
			t.EnableBox(b)
		}
	}
}

// IndexOfInteractionBox 返回box 所在的第一个位置索引
func (t *BoxTree) IndexOfInteractionBox(box *Box) (idx int) {
	i := -1
	for idx, val := range t.interactionBoxeslist {
		if val == box {
			i = idx
			break
		}
	}
	return i
}

// AddInteractionBox 添加交互层控件
func (t *BoxTree) AddInteractionBox(box *Box, parent *Box) {
	if t.IndexOfInteractionBox(box) > 0 {
		return
	}
	list := t.interactionBoxeslist
	if parent != nil {
		box.parent = parent
		children := box.parent.children
		box.parent.children = append(children, box)
	}
	t.interactionBoxeslist = append(list, box)

}

// RemoveInteractionBox 删除交互层控件
func (t *BoxTree) RemoveInteractionBox(box *Box) {
	idx := 0
	for idx < len(t.interactionBoxeslist) {
		if t.interactionBoxeslist[idx] == box {
			t.interactionBoxeslist = append(t.interactionBoxeslist[:idx], t.interactionBoxeslist[idx+1:]...)
			continue
		}
		idx++
	}
}

// ClearInteractionBoxes 清空交互层
func (t *BoxTree) ClearInteractionBoxes() {
	t.interactionBoxeslist = t.interactionBoxeslist[:1]
}

// GetAllPrevAndNext 找出所有兄弟节点
func (t *BoxTree) GetAllPrevAndNext(box *Box) []*Box {
	if box.parent == nil {
		return make([]*Box, 0, 0)
	}
	list := t.getChildren(box.parent, false)
	var i int
	for k, v := range list {
		if v == box {
			i = k
			break
		}
	}
	//去重自身
	result := append(list[:i], list[:(i+1)]...)
	return result
}

// GetBox 获取控件对象
func (t *BoxTree) GetBox(idx int) *Box {
	return t.GetBoxlist()[idx]
}

// ResetIsCorrect 重新计算合法性
func (t *BoxTree) ResetIsCorrect(box *Box) {
	box.isCorrect = !t.hitTestBoxToAllBoxes(box, t.GetAllPrevAndNext(box))
}

// IsPointInBox 判断点是否碰撞Box
func (t *BoxTree) IsPointInBox(x, y int, box *Box) bool {
	px, py := box.GetPosition()
	width := box.width
	height := box.height

	return x > px && x < px+width && y > py && y < py+height
}

// HitTestBoxToBox hittest
func (t *BoxTree) HitTestBoxToBox(box1, box2 *Box) bool {
	return t.hitTestBoundsToBounds(box1.GetBounds(), box2.GetBounds())
}

// 找出所有子节点 deep 是否深度遍历
func (t *BoxTree) getChildren(box *Box, deep bool) []*Box {
	children := make([]*Box, 0)
	boxeslist := t.GetBoxlist()
	for _, p := range boxeslist {
		if p.parent == box {
			children = append(children, p)
			if !deep {
				continue
			}
			children = append(children, t.getChildren(box, true)...)
		}
	}
	return children
}

func (t *BoxTree) parseBoolToUint8(value bool) uint8 {
	if value == true {
		return uint8(1)
	}
	return uint8(0)
}

func (t *BoxTree) hitTestBoundsToBounds(bounds1, bounds2 Bounds) bool {
	return intAbs((bounds1.x+bounds1.width/2)-(bounds2.x-bounds2.width/2)) < (bounds1.width+bounds2.width)/2 && intAbs((bounds1.y+bounds1.height/2)-(bounds2.y-bounds2.height/2)) < (bounds1.height+bounds2.height)/2
}

// 控件和所有元素hitTest
func (t *BoxTree) hitTestBoxToAllBoxes(box *Box, allBoxes []*Box) bool {
	result := false
	for _, i := range allBoxes {
		if box.isUsed && box.isCorrect && box.isSelected && t.HitTestBoxToBox(box, i) {
			result = true
			break
		}
	}
	return result
}

//获取一个元素的树形层级
func (t *BoxTree) getDeep(box *Box) int {
	// ROOT
	if box.parent == nil {
		return 0
	}

	deep := 0
	boxeslist := t.GetBoxlist()
	var parent *Box
	b := box
	for {
		parent = b.parent
		deep++
		if parent == boxeslist[ROOT] {
			break
		} else {
			b = parent
			continue
		}
	}
	return deep
}
