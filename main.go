package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ziutek/telnet"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//import (
//	"fmt"
//	"os/exec"
//)

//var cmd *exec.Cmd
//
//func main() {
//	var out []byte
//	cmd = exec.Command("ls")
//	out, err := cmd.Output()
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println(string(out))
//}

var (
	//server
	server = "127.0.0.1:11211" //mc连接地址

	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "my_counter",
			Help:      "This is my counter",
		})

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
			//ConstLabels: map[string]string{
			//	"path":"/api/test",
			//},
		})
)

func mcInit() {
	//create a handler
	conn, err := telnet.Dial("tcp", server)
	if err != nil {
		fmt.Println("conn faild")
		return
	}
	defer func(conn *telnet.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("defer conn close error ")
		}
	}(conn)
	buf := make([]byte, 30)
	var chunk []byte
	//Write 'Hello World' to the stream.
	_, err = conn.Write([]byte("stats slabs"))
	fmt.Println("write stats slabs")
	if err != nil {
		panic("Unable to write from stream.")
	}
	for {
		fmt.Println("start read")
		n, err := conn.Read(buf)
		fmt.Println(n, buf[:n])
		if err != nil {
			fmt.Println(err.Error())
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buf[:n]...)
		fmt.Println(chunk)
	}
	//return string(chunk)
	fmt.Println(string(chunk))

}

func main() {
	rand.Seed(time.Now().Unix())

	http.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(counter)
	prometheus.MustRegister(gauge)

	go func() {
		for {
			counter.Add(rand.Float64() * 5)
			gauge.Set(rand.Float64() * 15)
		}
	}()

	fmt.Println("Starting")
	mcInit()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
