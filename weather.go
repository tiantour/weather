package weather

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/tiantour/fetch"
	"github.com/tiantour/rsae"
)

var (
	// AppID app id
	AppID string

	// PrivateKey private key
	PrivateKey string

	weather = map[string]string{
		"00": "晴",
		"1":  "多云",
		"2":  "阴",
		"3":  "阵雨",
		"4":  "雷阵雨",
		"5":  "雷阵雨伴有冰雹",
		"6":  "雨夹雪",
		"7":  "小雨",
		"8":  "中雨",
		"9":  "大雨",
		"10": "暴雨",
		"11": "大暴雨",
		"12": "特大暴雨",
		"13": "阵雪",
		"14": "小雪",
		"15": "中雪",
		"16": "大雪",
		"17": "暴雪",
		"18": "雾",
		"19": "冻雨",
		"20": "沙尘暴",
		"21": "小到中雨",
		"22": "中到大雨",
		"23": "大到暴雨",
		"24": "暴雨到大暴雨",
		"25": "大暴雨到特大暴雨",
		"26": "小到中雪",
		"27": "中到大雪",
		"28": "大到暴雪",
		"29": "浮尘",
		"30": "扬沙",
		"31": "强沙尘暴",
		"53": "霾",
		"99": "无",
	}
	direction = map[string]string{
		"0": "无持续风向",
		"1": "东北风",
		"2": "东风",
		"3": "东南风",
		"4": "南风",
		"5": "西南风",
		"6": "西风",
		"7": "西北风",
		"8": "北风",
		"9": "旋转风",
	}
	rate = map[string]string{
		"0": "微风",
		"1": "3-4级",
		"2": "4-5级",
		"3": "5-6级",
		"4": "6-7级",
		"5": "7-8级",
		"6": "8-9级",
		"7": "9-10级",
		"8": "10-11级",
		"9": "11-12级",
	}
)

type (
	// Alarm alarm
	Alarm struct {
		W1  string `json:"w1,omitempty"`  // 省
		W2  string `json:"w2,omitempty"`  // 市
		W3  string `json:"w3,omitempty"`  // 县
		W4  string `json:"w4,omitempty"`  // 类别编号
		W5  string `json:"w5,omitempty"`  // 类别名称
		W6  string `json:"w6,omitempty"`  // 级别编号
		W7  string `json:"w7,omitempty"`  // 级别名称
		W8  string `json:"w8,omitempty"`  // 时间
		W9  string `json:"w9,omitempty"`  // 内容
		W10 string `json:"w10,omitempty"` // ID
	}
	// City city
	City struct {
		C1  string  `json:"c1,omitempty"`  // 区域编号
		C2  string  `json:"c2,omitempty"`  // 城市英文
		C3  string  `json:"c3,omitempty"`  // 城市中文
		C4  string  `json:"c4,omitempty"`  // 市英文
		C5  string  `json:"c5,omitempty"`  // 市中文
		C6  string  `json:"c6,omitempty"`  // 省英文
		C7  string  `json:"c7,omitempty"`  // 省中文
		C8  string  `json:"c8,omitempty"`  // 国家英文
		C9  string  `json:"c9,omitempty"`  // 国家中文
		C10 string  `json:"c10,omitempty"` // 城市级别
		C11 string  `json:"c11,omitempty"` // 区号
		C12 string  `json:"c12,omitempty"` // 邮编
		C13 float64 `json:"c13,omitempty"` // 经度
		C14 float64 `json:"c14,omitempty"` // 纬度
		C15 string  `json:"c15,omitempty"` // 海拔
		C16 string  `json:"c16,omitempty"` // 雷达
	}
	// Data data
	Data struct {
		F0 string      `json:"f0,omitempty"` // 时间
		F1 []*Forecast `json:"f1,omitempty"` // 天气
	}
	// Forecast forecast
	Forecast struct {
		Fa string `json:"fa,omitempty"` // 白天气象
		Fb string `json:"fb,omitempty"` // 晚上气象
		Fc string `json:"fc,omitempty"` // 白天温度
		Fd string `json:"fd,omitempty"` // 晚上温度
		Fe string `json:"fe,omitempty"` // 白天风向
		Ff string `json:"ff,omitempty"` // 晚上风向
		Fg string `json:"fg,omitempty"` // 白天风力
		Fh string `json:"fh,omitempty"` // 晚上风力
		Fi string `json:"fi,omitempty"` // 日出日落
	}
	// Index index
	Index struct {
		I1 string `json:"i1,omitempty"` // 简称
		I2 string `json:"i2,omitempty"` // 名称
		I3 string `json:"i3,omitempty"` // 别名
		I4 string `json:"i4,omitempty"` // 级别
		I5 string `json:"i5,omitempty"` // 内容
	}
	// Observe observe
	Observe struct {
		L1 string `json:"l1,omitempty"` // 温度
		L2 string `json:"l2,omitempty"` // 湿度
		L3 string `json:"l3,omitempty"` // 风力
		L4 string `json:"l4,omitempty"` // 风向
		L7 string `json:"l7,omitempty"` // 时间
	}
	// Weather weather
	Weather struct {
		A *Alarm   `json:"w,omitempty"` // 预警
		C *City    `json:"c,omitempty"` // 城市
		F *Data    `json:"f,omitempty"` // 常规
		L *Observe `json:"l,omitempty"` // 实况
		I []*Index `json:"i,omitempty"` // 指数
	}
)

