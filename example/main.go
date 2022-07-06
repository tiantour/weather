package main

import (
	"fmt"

	"github.com/tiantour/weather/sojson"
	"github.com/tiantour/weather/tianqiapi"
)

func main() {
	// use sojson
	x, err := sojson.NewWeather().Fetch("101030100")
	fmt.Println(x, err)

	// use tianqiapi
	tianqiapi.AppID = 0
	tianqiapi.AppSecret = "your appsecret"
	y, err := tianqiapi.NewWeather().Fetch("101030100")
	fmt.Println(y, err)
}
