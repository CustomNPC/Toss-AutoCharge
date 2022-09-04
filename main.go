package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	tossId string
	name   string
	amount string
}

func contains(s []float64, e float64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	var charged = make([]float64, 0)

	server := gin.Default()

	server.POST("/check", func(res *gin.Context) {
		requestData, err := ioutil.ReadAll(res.Request.Body)

		var jsonBody map[string]interface{}
		json.Unmarshal(requestData, &jsonBody)

		if jsonBody["tossId"] == nil || jsonBody["name"] == nil || jsonBody["amount"] == nil {
			res.JSON(400, gin.H{
				"success": false,
				"message": "One or more component is missing or invalid",
			})
			return
		}

		tossId := jsonBody["tossId"].(string)
		name := jsonBody["name"].(string)
		amount := jsonBody["amount"].(float64)

		req, err := http.NewRequest("POST", "https://api-gateway.toss.im:11099/api-public/v3/cashtag/transfer-feed/received/list?inputWord="+tossId, strings.NewReader("{}"))
		if err != nil {
			res.JSON(400, gin.H{
				"success": false,
				"message": "cannot create a request",
			})
			return
		}

		for k, v := range map[string]string{
			"Accept":             "application/json, text/plain, */*",
			"Accept-Language":    "en-US,en;q=0.9",
			"Cache-Control":      "no-cache",
			"Connection":         "keep-alive",
			"Content-Type":       "application/json",
			"Origin":             "https://toss.me",
			"Pragma":             "no-cache",
			"Referer":            "https://toss.me/",
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "cross-site",
			"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
			"X-Toss-Method":      "GET",
			"sec-ch-ua":          "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"",
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": "\"macOS\"",
		} {
			req.Header.Set(k, v)
		}

		response, err := http.DefaultClient.Do(req)
		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			res.JSON(400, gin.H{
				"success": false,
				"message": "cannot read body",
			})
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)

		if data["resultType"] != "SUCCESS" {
			res.JSON(400, gin.H{
				"success": false,
				"message": data["error"].(map[string]interface{})["reason"],
			})
			return
		}

		r, size := utf8.DecodeRuneInString(name)
		_name := string(r) + "*" + name[size*2:]

		for _, v := range data["success"].(map[string]interface{})["data"].([]interface{}) {
			value := v.(map[string]interface{})
			tranferId := value["cashtagTransferId"].(float64)
			containResult := contains(charged, tranferId)

			if value["amount"] == amount && (value["senderDisplayName"] == name || value["senderDisplayName"] == _name) && containResult == false {
				charged = append(charged, tranferId)
				res.JSON(200, gin.H{
					"success": true,
					"found":   true,
				})
				return
			}
		}
		res.JSON(200, gin.H{
			"success": true,
			"found":   false,
		})
		return
	})

	server.Run()
}
