package main

import (
	_ "bufio"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	_ "os"
	"strings"
	"time"
)

func LoginOperation(){

	if !haveLoginData{
		getAuthCreds()
		haveLoginData = true
	}
	var body io.Reader
	if csrfTokenNeeded{
		body = strings.NewReader(`loginOp=login&login_csrf=`+csrfToken+`&username=`+endUserName+`&password=`+endUserPassword+`&client=preferred`)

	}else{
		body = strings.NewReader(`loginOp=login&username=`+endUserName+`&password=`+endUserPassword+`&client=preferred`)
	}
	for {
		if disableSSLChecks{
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}

		req, err := http.NewRequest("POST", httpEndPoint, body)
		if err != nil {
			// handle err
		}
		req.Host = "localhost"
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Referer", httpEndPoint)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if csrfTokenNeeded{
			req.Header.Set("Cookie", "ZM_TEST=true; ZM_LOGIN_CSRF="+csrfToken)
		}


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

func getCSRFToken(){
	// Request the HTML page.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//res, err := http.Get(url)
	res, err := http.Get(httpEndPoint)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		nP := s.Get(0)
		if nP.Attr[1].Key=="name"{
			fmt.Println()
			if nP.Attr[1].Val =="login_csrf"{
				csrfToken = nP.Attr[2].Val
				csrfTokenNeeded = true
			}
		}
	})
}

