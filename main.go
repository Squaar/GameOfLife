package main

import (
	"code.google.com/p/goncurses"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	stdscr, err := goncurses.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer goncurses.End()

	board := NewBoard(20, 20, 50)
	board.Print(stdscr)

	for {
		board = board.Tick()
		board.Print(stdscr)
		stdscr.Refresh()
		stdscr.GetChar()
	}
}

type Board struct {
	grid      [][]rune
	w, h, pop int
}

func NewBoard(w, h, pop int) *Board {
	grid := make([][]rune, h)
	for i := range grid {
		grid[i] = make([]rune, w)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for i := 0; i < pop; i++ {
		x, y := rand.Intn(w), rand.Intn(h)
		if grid[y][x] == 'X' {
			i--
		} else {
			grid[y][x] = 'X'
		}
	}

	return &Board{grid, w, h, pop}
}

func (board *Board) Print(window *goncurses.Window) {
	for i := range board.grid {
		for j := range board.grid[i] {
			window.MovePrintf(i, j, "%c", board.grid[i][j])
		}
	}
}

func (board *Board) Tick() *Board {
	newBoard := NewBoard(board.w, board.h, 0)
	for i := range board.grid {
		for j := range board.grid[i] {
			adjacent := board.CountAdjacent(j, i)
			if (board.grid[i][j] == 'X' && (adjacent == 2 || adjacent == 3)) || (board.grid[i][j] == ' ' && adjacent == 3) {
				newBoard.grid[i][j] = 'X'
				newBoard.pop++
			}
		}
	}
	return newBoard
}

func (board *Board) CountAdjacent(x, y int) (count int) {
	count = 0
	if x-1 >= 0 && y-1 >= 0 && board.grid[y-1][x-1] == 'X' {
		count++
	}
	if y-1 >= 0 && board.grid[y-1][x] == 'X' {
		count++
	}
	if y-1 >= 0 && x+1 < board.w && board.grid[y-1][x+1] == 'X' {
		count++
	}
	if x+1 < board.w && board.grid[y][x+1] == 'X' {
		count++
	}
	if y+1 < board.h && x+1 < board.w && board.grid[y+1][x+1] == 'X' {
		count++
	}
	if y+1 < board.h && board.grid[y+1][x] == 'X' {
		count++
	}
	if y+1 < board.h && x-1 >= 0 && board.grid[y+1][x-1] == 'X' {
		count++
	}
	if x-1 >= 0 && board.grid[y][x-1] == 'X' {
		count++
	}
	return
}
