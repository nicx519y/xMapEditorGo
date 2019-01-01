package main

import "syscall/js"

const (
	// MOUSEDOWN 鼠标事件
	MOUSEDOWN = iota
	// MOUSEUP 鼠标事件
	MOUSEUP
	// MOUSEMOVE 鼠标事件
	MOUSEMOVE
	// CLICK 鼠标事件
	CLICK
	// DBCLICK 鼠标事件
	DBCLICK
	// MOUSEENTER 鼠标进入
	MOUSEENTER
	// MOUSEOLEAVE 鼠标离开
	MOUSEOLEAVE
)

// MouseEvent 鼠标事件对象
type MouseEvent struct {
	target    *Box
	eventType int
	mouseX    int
	mouseY    int
}

// EventHandler 事件回调函数
type EventHandler func(evt MouseEvent)

// EventAction 事件行为 描述一个事件行为 事件类型 和 callback
type EventAction struct {
	eventType int
	handler   EventHandler
}

// MouseEventManager 鼠标事件管理器
type MouseEventManager struct {
	boxTree         *BoxTree
	eventActionList map[*Box][]*EventAction
	eventTopBox     *Box
}

// NewMouseEventManager 构造函数
func NewMouseEventManager(boxTree *BoxTree) (manager *MouseEventManager) {
	manager = &MouseEventManager{}
	manager.eventActionList = make(map[*Box][]*EventAction)
	manager.eventTopBox = nil
	manager.boxTree = boxTree

	//系统鼠标事件接收
	sysMouseHandler := js.NewCallback(func(args []js.Value) {
		//chrome or firefox
		x := args[0].Get("offsetX").Int()
		y := args[0].Get("offsetY").Int()
		if x == 0 {
			x = args[0].Get("layerX").Int()
		}
		if y == 0 {
			y = args[0].Get("layerY").Int()
		}
		manager.dispatherEvents(args[0].Get("type").String(), x, y)
	})

	mainBox := js.Global().Get("document").Call("getElementById", "main-box")
	mainBox.Call("addEventListener", "mousedown", sysMouseHandler)
	mainBox.Call("addEventListener", "mouseup", sysMouseHandler)
	mainBox.Call("addEventListener", "mousemove", sysMouseHandler)
	mainBox.Call("addEventListener", "click", sysMouseHandler)
	mainBox.Call("addEventListener", "dblclick", sysMouseHandler)

	// defer sysMouseHandler.Release()
	return manager
}

// AddEventListener 添加事件监听器
func (t *MouseEventManager) AddEventListener(target *Box, eventType int, callback EventHandler) {
	_, ok := t.eventActionList[target]
	if !ok {
		t.eventActionList[target] = make([]*EventAction, 0)
	}
	action := EventAction{eventType, callback}
	targetActions := t.eventActionList[target]
	t.eventActionList[target] = append(targetActions, &action)
}

// RemoveEvents 清空事件
func (t *MouseEventManager) RemoveEvents(target *Box) {
	t.eventActionList[target] = make([]*EventAction, 0)
}

// DispatherEvents 接收系统事件 分发事件
func (t *MouseEventManager) dispatherEvents(eventType string, mx, my int) {

	list := t.getBubblingList(mx, my)
	bubblingLen := len(list)

	// mouseenter and mouseleave
	// if (bubblingLen == 0 && t.eventTopBox != nil) || list[bubblingLen-1] != t.eventTopBox {
	// 	t.dispatchBoxEvents(t.eventTopBox, MOUSEOLEAVE, mx, my)
	// }
	root := t.boxTree.GetBoxROOT()
	if bubblingLen > 0 && t.eventTopBox != list[bubblingLen-1] {
		if t.eventTopBox != nil && t.eventTopBox != root {
			t.dispatchBoxEvents(t.eventTopBox, MOUSEOLEAVE, mx, my)
		}
		t.eventTopBox = list[bubblingLen-1]
		if t.eventTopBox != root {
			t.dispatchBoxEvents(t.eventTopBox, MOUSEENTER, mx, my)
		}
	}

	if bubblingLen <= 0 {
		return
	}

	etype := 0
	switch eventType {
	case "mousedown":
		etype = MOUSEDOWN
		break
	case "mouseup":
		etype = MOUSEUP
		break
	case "mousemove":
		etype = MOUSEMOVE
		break
	case "click":
		etype = CLICK
		break
	case "dblclick":
		etype = DBCLICK
		break
	}
	// fmt.Println("事件类型：", eventType, "；冒泡序列：", list)
	for _, box := range list {
		t.dispatchBoxEvents(box, etype, mx, my)
		//如果禁止冒泡 跳出循环
		if !box.canBubble {
			break
		}
	}
}

//获取冒泡list
func (t *MouseEventManager) getBubblingList(x, y int) []*Box {
	boxeslist := t.boxTree.GetBoxlist()
	list := make([]*Box, 0)
	i := boxeslist[ROOT] //从根节点开始找

	if !t.boxTree.IsPointInBox(x, y, i) {
		return list
	}
	list = append(list, i)
	for {
		children := i.children
		//没有子节点
		if len(children) <= 0 {
			break
		}
		isChanged := false
		for k := len(children) - 1; k >= 0; k-- {
			if t.boxTree.IsPointInBox(x, y, children[k]) {
				list = append(list, children[k])
				i = children[k]
				isChanged = true
				break
			}
		}
		//子节点没有冒泡
		if !isChanged {
			break
		}
	}

	return list
}

//分发一个元素的事件
func (t *MouseEventManager) dispatchBoxEvents(target *Box, eventType, mx, my int) {
	_, ok := t.eventActionList[target]
	if !ok {
		return
	}
	for _, action := range t.eventActionList[target] {
		if eventType == action.eventType {
			action.handler(MouseEvent{target, eventType, mx, my})
		}
	}
}
