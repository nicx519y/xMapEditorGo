package main

// InteractionBasic 交互基类
type InteractionBasic struct {
}

// BoxHoverAction box hover交互
type BoxHoverAction struct {
	target         *Box
	interActionBox *Box
	engine         *Engine
}

// NewBoxHoverAction 构造函数
func NewBoxHoverAction(target *Box, engine *Engine) (action *BoxHoverAction) {
	action = &BoxHoverAction{}
	action.target = target
	action.engine = engine
	//用于交互的展示
	action.interActionBox = engine.boxTree.CreateSimpleBox(0, 0, 0, 0, "hoverborder")

	engine.mouseEvent.AddEventListener(action.target, MOUSEENTER, action.enterHandler)
	engine.mouseEvent.AddEventListener(action.target, MOUSEOLEAVE, action.leaveHandler)

	return action
}

func (t *BoxHoverAction) enterHandler(evt MouseEvent) {
	target := evt.target
	box := t.interActionBox
	box.x = target.x
	box.y = target.y
	box.width = target.width
	box.height = target.height
	t.engine.boxTree.AddInteractionBox(t.interActionBox, t.engine.boxTree.GetInteractionROOT())
	t.engine.render.PaintRectArea(Rect{box.x, box.y, box.width, box.height}, 1)
}

func (t *BoxHoverAction) leaveHandler(evt MouseEvent) {
	target := evt.target
	t.engine.boxTree.RemoveInteractionBox(t.interActionBox)
	t.engine.render.PaintRectArea(Rect{target.x, target.y, target.width, target.height}, 1)
}
