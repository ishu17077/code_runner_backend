package main

import (
	"fmt"

	CompositionRoot "github.com/ishu17077/code_runner_backend/composition_root"
)

func main() {
	if err := CompositionRoot.Start(); err != nil {
		fmt.Printf("Error starting the server: %s\n", err.Error())
		panic(err)
	}
}
