package sudoku

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

// Row slice of ints
type Row []int

// Sudoku slice of rows
type Sudoku []Row

type elibilityMap map[int]bool
type cell struct {
	rowIdx      int
	colIdx      int
	eligibleNos []int
}

type statusChannel struct {
	interimSolution Sudoku
	solved          bool
	err             error
}

func isComplete(r []int) bool {
	completionMap := make(map[int]int)
	if len(r) < 9 {
		return false
	}
	for _, el := range r {
		if _, found := completionMap[el]; found {
			completionMap[el]++
		} else {
			completionMap[el] = 1
		}
	}
	sequence := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, i := range sequence {
		if completionMap[i] != 1 {
			return false
		}
	}
	return true
}

func (em elibilityMap) init() {
	for i := 1; i <= 9; i++ {
		em[i] = true
	}
}

func (em elibilityMap) list() []int {
	var numberList []int
	for key, val := range em {
		if val {
			numberList = append(numberList, key)
		}
	}
	return numberList
}

func (s Sudoku) copy() Sudoku {
	copied := make(Sudoku, 0)
	done := make(chan int)
	go func() {
		for _, _Row := range s {
			myRow := make(Row, 0)
			for _, _col := range _Row {
				myRow = append(myRow, _col)
			}
			copied = append(copied, myRow)
		}
		done <- 0
	}()
	<-done
	return copied
}
func (s Sudoku) countUnfilledCells() int {
	unfilledCount := 0
	for _, row := range s {
		for _, col := range row {
			if col == 0 {
				unfilledCount++
			}
		}
	}
	return unfilledCount
}

func (s Sudoku) getNumbersFilledInRow(rowIdx int) []int {
	var filled []int
	row := s[rowIdx]
	for _, number := range row {
		if number != 0 {
			filled = append(filled, number)
		}
	}
	return filled
}

func (s Sudoku) getNumbersFilledInColumn(colIdx int) []int {
	var filled []int
	for _, row := range s {
		if row[colIdx] != 0 {
			filled = append(filled, row[colIdx])
		}
	}
	return filled
}

func (s Sudoku) getNumbersFilledInTile(rowIdx int, colIdx int) []int {
	var filled []int
	var startRow, endRow, startCol, endCol int

	if rowIdx < 3 {
		startRow = 0
		endRow = 2
	} else if rowIdx < 6 {
		startRow = 3
		endRow = 5
	} else {
		startRow = 6
		endRow = 8
	}
	if colIdx < 3 {
		startCol = 0
		endCol = 2
	} else if colIdx < 6 {
		startCol = 3
		endCol = 5
	} else {
		startCol = 6
		endCol = 8
	}

	for i := startRow; i <= endRow; i++ {
		for j := startCol; j <= endCol; j++ {
			if s[i][j] != 0 {
				filled = append(filled, s[i][j])
			}
		}
	}
	return filled
}

func (s Sudoku) solved() bool {
	// check if rows are complete
	for rowIdx, row := range s {
		filledInRow := s.getNumbersFilledInRow(rowIdx)
		if !isComplete(filledInRow) {
			return false
		}
		for i := 0; i < len(row); i++ {
			colIdx := i
			filledInColumn := s.getNumbersFilledInColumn(colIdx)
			if !isComplete(filledInColumn) {
				return false
			}
			filledInTile := s.getNumbersFilledInTile(rowIdx, colIdx)
			if !isComplete(filledInTile) {
				return false
			}
		}
	}
	return true
}

func (s Sudoku) mapEligibleNumbers(rowIdx int, colIdx int) []int {
	filledInRow := s.getNumbersFilledInRow(rowIdx)
	filledInColumn := s.getNumbersFilledInColumn(colIdx)
	filledInTile := s.getNumbersFilledInTile(rowIdx, colIdx)
	negativeList := append(filledInRow, filledInColumn...)
	negativeList = append(negativeList, filledInTile...)

	em := make(elibilityMap)
	em.init()
	// eliminate numbers in present in the negative list
	for _, number := range negativeList {
		em[number] = false
	}

	eligibleNos := em.list()
	return eligibleNos
}

