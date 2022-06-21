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

var pagesCount, chunkCount, useChunk, freeChunk int

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
	slice := strings.Split(string(out), string('\r'))

	for _, v := range slice {
		if findPage := strings.Contains(v, "total_pages"); findPage {
			fmt.Println(v)
			stp, _ := strconv.Atoi(strings.Split(v, " ")[1])
			pagesCount = stp + pagesCount

		}
		if findChunk := strings.Contains(v, "total_chunks"); findChunk {
			fmt.Println(v)
			stc, _ := strconv.Atoi(strings.Split(v, " ")[1])
			chunkCount = stc + chunkCount
		}
		if findFree := strings.Contains(v, "free_chunks"); findFree {
			fmt.Println(v)
			stf, _ := strconv.Atoi(strings.Split(v, " ")[1])
			freeChunk = stf + freeChunk
		}
		useChunk = chunkCount - freeChunk

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
	fmt.Println(pagesCount, float64(pagesCount))
	gaugeTotalPages.Set(float64(pagesCount))

	//get total chunks
	gaugeTotalChunks := prometheus.NewGauge(prometheus.GaugeOpts{ // 定义指标
		Name: "chunk_count",
		Help: "Test Gauge Custom help info",
	})
	if err := prometheus.Register(gaugeTotalChunks); err != nil {
		log.Fatal(err)
	}
	fmt.Println(chunkCount, float64(chunkCount))
	gaugeTotalChunks.Set(float64(chunkCount))

	//get use chunks
	gaugeUseChunks := prometheus.NewGauge(prometheus.GaugeOpts{ // 定义指标
		Name: "use_chunk",
		Help: "Test Gauge Custom help info",
	})
	if err := prometheus.Register(gaugeUseChunks); err != nil {
		log.Fatal(err)
	}
	useChunk = chunkCount - freeChunk
	fmt.Println()
	gaugeTotalPages.Set(float64(useChunk))

	//get free chunks
	gaugeFreeChunks := prometheus.NewGauge(prometheus.GaugeOpts{ // 定义指标
		Name: "free_chunk",
		Help: "Test Gauge Custom help info",
	})
	if err := prometheus.Register(gaugeFreeChunks); err != nil {
		log.Fatal(err)
	}
	gaugeTotalPages.Set(float64(freeChunk))

}
func main() {
	GetMcSlabs()
	Exporter()

	http.Handle("/metrics", promhttp.Handler()) // 服务！
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
