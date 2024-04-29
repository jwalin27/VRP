package main

import (
	"fmt"
	"os"

	"github.com/Vorto-interview/service"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_input_file>")
		os.Exit(1)
	}
	filePath := os.Args[1]

	loads, err := service.ParseAllLoads(filePath)
	if err != nil {
		fmt.Printf("failed to parse loads: %v", err)
	}

	solver := service.NewNearestNeighborSolver(loads)
	routes, err := solver.PlanRoutes()
	if err != nil {
		fmt.Printf("failed to plan routes: %v", err)
	}

	for _, r := range routes {
		r.PrintLoadNumbers()
	}
}
