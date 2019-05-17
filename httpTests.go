package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoginOperation(){

	if !haveLoginData{
		getAuthCreds()
		haveLoginData = true
	}

	for {
		body := strings.NewReader(`loginOp=login&login_csrf=b75b1d09-3e4c-43d9-a1c4-563b646436ec&username=john%40johnholder.net&password=123456&client=preferred`)
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		req, err := http.NewRequest("POST", "https://192.168.1.17/", body)
		if err != nil {
			// handle err
		}
		req.Host = "localhost"
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Referer", "https://192.168.1.17/")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", "ZM_TEST=true; ZM_LOGIN_CSRF=b75b1d09-3e4c-43d9-a1c4-563b646436ec")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err.Error())
		} else {
			//bodResp := resp.Body
			//fmt.Println(string(bodResp.Read()))
			resp.Body.Close()
		}
		c.Inc()
		time.Sleep(50*time.Millisecond)
	}
}

