package main

import (
	"math/rand"
	"syscall/js"
)

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
		engien.CreateNewBox(x, y, w, h, "")
	})
	defer addRectHandler.Release()
	doc.Call("getElementById", "add-rect-btn").Call("addEventListener", "click", addRectHandler)

	// 离开页面之前关闭线程
	destroyHandler := js.NewCallback(func(args []js.Value) {
		done <- 0
	})
	defer destroyHandler.Release()
	js.Global().Get("window").Call("addEventListener", "beforeunload", destroyHandler)

	<-done
}
