package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestProducer(t *testing.T) {
	partitions := map[int]chan string{
		0: make(chan string),
		1: make(chan string),
	}

	producer, err := NewProducer(partitions)
	if err != nil {
		t.Errorf("failed to init Producer: %s", err)
	}

	go producer.Start(bufio.NewReader(strings.NewReader("Cmd Inc")))

	if <-partitions[0] != "Cmd" {
		t.Errorf("producer should've written Cmd to partition 0")
	}

	if <-partitions[1] != "Inc" {
		t.Errorf("producer should've written Inc to partition 1")
	}
}

func TestNewProducerValidation(t *testing.T) {
	partitions := make(map[int]chan string)

	_, err := NewProducer(partitions)
	if err == nil {
		t.Errorf("failed to validate partitions")
	}
}
