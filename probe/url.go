package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UrlRes struct {
	Url            string
	Cert           float64
	Latency        float64
	HttpStatusCode float64
}

// UrlProbe 发起一个 url 拨测
func UrlProbe(urls []string, headers map[string]string, limit, timeout int64) []UrlRes {
	chLimit := make(chan bool, limit) //控制并发访问量
	chs := make([]chan UrlRes, len(urls))

	limitFunc := func(chLimit chan bool, ch chan UrlRes, url string) {
		urlProbe(url, headers, timeout, ch)
		<-chLimit
	}
	for i, url := range urls {
		chs[i] = make(chan UrlRes, 1)
		chLimit <- true
		go limitFunc(chLimit, chs[i], url)
	}
	result := []UrlRes{}
	for _, ch := range chs {
		res := <-ch
		result = append(result, res)
	}
	return result
}

func urlProbe(url string, headers map[string]string, timeout int64, ch chan UrlRes) {
	urlRes := urlCheck(url, headers, timeout)
	ch <- urlRes
	return
}

func urlCheck(url string, headers map[string]string, timeout int64) (res UrlRes) {
	res.Url = url
	now := time.Now()

	var statusCode int
	var err error
	statusCode, err = httpGet(url, headers, false, timeout)
	if err != nil && strings.Contains(err.Error(), "certificate") {
		now = time.Now()
		res.Cert = -1.0
		statusCode, err = httpGet(url, headers, true, timeout)
	}
	end := time.Now()
	d := end.Sub(now)

	if err != nil {
		res.Latency = -1.0
		return
	}

	res.Cert = 1

	rttStr := fmt.Sprintf("%.3f", float64(d.Nanoseconds())/1000000.0)
	rtt, _ := strconv.ParseFloat(rttStr, 64)
	res.Latency = rtt
	res.HttpStatusCode = float64(statusCode)
	return
}

func httpGet(url string, headers map[string]string, skipCert bool, timeout int64) (statusCode int, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipCert},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(timeout),
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	return
}
