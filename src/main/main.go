package main

import (
	"HystrixTest/src/app"
	"HystrixTest/src/hystrixConfig"
)

func main() {
	hystrixConfig.NewHystrix()
	app.SetUpRouter()
}
