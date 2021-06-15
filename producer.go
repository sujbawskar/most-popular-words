package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// Producer holds a number of stream partitions to which it writes strings.
// nextPartition indicates to which partition Producer should write the next word.
type Producer struct {
	numPartitions    int
	streamPartitions map[int]chan string
	nextPartition    int
}

// NewProducer takes in a map with at least one stream partition
// and returns a Producer that can be started given a Reader.
func NewProducer(streamPartitions map[int]chan string) (*Producer, error) {
	if len(streamPartitions) == 0 {
		return nil, fmt.Errorf("streamPartitions should have at least 1 partition")
	}

	return &Producer{numPartitions: len(streamPartitions), streamPartitions: streamPartitions, nextPartition: 0}, nil
}

// Start reads space-delimited strings from a Reader and writes each of them to
// the next stream partition in a round-robin fashion until it hits EOF.
func (p *Producer) Start(r io.Reader) {
	defer p.close()
	fmt.Println("starting producer")

	reader := bufio.NewReader(r)

	for {
		w, err := reader.ReadString(' ')
		if err != nil && err != io.EOF {
			log.Printf("failed to read word: %s", err)
			return
		}

		// When the string at the end of a line, ReadString will return
		// \n and the subsequent word if there are no spaces between them.
		words := strings.Split(w, "\n")
		for _, w := range words {
			if w2 := strings.Trim(w, " "); len(w2) > 0 {
				p.write(w2)
			}
		}

		// Break if we reached EOF
		if err != nil && err == io.EOF {
			break
		}
	}
}

func (p *Producer) write(word string) {
	partition := p.streamPartitions[p.nextPartition]
	partition <- word
	// Round-robin through the partitions
	p.nextPartition++

	if p.nextPartition >= p.numPartitions {
		p.nextPartition = 0
	}
}

func (p *Producer) close() {
	for _, s := range p.streamPartitions {
		close(s)
	}
}
