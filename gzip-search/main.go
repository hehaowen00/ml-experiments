package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"slices"
	"sync"
)

// var trainingSet = []string{
// 	"The sun sets over the horizon in a blaze of colors.",
// 	"The curious cat explores the garden with playful antics.",
// 	"A chef carefully crafts a delicious gourmet meal.",
// 	"The mountain peak is covered in a blanket of snow.",
// 	"Music fills the air as people dance under the starry sky.",
// }

// var testSet = []string{
// 	"The sky transforms into a canvas of vibrant hues as the sun sets.",
// 	"A playful cat frolics around the garden, investigating its surroundings.",
// 	"The skilled chef prepares a mouthwatering gourmet dish with precision.",
// 	"A snowy landscape surrounds the majestic peak of the mountain.",
// 	"Under the night sky, people dance to the rhythm of music in the air.",
// }

var trainingSet = []string{
	"The Renaissance: A Cultural and Artistic Rebirth",
	"Exploring the Wonders of Ancient Egyptian Civilization",
	"The Industrial Revolution: Transforming Societies and Economies",
	"Space Exploration: Unveiling the Mysteries Beyond Earth",
	"Advancements in Medicine: From Ancient Remedies to Modern Healthcare",
}

var testSet = []string{
	"A Cultural and Artistic Rebirth: Understanding the Renaissance",
	"Journey Through Ancient Egypt: Unraveling Its Civilization",
	"Transforming Societies: The Impact of the Industrial Revolution",
	"Unveiling Space Mysteries: The Era of Space Exploration",
	"From Ancient Remedies to Modern Healthcare: Evolution of Medicine",
}

func dist(a, b string) float64 {
	lenA := getEncodeLen(a)
	lenB := getEncodeLen(b)
	cx := getEncodeLen(a + b)
	dist := (cx - min(lenA, lenB)) / max(lenA, lenB)
	return dist
}

func getEncodeLen(s string) float64 {
	buf := new(bytes.Buffer)
	writer := gzip.NewWriter(buf)
	writer.Write([]byte(s))
	writer.Close()
	return float64(buf.Len())
}

type entry struct {
	dist float64
	q    string
}

func main() {
	for _, s := range testSet {
		var distances []entry = make([]entry, len(trainingSet))

		var wg sync.WaitGroup

		for i, q := range trainingSet {
			wg.Add(1)
			go func(i int, q string) {
				dist := dist(s, q)
				e := entry{
					dist: dist,
					q:    q,
				}
				distances[i] = e
				wg.Done()
			}(i, q)
		}

		wg.Wait()

		slices.SortFunc(distances, func(lhs, rhs entry) int {
			if lhs.dist > rhs.dist {
				return 1
			}
			if lhs.dist < rhs.dist {
				return -1
			}
			return 0
		})

		fmt.Printf("input: %s, output: %s\n", s, distances[0].q)
		fmt.Printf("distances: %v\n\n", distances)
	}
}
