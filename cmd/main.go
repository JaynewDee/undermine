package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/tawesoft/golib/v2/dialog"
)

const alarmPath string = "assets/alarm.wav"

//go:embed assets/alarm.wav
var audioFile embed.FS

func main() {
	var (
		note     string
		duration int
	)

	// Get input args --note <text>, --duration <minutes>
	n, d := parseInput(note, duration)

	// Thread is blocked for <duration> minutes
	sleepTick(n, d)

	// Display dialog and play sound
	alert(n)
}

// Main event loop
// clear -> print -> sleep
func sleepTick(note string, duration int) {
	clearScreen()

	start := time.Now()

	dur := time.Duration(time.Minute * time.Duration(duration))

	for {
		remaining := (dur - time.Since(start)).Minutes()

		printScreen(note, int(remaining+1))
		time.Sleep(time.Second)
		clearScreen()

		if time.Since(start) > dur {
			break
		}
	}
}

func parseInput(note string, duration int) (string, int) {
	flag.StringVar(&note, "note", "DO THE THING", "Reminder text content")
	flag.IntVar(&duration, "duration", 60, "Duration of sleep timer")

	flag.Parse()

	return note, duration
}

func useAlarm() {
	// read embedded wav data
	data, err := audioFile.ReadFile(alarmPath)
	if err != nil {
		fmt.Println("error encountred while reading audio file")
		log.Fatal(err)
	}

	reader := bytes.NewReader(data)

	// Decode the audio file
	streamer, format, err := wav.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// Initialize the speaker
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	// Play the audio
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the audio to finish playing
	<-done
}

func displayTitle() {
	titleStyle := color.New(color.FgCyan).Add(color.Underline)
	titleStyle.Println("REMINDER")
}

func displayNote(note string) {
	style := color.New(color.FgHiWhite, color.Italic)
	style.Printf("%v\n\n", note)
}

func displayDuration(dur int) {
	s := color.New(color.FgGreen, color.Bold)
	s.Printf("TIME UNTIL ::: ")
	sDur := color.New(color.FgHiRed, color.Underline)
	sDur.Printf("%dm\n", dur)
}

func printScreen(n string, d int) {
	displayTitle()
	displayNote(n)
	displayDuration(d)
}

func clearScreen() {
	screen.Clear()
	screen.MoveTopLeft()
}

func alert(n string) {
	useAlarm()
	dialog.Alert(n)
}
