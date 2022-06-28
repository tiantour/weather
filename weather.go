package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/tiantour/fetch"
)

var (
	cache     *ristretto.Cache
	AppID     int32  //  AppID
	AppSecret string // AppSecret
	Version   string // Version
)

func init() {
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters:        1e7,     // number of keys to track frequency of (10M).
		MaxCost:            1 << 30, // maximum cost of cache (1GB).
		BufferItems:        64,      // number of keys per Get buffer.
		IgnoreInternalCost: true,
	})
	if err != nil {
		log.Fatal(err)
	}
}

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
		Wea           string   `json:"wea,omitempty"`           // 天气情况1
		WeaImg        string   `json:"wea_img,omitempty"`       // 固定9种类型(您也可以根据wea字段自己处理)xue、lei、shachen、wu、bingbao、yun、yu、yin、qing
		Phrase        string   `json:"phrase,omitempty"`        // 天气情况短语
		Narrative     string   `json:"narrative,omitempty"`     // 天气情况描述
		Tem           string   `json:"tem,omitempty"`           // 实时温度
		Tem1          string   `json:"tem1,omitempty"`          // 高温
		Tem2          string   `json:"tem2,omitempty"`          // 低温
		Win           []string `json:"win,omitempty"`           // 风力
		WinSpeed      string   `json:"win_speed,omitempty"`     // 风力等级
		Humidity      string   `json:"humidity,omitempty"`      // 湿度
		Sunrise       string   `json:"sunrise,omitempty"`       // 日出
		Sunset        string   `json:"sunset,omitempty"`        // 日落
		Moonrise      string   `json:"moonrise,omitempty"`      // 月出
		Moonset       string   `json:"moonset,omitempty"`       // 月落
		MoonPhrase    string   `json:"moonPhrase,omitempty"`    // 月相
		Visibility    string   `json:"visibility,omitempty"`    // 能见度
		Pressure      string   `json:"pressure,omitempty"`      // 气压
		Air           string   `json:"air,omitempty"`           // 空气质量
		AirLevel      string   `json:"air_level,omitempty"`     // 空气质量等级
		AirTips       string   `json:"air_tips,omitempty"`      // 空气质量描述
		Rain          string   `json:"rain,omitempty"`          // 降雨概率
		UvIndex       string   `json:"uvIndex,omitempty"`       // 紫外线等级
		UvDescription string   `json:"uvDescription,omitempty"` // 紫外线等级描述
		Alarm         []*Alarm `json:"alarm,omitempty"`         // 预警

	}
	Alarm struct {
		AlarmType    string `json:"alarm_type,omitempty"`    // 预警类型
		AlarmLevel   string `json:"alarm_level,omitempty"`   // 预警级别
		AlarmTitle   string `json:"alarm_title,omitempty"`   // 预警标题
		AlarmContent string `json:"alarm_content,omitempty"` // 预警内容
	}
	Week struct {
		Date          string  `json:"date,omitempty"`            // 日期
		Week          string  `json:"week,omitempty"`            // 星期
		DayTem        string  `json:"day_tem,omitempty"`         // 日期温度
		DayWea        string  `json:"day_wea,omitempty"`         // 日期天气
		DayWeaImg     string  `json:"day_wea_img,omitempty"`     // 日期天气图片
		DayWeaIcon    string  `json:"day_wea_icon,omitempty"`    // 日期天气图标
		DayWin        string  `json:"day_win,omitempty"`         // 日期风力
		DayWinSpeed   string  `json:"day_win_speed,omitempty"`   // 日期风速
		NightTem      string  `json:"night_tem,omitempty"`       // 夜晚温度
		NightWea      string  `json:"night_wea,omitempty"`       // 夜晚天气
		NightWeaImg   string  `json:"night_wea_img,omitempty"`   // 夜晚天气图片
		NightWeaIcon  string  `json:"night_wea_icon,omitempty"`  // 夜晚天气图标
		NightWin      string  `json:"night_win,omitempty"`       // 夜晚风力
		NightWinSpeed string  `json:"night_win_speed,omitempty"` // 夜晚风速度
		Hours         []*Hour `json:"hours,omitempty"`           // 小时风力
	}
	Hour struct {
		Time     string `json:"time,omitempty"`      // 时间
		Icon     string `json:"icon,omitempty"`      // Icon
		Wea      string `json:"wea,omitempty"`       // 天气
		WeaImg   string `json:"wea_img,omitempty"`   // 天气图片
		Rain     string `json:"rain,omitempty"`      // 下雨
		Win      string `json:"win,omitempty"`       // 风力
		WinMeter string `json:"win_meter,omitempty"` // 风速
	}
)

// NewWeather new weather
func NewWeather() *Weather {
	return &Weather{}
}

func (w *Weather) FetchWithTTL(cityID string, cost int64, ttl time.Duration) (*Weather, error) {
	cache.Wait()
	result, ok := cache.Get(cityID)
	if ok {
		return result.(*Weather), nil
	}

	body, err := w.Fetch(cityID)
	if err != nil {
		return nil, err
	}

	_ = cache.SetWithTTL(cityID, body, cost, ttl)
	return body, err
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
