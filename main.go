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
	"log"
	"os"
	"strings"
	"sync"
)

func mapper(input <-chan string, output chan<- []string, wg *sync.WaitGroup) {
	for line := range input {
		output <- strings.Fields(line)
	}
	wg.Done()
}

func reducer(input <-chan []string, output chan<- int, wg *sync.WaitGroup) {
	for fields := range input {
		output <- len(fields)
	}
	wg.Done()
}

func printer(input <-chan int, wg *sync.WaitGroup) {
	count := 0
	for line := range input {
		fmt.Printf("%d \n", line)
		count += line
	}
	wg.Done()
	fmt.Printf("Total WC: %d\n", count)
}

func main() {
	input, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var mapWg sync.WaitGroup
	var reduceWg sync.WaitGroup
	var printWg sync.WaitGroup

	mapCh := make(chan string)
	reduceCh := make(chan []string)
	printCh := make(chan int)
	printWg.Add(1)
	go printer(printCh, &printWg)
	reduceWg.Add(2)
	go reducer(reduceCh, printCh, &reduceWg)
	go reducer(reduceCh, printCh, &reduceWg)
	mapWg.Add(2)
	go mapper(mapCh, reduceCh, &mapWg)
	go mapper(mapCh, reduceCh, &mapWg)

	defer input.Close()
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		mapCh <- scanner.Text()

	}
	close(mapCh)
	mapWg.Wait()
	close(reduceCh)
	reduceWg.Wait()
	close(printCh)
	printWg.Wait()
}
