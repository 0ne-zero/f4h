package payment

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func TorGETRequest(addr string, addr_port uint) (*http.Response, error) {
	// Tor = 127.0.0.1:9050
	// Privoxy = 127.0.0.1:8118
	// Privoxy convert HTTP to SOCKS5 and send it to Tor
	// Send request to Privoxy (localhost:8118) and privoxy send request to Tor (localhost:9050)
	// Privoxy address (proxy)
	privoxy_addr := "//127.0.0.1:8118"
	privoxy_url, err := url.Parse(privoxy_addr)
	if err != nil {
		return nil, err
	}
	// Create a Transport
	privoxy_transport := &http.Transport{Proxy: http.ProxyURL(privoxy_url)}
	// Create Client
	client := http.Client{Transport: privoxy_transport, Timeout: time.Second * 5}
	// Send GET request
	response, err := client.Get(fmt.Sprintf("%s:%d", addr, addr_port))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Return response (*http.Response)
	return response, nil
}
func ExtractBodyFromHttpResponse(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	return string(body), err
}
