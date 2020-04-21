package util

import (
	"encoding/json"
	"github.com/go-resty/resty"
	"log"
	"time"
)

// GET
func Get(url string, headers map[string]string, inter interface{}) error {
	if headers == nil {
		headers = make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded"
		headers["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"
	}

	resp, err := resty.New().SetTimeout(time.Minute * 1).R().
		SetHeaders(headers).
		Get(url)
	// log.Println("[Get]", resp.String())

	err = json.Unmarshal(resp.Body(), &inter)
	if err != nil {
		log.Print("[net][Get]", err, url)
	}
	return err
}

// POST
func Post(url string, headers map[string]string, params interface{}, inter interface{}) error {
	if headers == nil {
		headers = make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded"
		headers["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"
	}

	resp, err := resty.New().SetTimeout(time.Minute * 1).R().
		SetHeaders(headers).
		SetBody(params).
		Post(url)
	// log.Println("[PostForm]",resp.String())

	if inter != nil {
		err = json.Unmarshal(resp.Body(), &inter)
		if err != nil {
			log.Print("[net][PostForm]", resp.String(), err, url)
		}
	}
	return err
}

// POST 表单
func PostForm(url string, headers map[string]string, params map[string]string, inter interface{}) error {
	if headers == nil {
		headers = make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded"
		headers["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"
	}

	resp, err := resty.New().SetTimeout(time.Minute * 1).R().
		SetHeaders(headers).
		SetFormData(params).
		Post(url)

	if err != nil {
		log.Print("[PostForm]", err, "[url]", url)
	}
	// log.Println("[PostForm]",resp.String())

	if inter != nil {
		err = json.Unmarshal(resp.Body(), &inter)
		if err != nil {
			log.Print("[net][PostForm]", err, url)
		}
	}
	return err
}
