//go:build js && wasm

package main

import "syscall/js"

func fib(i int) int {
	if i == 0 || i == 1 {
		return 1
	}
	return fib(i-1) + fib(i-2)
}

func fibFunc(this js.Value, args []js.Value) interface{} {
	return js.ValueOf(fib(args[0].Int()))
}

func main() {
	// wasm调用Javascript函数
	alert := js.Global().Get("alert")
	alert.Invoke("Hello World!")

	// 注册函数供Javascript调用
	done := make(chan int, 0)
	js.Global().Set("fibFunc", js.FuncOf(fibFunc))
	<-done

}
