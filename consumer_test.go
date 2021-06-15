package main

import (
	"testing"
)

func TestConsumer(t *testing.T) {
	wordCountMap := WordCountMap{m: make(map[string]int)}
	myChan := make(chan string)
	consumer := Consumer{myChan, 1, &wordCountMap}
	wg.Add(1)
	go consumer.Start()
	myChan <- "Hi"
	myChan <- "there"
	close(myChan)
	wg.Wait()
	if wordCountMap.m["hi"] != 1 {
		t.Fail()
	}
	if wordCountMap.m["there"] != 1 {
		t.Fail()
	}
}
