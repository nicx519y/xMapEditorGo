package main

import "fmt"

// BoxState box的状态枚举
type BoxState int

const (
	// NORMAL 正常状态
	NORMAL BoxState = iota
	// HOVER 鼠标hover
	HOVER
	// SELECTED 选中
	SELECTED
	// MOVE 移动
	MOVE
	// STRETCH 拉伸
	STRETCH
	// HSTRETCH 水平拉伸
	HSTRETCH
	// VSTRETCH 垂直拉伸
	VSTRETCH
	// ROTATE 旋转
	ROTATE
)

// BoxEvent Box行为
type BoxEvent int

const (
	//IN 鼠标进入
	IN BoxEvent = iota
	// OUT 鼠标移开
	OUT
	// SELECT 选择
	SELECT
	// MOVESTART 开始移动
	MOVESTART
	// MOVEEND 结束移动
	MOVEEND
	// STRETCHSTART 开始拉伸
	STRETCHSTART
	// HSTRETCHSTART 开始水平拉伸
	HSTRETCHSTART
	// VSTRETCHSTART 开始垂直拉伸
	VSTRETCHSTART
	// STRETCHEND 结束拉伸
	STRETCHEND
	// HSTRETCHEND 结束水平拉伸
	HSTRETCHEND
	// VSTRETCHEND 结束垂直拉伸
	VSTRETCHEND
	// ROTATESTART 开始旋转
	ROTATESTART
	// ROTATEEND 结束旋转
	ROTATEEND
)

// BoxStateMachineFactroy 工厂方法
func BoxStateMachineFactroy(target *Box, engine *Engine) (machine *BoxStateMachine) {
	machine = NewBoxStateMachine(target, engine)

	states := make(map[BoxState]BoxStateInterface)
	states[NORMAL] = NewBoxNormalState(target, engine)
	states[HOVER] = NewBoxHoverState(target, engine)
	states[MOVE] = NewBoxMoveState(target, engine)
	machine.AddStates(states)

	events := make(map[BoxEvent]BoxState)
	events[IN] = HOVER
	events[OUT] = NORMAL
	events[MOVESTART] = MOVE
	events[MOVEEND] = HOVER
	machine.AddEvents(events)

	return machine
}

// BoxStateMachine 交互状态机
type BoxStateMachine struct {
	target       *Box
	engine       *Engine
	currentState BoxState
	// 描述每个状态
	states map[BoxState]BoxStateInterface
	// 描述行为 以及行为导致的状态跳转
	events map[BoxEvent]BoxState
}

// NewBoxStateMachine 构造函数
func NewBoxStateMachine(target *Box, engine *Engine) (stateMachine *BoxStateMachine) {
	stateMachine = &BoxStateMachine{}
	stateMachine.target = target
	stateMachine.engine = engine
	stateMachine.states = make(map[BoxState]BoxStateInterface)
	stateMachine.events = make(map[BoxEvent]BoxState)

	//默认Normal状态
	stateMachine.currentState = NORMAL

	return stateMachine
}

// AddStates 添加状态
func (t *BoxStateMachine) AddStates(states map[BoxState]BoxStateInterface) {
	for state, obj := range states {
		t.states[state] = obj
		t.states[state].RegisterEvents(t.eventsDispatchHandler)

		if state == NORMAL {
			obj.Start()
		}
	}
}

// AddEvents 添加行为
func (t *BoxStateMachine) AddEvents(events map[BoxEvent]BoxState) {
	for event, state := range events {
		t.events[event] = state
	}
}

// OpenState 状态跳转
func (t *BoxStateMachine) OpenState(state BoxState) {
	if t.currentState == state {
		return
	}
	t.closeCurrentState()
	t.currentState = state
	_, hasState := t.states[state]
	if hasState {
		t.states[state].Start()
	}

	// log for debug
	switch state {
	case NORMAL:
		fmt.Println("状态跳转：normal")
	case HOVER:
		fmt.Println("状态跳转：hover")
	case SELECTED:
		fmt.Println("状态跳转：selected")
	case MOVE:
		fmt.Println("状态跳转：move")
	case STRETCH:
		fmt.Println("状态跳转：stretch")
	case HSTRETCH:
		fmt.Println("状态跳转：hstretch")
	case VSTRETCH:
		fmt.Println("状态跳转：vstretch")
	case ROTATE:
		fmt.Println("状态跳转：rotate")
	default:
		break
	}
}

func (t *BoxStateMachine) closeCurrentState() {
	_, hasState := t.states[t.currentState]
	if hasState {
		t.states[t.currentState].Stop()
	}
}

func (t *BoxStateMachine) eventsDispatchHandler(event BoxEvent) {
	_, ok := t.events[event]
	if ok {
		t.OpenState(t.events[event])
	}
}
