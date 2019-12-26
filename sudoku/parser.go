package sudoku

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ParseFromFile parses file to read sudoku input
func ParseFromFile(filePath string) Sudoku {
	var parsedSudoku Sudoku
	f, _ := os.Open(filePath)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var parsedRow Row
		line := scanner.Text()
		inputRow := strings.Split(line, ",")
		for _, el := range inputRow {
			intVal, _ := strconv.Atoi(el)
			parsedRow = append(parsedRow, intVal)
		}
		parsedSudoku = append(parsedSudoku, parsedRow)
	}

	return parsedSudoku
}
