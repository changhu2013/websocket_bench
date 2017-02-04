package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"strconv"

	"github.com/gorilla/websocket"
)

func readMessage(c *websocket.Conn) {
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}

func sendMessage(c *websocket.Conn, m string, s int) {
	defer c.Close()
	for {
		err := c.WriteMessage(websocket.TextMessage, []byte(m))
		if err != nil {
			log.Println("write:", err)
			return
		}
		log.Printf("write %s", m)

		ss, _ := time.ParseDuration(strconv.Itoa(s) + "s")
		time.Sleep(ss)
	}
}

func doConnect(h string, m string, s int, wg *sync.WaitGroup) {
	c, _, err := websocket.DefaultDialer.Dial(h, nil)

	if err != nil {
		log.Println("dial:", err)
	}

	go readMessage(c)
	go sendMessage(c, m, s)

	wg.Done()
}

func connect(c int, k int, h string, m string, s int) {

	if c == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(c)

	for i := 0; i < c; i = i + 1 {
		go doConnect(h, m, s, &wg)
	}

	kk, _ := time.ParseDuration(strconv.Itoa(k) + "s")
	time.Sleep(kk)

	wg.Wait()
}

func bench(a int, c int, k int, h string, m string, s int) {

	b := a / c
	p := a % c

	cc := 0

	for i := 0; i < b; i = i + 1 {
		connect(c, k, h, m, s)
		cc = cc + c
		log.Printf("connected %d", cc)
	}

	connect(p, k, h, m, s)

	cc = cc + p
	log.Printf("connected %d", cc)
}

func main() {

	var a int
	var c int
	var k int
	var s int
	var h string
	var m string

	flag.IntVar(&a, "a", 1, "总数")
	flag.IntVar(&c, "c", 1, "并发数")
	flag.IntVar(&k, "k", 5, "分批创建连接时间间隔 单位:秒 默认为5秒")
	flag.IntVar(&s, "s", 6, "发送消息间隔时间 单位:秒 默认为6秒")
	flag.StringVar(&h, "h", "127.0.0.1", "服务地址 eg: ws://127.0.0.1:6600?sessionid=11011011000_KGaKhmGrAF")
	flag.StringVar(&m, "m", "", "要发送的消息")

	flag.Parse()

	bench(a, c, k, h, m, s)

	<-make(chan struct{})
}
