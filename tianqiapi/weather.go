package tianqiapi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tiantour/fetch"
)

var (
	AppID     int32  //  AppID
	AppSecret string // AppSecret
	Version   string // Version
)

type (
	// Weather weather
	Weather struct {
		ErrCode    int32   `json:"errcode,omitempty"`     // 错误代码
		ErrMsg     string  `json:"errmsg,omitempty"`      // 错误消息
		UpdateTime string  `json:"update_time,omitempty"` // 更新时间
		CityID     string  `json:"citydid,omitempty"`     // 城市编号
		City       string  `json:"city,omitempty"`        // 城市名称
		CityEN     string  `json:"cityEN,omitempty"`      // 城市英文名称
		Country    string  `json:"country,omitempty"`     // 国家名称
		CountryEN  string  `json:"countryEN,omitempty"`   // 国家英文名称
		Data       []*Data `json:"data,omitempty"`        // 每日数据
	}
	Data struct {
		Day           string      `json:"day,omitempty"`           // 天
		Date          string      `json:"date,omitempty"`          // 日期
		Week          string      `json:"week,omitempty"`          // 星期
		Wea           string      `json:"wea,omitempty"`           // 天气情况1
		WeaImg        string      `json:"wea_img,omitempty"`       // 固定9种类型(您也可以根据wea字段自己处理)xue、lei、shachen、wu、bingbao、yun、yu、yin、qing
		WeaDay        string      `json:"wea_day,omitempty"`       // 天气情况1
		WeaDayImg     string      `json:"wea_day_img,omitempty"`   // 固定9种类型(您也可以根据wea字段自己处理)xue、lei、shachen、wu、bingbao、yun、yu、yin、qing
		WeaNight      string      `json:"wea_night,omitempty"`     // 天气情况1
		WeaNightImg   string      `json:"wea_night_img,omitempty"` // 固定9种类型(您也可以根据wea字段自己处理)xue、lei、shachen、wu、bingbao、yun、yu、yin、qing
		Tem           string      `json:"tem,omitempty"`           // 实时温度
		Tem1          string      `json:"tem1,omitempty"`          // 高温
		Tem2          string      `json:"tem2,omitempty"`          // 低温
		Humidity      string      `json:"humidity,omitempty"`      // 湿度
		Visibility    string      `json:"visibility,omitempty"`    // 能见度
		Pressure      string      `json:"pressure,omitempty"`      // 气压
		Win           []string    `json:"win,omitempty"`           // 风力
		WinSpeed      string      `json:"win_speed,omitempty"`     // 风力等级
		WinMeter      string      `json:"win_meter,omitempty"`     // 风力等级
		Sunrise       interface{} `json:"sunrise,omitempty"`       // 日出
		Sunset        interface{} `json:"sunset,omitempty"`        // 日落
		Air           string      `json:"air,omitempty"`           // 空气质量
		AirLevel      string      `json:"air_level,omitempty"`     // 空气质量等级
		AirTips       string      `json:"air_tips,omitempty"`      // 空气质量描述
		Phrase        string      `json:"phrase,omitempty"`        // 天气情况短语
		Narrative     string      `json:"narrative,omitempty"`     // 天气情况描述
		Moonrise      interface{} `json:"moonrise,omitempty"`      // 月出
		Moonset       interface{} `json:"moonset,omitempty"`       // 月落
		MoonPhrase    interface{} `json:"moonPhrase,omitempty"`    // 月相
		Rain          string      `json:"rain,omitempty"`          // 降雨概率
		UvIndex       string      `json:"uvIndex,omitempty"`       // 紫外线等级
		UvDescription string      `json:"uvDescription,omitempty"` // 紫外线等级描述
		Alarm         []*Alarm    `json:"alarm,omitempty"`         // 预警
	}
	Alarm struct {
		AlarmType    string `json:"alarm_type,omitempty"`    // 预警类型
		AlarmLevel   string `json:"alarm_level,omitempty"`   // 预警级别
		AlarmTitle   string `json:"alarm_title,omitempty"`   // 预警标题
		AlarmContent string `json:"alarm_content,omitempty"` // 预警内容
	}
)

// NewWeather new weather
func NewWeather() *Weather {
	return &Weather{}
}

// Fetch fetch weather
func (w *Weather) Fetch(cityID string) (*Weather, error) {
	path := fmt.Sprintf("https://v0.yiketianqi.com/api?unescape=1&version=%s&appid=%d&appsecret=%s&ext&cityid=%s", Version, AppID, AppSecret, cityID)
	body, err := fetch.Cmd(&fetch.Request{
		Method: "GET",
		URL:    path,
	})
	if err != nil {
		return nil, err
	}

	var result Weather
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, errors.New(result.ErrMsg)
	}
	return &result, nil
}
