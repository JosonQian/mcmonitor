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
)

var pages_count, chunk_count, use_chunk, free_chunk int

func GetMcSlabs() {
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
	fmt.Println(slice)
	for _, v := range slice {
		if findPage := strings.Contains(v, "total_pages"); findPage {
			fmt.Println(v)
			sti, _ := strconv.Atoi(strings.Split(v, " ")[1])
			pages_count += sti

		}
		if findChunk := strings.Contains(v, "total_chunks"); findChunk {
			fmt.Println(v)
			sti, _ := strconv.Atoi(strings.Split(v, " ")[1])
			chunk_count += sti
		}
		if findFree := strings.Contains(v, "free_chunks"); findFree {
			fmt.Println(v)
			sti, _ := strconv.Atoi(strings.Split(v, " ")[1])
			free_chunk += sti
		}
		use_chunk = chunk_count - free_chunk
		fmt.Println(pages_count, chunk_count, use_chunk, free_chunk)
	}
}

func Exporter() {
	//get total pages
	gaugeTotalPages := prometheus.NewGauge(prometheus.GaugeOpts{ // 定义指标
		Name: "pages_count",
		Help: "Test Gauge Custom help info",
	})
	if err := prometheus.Register(gaugeTotalPages); err != nil {
		log.Fatal(err)
	}
	gaugeTotalPages.Set(float64(pages_count))

	//get total chunks
	gaugeTotalChunks := prometheus.NewGauge(prometheus.GaugeOpts{ // 定义指标
		Name: "chunk_count",
		Help: "Test Gauge Custom help info",
	})
	if err := prometheus.Register(gaugeTotalChunks); err != nil {
		log.Fatal(err)
	}
	gaugeTotalChunks.Set(float64(chunk_count))

	//get use chunks
	gaugeUseChunks := prometheus.NewGauge(prometheus.GaugeOpts{ // 定义指标
		Name: "use_chunk",
		Help: "Test Gauge Custom help info",
	})
	if err := prometheus.Register(gaugeUseChunks); err != nil {
		log.Fatal(err)
	}
	use_chunk = chunk_count - free_chunk
	gaugeTotalPages.Set(float64(use_chunk))

	//get free chunks
	gaugeFreeChunks := prometheus.NewGauge(prometheus.GaugeOpts{ // 定义指标
		Name: "free_chunk",
		Help: "Test Gauge Custom help info",
	})
	if err := prometheus.Register(gaugeFreeChunks); err != nil {
		log.Fatal(err)
	}
	gaugeTotalPages.Set(float64(free_chunk))

}
func main() {
	GetMcSlabs()
	Exporter()

	http.Handle("/metrics", promhttp.Handler()) // 服务！
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
