package main

import (
	"strings"

	"github.com/aaaton/golem/v4"
)

type Document struct {
	ID      int
	Title   string
	Summary string
}

func NewDocument(id int, line string) *Document {
	parts := strings.Split(line, "|||")
	title := strings.TrimSpace(parts[0])
	summary := strings.TrimSpace(parts[1])

	return &Document{
		ID:      id,
		Title:   title,
		Summary: summary,
	}
}

type DocumentIndex struct {
	ID         int
	Tokens     map[string]int
	TF         map[string]float64
	TokenCount int
}

func ParseDocument(lemmatizer *golem.Lemmatizer, id int, dataset string, punctuation, stopList []string) *DocumentIndex {
	tokens := map[string]int{}

	words := processDataset(lemmatizer, dataset, punctuation, stopList)

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

	return &DocumentIndex{
		ID:         id,
		Tokens:     tokens,
		TF:         tf,
		TokenCount: count,
	}
}

func processDataset(lemmatizer *golem.Lemmatizer, data string, punctuation, stopList []string) []string {
	// to lowercase
	data = strings.ToLower(data)

	// remove punctuation
	for _, p := range punctuation {
		data = strings.ReplaceAll(data, p, " ")
	}

	acc := []string{}

	stopListTable := map[string]struct{}{}
	for _, s := range stopList {
		stopListTable[s] = struct{}{}
	}

	// remove stop words
	for _, x := range strings.Split(data, " ") {
		x = strings.TrimSpace(x)

		if _, ok := stopListTable[x]; ok {
			continue
		}

		for _, p := range punctuation {
			x = strings.ReplaceAll(x, p, " ")
		}

		// lemmatize
		w := lemmatizer.LemmaLower(x)
		acc = append(acc, w)
	}

	return acc
}
