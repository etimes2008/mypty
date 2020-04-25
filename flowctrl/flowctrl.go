package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func CopyRate(dst io.Writer, src io.Reader, bps int64) (written int64, err error) {
	throttle := time.NewTicker(time.Second)
	defer throttle.Stop()

	var n int64
	for {
		n, err = io.CopyN(dst, src, bps)
		if n > 0 {
			written += n
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		<-throttle.C // rate limit our flows
	}
	return written, err
}

func main() {
	// var fillInterval = time.Millisecond * 10
	// var capacity = 100
	// var tokenBucket = make(chan struct{}, capacity)
	// fillToken := func() {
	// 	ticker := time.NewTicker(fillInterval)
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			select {
	// 			case tokenBucket <- struct{}{}:
	// 			default:
	// 			}
	// 			fmt.Println("current token cnt:", len(tokenBucket), time.Now())
	// 		}
	// 	}
	// }
	// go fillToken()

	FolwInfos := make(map[string]int, 12*24)

	file, err := os.Open("flow")
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		} else {
			// fmt.Println(line)
			flowInfo := strings.Fields(line)
			tmp, _ := strconv.ParseFloat(flowInfo[3], 64)
			FolwInfos[flowInfo[1][:5]] = int(tmp) / 1024 / 8
			fmt.Println(flowInfo[1][:5], flowInfo[1][0:2], flowInfo[1][3:5], flowInfo[3], FolwInfos[flowInfo[1][:5]])
		}
	}

	// timeNow := time.Now()
	// fmt.Println(timeNow.Format("15:04"), timeNow.Minute()%5, FolwInfos[timeNow.Format("15:04")])

	throttle := time.NewTicker(time.Second)
	defer throttle.Stop()

	var flow int64
	// kps := 1024 * rand.Intn(10)
	var kps int32
	kps = int32(rand.Intn(500))
	// var data [1024]byte
	go func() {
		// for {
		// 	kps = int32(1024 * rand.Intn(10))
		// 	kpsA := atomic.LoadInt32(&kps)
		// 	atomic.CompareAndSwapInt32(&kps, kpsA, int32(1024*rand.Intn(10)))
		// 	println(kps)
		// 	time.Sleep(time.Second)
		// }
		preMinute := -1
		for {
			timeNow := time.Now()
			curMinute := timeNow.Minute()
			if curMinute%5 == 0 && preMinute != curMinute {
				preMinute = curMinute
				fmt.Println(timeNow.Format("15:04"), timeNow.Minute()%5, FolwInfos[timeNow.Format("15:04")])
				kps = int32(FolwInfos[timeNow.Format("15:04")])
				kpsA := atomic.LoadInt32(&kps)
				atomic.CompareAndSwapInt32(&kps, kpsA, kps)
			}
			time.Sleep(time.Second)
		}
	}()

	uClinet, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		// IP:   net.IPv4(127, 0, 0, 1),
		IP:   net.IPv4(61, 164, 110, 198),
		Port: 8080,
	})
	if err != nil {
		fmt.Println("connect fail !", err)
		return
	}
	defer uClinet.Close()
	data := make([]byte, 1024)
	for {
		kpsA := atomic.LoadInt32(&kps)
		println(1, kpsA)
		kpsA = kpsA + kpsA*int32((rand.Intn(400)-200))/1000
		for i := 0; i < int(kpsA); i++ {
			// rand.Read(data) //生成随机字符串
			uClinet.Write(data)
			flow += int64(len(data))
		}
		log.Println("===>", len(data)*int(kps), flow)
		<-throttle.C // rate limit our flows
	}

	// time.Sleep(time.Hour)
}
