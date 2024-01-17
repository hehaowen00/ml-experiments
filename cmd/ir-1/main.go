package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/vmihailenco/msgpack/v5"
)

func main() {
	log.SetFlags(log.Lshortfile)
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(err)
	}

	f1, err := os.ReadFile("./stopwords.txt")
	if err != nil {
		log.Fatalln("error reading stop words list -", err)
	}

	stopList := strings.Split(string(f1), "\n")

	f1, err = os.ReadFile("./punctuation.txt")
	if err != nil {
		log.Fatalln("error reading punctuation list -", err)
	}

	punctuation := strings.Split(string(f1), "")
	if err != nil {
		log.Fatalln("error reading punctuation list -", err)
	}

	f, err := os.Open("data.txt")
	if err != nil {
		log.Fatalln("error opening data file -", err)
	}

	err = os.MkdirAll("./docs/", 0777)
	if err != nil {
		log.Fatalln("error creating directory -", err)
	}

	err = os.MkdirAll("./index/", 0777)
	if err != nil {
		log.Fatalln("error creating directory -", err)
	}

	// err = os.MkdirAll("./tf/", 0777)
	// if err != nil {
	// 	log.Fatalln("error creating directory -", err)
	// }

	err = os.MkdirAll("./tfidf/", 0777)
	if err != nil {
		log.Fatalln("error creating directory -", err)
	}

	documentCount := 0

	if false {
		scanner := bufio.NewScanner(f)

		const bufSize = 1024 * 1024
		maxBuffer := make([]byte, bufSize)
		scanner.Buffer(maxBuffer, bufSize)

		var count int = 1
		var inc int = 1
		var docs []*Document

		for scanner.Scan() {
			doc := NewDocument(count, scanner.Text())
			documentCount++

			docs = append(docs, doc)
			count += 1

			if count%100000 == 0 {
				data, err := msgpack.Marshal(docs)
				if err != nil {
					log.Fatalln("error marshalling data -", err)
				}

				err = os.WriteFile(fmt.Sprintf("./docs/%d.msgpack", inc), data, 0777)
				if err != nil {
					log.Fatalln("error writing file -", err)
				}

				docs = nil
				inc += 1
			}
		}

		if len(docs) > 0 {
			data, err := msgpack.Marshal(docs)
			if err != nil {
				log.Fatalln("error marshalling data -", err)
			}

			err = os.WriteFile(fmt.Sprintf("./docs/%d.msgpack", inc), data, 0777)
			if err != nil {
				log.Fatalln("error writing file -", err)
			}
		}
	}

	if false {
		files, err := os.ReadDir("./docs/")
		if err != nil {
			log.Fatalln("error reading directory -", err)
		}

		idf := map[string]float64{}
		n := 0

		for _, f := range files {
			// tf := map[string]int{}

			data, err := os.ReadFile(fmt.Sprintf("./docs/%s", f.Name()))
			if err != nil {
				log.Fatalln("error reading file -", err)
			}

			var docs []*Document
			err = msgpack.Unmarshal(data, &docs)
			if err != nil {
				log.Fatalln("error unmarshalling file -", err)
			}

			indexes := make([]*DocumentIndex, len(docs))

			log.Println("processing", f.Name())

			wg := sync.WaitGroup{}

			for i, d := range docs {
				wg.Add(1)
				go func(i int, doc *Document) {
					index := ParseDocument(lemmatizer, doc.ID, doc.Summary, punctuation, stopList)
					indexes[i] = index
					// _, df := termFrequency(lemmatizer, d, punctuation, stopList)
					defer wg.Done()
				}(i, d)
			}

			wg.Wait()

			data, err = msgpack.Marshal(indexes)
			if err != nil {
				log.Fatalln("error marshalling data -", err)
			}

			err = os.WriteFile(fmt.Sprintf("./index/%s", f.Name()), data, 0777)
			if err != nil {
				log.Fatalln("error writing file -", err)
			}

			log.Println("processing corpus")

			for _, index := range indexes {
				for k, _ := range index.Tokens {
					idf[k] += 1
				}
				n++
			}

			// data, err = msgpack.Marshal(idf)
			// if err != nil {
			// 	log.Fatalln("error marshalling data -", err)
			// }

			// err = os.WriteFile(fmt.Sprintf("./tf/%s", f.Name()), data, 0777)
			// if err != nil {
			// 	log.Fatalln("error writing file -", err)
			// }
		}

		for k, v := range idf {
			idf[k] = math.Log(float64(n) / (float64(v)))
		}

		data, err := msgpack.Marshal(idf)
		if err != nil {
			log.Fatalln("error marshalling data -", err)
		}

		err = os.WriteFile("./index/idf.msgpack", data, 0777)
		if err != nil {
			log.Fatalln("error writing file -", err)
		}
	}

	tfidf := map[int]map[string]float64{}

	if false {
		idf := map[string]float64{}

		data, err := os.ReadFile("./index/idf.msgpack")
		if err != nil {
			log.Fatalln("error reading file -", err)
		}

		err = msgpack.Unmarshal(data, &idf)
		if err != nil {
			log.Fatalln("error unmarshalling file -", err)
		}

		files, err := os.ReadDir("index")
		if err != nil {
			log.Fatalln("error reading directory -", err)
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			if f.Name() == "idf.msgpack" {
				continue
			}

			log.Println("processing", f.Name())

			data, err = os.ReadFile(fmt.Sprintf("./index/%s", f.Name()))
			if err != nil {
				log.Fatalln("error reading file -", err)
			}

			var docs []*DocumentIndex
			err = msgpack.Unmarshal(data, &docs)
			if err != nil {
				log.Fatalln("error unmarshalling file -", err)
			}

			for _, index := range docs {
				for k, v := range index.TF {
					index.TF[k] = v * idf[k]
					if index.TF[k] < 0.000001 {
						delete(index.TF, k)
					}
				}
				tfidf[index.ID] = index.TF
			}

			// data, err = msgpack.Marshal(docs)
			// if err != nil {
			// 	log.Fatalln("error marshalling data -", err)
			// }

			// err = os.WriteFile(fmt.Sprintf("./tfidf/%s", f.Name()), data, 0777)
			// if err != nil {
			// 	log.Fatalln("error writing file -", err)
			// }
		}

		f, err := os.OpenFile("./tfidf/tfidf.msgpack", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			log.Fatalln("error opening file -", err)
		}

		enc := msgpack.NewEncoder(f)
		err = enc.Encode(tfidf)
		if err != nil {
			log.Println("error encoding -", err)
		}

		// tfidfData, err := msgpack.Marshal(tfidf)
		// if err != nil {
		// 	log.Fatalln("error marshalling data -", err)
		// }

		// err = os.WriteFile("./tfidf/tfidf.msgpack", tfidfData, 0777)
		// if err != nil {
		// 	log.Fatalln("error writing file -", err)
		// }
	}

	// tfidf := map[int]map[string]float64{}
	if false {
		files, err := os.ReadDir("tfidf")
		if err != nil {
			log.Fatalln("error reading directory -", err)
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			log.Println("processing", f.Name())

			data, err := os.ReadFile(fmt.Sprintf("./tfidf/%s", f.Name()))
			if err != nil {
				log.Fatalln("error reading file -", err)
			}

			var docs []*DocumentIndex

			err = msgpack.Unmarshal(data, &docs)
			if err != nil {
				log.Fatalln("error unmarshalling file -", err)
			}

			for _, index := range docs {
				fmt.Println(index.ID)
				tfidf[index.ID] = index.TF
				fmt.Println(len(tfidf))
				break
			}
		}

		fmt.Println(len(tfidf))

		tfidfData, err := msgpack.Marshal(tfidf)
		if err != nil {
			log.Fatalln("error marshalling data -", err)
		}

		err = os.WriteFile("./tfidf/tfidf.msgpack", tfidfData, 0777)
		if err != nil {
			log.Fatalln("error writing file -", err)
		}
	}

	f, err = os.OpenFile("./tfidf/tfidf.msgpack", os.O_RDONLY, 0777)
	if err != nil {
		log.Fatalln("error opening file -", err)
	}

	dec := msgpack.NewDecoder(f)
	err = dec.Decode(&tfidf)
	if err != nil {
		log.Println("error decoding -", err)
	}

	// readline
	// read input from stdin
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter text: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	fmt.Println("searching for", text)

	words := processDataset(lemmatizer, text, punctuation, stopList)

	tokens := map[string]int{}

	for _, w := range words {
		if len(w) == 0 {
			continue
		}

		tokens[w]++
	}

	count := len(words)

	tf := map[string]float64{}

	for k, v := range tokens {
		tf[k] = float64(v) / float64(count)
	}

	idf := map[string]float64{}
	for k, v := range tf {
		idf[k] = math.Log(float64(len(tfidf)) / (1 + float64(v)))
	}

	queryTfIdf := map[string]float64{}
	for k, v := range tf {
		queryTfIdf[k] = v * idf[k]
	}

	res := []float64{}
	for _, v := range queryTfIdf {
		res = append(res, v)
	}

	type Entry struct {
		Doc   int
		Score float64
	}

	results := make(map[int]Entry, len(tfidf))
	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}

	for doc := range tfidf {
		wg.Add(1)
		go func(doc int) {
			docV := []float64{}
			defer wg.Done()

			for k, _ := range queryTfIdf {
				docV = append(docV, tfidf[doc][k])
			}

			cos, err := cosine(docV, res)
			if err != nil {
				return
			}

			if cos < 0.2 {
				return
			}

			mutex.Lock()
			results[doc] = Entry{Doc: doc, Score: cos}
			mutex.Unlock()
		}(doc)
	}

	wg.Wait()

	ranked := make([]Entry, 0, len(results))

	for _, e := range results {
		ranked = append(ranked, e)
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	clear(results)

	fmt.Println("results", len(ranked))

	files, err := os.ReadDir("./docs/")
	if err != nil {
		log.Fatalln("error reading directory -", err)
	}

	if len(files) == 0 {
		log.Fatalln("no files found")
	}

	for _, f := range files {
		data, err := os.ReadFile(fmt.Sprintf("./docs/%s", f.Name()))
		if err != nil {
			log.Fatalln("error reading file -", err)
		}

		var docs []*Document
		err = msgpack.Unmarshal(data, &docs)
		if err != nil {
			log.Fatalln("error unmarshalling file -", err)
		}

		// last 100

		if len(ranked) > 100 {
			ranked = ranked[:100]
		}

		for _, e := range ranked {
			if e.Score < 0.2 {
				continue
			}
			id := e.Doc
			if id == 0 {
				continue
			}
			for _, doc := range docs {
				if doc.ID == id {
					fmt.Println(doc)
				}
			}
		}
	}
}
