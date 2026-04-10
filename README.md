# llm-tokencount

Go library for approximate token counting and cost estimation for LLM APIs. Supports OpenAI, Anthropic, and Google models.

## Installation

```bash
go get github.com/stef41/llm-tokencount
```

## Usage

```go
package main

import (
    "fmt"
    tc "github.com/stef41/llm-tokencount"
)

func main() {
    tokens := tc.CountTokens("Hello, how are you?")
    fmt.Printf("Tokens: %d\n", tokens)

    est, _ := tc.EstimateCost("gpt-4o", tokens, 100)
    fmt.Printf("Cost: $%.6f\n", est.TotalCost)

    fits, remaining := tc.FitsContext("gpt-4", tokens)
    fmt.Printf("Fits: %v, Remaining: %d\n", fits, remaining)
}
```

## License

MIT
