package audio

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func Play(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch radio stream: %v\n", err)
		return
	}
	defer resp.Body.Close()

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode MP3 stream: %v\n", err)
		return
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize speaker: %v\n", err)
		return
	}

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		fmt.Println("Stream finished playing")
	})))
	select {}
}
