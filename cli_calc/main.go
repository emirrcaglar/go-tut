package main

import "fmt"

func add(x, y int) int {
	return x + y
}

func subtract(x, y int) int {
	return x - y
}

func multiply(x, y int) int {
	return x * y
}

func divide(x, y int) (n int, err error) {
	if y == 0 {
		err = fmt.Errorf("cannot divide by zero")
		return
	}
	return x / y, nil
}

func main() {
	var operation string

	fmt.Print("Type your operation: ")
	fmt.Scan((&operation))

	var allowedOps = map[string]bool{
		"add":      true,
		"subtract": true,
		"multiply": true,
		"divide":   true,
	}
	if allowedOps[operation] {

		var x, y int

		fmt.Println("Type two numbers: ")
		fmt.Println("Careful: the order matters")
		fmt.Print("Number 1: ")
		fmt.Scan(&x)
		fmt.Print("Number 2: ")
		fmt.Scan(&y)

		var result int
		var err error

		switch operation {

		case "add":
			result = add(x, y)
		case "subtract":
			result = subtract(x, y)
		case "multiply":
			result = multiply(x, y)
		case "divide":
			result, err = divide(x, y)
		}
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Result: ", result)

	}

}
