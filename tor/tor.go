package tor

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/hailiang/gosocks"
)

// Your torrc should have these lines
// ControlPort 9051
// HashedControlPassword 16:CB7707079B9712C860BB052D2D6A96323211DD16D50A170E6ADD10BEFD

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

func controlCommand(commands []string) (result string, err error) {
	conn, err := net.Dial("tcp", "127.0.0.1:9051")
	if err != nil {
		return
	}

	fmt.Fprintf(conn, "AUTHENTICATE \"hM4LoBWbaVCz4rb\"\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	parts := strings.Split(status, " ")
	if parts[0] != "250" || err != nil {
		return status, err
	}

	for _, cmd := range commands {
		fmt.Fprintf(conn, cmd+"\r\n")
		status, err = bufio.NewReader(conn).ReadString('\n')
		result = status
	}

	return
}

func NewIP() (newip string, err error) {
	cmds := []string{
		"SIGNAL NEWNYM",
		"GETINFO stream-status",
	}
	result, err := controlCommand(cmds)
	if err != nil {
		return
	}
	parts := strings.Split(result, " ")
	if len(parts) == 4 {
		newip = parts[3]
	} else {
		err = errors.New("No IP found")
	}

	return
}
