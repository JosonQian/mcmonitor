package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var pages_count, chunk_count, use_chunk, free_chunk int

func getMcSlabs() {
	var cmd *exec.Cmd
	var out []byte
	cmd = exec.Command("bash", "-c", "echo 'stats slabs'|nc 127.0.0.1 11211 |grep -E 'total_pages|total_chunks|free_chunks'|grep -v free_chunks_end|awk -F ':' '{print $2}'")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("output", string(out))
	slice := strings.Split(string(out), string('\n'))
	for _, v := range slice {
		if findPage := strings.Contains(v, "a"); findPage {
			fmt.Println(v)
			sti, _ := strconv.Atoi(strings.Split(v, " ")[1])
			pages_count += sti

		}
		if findChunk := strings.Contains(v, "b"); findChunk {
			fmt.Println(v)
			sti, _ := strconv.Atoi(strings.Split(v, " ")[1])
			chunk_count += sti
		}
		if findFree := strings.Contains(v, "b"); findFree {
			fmt.Println(v)
			sti, _ := strconv.Atoi(strings.Split(v, " ")[1])
			free_chunk += sti
		}
		use_chunk = chunk_count - free_chunk
		fmt.Println(pages_count, chunk_count, use_chunk, free_chunk)
	}
}
func recordMetrics(c *prometheus.GaugeVec) {
	go func() {
		for {
			c.With(prometheus.Labels{"total_pages": "Red"}).Set(float64(pages_count)) // 写入指标
			c.With(prometheus.Labels{"total_chunk": "Red"}).Set(float64(chunk_count))
			c.With(prometheus.Labels{"use_chunk": "Red"}).Set(float64(use_chunk))
			c.With(prometheus.Labels{"free_chunk": "Red"}).Set(float64(free_chunk))
			time.Sleep(time.Second * 2)
		}
	}()
}

func main() {
	MyGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{ // 定义指标
		Name: "test_gauge",
		Help: "Test Gauge Custom help info",
	},
		[]string{"total_pages", "total_chunk", "use_chunk", "free_chunk"},
	)

	if err := prometheus.Register(MyGauge); err != nil {
		log.Fatal(err)
	}
	recordMetrics(MyGauge)
	http.Handle("/metrics", promhttp.Handler()) // 服务！
	log.Fatal(http.ListenAndServe("0.0.0.0:9109", nil))
}
