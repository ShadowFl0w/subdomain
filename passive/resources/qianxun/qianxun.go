package qianxun

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"subdomain/libs"
	"subdomain/utils"
)

func Qianxun(domain string, opt *libs.Options) {

	log.Println("==Qianxun==")
	data := "ecmsfrom=&show=&num=&classid=0&keywords=" + domain
	num := 1

	for {
		url := "https://www.dnsscan.cn/dns.html?keywords=" + domain + "&page=" + fmt.Sprint(num)
		resp, err := HTTPRequest("POST", url, strings.NewReader(data))
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.Status == "200 OK" {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			respbody := string(body)
			//fmt.Println(respbody)
			result := utils.GetSubomainsNot(respbody, domain)
			//fmt.Println(result)
			if result == nil {
				break
			}
			for _, record := range result {

				log.Println("qianxun: ", record)
				opt.PRwmutex.Lock()
				opt.Presults[record] = struct{}{}
				opt.PRwmutex.Unlock()
			}

			lastPageString := `<li class="disabled"><span>&raquo;</span></li>`
			if strings.Contains(respbody, lastPageString) {
				resp.Body.Close()
				break
			}

			resp.Body.Close()
			num++
		} else {
			_, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			break
		}

	}

	//最后一页
	url := "https://www.dnsscan.cn/dns.html?keywords=" + domain + "&page=" + fmt.Sprint(num)
	resp, err := HTTPRequest("POST", url, strings.NewReader(data))
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	respbody := string(body)
	// fmt.Println(respbody)
	resultLast := utils.GetSubomains(respbody, domain)
	for _, record := range resultLast {

		log.Println("qianxun: ", record)
		opt.PRwmutex.Lock()
		opt.Presults[record] = struct{}{}
		opt.PRwmutex.Unlock()
	}
	resp.Body.Close()

}

func HTTPRequest(method string, requestUrl string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		log.Println(err)
	}
	// proxy, _ := url.Parse("http://127.0.0.1:8080")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Proxy:           http.ProxyURL(proxy),
	}
	client := &http.Client{Transport: tr}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if err != nil {
		log.Println(err)
	}

	return resp, nil

}
