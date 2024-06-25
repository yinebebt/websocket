package main

import "fmt"

func main() {
	items := [][2]byte{{1, 2}, {3, 4}, {5, 6}}
	a := [][]byte{}

	for _, i := range items {
		a = append(a, i[:])
	}

	fmt.Println(a)
}
