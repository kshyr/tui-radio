package main

import (
	"github.com/kshyr/tui-radio/internal/client"
)

func main() {
	c := client.DefaultClient()
	c.Run()
}
