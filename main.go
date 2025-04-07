package main

import (
	"context"

	"github.com/axone-protocol/axone-mcp/cmd"
)

func main() {
	ctx := context.Background()
	cmd.Execute(ctx)
}
