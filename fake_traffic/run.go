package fake_traffic

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func Run(fakeIPs []string) {
	for {
		randomIndex := rand.Intn(len(fakeIPs))
		proxy := fakeIPs[randomIndex]
		go get(proxy)
		time.Sleep(time.Second * 2)
	}
}

func get(proxyIP string) {
	u, err := url.Parse("http://" + proxyIP)
	if err != nil {
		log.Printf("Url parser error: %s\n", err.Error())
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
	}

	sites := []string{
		"https://google.com",
		"https://microsoft.com",
		"https://bing.com",
		"https://youtube.com",
		"https://tiktok.com",
		"https://netflix.com",
		"https://samsung.com",
		"https://fandom.com",
	}

	randomIndex := rand.Intn(len(sites))
	targetURL := sites[randomIndex]

	resp, err := client.Get(targetURL)
	if err != nil {
		fmt.Printf("get err: %v\n", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("readAll error: %v\n", err)
		return
	}
	return
}
