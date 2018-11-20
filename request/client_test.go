package request

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	response, err := Get("https://life.qq.com/api/activity/detail?id=773947310848622080")
	if err != nil {
		fmt.Println(err.Error())
	}
	contentStr := response.ToContentStr()
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	json.Unmarshal(response.Content, &result)
	fmt.Println(result.Code)
	fmt.Println(result.Msg)
	fmt.Println(contentStr)
}

func TestPost(t *testing.T)  {
	headers := map[string] string {
		"Content-Type": "application/x-www-form-urlencoded",
	}
	postData := map[string]interface{} {
		"chlid": "news_news_bj",
	}
	response, err := Post("https://life.qq.com/api/activity/get_good_act_list?cachedCount=0", postData, headers)
	if err != nil {
		fmt.Println(err.Error())
	}
	contentStr := response.ToContentStr()
	fmt.Println(contentStr)
}

func TestPostForm(t *testing.T)  {
	requestUrl := "http://captcha.qq.com/mcheck"
	postData := map[string] interface {} {
		"aid": "9288506",
		"AppSecretKey": "1aXd6Sj_iUbBpKd8u99joPA",
		"SceneType": "3",
		"AccountType": "2",
		"AccountId": "oelUQ5XXDC08aJqMOgRyUEsP2hwQ",
		"AccountAppid": "wx304df47a4abc2a7f",
		"UserIP": "121.206.165.104",
	}
	headers := map[string] string {
		"Content-Type": "application/x-www-form-urlencoded",
	}
	response, _ := Post(requestUrl, postData, headers)
	fmt.Println(response.ToContentStr())
}