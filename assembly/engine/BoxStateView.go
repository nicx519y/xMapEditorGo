package main

// BoxStateViewInterface 交互层视图
type BoxStateViewInterface interface {
	Target() *Box
	Render()
	Refresh()
	Close()
}

// BoxBorderView 交互层边框视图
type BoxBorderView struct {
	target            *Box
	interactionTarget *Box
	engine            *Engine
}

// NewBoxBorderView 构造函数
func NewBoxBorderView(target *Box, styleClass string, engine *Engine) (view *BoxBorderView) {
	view = &BoxBorderView{}
	view.target = target
	view.engine = engine
	// view.interactionTarget = NewBox(0, 0, 0, 0, styleClass)
	return view
}

// Target 获取目标
func (t *BoxBorderView) Target() *Box {
	return t.target
}

// Render 渲染视图
func (t *BoxBorderView) Render() {
	t.Refresh()
	if t.interactionTarget != nil {
		t.engine.boxTree.AddInteractionBox(t.interactionTarget, t.engine.boxTree.GetInteractionROOT())
		t.engine.render.PaintBox(t.target, 1)
	}
}

// Refresh 刷新
func (t *BoxBorderView) Refresh() {
	if t.interactionTarget != nil {
		px, py := t.target.GetPosition()
		t.interactionTarget.x = px
		t.interactionTarget.y = py
		t.interactionTarget.width = t.target.width
		t.interactionTarget.height = t.target.height
	}
}

// Close 关闭渲染
func (t *BoxBorderView) Close() {
	if t.interactionTarget != nil {
		t.engine.boxTree.RemoveInteractionBox(t.interactionTarget)
		t.engine.render.PaintBox(t.target, 1)
	}
}
