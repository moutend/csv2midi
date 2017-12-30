package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
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

func parseFile(file io.Reader) ([]event.Event, error) {
	var events []event.Event
	var line int
	var fields []string
	var e event.Event
	var err error

	reader := csv.NewReader(file)
	for {
		fields, err = reader.Read()
		if err != nil {
			break
		}
		e, err = parseFields(fields)
		if err != nil {
			break
		}
		events = append(events, e)
		line++
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		err = fmt.Errorf("line %v: %v", line+1, err)
	}

	return events, err
}

func parseFields(fields []string) (event.Event, error) {
	if len(fields) < 1 {
		return nil, fmt.Errorf("specify delta time at first column")
	}
	i, err := strconv.Atoi(fields[0])
	if err != nil {
		err = fmt.Errorf("delta time %v is invalid (%v)", fields[0], err)
		return nil, err
	}
	deltaTime, err := deltatime.New(i)
	if err != nil {
		return nil, err
	}
	if len(fields) == 1 {
		err = fmt.Errorf("specify event type at second column")
		return nil, err
	}
	eventName := strings.ToLower(fields[1])
	switch eventName {
	case "on", "note on", "note-on", "note_on", "off", "note off", "note-off", "note_off":
		if len(fields) == 2 {
			return nil, fmt.Errorf("specify note name at third column")
		}
		note, err := constant.ParseNote(fields[2])
		if err != nil {
			return nil, err
		}
		if len(fields) == 3 {
			return nil, fmt.Errorf("specify velocity at fourth column")
		}
		velocity, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(eventName, "on") {
			return event.NewNoteOnEvent(deltaTime, 0, note, uint8(velocity))
		} else {
			return event.NewNoteOffEvent(deltaTime, 0, note, uint8(velocity))
		}
	case "cc", "control change", "control-change", "control_change":
		if len(fields) == 2 {
			return nil, fmt.Errorf("specify controller name at third column")
		}
		cc, err := constant.ParseControlName(fields[2])
		if err != nil {
			return nil, err
		}
		if len(fields) == 3 {
			return nil, fmt.Errorf("specify value at fourth column")
		}
		value, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, err
		}
		return event.NewControllerEvent(deltaTime, 0, cc, uint8(value))
	case "bend", "pitch bend", "pitch-bend", "pitch_bend":
		if len(fields) == 2 {
			return nil, fmt.Errorf("specify")
		}
		pb, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}
		return event.NewPitchBendEvent(deltaTime, 0, uint16(pb))
	}
	return nil, fmt.Errorf("unknown event type '%v'", eventName)
}

func main() {
	if err := run(os.Args); err != nil {
		log.New(os.Stderr, "error: ", 0).Fatal(err)
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
	file, err := os.Open(inputFilename)
	if err != nil {
		return err
	}

	events, err := parseFile(file)
	if err != nil {
		return err
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
