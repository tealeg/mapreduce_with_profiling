// In this mentoring session, we take the mapreduce program we created
// (https://github.com/tealeg/mapreduce) and do some performance
// testing and profiling on it.
//
// The goal here is to provide an introduction to advanced tooling
// available in Go for this kind of investigative work.

// The steps we'll take:
// 1. Create "Benchmarks" using the Go test packages Benchmark functionality.
// 2. Insert code to start a Go profiler and present an HTTP interface to interact with it.

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type WordCounter struct {
	mapWg       sync.WaitGroup
	reduceWg    sync.WaitGroup
	aggregateWg sync.WaitGroup
	mapCh       chan string
	reduceCh    chan []string
	aggregateCh chan int
	count       int
}

func NewWordCounter() *WordCounter {
	return &WordCounter{}
}

func (wc *WordCounter) mapper() {
	for line := range wc.mapCh {
		wc.reduceCh <- strings.Fields(line)
	}
	wc.mapWg.Done()
}

func (wc *WordCounter) reducer() {
	for fields := range wc.reduceCh {
		wc.aggregateCh <- len(fields)
	}
	wc.reduceWg.Done()
}

func (wc *WordCounter) agregator() {
	var count int = 0
	for line := range wc.aggregateCh {
		count += line
	}
	wc.count = count
	wc.aggregateWg.Done()
}

func (wc *WordCounter) Count(input io.Reader) int {

	wc.mapCh = make(chan string)
	wc.reduceCh = make(chan []string)
	wc.aggregateCh = make(chan int)
	wc.count = 0

	wc.aggregateWg.Add(1)
	go wc.agregator()
	wc.reduceWg.Add(2)
	go wc.reducer()
	go wc.reducer()
	wc.mapWg.Add(2)
	go wc.mapper()
	go wc.mapper()

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		wc.mapCh <- text
	}
	close(wc.mapCh)
	wc.mapWg.Wait()
	close(wc.reduceCh)
	wc.reduceWg.Wait()
	close(wc.aggregateCh)
	wc.aggregateWg.Wait()
	return wc.count
}

func main() {
	input, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	wc := NewWordCounter()

	count := wc.Count(input)

	fmt.Printf("Total WC: %d\n", count)
}
