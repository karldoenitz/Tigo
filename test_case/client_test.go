package test_case

import (
	"encoding/json"
	"fmt"
	"github.com/karldoenitz/Tigo/request"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	response, err := request.Get("https://life.qq.com/api/activity/detail?id=773947310848622080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	contentStr := response.ToContentStr()
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	_ = json.Unmarshal(response.Content, &result)
	fmt.Println(result.Code)
	fmt.Println(result.Msg)
	fmt.Println(contentStr)
}

func TestPost(t *testing.T) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	postData := map[string]interface{}{
		"chlid": "news_news_bj",
	}
	response, err := request.Post("https://life.qq.com/api/activity/get_good_act_list?cachedCount=0", postData, headers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	contentStr := response.ToContentStr()
	fmt.Println(contentStr)
}

func TestMakeRequest(t *testing.T) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	postData := "chlid=news_news_bj"
	response, err := request.MakeRequest("POST", "https://life.qq.com/api/activity/get_good_act_list?cachedCount=0", strings.NewReader(postData), headers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	contentStr := response.ToContentStr()
	fmt.Println(contentStr)
}

func TestMap2Xml(t *testing.T) {
	testMap := map[string]interface{}{
		"k1": "v1",
		"k2": []string{"1", "2", "3"},
		"k3": map[string]int{
			"kk1": 1,
			"kk2": 2,
		},
	}
	result := request.Map2Xml(testMap)
	t.Logf("xml result => %s", result)
}
