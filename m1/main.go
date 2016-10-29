package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type shiftedLine struct {
	line      int
	startChar int
}

// Module 1: Input

func input(filename string) ([]byte, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	chars := make([]byte, stat.Size())
	_, err = io.ReadFull(file, chars)
	if err != nil {
		return nil, nil, err
	}
	lines := []int{0}
	for i, char := range chars {
		if char == '\n' {
			chars[i] = ' '
			lines = append(lines, i+1)
		}
	}
	return chars, lines, nil
}

// Module 2: Circular Shift

func circularShifts(chars []byte, lines []int) []shiftedLine {
	shifts := []shiftedLine{}
	for line, startChar := range lines {
		var endChar int // exclusive
		if line+1 < len(lines) {
			endChar = lines[line+1] - 1
		} else {
			endChar = len(chars)
		}
		wordBreak := true
		for i := startChar; i < endChar; i++ {
			if wordBreak {
				shifts = append(shifts, shiftedLine{line, i})
			}
			wordBreak = (chars[i] == ' ')
		}
	}
	return shifts
}

// Module 3: Alphabetizing

func alphabetize(chars []byte, lines []int, unsorted []shiftedLine) []shiftedLine {
	// Copy unsorted so we can sort it in place.
	sorted := make([]shiftedLine, 0, len(unsorted))
	for _, shift := range unsorted {
		sorted = append(sorted, shift)
	}

	// strcmp for first word found at startChar
	less := func(a, b shiftedLine) bool {
		i := a.startChar
		j := b.startChar
		for {
			if chars[i] < chars[j] {
				return true
			}
			if chars[j] < chars[i] {
				return false
			}
			if i == len(chars) || chars[i] == ' ' {
				return true
			}
			if j == len(chars) || chars[j] == ' ' {
				return false
			}
			i++
			j++
		}
		return true
	}

	var quickSort func(left, right int)
	quickSort = func(left, right int) {
		if right-left <= 1 {
			return
		}
		pivot := left
		// Invariants: elements left of pivot are less than pivot, i > pivot
		for i := pivot + 1; i < right; i++ {
			if less(sorted[i], sorted[pivot]) {
				if i == pivot+1 {
					sorted[pivot], sorted[i] = sorted[i], sorted[pivot]
				} else {
					sorted[pivot], sorted[pivot+1], sorted[i] = sorted[i], sorted[pivot], sorted[pivot+1]
				}
				pivot++
			}
		}
		quickSort(left, pivot)
		quickSort(pivot+1, right)
	}
	quickSort(0, len(sorted))
	return sorted
}

// Module 4: Output

func output(w io.Writer, chars []byte, lines []int, shifts []shiftedLine) {
	for _, shift := range shifts {
		lineStart := lines[shift.line]
		var lineEnd int // exclusive
		if shift.line+1 < len(lines) {
			lineEnd = lines[shift.line+1] - 1
		} else {
			lineEnd = len(chars)
		}
		if lineStart < shift.startChar {
			fmt.Fprintf(w, "%s | %s [%d]\n",
				string(chars[shift.startChar:lineEnd]),
				string(chars[lineStart:shift.startChar-1]),
				shift.line+1)
		} else {
			fmt.Fprintf(w, "%s [%d]\n",
				string(chars[shift.startChar:lineEnd]),
				shift.line+1)
		}
	}
}

// Module 5: Master Control

func main() {
	filename := "input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	chars, lines, err := input(filename)
	if err != nil {
		log.Fatalf("Error in input(%v): %v", filename, err)
	}
	unsortedShifts := circularShifts(chars, lines)
	shifts := alphabetize(chars, lines, unsortedShifts)
	output(os.Stdout, chars, lines, shifts)
}