// NewWeather new weather
func NewWeather() *Weather {
	return &Weather{}
}

// Alarm get weather index
func (w *Weather) Alarm(placeID int) (*Weather, error) {
	body, err := w.Body(placeID, "alarm")
	if err != nil {
		return nil, err
	}
	result := Weather{}
	err = json.Unmarshal(body, &result)
	return &result, err
}

// Forecast get weather forecast
func (w *Weather) Forecast(placeID int) (*Weather, error) {
	body, err := w.Body(placeID, "forecast")
	if err != nil {
		return nil, err
	}
	result := Weather{}
	err = json.Unmarshal(body, &result)
	for k, v := range result.F.F1 {
		v = w.transform(v)
		result.F.F1[k] = v
	}
	return &result, err
}

// Index get weather index
func (w *Weather) Index(placeID int) (*Weather, error) {
	body, err := w.Body(placeID, "index")
	if err != nil {
		return nil, err
	}
	result := Weather{}
	err = json.Unmarshal(body, &result)
	return &result, err
}

// Observe get weather observe
func (w *Weather) Observe(placeID int) (*Weather, error) {
	body, err := w.Body(placeID, "observe")
	if err != nil {
		return nil, err
	}
	result := Weather{}
	err = json.Unmarshal(body, &result)
	return &result, err
}

// transform transform
func (w *Weather) transform(args *Forecast) *Forecast {
	args.Fa = weather[args.Fa]
	args.Fb = weather[args.Fb]
	args.Fe = direction[args.Fe]
	args.Ff = direction[args.Ff]
	args.Fg = rate[args.Fg]
	args.Fh = rate[args.Fh]
	return args
}

// Body get weather body
func (w *Weather) Body(placeID int, action string) ([]byte, error) {
	date := time.Now().Format("200601021504")
	path := fmt.Sprintf("http://open.weather.com.cn/data/?areaid=%d&type=%s_v&date=%s", placeID, action, date)
	return fetch.Cmd(fetch.Request{
		Method: "GET",
		URL: fmt.Sprintf("%s&appid=%s&key=%s",
			path, AppID[:6],
			w.sign(path),
		),
	})
}

// sign get weathr path sign
func (w *Weather) sign(path string) string {
	publickKey := fmt.Sprintf("%s&appid=%s", path, AppID)
	body := rsae.NewSHA().HmacSha1(publickKey, PrivateKey)
	sign := rsae.NewBase64().Encode(body)
	return url.QueryEscape(sign)
}
