package main

import (
	"fmt"
	"image"
	"math"
	"syscall/js"
	"time"

	"github.com/fogleman/gg"
)

const (
	// BOX 控件
	BOX = iota
	// INTERACTION 交互层
	INTERACTION
)

// Rect 矩形
type Rect struct {
	x, y, width, height int
}

// Viewport 渲染视口
type Viewport struct {
	x       int
	y       int
	width   int
	height  int
	zoom    float64
	context *gg.Context
}

// RenderEngine 渲染器
type RenderEngine struct {
	boxTree    *BoxTree
	styleSheet *StyleSheetManager
	vp         *Viewport
}

// NewRenderEngine 渲染器构造函数
func NewRenderEngine(boxTree *BoxTree, styleSheet *StyleSheetManager) (engine *RenderEngine) {
	engine = &RenderEngine{}
	engine.boxTree = boxTree
	engine.styleSheet = styleSheet
	return engine
}

// PaintBox 简易方法 绘制一个控件
func (t *RenderEngine) PaintBox(box *Box, zoom float64) {
	bounds := box.GetBounds()
	// fmt.Println(bounds)
	go t.PaintRectArea(&Rect{bounds.x, bounds.y, bounds.width, bounds.height}, zoom)
}

// PaintRectArea 绘制一个区域
func (t *RenderEngine) PaintRectArea(rect *Rect, zoom float64) {
	go t.paintRect(rect, zoom)
}

// PaintToScreen 绘制到屏幕
func (t *RenderEngine) PaintToScreen() {
	//绘制
	tvp := t.vp
	c := tvp.context
	rgba := c.Image().(*image.RGBA)
	width := c.Width()
	height := c.Height()

	print := js.NewCallback(func(args []js.Value) {
		if width > 0 && height > 0 {
			js.Global().Get("window").Call("printer", js.TypedArrayOf(rgba.Pix), js.ValueOf(tvp.x), js.ValueOf(tvp.y), js.ValueOf(width), js.ValueOf(height))
		}
	})
	js.Global().Call("requestAnimationFrame", print)

}

// PaintRectArea 绘制一个矩形区域 矩形是缩放前的坐标
func (t *RenderEngine) paintRect(rect *Rect, zoom float64) {
	t.vp = &Viewport{}
	vp := t.vp
	vp.x = round(math.Ceil(float64(rect.x) * zoom))
	vp.y = round(math.Ceil(float64(rect.y) * zoom))
	vp.width = round(math.Ceil(float64(rect.width) * zoom))
	vp.height = round(math.Ceil(float64(rect.height) * zoom))
	vp.zoom = zoom
	vp.context = gg.NewContext(vp.width, vp.height)
	tm := time.Now()
	//填写viewport数据
	t.paintViewport()
	//绘制到屏幕
	t.PaintToScreen()
	fmt.Println("渲染数据准备：", time.Now().Sub(tm))
}

// 绘制一个 视口
func (t *RenderEngine) paintViewport() {
	boxeslist := t.boxTree.GetBoxlist()
	interactionBoxeslist := t.boxTree.GetInteractionBoxeslist()
	if len(boxeslist) > 1 {
		// 绘制控件
		t.renderBox(boxeslist[ROOT], BOX)
	}
	if len(interactionBoxeslist) > 1 {
		// 绘制交互层
		t.renderBox(interactionBoxeslist[ROOT], INTERACTION)
	}
}

//绘制一个容器里的控件
func (t *RenderEngine) renderBoxesInContainer(container *Box, layer int) {

	var list []*Box
	if layer == BOX {
		list = t.boxTree.GetBoxlist()
	} else if layer == INTERACTION {
		list = t.boxTree.GetInteractionBoxeslist()
	}

	for _, v := range list {
		if v.parent == container {
			t.renderBox(v, layer)
		} else {
			continue
		}
	}
}

//绘制控件
func (t *RenderEngine) renderBox(box *Box, layer int) {

	var list []*Box
	if layer == BOX {
		list = t.boxTree.GetBoxlist()
	} else if layer == INTERACTION {
		list = t.boxTree.GetInteractionBoxeslist()
	}

	//如果是 根节点 直接绘制子节点
	if box == list[ROOT] {
		t.renderBoxesInContainer(box, layer)
		return
	}

	//选中的控件或者不可用的控件不渲染
	if !box.isUsed || box.isSelected {
		return
	}

	//绘制
	t.drawBox(box)

	//递归 如果此box是容器，继续绘制里面的元素
	t.renderBoxesInContainer(box, layer)
}

// 按照物理尺寸填充一个矩形区域
func (t *RenderEngine) drawBox(box *Box) {
	style := t.styleSheet.GetStyle(box.styleClass)

	x := box.x - t.vp.x
	y := box.y - t.vp.y
	// cx, cy := box.GetCenterPoint()
	cx, cy := box.width/2+box.x, box.height/2+box.y
	cx -= t.vp.x
	cy -= t.vp.y

	// fmt.Println(cx, cy)
	context := t.vp.context
	context.Push()
	// context.RotateAbout(math.Pi/4, 100, 100)
	context.SetFillStyle(gg.NewSolidPattern(style.backgroundColor))
	context.RotateAbout(box.angle, float64(cx), float64(cy))
	context.DrawRectangle(float64(x), float64(y), float64(box.width), float64(box.height))
	context.Fill()
	context.Pop()
}

//获取跟视口的交集
func (t *RenderEngine) intersectionRect(rect1 Rect) Rect {
	var r Rect
	viewport := t.vp
	r.x = intMax(rect1.x, viewport.x)
	r.y = intMax(rect1.y, viewport.y)

	dx1 := rect1.x + rect1.width
	dy1 := rect1.y + rect1.height
	dx2 := viewport.x + viewport.width
	dy2 := viewport.y + viewport.height

	dx := intMin(dx1, dx2)
	dy := intMin(dy1, dy2)

	r.width = intMax(dx-r.x, 0)
	r.height = intMax(dy-r.y, 0)

	r.x = r.x - viewport.x
	r.y = r.y - viewport.y

	return r
}
