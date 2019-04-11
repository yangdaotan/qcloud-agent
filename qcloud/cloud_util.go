package qcloud

import (
	"sort"
	"fmt"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"io/ioutil"
	"math/rand"
	"time"
	"net/url"
	"errors"
)


// 公共参数：Action, SecretId, Timestamp, Nonce, Region
// 方法：
//	1. 所有参数做升序
//	2. 拼接请求字符串： Action=xxx&Nonce=xxx...
//	3. 拼接签名字符串：
// 		请求方法+domain+/v2/index.php?Action=xxx&Nonce=xxx...
//	4. 签名：base64_encode(hash_hmac('sha1', $srcStr, $secretKey, true))
//  5. url编码
func GenerateSignature(method string, module Module, hostPath string, params map[string]interface{}) string {
	keys := []string{}
	for key := range params {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))

	srcStr := ""
	for index, sortKey := range keys {
		srcStr += fmt.Sprintf("%v=%v", sortKey, params[sortKey])
		if index != len(keys)-1 {
			srcStr += "&"
		}
	}

	spliceStr := fmt.Sprintf("%v%v.%v%v?%v", method, module, DomainBase, hostPath, srcStr)

	mac := hmac.New(sha1.New, []byte(SecretKey))
	mac.Write([]byte(spliceStr))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}


func GetCommonParamsMap(action Action, region Region) map[string]interface{}  {
	nonce := rand.Intn(100000)

	params := make(map[string]interface{})

	params["Action"] = action
	params["Nonce"] = nonce
	params["Region"] = region
	params["SecretId"] = SecretId
	params["Timestamp"] = time.Now().Unix()

	return params
}

func Get(domain string, params map[string]interface{}) string {

	var values = url.Values{}
	for key, value := range params {
		values.Add(key, fmt.Sprintf("%v",value))
	}
	paramsStr := values.Encode()

	domain = domain + "?" + paramsStr
	resp, err := http.Get(domain)

	if err != nil {
		fmt.Printf("curl error %v\n", domain)
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return  string(body)
}

func CheckSuccess(resp string) (error, interface{}) {
	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	code, _ := dat["code"].(int)

	if code != SUCCESS {
		return errors.New(fmt.Sprintf("%v:%v", code, dat["message"])), nil
	}

	return nil, dat["data"]
}