package weather

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/tiantour/conf"
	"github.com/tiantour/imago"
	"github.com/tiantour/requests"
)

// Weather weather
var (
	Weather     = &weather{}
	appID       = conf.Options.Weather.ID
	privateKey  = conf.Options.Weather.Key
	weatherInfo = map[string]string{
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
	windDirection = map[string]string{
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
	windSize = map[string]string{
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

type weather struct{}

// GetWeatherData 转换
func (w *weather) Data(weatherArea, weatherType string) (weatherData map[string]interface{}, err error) {
	requestURL, requestData, requestHeader := requests.Options()
	requestURL = w.url(weatherArea, weatherType)
	body, err := requests.Get(requestURL, requestData, requestHeader)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &weatherData)
	if weatherType == "forecast_v" {
		weatherData = w.tranform(weatherData)
	}
	return
}

// getWeatherURL 获取天气地址
func (w *weather) url(weatherArea, weatherType string) (weatherURL string) {
	timeString := time.Now().Format("200601021504")
	tempURL := "http://open.weather.com.cn/data/?areaid=" + weatherArea + "&type=" + weatherType + "&date=" + timeString
	publicKey := tempURL + "&appid=" + appID
	signatureKey := imago.Crypto.Base64Encode(imago.Crypto.HmacSha1(publicKey, privateKey))
	weatherURL = tempURL + "&appid=" + appID[:6] + "&key=" + url.QueryEscape(signatureKey)
	return
}

// tranformWeatherData 转换
func (w *weather) tranform(inputData map[string]interface{}) (outPutData map[string]interface{}) {
	forecast := inputData["f"].(map[string]interface{})["f1"].([]interface{})
	for _, v := range forecast {
		t := v.(map[string]interface{})
		faItem := t["fa"].(string)
		if faItem != "" {
			t["fa"] = weatherInfo[faItem]
		}
		feItem := t["fe"].(string)
		if faItem != "" {
			t["fe"] = windDirection[feItem]
		}
		fgItem := t["fg"].(string)
		if fgItem != "" {
			t["fg"] = windSize[fgItem]
		}
		t["fb"] = weatherInfo[t["fb"].(string)]
		t["ff"] = windDirection[t["ff"].(string)]
		t["fh"] = windSize[t["fh"].(string)]
		temp := strings.Split(t["fi"].(string), "|")
		t["fi"] = temp[0]
		t["fj"] = temp[1]
	}
	outPutData = inputData
	return
}
