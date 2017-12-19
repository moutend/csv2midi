package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	midi "github.com/moutend/go-midi"
	"github.com/moutend/go-midi/constant"
	"github.com/moutend/go-midi/deltatime"
	"github.com/moutend/go-midi/event"
)

var (
	DeltaTimeRandomizer *randomizer
	VelocityRandomizer  *randomizer
)

type randomizer struct {
	factor   int
	position int
}

func (r *randomizer) Randomize(input int) int {
	if input <= r.factor || r.factor <= 0 {
		return input
	}

	offset := 0
	for {
		offset = rand.Intn(r.factor)
		if (offset + input) < 0 {
			continue
		} else {
			break
		}
	}
	if r.position <= -r.factor || r.factor <= r.position {
		offset = -offset
	}

	r.position += offset

	return input + offset
}

func parseLine(line string) (event.Event, error) {
	s := strings.Split(line, ",")
	if len(s) != 4 {
		return nil, fmt.Errorf("line must be contain: delta time, type, note and velocity: %v", line)
	}
	dt, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, err
	}
	deltaTime, err := deltatime.New(DeltaTimeRandomizer.Randomize(dt))
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
	var err error

	if err = run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var randomizeDeltaTimeFlag int
	var randomizeVelocityFlag int

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.IntVar(&randomizeDeltaTimeFlag, "delta-time", 0, "randomize delta time with given factor")
	f.IntVar(&randomizeDeltaTimeFlag, "d", 0, "alias of --delta-time")
	f.IntVar(&randomizeVelocityFlag, "velocity", 0, "randomize velocity with given factor")
	f.IntVar(&randomizeVelocityFlag, "v", 0, "alias of --velocity")
	f.Parse(args[1:])
	if randomizeVelocityFlag < 0 || randomizeDeltaTimeFlag < 0 {
		return fmt.Errorf("randomize factor must be positive integer")
	}
	DeltaTimeRandomizer = &randomizer{factor: randomizeDeltaTimeFlag}
	VelocityRandomizer = &randomizer{factor: randomizeVelocityFlag}

	if len(f.Args()) < 1 {
		return nil
	}

	inputFilename := f.Args()[0]
	file, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		return err
	}

	events := []event.Event{}
	lines := strings.Split(string(file), "\n")
	lines = lines[0 : len(lines)-1]
	for _, line := range lines {
		e, err := parseLine(line)
		if err != nil {
			return err
		}
		events = append(events, e)
	}

	var defaultTimeDivision int = 960
	remainingDeltaTime, err := deltatime.New(defaultTimeDivision - DeltaTimeRandomizer.position)
	if err != nil {
		return err
	}
	endOfTrack, _ := event.NewEndOfTrackEvent(remainingDeltaTime)
	events = append(events, endOfTrack)

	track := midi.NewTrack(events...)
	m := midi.MIDI{}
	m.TimeDivision().SetBPM(defaultTimeDivision)
	m.Tracks = append(m.Tracks, track)

	outputFilename := inputFilename
	if strings.HasSuffix(inputFilename, ".csv") {
		outputFilename = inputFilename[0:len(inputFilename)-4] + ".mid"
	} else {
		outputFilename += ".mid"
	}
	return ioutil.WriteFile(outputFilename, m.Serialize(), 0644)
}
