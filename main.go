package main

import (
	"fmt"
	"log"
	"time"

	"github.com/h-w-s/brute-doku/sudoku"
)

func main() {
	fmt.Println("Start solving")
	start := time.Now()
	sudokuIn := sudoku.ParseFromFile("./sample.txt")
	fmt.Println(sudokuIn)
	sudokuOut, solved, interations, err := sudokuIn.Solve()
	if err != nil {
		log.Fatal("ERR:", err.Error())
	}
	if !solved {
		log.Fatal("Unsolved oops")
	}
	fmt.Println("solution", sudokuOut)
	elapsed := time.Since(start)
	fmt.Println("Total time taken: ", elapsed)
	fmt.Println("Total interations: ", interations)
}
