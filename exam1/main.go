package main

import (
	"fmt"
)

func printStarPattern(x int) {
	for i := 1 ; i <= x ; i++ {
		fmt.Println(printStar("*", i))
	}
	for i := x-1 ; i > 0 ; i-- {
		fmt.Println(printStar("*", i))
	}
}

func printStar(s string, count int) string {
	star := ""
	for i := 1 ; i<= count ; i++ {
		star += s
	}
	return star
}


func main() {
	var x int
	fmt.Print("Enter a number: ")
	fmt.Scan(&x)
	printStarPattern(x)
}