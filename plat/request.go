package plat

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
	}

	resp, err := resty.New().SetTimeout(time.Minute * 1).R().
		SetHeaders(headers).
		Get(url)
	// log.Print("[Get]",resp.String())

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
	}

	resp, err := resty.New().SetTimeout(time.Minute * 1).R().
		SetHeaders(headers).
		SetBody(params).
		Post(url)

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
	}

	resp, err := resty.New().SetTimeout(time.Minute * 1).R().
		SetHeaders(headers).
		SetFormData(params).
		Post(url)

	if err != nil {
		log.Print("[PostForm]", err, "[url]", url)
	}
	// log.Print("[PostForm]",resp.String())

	if inter != nil {
		err = json.Unmarshal(resp.Body(), &inter)
		if err != nil {
			log.Print("[net][PostForm]", err, url)
		}
	}
	return err
}
