package main

import (
	"fmt"
	"os"
	"strings"

	tc "github.com/stef41/llm-tokencount"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: llm-tokencount <text> [model]\n")
		fmt.Fprintf(os.Stderr, "       llm-tokencount --models\n")
		os.Exit(1)
	}

	if os.Args[1] == "--models" {
		for _, m := range tc.ListModels() {
			fmt.Printf("%-20s %-10s %d tokens  $%.2f/$%.2f per 1M\n",
				m.Name, m.Provider, m.MaxTokens, m.InputCostPer1M, m.OutputCostPer1M)
		}
		return
	}

	text := strings.Join(os.Args[1:], " ")
	model := "gpt-4o"

	// Check if last arg is a known model
	if len(os.Args) > 2 {
		if _, ok := tc.GetModel(os.Args[len(os.Args)-1]); ok {
			model = os.Args[len(os.Args)-1]
			text = strings.Join(os.Args[1:len(os.Args)-1], " ")
		}
	}

	tokens := tc.CountTokens(text)
	fmt.Printf("Tokens: %d\n", tokens)

	est, err := tc.EstimateCost(model, tokens, tokens/2)
	if err == nil {
		fmt.Printf("Model:  %s\n", model)
		fmt.Printf("Cost:   $%.6f (input: $%.6f, output: $%.6f)\n", est.TotalCost, est.InputCost, est.OutputCost)
	}

	fits, remaining := tc.FitsContext(model, tokens)
	if fits {
		fmt.Printf("Context: fits (%d remaining)\n", remaining)
	} else {
		fmt.Printf("Context: EXCEEDS by %d tokens\n", -remaining)
	}
}
