package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("exit from main().")
	os.Exit(0) // want "os.Exit exists in main body"
}
