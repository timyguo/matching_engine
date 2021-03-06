package main

import (
	"flag"
	"github.com/fmstephe/fstrconv"
	"github.com/fmstephe/matching_engine/cbuf"
	"github.com/fmstephe/matching_engine/itch"
	"github.com/fmstephe/matching_engine/matcher"
	"github.com/fmstephe/matching_engine/trade"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

const (
	StockId = uint32(1)
)

var (
	filePath   = flag.String("f", "", "Relative path to an ITCH file providing test data")
	profile    = flag.String("p", "", "Write out a profile of this application, 'cpu' and 'mem' supported")
	orderNum   = flag.Int("o", 100*1000, "The number of orders to generate. Ignored if -f is provided")
	delDelay   = flag.Int("d", 100, "The number of orders generated before we begin deleting existing orders")
	perfRand   = rand.New(rand.NewSource(1))
	orderMaker = trade.NewOrderMaker()
)

func main() {
	flag.Parse()
	orders := getData()
	orderCount := fstrconv.Itoa64Comma(int64(len(orders)))
	println(orderCount, "Orders Built")
	buffer := cbuf.New(len(orders))
	m := matcher.NewMatcher(*delDelay*2, buffer)
	startProfile()
	defer endProfile()
	start := time.Now().UnixNano()
	for i := range orders {
		m.Submit(&orders[i])
	}
	println("Buffer Writes: ", buffer.Writes())
	total := time.Now().UnixNano() - start
	println("Nanos\t", fstrconv.Itoa64Comma(total))
	println("Micros\t", fstrconv.Itoa64Comma(total/1000))
	println("Millis\t", fstrconv.Itoa64Comma(total/(1000*1000)))
	println("Seconds\t", fstrconv.Itoa64Comma(total/(1000*1000*1000)))
}

func startProfile() {
	if *profile == "cpu" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func endProfile() {
	if *profile == "cpu" {
		pprof.StopCPUProfile()
	}
	if *profile == "mem" {
		f, err := os.Create("mem.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
	}
}

func getData() []trade.OrderData {
	orders, err := orderMaker.RndTradeSet(*orderNum, *delDelay, 1000, 1500)
	if err != nil {
		panic(err.Error())
	}
	return orders
	/*
		if *filePath == "" {
			return mkRandomData()
		}
		return getItchData()
	*/
}

func getItchData() []*trade.Order {
	ir := itch.NewItchReader(*filePath)
	orders, err := ir.ReadAll()
	if err != nil {
		panic(err.Error())
	}
	return orders
}
