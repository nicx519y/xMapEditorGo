package main

import "fmt"

// BoxStateInterface 交互状态对象
type BoxStateInterface interface {
	Target() *Box                                 //目标box
	IsRunning() bool                              //判断是否在状态内运行
	Views() []BoxStateViewInterface               //获取视图列表
	Stop()                                        //停止状态
	Start()                                       //开始状态
	RegisterEvents(handlers func(event BoxEvent)) //注册行为
}

// BoxBasicState 状态基类
type BoxBasicState struct {
	target        *Box
	engine        *Engine
	isRunning     bool
	views         []BoxStateViewInterface
	eventsHandler func(evt BoxEvent)
}

// Target 获取目标box
func (t *BoxBasicState) Target() *Box {
	return t.target
}

// IsRunning 是否在运行
func (t *BoxBasicState) IsRunning() bool {
	return t.isRunning
}

// Views 获取视图列表
func (t *BoxBasicState) Views() []BoxStateViewInterface {
	return t.views
}

// Stop 停止状态
func (t *BoxBasicState) Stop() {
	t.isRunning = false
	for _, view := range t.views {
		view.Close()
	}
}

// Start 开始状态
func (t *BoxBasicState) Start() {
	t.isRunning = true
	for _, view := range t.views {
		view.Render()
	}
}

// RegisterEvents 注册监听器
func (t *BoxBasicState) RegisterEvents(handler func(event BoxEvent)) {
	t.eventsHandler = handler
}

// initState 初始化
func (t *BoxBasicState) initState(target *Box, engine *Engine) {
	t.target = target
	t.engine = engine
	t.isRunning = false
	t.views = make([]BoxStateViewInterface, 0, 1)
}

// 添加鼠标事件监听
func (t *BoxBasicState) addEventListener(target *Box, EventType MouseEventType, handler EventHandler) (listener *EventListener) {
	return t.engine.mouseEvent.AddEventListener(target, EventType, func(evt MouseEvent) {
		if t.isRunning {
			handler(evt)
		}
	})
}

// 删除鼠标事件监听
func (t *BoxBasicState) removeEventListener(target *Box, listener *EventListener) {
	t.engine.mouseEvent.RemoveEventListener(target, listener)
}

/////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  BoxNormalState start ///////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////

// BoxNormalState 正常状态
type BoxNormalState struct {
	BoxBasicState
	enterListener *EventListener
}

// NewBoxNormalState 构造函数
func NewBoxNormalState(target *Box, engine *Engine) (state *BoxNormalState) {
	state = &BoxNormalState{}
	state.initState(target, engine)
	state.isRunning = false
	return state
}

// Stop 停止状态
func (t *BoxNormalState) Stop() {
	t.isRunning = false
	for _, view := range t.views {
		view.Close()
	}
	t.removeEventListener(t.target, t.enterListener)
}

// Start 开始状态
func (t *BoxNormalState) Start() {
	t.isRunning = true
	for _, view := range t.views {
		view.Render()
	}
	t.enterListener = t.addEventListener(t.target, MOUSEENTER, func(evt MouseEvent) {
		if t.eventsHandler != nil {
			t.eventsHandler(IN)
		}
	})
}

/////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  BoxHoverState start ////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////

// BoxHoverState Hover状态
type BoxHoverState struct {
	BoxBasicState
	leaveListener     *EventListener
	dragstartListener *EventListener
}

// NewBoxHoverState 构造函数
func NewBoxHoverState(target *Box, engine *Engine) (state *BoxHoverState) {
	state = &BoxHoverState{}
	state.initState(target, engine)
	state.views = append(state.views, NewBoxBorderView(target, "hoverborder", engine))
	return state
}

// Stop 停止状态
func (t *BoxHoverState) Stop() {
	t.isRunning = false
	for _, view := range t.views {
		view.Close()
	}
	t.removeEventListener(t.target, t.dragstartListener)
	t.removeEventListener(t.target, t.leaveListener)
}

// Start 开始状态
func (t *BoxHoverState) Start() {
	t.isRunning = true
	for _, view := range t.views {
		view.Render()
	}
	t.dragstartListener = t.addEventListener(t.target, DRAGSTART, func(evt MouseEvent) {
		if t.eventsHandler != nil {
			t.eventsHandler(MOVESTART)
		}
	})
	t.leaveListener = t.addEventListener(t.target, MOUSELEAVE, func(evt MouseEvent) {
		if t.eventsHandler != nil {
			t.eventsHandler(OUT)
		}
	})
}

/////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  BoxMoveState start /////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////

// BoxMoveState 平移状态
type BoxMoveState struct {
	BoxBasicState
	dragListener    *EventListener
	dragendListener *EventListener
}

// NewBoxMoveState 构造函数
func NewBoxMoveState(target *Box, engine *Engine) (state *BoxMoveState) {
	state = &BoxMoveState{}
	state.initState(target, engine)
	state.views = append(state.views, NewBoxBorderView(target, "hoverborder", engine))
	return state
}

// Stop 停止状态
func (t *BoxMoveState) Stop() {
	t.isRunning = false
	for _, view := range t.views {
		view.Close()
	}
	t.removeEventListener(t.target, t.dragListener)
	t.removeEventListener(t.target, t.dragendListener)
}

// Start 开始状态
func (t *BoxMoveState) Start() {
	t.isRunning = true
	for _, view := range t.views {
		view.Render()
	}

	// drag 刷新target视图
	t.dragListener = t.addEventListener(t.target, DRAG, func(evt MouseEvent) {
		target := t.target
		ob := target.GetBounds()
		target.x = evt.mouseX - evt.data.x
		target.y = evt.mouseY - evt.data.y
		b := target.GetBounds()
		for _, view := range t.views {
			view.Refresh()
		}
		rect := &Rect{}
		rect.x = intMin(ob.x, b.x)
		rect.y = intMin(ob.y, b.y)
		rect.width = intMax(ob.x+ob.width, b.x+b.width) - rect.x
		rect.height = intMax(ob.y+ob.height, b.y+b.height) - rect.y
		t.engine.render.PaintRectArea(rect, 1)
	})

	t.dragendListener = t.addEventListener(t.target, DRAGEND, func(evt MouseEvent) {
		fmt.Println("drag end")
		if t.eventsHandler != nil {
			t.eventsHandler(MOVEEND)
		}
	})
}
