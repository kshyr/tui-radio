package main

import (
	"github.com/kshyr/tui-radio/internal/tui"
)

const (
	ExampleStationURL = "https://icecast.walmradio.com:8443/classic"
)

func main() {
	c := tui.DefaultClient()
	c.Run()
}
