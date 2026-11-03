// Command authopusctl is the Authopus administrative CLI entry point.
package main

import (
	"fmt"

	"github.com/tarangchikhalia/authopus/internal/app"
)

func main() {
	fmt.Println(app.Greeting())
}
