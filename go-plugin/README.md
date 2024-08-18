# Golang 原生插件

# 将模块添加到工作空间
go work use .

# 插件编译
go build -buildmode=plugin -o plugin.so plugin.go

# 运行插件主程序
go run cmd/*.go -from 3333 -to google.com

# 执行测试
curl localhost:3333