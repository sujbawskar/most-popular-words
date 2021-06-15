package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Consumer struct {
	streamPartitions chan string
	numPartition     int
	wordCountMap     *WordCountMap
}

// Start consumer
func (c *Consumer) Start() {
	msg := ""
	opened := true
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	// Checks if channel is opened
	for opened {
		msg, opened = <-c.streamPartitions
		msg = strings.ToLower(msg)
		msg = reg.ReplaceAllString(msg, "")
		// Locks the word count map for reading and writing
		c.wordCountMap.Lock()
		count := c.wordCountMap.m[msg]
		c.wordCountMap.m[msg] = count + 1
		c.wordCountMap.Unlock()
	}
	fmt.Println("Finished ", c.numPartition)
	wg.Done()
}
