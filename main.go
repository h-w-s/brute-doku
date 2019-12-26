package main

import (
	"fmt"
	"time"

	"github.com/h-w-s/brute-doku/sudoku"
)

func main() {
	fmt.Println("Start solving")
	start := time.Now()
	sudokuIn := sudoku.ParseFromFile("./sample.txt")
	fmt.Println(sudokuIn)
	elapsed := time.Since(start)
	fmt.Println("Total time taken: ", elapsed)
}
