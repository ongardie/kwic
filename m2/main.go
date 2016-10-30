package main

import (
	"io"
	"os"
)

type lineHolder interface {
	// word returns the requested character of a word in a line.
	word(line, word, char int) byte
	// words returns the number of words in a line.
	words(line int) int
	// lines returns the total number of lines.
	lines() int
	// chars returns the number of characters in a word.
	chars(line, word int) int
}

// Module 1: Line Storage

type lineStorage struct{}

func (storage *lineStorage) word(line, word, char int) byte {
	return 0
}

// setWord adds a character to the last word, a new word on the last line, or a
// new word on a new line.
func (storage *lineStorage) setWord(line, word, char int, value byte) {
}

func (storage *lineStorage) words(line int) int {
	return 0
}

func (storage *lineStorage) lines() int {
	return 0
}

// deleteWord and deleteLine are unused and not implemented.

func (storage *lineStorage) chars(line, word int) int {
	return 0
}

// Module 2: Input

func input(filename string, storage *lineStorage) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// Module 3: Circular Shifter

type circularShifter struct{}

func setupCircularShifter(storage *lineStorage) *circularShifter {
	return nil
}

func (shifter *circularShifter) word(line, word, char int) byte {
	return 0
}

func (shifter *circularShifter) words(line int) int {
	return 0
}

func (shifter *circularShifter) lines() int {
	return 0
}

func (shifter *circularShifter) chars(line, word int) int {
	return 0
}

// Module 4: Alphabetizer

type alphabetizer struct{}

func alphabetize(shifter *circularShifter) *alphabetizer {
	return nil
}

// ith returns the index of the circular shift that comes i-th in alphabetical
// order.
func (alpha *alphabetizer) ith(i int) int {
	return 0
}

func normalizeChar(char byte) byte {
	return char
}

func wordsEqual(lines *lineHolder, line1, word1, line2, word2 int) bool {
	return false
}

func wordsStrictlyLess(lines *lineHolder, line1, word1, line2, word2 int) bool {
	return false
}

func linesEqual(lines *lineHolder, line1, line2 int) bool {
	return false
}

func linesStrictlyLess(lines *lineHolder, line1, line2 int) bool {
	return false
}

// Module 5: Output

func output(w io.Writer, shifter *circularShifter, i int) {
}

// Module 6: Master Control

func main() {
}
