package tor

import (
  "net/http"
  "github.com/hailiang/gosocks"
  "io/ioutil"
)

func prepareProxyClient() *http.Client {
  	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:9050")
  	transport := &http.Transport{Dial: dialSocksProxy}
  	return &http.Client{Transport: transport}
}

func httpGet(httpClient *http.Client, url string) (resp *http.Response, err error) {
  	req, err := http.NewRequest("GET", url, nil)
  	req.Header.Set("User-Agent", "curl/7.21.4 (universal-apple-darwin11.0) libcurl/7.21.4 OpenSSL/0.9.8x zlib/1.2.5")
  	resp, err = httpClient.Do(req)
  	return
}

func httpGetBody(httpClient *http.Client, url string) (body string, err error) {
  	resp, err := httpGet(httpClient, url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
  	bodyb, err := ioutil.ReadAll(resp.Body)
  	body = string(bodyb)
  	return
}

func httpGetStatus(httpClient *http.Client, url string) (status string, err error) {
  	resp, err := httpGet(httpClient, url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	status = resp.Status
  	return
}

func TorGet(url string) (body string, err error) {
        clientPtr := prepareProxyClient()
        body, err = httpGetBody(clientPtr, url)
	return
}

func TorGetStatus(url string) (status string, err error) {
        clientPtr := prepareProxyClient()
        status, err = httpGetStatus(clientPtr, url)
	return
}
