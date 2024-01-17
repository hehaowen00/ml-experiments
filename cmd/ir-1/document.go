package main

import (
	"strings"
	"unicode"
)

type Document struct {
	data []string
}

func NewDocument(line string) *Document {
	parts := strings.Split(line, "|||")
	title := strings.ToLower(strings.TrimSpace(parts[0]))
	summary := strings.ToLower(strings.TrimSpace(parts[1]))

	return &Document{
		data: []string{title, summary},
	}
}

func (doc *Document) Title() string {
	return doc.data[0]
}

func (doc *Document) Summary() string {
	return doc.data[1]
}

type DocumentIndex struct {
	tf map[string]int
}

func ParseDocument(doc *Document) *DocumentIndex {
	tf := map[string]int{}

	var acc []rune
	summary := []rune(doc.Summary())

	for i := 0; i < len(summary); i++ {
		if summary[i] == rune('(') {
			continue
		}
		if summary[i] == rune(')') {
			continue
		}

		if summary[i] == rune(',') {
			continue
		}

		if summary[i] == rune('"') {
			continue
		}

		if summary[i] == rune('-') {
			tf[string(acc)] += tf[string(acc)] + 1
			acc = nil
			continue
		}

		acc = append(acc, summary[i])
		if unicode.IsSpace(summary[i]) {
			if string(acc) == "..." {
				acc = nil
				continue
			}

			tf[string(acc)] += tf[string(acc)] + 1
			acc = nil
		}
	}

	if len(acc) > 0 {
		tf[string(acc)] += tf[string(acc)] + 1
	}

	return &DocumentIndex{
		tf: tf,
	}
}
