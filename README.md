# weather

weather skd for tianqiapi.com v91

# how to use

```
package main

import (
	"fmt"
	"time"

	"github.com/tiantour/weather"
)

func main() {
	weather.AppID = "you appid"
	weather.AppSecret = "you appsecret"
	weather.Version = "v91"

	result, err := weather.NewWeather().FetchWithTTL(cityID, 1, 7200*time.Second)
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range result.Data {
		fmt.Println(k, v)
	}
}

```