package main

import (
	"math"
	"math/rand"
	"syscall/js"
)

// var n int

func main() {

	done := make(chan int, 0)
	doc := js.Global().Get("document")

	engien := NewEngine()
	// 添加矩形
	addRectHandler := js.NewCallback(func(args []js.Value) {
		x := rand.Intn(500)
		y := rand.Intn(600)
		w := rand.Intn(300)
		h := rand.Intn(300)
		// x := 300
		// y := 200
		// w := 250
		// h := 50
		// var angle float64
		// if n%2 == 0 {
		// 	angle = math.Pi / 2
		// } else {
		// 	angle = float64(0)
		// }
		// n++
		// angle := float64(0)
		angle := rand.Float64() * 2 * math.Pi
		engien.CreateNewBox(x, y, w, h, angle, "")
		// engien.CreateNewBox(0, 0, 200, 200, "")
	})
	defer addRectHandler.Release()
	doc.Call("getElementById", "add-rect-btn").Call("addEventListener", "click", addRectHandler)

	// 离开页面之前关闭线程
	destroyHandler := js.NewCallback(func(args []js.Value) {
		done <- 0
	})
	defer destroyHandler.Release()
	js.Global().Get("window").Call("addEventListener", "beforeunload", destroyHandler)

	// const S = 1024
	// dc := gg.NewContext(S, S)
	// dc.SetRGBA(0, 0, 0, 0.1)
	// for i := 0; i < 360; i += 15 {
	// 	dc.Push()
	// 	dc.RotateAbout(gg.Radians(float64(i)), S/2, S/2)
	// 	dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
	// 	dc.Fill()
	// 	dc.Pop()
	// }
	// dc.SavePNG("out.png")
	// // fmt.Println(dc.Image())

	<-done
}
