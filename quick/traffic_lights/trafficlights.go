package main

import (
	"time"
	"fmt"
)

type Direction string
type Signal string

type Event struct {
	direction Direction
	state     Signal
}

const (
	STOP                     = Direction("STOP")
	NORTH, SOUTH, EAST, WEST = Direction("N"), Direction("S"), Direction("E"), Direction("W")

	SECOND_10, SECOND_5, SECOND_0 = 10 * time.Second, 5 * time.Second, 0 * time.Second
	MINUTE_1                      = 1 * time.Minute

	RED, YELLOW, GREEN = Signal("RED"), Signal("YELLOW"), Signal("GREEN")
)

func warden(msgbus chan *Event) {
	northBus, southBus, eastBus, westBus := instantiateProcesses(msgbus)
	for {
		got_event := <-msgbus
		switch got_event.direction {
		case NORTH:
			northBus <- got_event
		case SOUTH:
			southBus <- got_event
		case EAST:
			eastBus <- got_event
		case WEST:
			westBus <- got_event
		default:
			//close(msgbus) race condition - why?
			break
		}
	}
}
func instantiateProcesses(msgBus chan *Event) (northBus, southBus, eastBus, westBus chan *Event) {
	northBus = make(chan *Event)
	southBus = make(chan *Event)
	eastBus = make(chan *Event)
	westBus = make(chan *Event)
	go genericBus(NORTH, northBus, msgBus)
	go genericBus(SOUTH, southBus, msgBus)
	go genericBus(EAST, eastBus, msgBus)
	go genericBus(WEST, westBus, msgBus)
	return
}
func genericBus(owner Direction, inbus <-chan *Event, outbus chan<- *Event) {
	for {
		owner_event := <-inbus
		oldState := *owner_event
		var duration time.Duration
		switch oldState.state {
		case RED:
			owner_event.state = GREEN
			duration = SECOND_10
		case GREEN:
			owner_event.state = YELLOW
			duration = SECOND_5
		case YELLOW:
			owner_event.state = RED
			duration = SECOND_0
			switch oldState.direction {
			case NORTH:
				owner_event.direction = EAST
			case SOUTH:
				owner_event.direction = WEST
			case EAST:
				owner_event.direction = NORTH
			case WEST:
				owner_event.direction = SOUTH
			}
		}
		t := time.Now()
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d   Changing %v from %v to %v\n",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(),
			owner, oldState.state, owner_event.state)
		time.Sleep(duration)
		outbus <- owner_event
	}
}

func startSimulation(msgbus chan *Event) {
	msgbus <- &Event{NORTH, RED}
	msgbus <- &Event{SOUTH, RED}
	defer func() {
		msgbus <- &Event{direction: STOP}
		msgbus <- &Event{direction: STOP}
	}()
	time.Sleep(MINUTE_1)
}

func main() {
	msgbus := make(chan *Event, 2)
	go warden(msgbus)
	startSimulation(msgbus)
}
