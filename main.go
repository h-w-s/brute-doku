package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Start solving")
	start := time.Now()
	elapsed := time.Since(start)
	fmt.Println("Total time taken: ", elapsed)
}
