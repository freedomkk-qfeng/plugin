package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	vers := flag.Bool("v", false, "display the version.")
	help := flag.Bool("h", false, "print this help.")
	conf := flag.String("f", "config.yml", "specify configuration file.")
	flag.Parse()

	if *vers {
		fmt.Println("version:", Version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	config := Parse(*conf)

	pingRes := PingProbe(config.Ping, config.Probe.Limit, config.Probe.Timeout)

	urlRes := UrlProbe(config.Url, config.Probe.Headers, config.Probe.Limit, config.Probe.Timeout)

	L := []*MetricValue{}
	ts := time.Now().Unix()
	for _, res := range pingRes {
		tag := "ip=" + res.Ip
		L = append(L, GaugeValue(ts, "pingProbe.latency", res.Ping, tag))
	}
	for _, res := range urlRes {
		tag := "url=" + res.Url
		L = append(L, GaugeValue(ts, "urlProbe.latency", res.Latency, tag))
		L = append(L, GaugeValue(ts, "urlProbe.cert", res.Cert, tag))
		L = append(L, GaugeValue(ts, "urlProbe.statusCode", res.HttpStatusCode, tag))
	}
	if len(L) > 0 {
		js, _ := json.Marshal(L)
		fmt.Println(string(js))
	}

}
