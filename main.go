package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: aux4-secret-aux4 <command> [options]\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		runGet(os.Args[2:])
	case "create":
		runCreate(os.Args[2:])
	case "set":
		runSet(os.Args[2:])
	case "list":
		runList(os.Args[2:])
	case "search":
		runSearch(os.Args[2:])
	case "remove":
		runRemove(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "Available commands: get, create, set, list, search, remove\n")
		os.Exit(1)
	}
}
