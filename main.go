package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

/////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                             //
// This is a more sophisticated version of the good ol' word count problem. In this version,   //
// we have a producer that reads a number of words from a file on disk and sends them          //
// to one of N partitions, simulating the semantics implemented by Kafka streams.              //
//                                                                                             //
// Your goal is to write some code to consume the data coming through each of the partitions   //
// and count the occurrence of each word across all partitions.                                //
//                                                                                             //
// Feel free to make changes to the main function below in order to accomodate your solution,  //
// but we ask that you follow these guidelines:                                                //
//                                                                                             //
// 1. Don't modify producer.go                                                                 //
// 2. Don't use any external libraries;                                                        //
// 3. Do state any relevant assumptions or limitations.                                        //
// 4. Running make run/test should successfully run/test your solution.                        //
//                                                                                             //
/////////////////////////////////////////////////////////////////////////////////////////////////

const (
	filepath      = "world192.txt"
	numPartitions = 3
)

// TODO: implement some sort of consumer code that
// reads data from all partitions in streamPartitions.
var wg sync.WaitGroup

type WordCountMap struct {
	m map[string]int
	sync.RWMutex
}
type by_freq []word_struct

func (a by_freq) Len() int           { return len(a) }
func (a by_freq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a by_freq) Less(i, j int) bool { return a[i].freq > a[j].freq }

type word_struct struct {
	freq int
	word string
}

func main() {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var cons [numPartitions]Consumer
	streamPartitions := make(map[int]chan string)
	wordCountMap := WordCountMap{m: make(map[string]int)}
	// Initialize consumers
	for i := 0; i < numPartitions; i++ {
		streamPartitions[i] = make(chan string)
		consumer := Consumer{streamPartitions[i], i, &wordCountMap}
		cons[i] = consumer
	}

	// init producer with a number of stream partitions to which
	// the producer will write the words in the file at filepath.
	// Start consumers
	producer, err := NewProducer(streamPartitions)
	for i := 0; i < numPartitions; i++ {
		wg.Add(i)
		go cons[i].Start()
	}
	if err != nil {
		log.Fatalf("failed to init Producer: %s", err)
	}
	go producer.Start(f)
	wg.Wait()
	fmt.Println("Cmd Inc. 2021")
	pws := new(word_struct)
	struct_slice := make([]word_struct, len(wordCountMap.m))
	ix := 0
	for key, val := range wordCountMap.m {
		pws.freq = val
		pws.word = key
		// test, %+v shows field names
		//fmt.Printf("%v %v  %+v\n", pws.freq, pws.word, *pws)
		struct_slice[ix] = *pws
		ix++
	}
	sort.Sort(by_freq(struct_slice))
	for ix := 0; ix < 100; ix++ {
		fmt.Printf("%v\n", struct_slice[ix])
	}
}
