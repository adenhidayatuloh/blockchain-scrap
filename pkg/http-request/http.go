package httprequest

import (
	"blockchain-scrap/pkg/errs"
	"bytes"
	"fmt"
	"io"
	"net/http"

	_ "github.com/joho/godotenv"
)

// func ProcessRequest(url string) ([]byte, errs.MessageErr) {
// 	body, err := ProcessJSONRequest("GET", url, nil, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return body, nil
// }

// func ProcessPostRequest(url string, req []byte) ([]byte, errs.MessageErr) {
// 	body, err := ProcessJSONRequest("POST", url, req, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return body, nil
// }

func ProcessJSONRequest(method, url string, payload []byte, headers map[string]string) ([]byte, errs.MessageErr) {
	var req *http.Request
	var err error

	switch method {
	case "POST":
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
		if err != nil {
			return nil, errs.NewInternalServerError("failed to create POST request: " + err.Error())
		}

	case "GET":
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, errs.NewInternalServerError("failed to create GET request: " + err.Error())
		}

	default:
		return nil, errs.NewBadRequest("unsupported HTTP method: " + method)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.NewInternalServerError("request failed: " + err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errs.NewInternalServerError("failed to read response body: " + err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return body, errs.NewInternalServerError(fmt.Sprintf("unexpected status code %d: %s", res.StatusCode, string(body)))
	}

	return body, nil
}
