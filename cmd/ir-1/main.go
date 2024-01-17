package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

func main() {
	f, err := os.Open("raw.txt")
	if err != nil {
		log.Fatalln("error opening data file -", err)
	}

	err = os.MkdirAll("./tf/", 0777)
	if err != nil {
		log.Fatalln("error creating directory -", err)
	}

	scanner := bufio.NewScanner(f)

	const bufSize = 1024 * 1024
	maxBuffer := make([]byte, bufSize)
	scanner.Buffer(maxBuffer, bufSize)

	// var docs []*DocumentIndex

	var wg sync.WaitGroup

	var count int

	for scanner.Scan() {
		wg.Add(1)
		doc := NewDocument(scanner.Text())
		go func(doc *Document, count int) {
			index := ParseDocument(doc)

			data, err := msgpack.Marshal(index)
			if err != nil {
				log.Fatalln("error marshalling data -", err)
			}

			err = os.WriteFile(fmt.Sprintf("%d.msgpack", count), data, 0777)
			if err != nil {
				log.Fatalln("error writing file -", err)
			}

			wg.Done()
		}(doc, count)

		count += 1
	}

	wg.Wait()

}
