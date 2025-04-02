package main

import (
	"R2-D2/cmd"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Check if any arguments were provided
	if len(os.Args) > 1 {
		// Use Cobra CLI normally if arguments are provided
		if err := cmd.Execute(); err != nil {
			panic(err)
		}
		return
	}

	// Otherwise, start REPL mode
	fmt.Println("Welcome to R2-D2 Todo REPL")
	fmt.Println("Type 'help' to see available commands, or 'exit' to quit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" || input == "quit" {
			break
		}

		// Execute the command through Cobra
		args := strings.Fields(input)
		if len(args) > 0 {
			// Save original os.Args
			originalArgs := os.Args
			// Set new args for Cobra to process
			os.Args = append([]string{"R2-D2"}, args...)

			// Execute the command
			err := cmd.Execute()
			if err != nil {
				fmt.Println("Error:", err)
			}

			// Restore original args
			os.Args = originalArgs
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
