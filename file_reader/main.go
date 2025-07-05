package main

import (
	"fmt"
	"os"
)

func take_file_path() (path string) {
	fmt.Print("please write a file path: ")
	fmt.Scan(&path)
	return path
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	path := take_file_path()

	dat, err := os.ReadFile(path)
	check(err)
	fmt.Print(dat)

	f, err := os.Open(path)
	check(err)

	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check(err)
	fmt.Println()
	fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))
}
