package main

import (
	"testing"
)

func Test_PingProbe(t *testing.T) {
	ips := []string{"114.114.114.114", "114.114.115.115", "39.156.69.79"}
	var timeout int64 = 2
	var limit int64 = 1
	res := PingProbe(ips, limit, timeout)
	t.Log("rtt:", res)
}

func Test_UrlProbe(t *testing.T) {
	url := []string{"https://bbs.ngacn.cc", "https://www.163.com", "https://www.baidu.com"}
	headers := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; Android 6.0.1; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Mobile Safari/537.36 Edg/87.0.664.66",
	}

	var timeout int64 = 2
	var limit int64 = 3
	res := UrlProbe(url, headers, limit, timeout)
	t.Log(res)
}
