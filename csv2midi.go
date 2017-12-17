package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	midi "github.com/moutend/go-midi"
	"github.com/moutend/go-midi/constant"
	"github.com/moutend/go-midi/deltatime"
	"github.com/moutend/go-midi/event"
)

func makeRandomOffsets(width, length int) []int {
	position := 0
	offsets := []int{}
	rand.Seed(time.Now().Unix())

	for i := 0; i < length-1; i++ {
		r := rand.Intn(width)
		p := position + r
		if p < -width || width < p {
			r = -r
		}
		offsets = append(offsets, r)
		position += r
	}
	if position > 0 {
		offsets = append(offsets, -position)
	} else {
		offsets = append(offsets, position)
	}

	return offsets
}

func newDeltaTime(s string) (*deltatime.DeltaTime, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return deltatime.New(i)
}

func parseLine(line string) (event.Event, error) {
	s := strings.Split(line, ",")
	if len(s) != 4 {
		return nil, fmt.Errorf("line must be contain: delta time, type, note and velocity: %v", line)
	}
	deltaTime, err := newDeltaTime(s[0])
	if err != nil {
		return nil, err
	}
	note, err := constant.ParseNote(s[2])
	if err != nil {
		return nil, err
	}
	velocity, err := strconv.Atoi(s[3])
	if err != nil {
		return nil, err
	}
	eventName := strings.ToLower(s[1])
	switch eventName {
	case "on":
		return event.NewNoteOnEvent(deltaTime, 0, note, uint8(velocity))
	case "off":
		return event.NewNoteOffEvent(deltaTime, 0, note, uint8(velocity))
	default:
		return nil, fmt.Errorf("invalid line '%v'", line)
	}
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var events []event.Event
	position := 0
	lines := strings.Split(string(file), "\n")
	lines = lines[0 : len(lines)-1]
	for _, line := range lines {
		e, err := parseLine(line)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, e)
		position += int(e.DeltaTime().Quantity().Uint32())
	}
	var defaultTimeDivision int = 960
	remainingDeltaTime, err := deltatime.New(position % defaultTimeDivision)
	if err != nil {
		log.Fatal(err)
	}
	endOfTrack, _ := event.NewEndOfTrackEvent(remainingDeltaTime)
	events = append(events, endOfTrack)

	track := midi.NewTrack(events...)
	m := midi.MIDI{}
	m.TimeDivision().SetBPM(defaultTimeDivision)
	m.Tracks = append(m.Tracks, track)
	fmt.Println(track)
	if err := ioutil.WriteFile("output.mid", m.Serialize(), 0644); err != nil {
		log.Fatal(err)
	}
	return
}