func (s Sudoku) fill(rowIdx int, colIdx int, toFill int) {
	done := make(chan int)
	go func(s Sudoku) {
		s[rowIdx][colIdx] = toFill
		done <- toFill
	}(s)
	<-done
}

func (s Sudoku) fillElibileNumber(c cell) (int, error) {
	filled := 0
	if len(c.eligibleNos) == 9 {
		return 0, errors.New("incorrect Sudoku")
	}
	// fill only if there is a single eligible number
	if len(c.eligibleNos) == 1 {
		toFill := c.eligibleNos[0]
		s.fill(c.rowIdx, c.colIdx, toFill)
		filled = toFill
	} else {
		// log.Println("multiple numbers eligible:no fill")
	}
	return filled, nil
}

// Solve fills the sudoku first by naive approach then by brute forcing
func (s Sudoku) Solve() (Sudoku, bool, int, error) {
	toSolve := s.copy()
	solved := toSolve.solved()
	iterations := 0
	for !solved {
		unfilledCount := toSolve.countUnfilledCells()
		iterations++
		// fill all cells wil single options
		for rowIdx, row := range toSolve {
			for colIdx, col := range row {
				if col == 0 {
					numbers := toSolve.mapEligibleNumbers(rowIdx, colIdx)
					c := cell{rowIdx, colIdx, numbers}
					filled, err := toSolve.fillElibileNumber(c)
					if filled > 0 {
						// log.Printf("Filled %d,%d : %d \n", rowIdx, colIdx, filled)
					}
					if err != nil {
						log.Fatal("Fill error", err)
					}
				}
			}
		}

		if toSolve.solved() {
			// done
			break
		}

		currentUnfilled := toSolve.countUnfilledCells()
		fmt.Println("Unfilled count ----->", currentUnfilled)
		if currentUnfilled >= unfilledCount {
			// naive approach has reached a dead end, use brute force
			var cellWithLeastOptions cell
			maxOptions := 9
			fmt.Println("start brute force")
			for rowIdx, row := range toSolve {
				for colIdx, col := range row {
					if col == 0 {
						eligibleNumbers := toSolve.mapEligibleNumbers(rowIdx, colIdx)
						candidateCount := len(eligibleNumbers)
						if candidateCount > 0 {
							if candidateCount <= maxOptions {
								cellWithLeastOptions = cell{rowIdx, colIdx, eligibleNumbers}
								maxOptions = candidateCount
							}
						}
					}
				}
			}

			// for handling concurrent brute forcing
			solvedStatusChannel := make(chan statusChannel)
			wg := new(sync.WaitGroup)
			fmt.Println("Eligible options", cellWithLeastOptions.eligibleNos)
			for _, number := range cellWithLeastOptions.eligibleNos {
				fmt.Printf("Testing %d at %d,%d\n", number, cellWithLeastOptions.rowIdx, cellWithLeastOptions.colIdx)
				wg.Add(1)

				go func(inSudoku Sudoku, rowIdx int, colIdx int, toFill int, wg *sync.WaitGroup, c *chan statusChannel) {
					defer wg.Done()
					outSudoku := inSudoku.copy()
					outSudoku.fill(rowIdx, colIdx, toFill)
					interimSudoku, solved, _, err := outSudoku.Solve()
					*c <- statusChannel{interimSudoku, solved, err}
				}(toSolve, cellWithLeastOptions.rowIdx, cellWithLeastOptions.colIdx, number, wg, &solvedStatusChannel)
			}

			// wait for all go routines to execute
			go func(wg *sync.WaitGroup, c chan statusChannel) {
				wg.Wait()
				close(c)
			}(wg, solvedStatusChannel)

			// check all interim solutions to check if filled
			for interim := range solvedStatusChannel {
				sudoku := interim.interimSolution
				solved := interim.solved
				err := interim.err

				if err != nil {
					log.Fatal("channel error", err.Error())
				}

				if solved {
					log.Println("Solved")
					return sudoku, solved, iterations, nil
				} else {
					// not solved, but the guess is correct. try from beginning
					toSolve = sudoku.copy()
					break
				}
			}
		}

	}
	return toSolve, toSolve.solved(), iterations, nil
}
