package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const gridSize = 9

type Sudoku struct {
	grid [gridSize][gridSize]int
}

func (s *Sudoku) printBoard() {
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if s.grid[i][j] == 0 {
				fmt.Print("■ ")
			} else {
				fmt.Print(s.grid[i][j], " ")
			}
		}
		fmt.Println()
	}
}

func (s *Sudoku) isValidMove(row, col, num int) bool {
	for i := 0; i < gridSize; i++ {
		if s.grid[row][i] == num || s.grid[i][col] == num {
			return false
		}
	}

	startRow, startCol := (row/3)*3, (col/3)*3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.grid[startRow+i][startCol+j] == num {
				return false
			}
		}
	}
	return true
}

func (s *Sudoku) makeMove(row, col, num int) bool {
	if row >= 0 && row < gridSize && col >= 0 && col < gridSize && s.grid[row][col] == 0 {
		if s.isValidMove(row, col, num) {
			s.grid[row][col] = num
			return true
		}
	}
	return false
}

func generateSudoku() Sudoku {
	rand.Seed(time.Now().UnixNano())
	s := Sudoku{}
	for i := 0; i < 10; i++ {
		row, col, num := rand.Intn(gridSize), rand.Intn(gridSize), rand.Intn(9)+1
		if s.isValidMove(row, col, num) {
			s.grid[row][col] = num
		}
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for level := 1; level <= 10; level++ {
		fmt.Printf("Level %d\n", level)
		sudoku := generateSudoku()
		sudoku.printBoard()
		for {
			fmt.Print("Enter your move (row column number): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			parts := strings.Split(input, " ")
			if len(parts) != 3 {
				fmt.Println("Please enter in the correct format: row column number")
				continue
			}

			row, err1 := strconv.Atoi(parts[0])
			col, err2 := strconv.Atoi(parts[1])
			num, err3 := strconv.Atoi(parts[2])

			if err1 != nil || err2 != nil || err3 != nil || row < 1 || row > 9 || col < 1 || col > 9 || num < 1 || num > 9 {
				fmt.Println("Invalid input, please try again.")
				continue
			}

			if sudoku.makeMove(row-1, col-1, num) {
				sudoku.printBoard()
			} else {
				fmt.Println("Invalid move.")
			}
		}
	}
}
