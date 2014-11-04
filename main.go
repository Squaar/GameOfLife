package main

import (
	"bytes"
	"code.google.com/p/goncurses"
	"crypto/md5"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// con := flag.Bool("c", false, "Run continuously instead of pressing a key for each step")
	width := flag.Int("w", 20, "Width of the game board.")
	height := flag.Int("h", 20, "Height of the game board.")
	pop := flag.Int("p", 50, "The starting population of the game board.")
	// seed := flag.Int64("s", time.Now().Unix(), "The seed to generate the board with.")
	seed := flag.String("s", "", "The seed to generate the board with.")
	flag.Parse()

	stdscr, err := goncurses.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer goncurses.End()

	if *seed != "" {
		rand.Seed(hash(*seed))
	} else {
		rand.Seed(time.Now().Unix())
	}

	board := NewBoard(*width, *height, *pop)
	board.Print(stdscr)

	for {
		board = board.Tick()
		board.Print(stdscr)
		stdscr.Refresh()
		if stdscr.GetChar() == 27 {
			goncurses.End()
			os.Exit(0)
		}
	}
}

func hash(s string) int64 {
	sum := md5.Sum([]byte(s))
	return bytes2int64(sum[:8]) ^ bytes2int64(sum[8:])
}

func bytes2int64(b []byte) (n int64) {
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &n)
	return
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
