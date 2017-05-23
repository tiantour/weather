# weather

get weather.com.cn api with go

### Index

```
package main

import (
	"fmt"

	"github.com/tiantour/weather"
)

func init() {
	weather.AppID = "your appid"
	weather.PrivateKey = "your private key"
}

func main() {
	area := "101251201"
	index, err := weather.NewWeather().Index(area)
	fmt.Print(index, err)
}
```


### Forecast

```
package main

import (
	"fmt"

	"github.com/tiantour/weather"
)

func init() {
	weather.AppID = "your appid"
	weather.PrivateKey = "your private key"
}

func main() {
	area := "101251201"
	forecast, err := weather.NewWeather().Forecast(area)
	fmt.Print(forecast, err)
}
```