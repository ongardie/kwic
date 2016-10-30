package main

import (
	"io"
	"log"
	"os"
)

type lineHolder interface {
	// char returns the requested character of a word in a line.
	char(line, word, char int) byte
	// lines returns the total number of lines.
	lines() int
	// words returns the number of words in a line.
	words(line int) int
	// chars returns the number of characters in a word.
	chars(line, word int) int
}

// Module 1: Line Storage

type lineStorage struct {
	array [][][]byte
}

func (storage *lineStorage) char(line, word, char int) byte {
	return storage.array[line-1][word-1][char-1]
}

func (storage *lineStorage) lines() int {
	return len(storage.array)
}

func (storage *lineStorage) words(line int) int {
	return len(storage.array[line-1])
}

func (storage *lineStorage) chars(line, word int) int {
	return len(storage.array[line-1][word-1])
}

// setWord adds a character to the last word, a new word on the last line, or a
// new word on a new line.
func (storage *lineStorage) setWord(line, word, char int, value byte) {
	lines := storage.lines()
	if line < lines || line > lines+1 {
		panic("Line not last or just past last (ERLSBL)")
	}
	words := 0
	if line == lines {
		words = storage.words(line)
	}
	if word < words || word > words+1 {
		panic("Word not last or just past last (ERLSBW)")
	}
	chars := 0
	if line == lines && word == words {
		chars = storage.chars(line, word)
	}
	if char != chars+1 {
		panic("Char not just past last (ERLSBC)")
	}
	if line == lines+1 {
		storage.array = append(storage.array, nil)
	}
	if word == words+1 {
		storage.array[line-1] = append(storage.array[line-1], nil)
	}
	storage.array[line-1][word-1] = append(storage.array[line-1][word-1], value)
}

// deleteWord and deleteLine are unused and not implemented.

// Module 2: Input

func input(filename string, storage *lineStorage) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, 1)
	line, word, char := 1, 1, 1
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}
		if buf[0] == ' ' {
			word++
			char = 1
		} else if buf[0] == '\n' {
			line++
			word, char = 1, 1
		} else {
			storage.setWord(line, word, char, buf[0])
			char++
		}
	}
	return nil
}

// Module 3: Circular Shifter

func newCircularShifter(storage *lineStorage) lineHolder {
	shifted := &lineStorage{}
	line := 1
	for inputLine := 1; inputLine <= storage.lines(); inputLine++ {
		words := storage.words(inputLine)
		for inputStartWord := 1; inputStartWord <= words; inputStartWord++ {
			for word := 1; word <= words; word++ {
				inputWord := inputStartWord + word - 1
				if inputWord > words {
					inputWord -= words
				}
				for char := 1; char <= storage.chars(inputLine, inputWord); char++ {
					shifted.setWord(line, word, char, storage.char(inputLine, inputWord, char))
				}
			}
			line++
		}
	}
	return shifted
}

// Module 4: Alphabetizer

type alphabetizer struct {
	perm []int
}

func alphabetize(lines lineHolder) *alphabetizer {
	perm := make([]int, lines.lines())
	for i := range perm {
		perm[i] = i + 1
	}
	var quickSort func(left, right int)
	quickSort = func(left, right int) {
		if right-left <= 1 {
			return
		}
		pivot := left
		// Invariants: line[perm[elements left of pivot]] < line[perm[pivot]], i > pivot
		for i := pivot + 1; i < right; i++ {
			if linesLess(lines, perm[i], perm[pivot]) {
				if i == pivot+1 {
					perm[pivot], perm[i] = perm[i], perm[pivot]
				} else {
					perm[pivot], perm[pivot+1], perm[i] = perm[i], perm[pivot], perm[pivot+1]
				}
				pivot++
			}
		}
		quickSort(left, pivot)
		quickSort(pivot+1, right)
	}
	quickSort(0, len(perm))
	return &alphabetizer{perm}
}

// ith returns the index of the circular shift that comes i-th in alphabetical
// order.
func (alpha *alphabetizer) ith(i int) int {
	return alpha.perm[i-1]
}

func normalizeChar(char byte) byte {
	if char >= 'A' && char <= 'Z' {
		return (char - 'A') * 2
	}
	if char >= 'a' && char <= 'z' {
		return (char-'a')*2 + 1
	}
	return 0
}

func wordsLess(lines lineHolder, line1, word1, line2, word2 int) bool {
	chars1 := lines.chars(line1, word1)
	chars2 := lines.chars(line2, word2)
	char := 1
	for {
		if char > chars1 && char <= chars2 {
			return true
		}
		if char > chars2 {
			return false
		}
		n1 := normalizeChar(lines.char(line1, word1, char))
		n2 := normalizeChar(lines.char(line2, word2, char))
		if n1 < n2 {
			return true
		} else if n1 > n2 {
			return false
		}
		char++
	}
}

func linesLess(lines lineHolder, line1, line2 int) bool {
	words1 := lines.words(line1)
	words2 := lines.words(line2)
	word := 1
	for {
		if word > words1 && word <= words2 {
			return true
		}
		if word > words2 {
			return false
		}
		if wordsLess(lines, line1, word, line2, word) {
			return true
		}
		if wordsLess(lines, line2, word, line1, word) {
			return false
		}
		word++
	}
}

// wordsEqual and linesEqual are unused and not implemented.

// Module 5: Output

func output(w io.Writer, lines lineHolder, line int) {
	for word := 1; word <= lines.words(line); word++ {
		for char := 1; char <= lines.chars(line, word); char++ {
			w.Write([]byte{lines.char(line, word, char)})
		}
		if word < lines.words(line) {
			w.Write([]byte{' '})
		}
	}
	w.Write([]byte{'\n'})
}

// Module 6: Master Control

func main() {
	filename := "input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	storage := &lineStorage{}
	err := input(filename, storage)
	if err != nil {
		log.Fatalf("Error in input(%v): %v", filename, err)
	}
	shifter := newCircularShifter(storage)
	alpha := alphabetize(shifter)
	for line := 1; line <= shifter.lines(); line++ {
		output(os.Stdout, shifter, alpha.ith(line))
	}
}
