# WebAssembly
具有紧凑的二进制格式，可以接近原生的性能。

Go语言内置的 syscall/js 包，可以在Go语言中直接调用JavaScript函数，包括DOM树相关操作。

## 编译wasm
GOOS=js GOARCH=wasm go build -o static/main.wasm

## 前端需要依赖wasm_exec.js
将其拷贝到static文件夹
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static