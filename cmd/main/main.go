package main

import (
	"app"
	"fmt"
)

func init() {
	app.Import()
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error:", err)
		}

		// print stack trace
		println(app.ColorRed)
		println("-----------------------------------------------------------------------")
		println("End of program!")
		println(app.ColorReset)
	}()

	fmt.Println("Welcome to the Command Line Interface of Application\n" +
		"\r\nSelect a command:" +
		"\r\n1. Create a new project" +
		"\r\n2. Create a new project" +
		"\r\n3. Create a new project" +
		"\r\n4. Create a new project" +
		"\r\n5. Create a new project" +
		"\r\n6. Create a new project",
	)

	var input string
	fmt.Print("\r\nPlease enter a command: ")
	fmt.Scanln(&input)
	// TODO: add a switch case for the input
	switch input {
	case "1":
		fmt.Printf("-> Execute a command: (%s)\n", input)
	case "2":
		fmt.Printf("-> Execute a command: (%s)\n", input)
	default:
		fmt.Printf("-> Command `%s` is not defined.\nPlease enter a command from the list (1 >> 2).\n", input)
	}

}
