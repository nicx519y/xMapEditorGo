package main

// ROOT 根
const ROOT = 0

//Box 控件
type Box struct {
	x          int
	y          int
	width      int
	height     int
	styleClass string
	parent     *Box
	isSelected bool
	isCorrect  bool
	isUsed     bool
	canBubble  bool   //是否继续冒泡
	children   []*Box //子节点 按照z-index排序
}

//Bounds 元素外框
type Bounds struct {
	x      int
	y      int
	width  int
	height int
}

// Position 位置数据
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
	boxROOT := Box{0, 0, width, height, "", nil, false, true, true, true, make([]*Box, 0)}
	tree.AddBox(&boxROOT, nil)
	interactionROOT := Box{0, 0, width, height, "", nil, false, true, true, true, make([]*Box, 0)}
	tree.AddInteractionBox(&interactionROOT, nil)
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

// CreateSimpleBox 创建一个简易的Box
func (t *BoxTree) CreateSimpleBox(x, y, width, height int, styleClass string) *Box {
	return &Box{x, y, width, height, styleClass, nil, false, true, true, true, make([]*Box, 0)}
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

// AddInteractionBox 添加交互层控件
func (t *BoxTree) AddInteractionBox(box *Box, parent *Box) {
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
	list := t.interactionBoxeslist
	idx := 0
	for k, v := range list {
		if v == box {
			idx = k
			break
		}
	}
	t.interactionBoxeslist = append(list[:idx], list[idx+1:]...)
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

// GetPosition 获取绝对位置
func (t *BoxTree) GetPosition(box *Box) Position {
	b := box
	x := box.x
	y := box.y
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
	return Position{x, y}
}

// IsPointInBox 判断点是否碰撞Box
func (t *BoxTree) IsPointInBox(x, y int, box *Box) bool {
	pos := t.GetPosition(box)
	width := box.width
	height := box.height

	return x > pos.x && x < pos.x+width && y > pos.y && y < pos.y+height
}

// HitTestBoxToBox hittest
func (t *BoxTree) HitTestBoxToBox(box1, box2 *Box) bool {
	return t.hitTestBoundsToBounds(t.getBoxBounds(box1), t.getBoxBounds(box2))
}

func (t *BoxTree) getBoxBounds(box *Box) Bounds {
	return Bounds{
		box.x,
		box.y,
		box.width,
		box.height,
	}
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
