// Command authopus is the Authopus server entry point.
package main

import (
	"fmt"

	"github.com/tarangchikhalia/authopus/internal/app"
)

func main() {
	fmt.Println(app.Greeting())
}
