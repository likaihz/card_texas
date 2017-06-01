package main

import (
	// "flag"
	"fmt"
	//"os"
)

func main() {
	var c [5]int

	for i := 0; i < 5; i++ {
		fmt.Scanf("%d\n", &c[i])
	}
	for i := 0; i < 5; i++ {
		fmt.Println(c[i])
	}
}
