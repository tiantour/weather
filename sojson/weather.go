package sojson

import (
	"errors"
	"fmt"

	"github.com/duke-git/lancet/v2/netutil"
)

type (
	// Weather weather
	Weather struct {
		Status   int32     `json:"status,omitempty"`   // 错误代码
		Message  string    `json:"message,omitempty"`  // 错误消息
		Date     string    `json:"date,omitempty"`     // 更新日期
		Time     string    `json:"time,omitempty"`     // 更新时间
		CityInfo *CityInfo `json:"cityInfo,omitempty"` // 城市信息
		Data     *Data     `json:"data,omitempty"`     // 每日数据
	}
	CityInfo struct {
		City       string `json:"city,omitempty"`       // 城市名称
		CityID     string `json:"cityId,omitempty"`     // 城市编号
		Parent     string `json:"parent,omitempty"`     // 上一级
		UpdateTime string `json:"updateTime,omitempty"` // 更新时间
	}
	Data struct {
		Humidity  string  `json:"shidu,omitempty"`     // 湿度
		Tem       string  `json:"wendu,omitempty"`     // 温度
		PM10      float32 `json:"pm10,omitempty"`      // pm1.0
		PM25      float32 `json:"pm25,omitempty"`      // pm2.5
		Air       string  `json:"quality,omitempty"`   // 空气指数
		Cold      string  `json:"ganmao,omitempty"`    // 感冒指数
		Yesterday *Date   `json:"yesterday,omitempty"` // 昨天
		Forecast  []*Date `json:"forecast,omitempty"`  // 未来
	}
	Date struct {
		Date      string `json:"ymd,omitempty"`     // 日期
		Day       string `json:"date,omitempty"`    // 天
		Week      string `json:"week,omitempty"`    // 周
		Sunrise   string `json:"sunrise,omitempty"` // 日出
		Sunset    string `json:"sunset,omitempty"`  // 日落
		High      string `json:"high,omitempty"`    // 高温
		Low       string `json:"low,omitempty"`     // 低温
		Air       int32  `json:"aqi,omitempty"`     // 空气质量
		Wind      string `json:"fx,omitempty"`      // 风向
		WindSpeed string `json:"fl,omitempty"`      // 风速
		Type      string `json:"type,omitempty"`    // 天气
		Notice    string `json:"notice,omitempty"`  // 提示
	}
)

// NewWeather new weather
func NewWeather() *Weather {
	return &Weather{}
}

// Fetch fetch weather
func (w *Weather) Fetch(cityID string) (*Weather, error) {
	client := netutil.NewHttpClient()
	resp, err := client.SendRequest(&netutil.HttpRequest{
		RawURL: fmt.Sprintf("http://t.weather.sojson.com/api/weather/city/%s", cityID),
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	result := Weather{}
	err = client.DecodeResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != 200 {
		return nil, errors.New(result.Message)
	}
	return &result, nil
}
