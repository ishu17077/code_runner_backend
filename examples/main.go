package main

import "fmt"

func main() {
	var num int
	fmt.Scanf("%d", &num)
	if num%2 == 0 {
		fmt.Printf("Yes")
	} else {
		fmt.Printf("No")
	}
}
