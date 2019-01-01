package main

import (
	"syscall/js"
	"time"
)

// Engine 引擎定义
type Engine struct {
	boxTree    *BoxTree
	mouseEvent *MouseEventManager
	styleSheet *StyleSheetManager
	render     *RenderEngine
}

// NewEngine 构造函数
func NewEngine() (engine *Engine) {
	engine = &Engine{}
	initStage := js.NewCallback(func(args []js.Value) {
		engine.boxTree = NewBoxTree(args[0].Int(), args[1].Int())
		engine.mouseEvent = NewMouseEventManager(engine.boxTree)
		engine.styleSheet = NewStyleSheetManager()
		engine.render = NewRenderEngine(engine.boxTree, engine.styleSheet)
	})
	js.Global().Get("window").Call("isReady", initStage)
	return engine
}

// CreateNewBox 创建一个新的控件
func (t *Engine) CreateNewBox(x, y, width, height int, styleClass string) {

	styleName := time.Now().String()
	style := t.styleSheet.GetRandStyle()
	t.styleSheet.AddStyle(styleName, style)

	box := t.boxTree.CreateSimpleBox(x, y, width, height, styleName)
	box.parent = t.boxTree.GetBoxROOT()
	t.boxTree.AddBox(box)

	t.render.PaintRectArea(Rect{box.x, box.y, box.width, box.height}, 1)
	NewBoxHoverAction(box, t)
}
