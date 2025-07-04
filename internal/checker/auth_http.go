package checker

import (
	"net/http"
	"time"
)

var apiURL = "http://127.0.0.1:3000/api/check?uuid="

func SetAuthAPI(url string) {
	apiURL = url
}

func CheckUUIDViaAPI(uuid string) bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(apiURL + uuid)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
} 