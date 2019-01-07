package main

import (
	"fmt"
	"syscall/js"
)

// MouseEventType 鼠标事件类型
type MouseEventType int

const (
	// MOUSEDOWN 鼠标事件
	MOUSEDOWN MouseEventType = iota
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
	// MOUSELEAVE 鼠标离开
	MOUSELEAVE
	// DRAG 拖拽
	DRAG
	// DRAGSTART 拖拽开始
	DRAGSTART
	// DRAGEND 拖拽结束
	DRAGEND
)

const (
	// NONE 无状态
	NONE = iota
	// BEFOREDRAG 拖拽之前
	BEFOREDRAG
	// DRAGING 拖拽中
	DRAGING
)

// MouseEvent 鼠标事件对象
type MouseEvent struct {
	target    *Box
	eventType MouseEventType
	mouseX    int //绝对坐标
	mouseY    int
	data      *Position
}

// EventHandler 事件回调函数
type EventHandler func(evt MouseEvent)

// EventListener 事件行为监听器 描述一个事件行为 事件类型 和 callback
type EventListener struct {
	eventType MouseEventType
	handler   EventHandler
}

// DragStartState 记录drag的状态
type DragStartState struct {
	target    *Box
	state     int
	x         int
	y         int
	dPosition *Position
}

// MouseEventManager 鼠标事件管理器
type MouseEventManager struct {
	boxTree         *BoxTree
	eventActionList map[*Box][]*EventListener
	eventTopBox     *Box
	dragState       *DragStartState
}

// NewMouseEventManager 构造函数
func NewMouseEventManager(boxTree *BoxTree) (manager *MouseEventManager) {
	manager = &MouseEventManager{}
	manager.eventActionList = make(map[*Box][]*EventListener)
	manager.eventTopBox = nil
	manager.boxTree = boxTree
	manager.dragState = &DragStartState{}

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
func (t *MouseEventManager) AddEventListener(target *Box, eventType MouseEventType, callback EventHandler) (listener *EventListener) {
	_, ok := t.eventActionList[target]
	if !ok {
		t.eventActionList[target] = make([]*EventListener, 0)
	}
	listener = &EventListener{eventType, callback}
	targetActions := t.eventActionList[target]
	t.eventActionList[target] = append(targetActions, listener)
	return listener
}

// RemoveEventListener  删除时间监听器
func (t *MouseEventManager) RemoveEventListener(target *Box, listener *EventListener) {
	_, ok := t.eventActionList[target]
	if !ok {
		return
	}
	idx := 0
	for idx < len(t.eventActionList[target]) {
		if t.eventActionList[target][idx] == listener {
			t.eventActionList[target] = append(t.eventActionList[target][:idx], t.eventActionList[target][idx+1:]...)
			continue
		}
		idx++
	}
}

// RemoveEvents 清空事件
func (t *MouseEventManager) RemoveEvents(target *Box) {
	delete(t.eventActionList, target)
}

// DispatherEvents 接收系统事件 分发事件
func (t *MouseEventManager) dispatherEvents(eventType string, mx, my int) {

	list := t.getBubblingList(mx, my)
	bubblingLen := len(list)

	if bubblingLen <= 0 {
		return
	}

	root := t.boxTree.GetBoxROOT()

	// MouseEnter和MouseLeave事件
	if bubblingLen > 0 && t.eventTopBox != list[bubblingLen-1] {
		if t.eventTopBox != nil && t.eventTopBox != root {
			t.dispatchBoxEvents(t.eventTopBox, MOUSELEAVE, mx, my, nil)
		}
		t.eventTopBox = list[bubblingLen-1]
		if t.eventTopBox != root {
			t.dispatchBoxEvents(t.eventTopBox, MOUSEENTER, mx, my, nil)
		}
	}

	var etype MouseEventType
	dragState := t.dragState
	switch eventType {
	case "mousedown":
		etype = MOUSEDOWN
		px, py := list[bubblingLen-1].GetPosition()
		//如果存在冒泡队列，取最上面一个作为拖拽对象记录
		dragState.target = list[bubblingLen-1]
		dragState.state = BEFOREDRAG
		dragState.x = mx
		dragState.y = my
		dragState.dPosition = &Position{mx - px, my - py}
		break
	case "mouseup":
		etype = MOUSEUP
		fmt.Println("mouseup", dragState.state)
		if dragState.state == DRAGING {
			t.dispatchBoxEvents(dragState.target, DRAGEND, mx, my, dragState.dPosition)
		}
		if dragState.state != NONE {
			dragState.state = NONE
		}
		break
	case "mousemove":
		etype = MOUSEMOVE
		if dragState.state == BEFOREDRAG {
			dragState.state = DRAGING
			t.dispatchBoxEvents(dragState.target, DRAGSTART, dragState.x, dragState.y, dragState.dPosition)
		}
		if dragState.state == DRAGING {
			t.dispatchBoxEvents(dragState.target, DRAG, mx, my, dragState.dPosition)
		}
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
		t.dispatchBoxEvents(box, etype, mx, my, nil)
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
func (t *MouseEventManager) dispatchBoxEvents(target *Box, eventType MouseEventType, mx, my int, data *Position) {
	_, ok := t.eventActionList[target]
	if !ok {
		return
	}
	for _, action := range t.eventActionList[target] {
		if eventType == action.eventType {
			action.handler(MouseEvent{target, eventType, mx, my, data})
		}
	}
}
