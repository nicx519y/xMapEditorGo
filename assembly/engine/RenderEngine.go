package main

import (
	"fmt"
	"math"
	"syscall/js"
	"time"
)

const (
	// BOX 控件
	BOX = iota
	// INTERACTION 交互层
	INTERACTION
)

// Viewport 视口
type Viewport struct {
	x      int     //视口起点
	y      int     //视口起点
	width  int     //视口宽度
	height int     //视口高度
	zoom   float64 //缩放比
	pixels []uint32
}

//Rect 矩形
type Rect struct {
	x, y, width, height int
}

// RenderEngine 渲染器
type RenderEngine struct {
	boxTree    *BoxTree
	styleSheet *StyleSheetManager
	vp         *Viewport
}

// NewRenderEngine 渲染器构造函数
func NewRenderEngine(boxTree *BoxTree, styleSheet *StyleSheetManager) (engine *RenderEngine) {
	vp := Viewport{0, 0, 0, 0, 1, make([]uint32, 0)}
	engine = &RenderEngine{}
	engine.vp = &vp
	engine.boxTree = boxTree
	engine.styleSheet = styleSheet

	return engine
}

// GetViewport 获取viewport
func (t *RenderEngine) GetViewport() *Viewport {
	return t.vp
}

// ClearViewport 清空viewport数据
func (t *RenderEngine) ClearViewport() {
	t.vp.width = 0
	t.vp.height = 0
	t.vp.pixels = make([]uint32, 0)
}

// PaintRectArea 绘制一个矩形区域 矩形是缩放前的坐标
func (t *RenderEngine) PaintRectArea(rect Rect, zoom float64) {
	vp := t.vp
	vp.x = int(math.Ceil(float64(rect.x) * zoom))
	vp.y = int(math.Ceil(float64(rect.y) * zoom))
	vp.width = int(math.Ceil(float64(rect.width) * zoom))
	vp.height = int(math.Ceil(float64(rect.height) * zoom))
	vp.zoom = zoom
	vp.pixels = make([]uint32, vp.width*vp.height, vp.width*vp.height)
	tm := time.Now()
	//填写viewport数据
	t.paintViewport()
	//绘制到屏幕
	t.PaintToScreen()
	fmt.Println("渲染数据准备：", time.Now().Sub(tm))
}

// PaintToScreen 绘制到屏幕
func (t *RenderEngine) PaintToScreen() {
	//绘制
	tvp := Viewport{t.vp.x, t.vp.y, t.vp.width, t.vp.height, t.vp.zoom, t.vp.pixels}
	print := js.NewCallback(func(args []js.Value) {
		if tvp.width > 0 && tvp.height > 0 {
			js.Global().Get("window").Call("printer", js.TypedArrayOf(tvp.pixels), js.ValueOf(tvp.x), js.ValueOf(tvp.y), js.ValueOf(tvp.width), js.ValueOf(tvp.height))
		}
	})
	js.Global().Call("requestAnimationFrame", print)
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

	//获取绝对坐标
	position := t.boxTree.GetPosition(box)
	//绘制
	rect := Rect{position.x, position.y, box.width, box.height}
	t.fillRect(rect, t.styleSheet.GetStyle(box.styleClass))

	//递归 如果此box是容器，继续绘制里面的元素
	t.renderBoxesInContainer(box, layer)
}

// 按照物理尺寸填充一个矩形区域
func (t *RenderEngine) fillRect(rect Rect, style *Style) {

	var mrect Rect
	viewport := t.vp

	//如果边框粗细为0 并且背景透明 直接返回
	if style.borderWeight <= 0 && style.bgTransparent == true {
		return
	}

	mrect.x = int(math.Ceil(float64(rect.x) * viewport.zoom))
	mrect.y = int(math.Ceil(float64(rect.y) * viewport.zoom))
	mrect.width = int(math.Ceil(float64(rect.width) * viewport.zoom))
	mrect.height = int(math.Ceil(float64(rect.height) * viewport.zoom))
	// 求矩形和视口的交集用于绘制
	nrect := t.intersectionRect(mrect)

	vpx1 := nrect.x
	vpy1 := nrect.y
	vpx2 := vpx1 + nrect.width
	vpy2 := vpy1 + nrect.height

	vbx1 := vpx1 + style.borderWeight
	vbx2 := vpx2 - style.borderWeight
	vby1 := vpy1 + style.borderWeight
	vby2 := vpy2 - style.borderWeight

	boc := style.borderColor
	bgc := style.backgroundColor
	bgt := style.bgTransparent

	for i := vpy1; i < vpy2; i++ {
		for j := vpx1; j < vpx2; j++ {
			if i >= vby1 && i < vby2 && j >= vbx1 && j < vbx2 {
				if !bgt {
					// 绘制底色
					viewport.pixels[j+i*viewport.width] = bgc // a b g r
				} else {
					continue
				}
			} else {
				// 绘制边框
				viewport.pixels[j+i*viewport.width] = boc // a b g r
			}
		}
	}
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
