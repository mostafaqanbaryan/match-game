package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

type Board struct {
	cells     []Cell
	selection []int
	cursor    int
	cols      int
}

func (b Board) shuffle() Board {
	total := len(b.cells)
	rand.Shuffle(total, func(i int, j int) {
		b.cells[i], b.cells[j] = b.cells[j], b.cells[i]
	})
	return b
}

func (b Board) selectCell() Board {
	b.selection = append(b.selection, b.cursor)
	return b
}

func (b Board) isCellSelected(index int) bool {
	for _, selectedIndex := range b.selection {
		if index == selectedIndex {
			return true
		}
	}
	return false
}

func (b Board) newSelection() Board {
	b.selection = make([]int, 0)
	return b
}

func (b Board) isCellAlreadyMatched() bool {
	return b.cells[b.cursor].matched
}

func (b Board) isSelectionMatched() bool {
	return len(b.selection) == 2 && b.cells[b.selection[0]].value == b.cells[b.selection[1]].value
}

func (b Board) isCompleted() bool {
	for _, cell := range b.cells {
		if !cell.matched {
			return false
		}
	}
	return true
}

func (b Board) match() Board {
	if b.isSelectionMatched() {
		b.cells[b.selection[0]].matched = true
		b.cells[b.selection[1]].matched = true
	}
	return b.newSelection()
}

func (b Board) view() string {
	s := ""
	s += "    ╭"
	s += strings.Repeat("-", 15)
	s += "╮"
	for row := 1; row <= b.cols; row++ {
		if row > 1 {
			s += "\n"
			s += "    ├─" + strings.Repeat("──┼─", 3) + "──┤"
		}
		s += "\n    │"
		for col := 1; col <= b.cols; col++ {
			index := ((row - 1) * b.cols) + col - 1
			cell := b.cells[index]
			value := "X"
			if cell.matched || b.isCellSelected(index) {
				value = cell.value
			}
			rPadding := 1
			lPadding := 0
			if len(value) == 1 {
				lPadding = 1
				rPadding = 1
			}
			if b.cursor == index {
				value = fmt.Sprintf("[%s]", value)
				lPadding = 0
				rPadding = 0
			}
			s += fmt.Sprintf("%*s%s%*s│", lPadding, "", value, rPadding, "")
		}
	}
	s += "\n"
	s += "    ╰"
	s += strings.Repeat("-", 15)
	s += "╯"
	s += "\n"
	return s
}

func NewBoard(level Level) Board {
	total := int(level)
	cells := make([]Cell, total)
	for i := 0; i < total; i++ {
		index := i
		if i >= total/2 {
			index -= total / 2
		}
		cells[i] = Cell{
			value:   fmt.Sprintf("%d", index+1),
			matched: false,
		}
	}

	b := Board{
		cells:  cells,
		cols:   int(math.Sqrt(float64(total))),
		cursor: 0,
	}
	b.newSelection()
	return b.shuffle()
}
