package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type simulator struct{}

type eventPayload map[string]any

type event struct {
	Event     string       `json:"event"`
	UserID    string       `json:"user_id"`
	Timestamp time.Time    `json:"timestamp"`
	Params    eventPayload `json:"params,omitempty"`
}

type choose_tier struct {
	name   string
	weight int
}

func chooseRand(choices []choose_tier) string {
	return choose(rand.Intn(100), choices)
}

func choose(weight int, choices []choose_tier) string {
	total := 0
	for _, v := range choices {
		total += v.weight
	}

	mod_weight := weight % total
	run := 0
	for _, v := range choices {
		run += v.weight
		if mod_weight <= run {
			return v.name
		}
	}

	return ""
}

var payment_tiers = []choose_tier{
	{"none", 70},
	{"premium", 30},
}

var stages = []choose_tier{
	{"", 70},
	{"premium", 30},
}

var battle_arena = []choose_tier{
	{"Last Player Standing", 10},
	{"Rumble", 10},
	{"Score Attack", 10},
	{"Team Knockout", 10},
}

var locations = []choose_tier{
	{"Keep it Low", 10},
	{"Pit 1", 10},
	{"Pit 2", 10},
	{"Pit 3", 10},
	{"Bounce House", 10},
	{"Plateau", 10},
	{"Pieces", 10},
	{"Tractor", 10},
	{"Pit 4", 10},
	{"Saw Mill", 10},
	{"Pit 6", 10},
	{"Pinch", 10},
	{"Burping Burt", 10},
	{"Fracking Frank's Place", 10},
}

type session_context struct {
	userId  string
	events  []event
	simTime float32
	bucket  int
}

func (s *session_context) addEvent(eventType string, params eventPayload) {
	s.addEventT(rand.Intn(30)+5, eventType, params)
}

func (s *session_context) addEventT(duration int, eventType string, params eventPayload) {
	s.simTime += float32(duration)
	e := event{
		Event:     eventType,
		UserID:    s.userId,
		Timestamp: time.Now().Add(time.Duration(s.simTime) * time.Second),
		Params:    params,
	}

	s.events = append(s.events, e)
}

func (s *session_context) sim_identify() {
	paymentTier := choose(s.bucket, payment_tiers)
	s.addEvent("$identify", eventPayload{
		"name":         "kneehat",
		"payment_tier": paymentTier,
	})
}

func (s *session_context) sim_settings() {
	if rand.Intn(100) < 30 {
		s.addEvent("$uiScreen", eventPayload{
			"name": "settings",
			"tab":  "game config",
		})
		s.addEvent("$uiScreen", eventPayload{
			"name": "settings",
			"tab":  "controls",
		})
		s.addEvent("$uiScreen", eventPayload{
			"name": "welcome",
		})
	}
}

func (s *session_context) sim_tutorial() {
	if rand.Intn(100) < 10 {
		s.addEvent("$tutorialBegin", nil)
		s.addEvent("$tutorialStep", eventPayload{
			"step": 1,
		})
		s.addEvent("$tutorialStep", eventPayload{
			"step": 2,
		})
		s.addEvent("$tutorialStep", eventPayload{
			"step": 3,
		})
		s.addEvent("$tutorialEnd", nil)
		s.addEvent("$uiScreen", eventPayload{
			"name": "welcome",
		})
	}
}

func (s *session_context) sim_customise(homeScreen string) {
	if rand.Intn(100) < 30 {
		// sim the character select screen
		s.addEvent("$uiScreen", eventPayload{
			"name": "customise",
		})
		s.addEvent("$uiScreen", eventPayload{
			"name": "change character",
			"tab":  "regular characters",
		})
		s.addEvent("$uiScreen", eventPayload{
			"name": homeScreen,
		})
	}
}

func (s *session_context) sim_game_battle_arena(homeScreen string) {
	s.sim_customise(homeScreen)

	_type := chooseRand(battle_arena)
	length := rand.Intn(4)*2 + 1

	s.addEvent("$uiScreen", eventPayload{
		"name": "game mode",
	})
	s.addEvent("$uiScreen", eventPayload{
		"name": homeScreen,
	})
	s.addEvent("$gameBegin", eventPayload{
		"type":   _type,
		"length": length,
	})

	numPlayers := rand.Intn(2) + 2

	for i := 0; i < length; i++ {
		location := chooseRand(locations)
		s.addEventT(3, "levelBegin", eventPayload{
			"name": location,
		})

		player := []string{}
		score := []int{}
		// death := []int{}
		// banana := []int{}

		for i := 0; i < numPlayers; i++ {
			if i == 0 {
				player = append(player, s.userId)
			} else {
				player = append(player, fmt.Sprintf("player%d", i+1))
			}

			score = append(score, rand.Intn(600)*50)

			//    numDeath := rand.IntN(5)
			// numBanana := 0
			// if numDeath > 0 {
			// 	numBanana = rand.IntN(numDeath)
			// }
			// death = append(death, numDeath)
			// banana = append(banana, numBanana)
		}

		s.addEventT(rand.Intn(30)+30, "levelEnd", eventPayload{
			"player": player,
			"score":  score,
			// "death":  death,
			// "banana": banana,
		})
	}

	s.addEventT(1, "$gameEnd", nil)
	s.addEvent("$uiScreen", eventPayload{
		"name": homeScreen,
	})
}

func (s *session_context) sim_home() {
	s.addEvent("$uiScreen", eventPayload{
		"name": "home",
	})

	switch val := rand.Intn(100); {
	case val < 50:
		numGames := rand.Intn(4) + 1
		for i := 0; i < numGames; i++ {
			s.sim_game_battle_arena("home")
		}
	default:
		s.addEvent("$uiScreen", eventPayload{
			"name": "play online",
		})
		s.addEvent("$uiScreen", eventPayload{
			"name": "create lobby",
		})
		s.addEvent("$uiScreen", eventPayload{
			"name": "al's imports",
		})
		numGames := rand.Intn(4) + 1
		for i := 0; i < numGames; i++ {
			s.sim_game_battle_arena("al's imports")
		}
		s.addEvent("$uiScreen", eventPayload{
			"name": "home",
		})
	}
}

func (s *session_context) begin() {
	s.bucket = rand.Intn(100)
	s.userId = fmt.Sprintf("STEAM#%d", s.bucket)
	s.sim_identify()
	s.addEvent("$sessionBegin", eventPayload{
		"steam_branch": "beta",
		"version":      "0.1.4937",
	})
	s.addEvent("$uiScreen", eventPayload{
		"name": "welcome",
	})
}

func (s *session_context) end() {
	s.addEvent("$sessionEnd", nil)
}

func (s *session_context) serialise() error {
	r, err := json.MarshalIndent(s.events, "", "  ")
	// r, err := json.Marshal(s.events)
	if err != nil {
		return err
	}
	fmt.Printf("%s", string(r))
	return nil
}

func simulateForProject(numSessions int) error {
	ctx := session_context{}
	for i := 0; i < numSessions; i++ {
		ctx.begin()
		ctx.sim_settings()
		ctx.sim_tutorial()
		ctx.sim_home()
		ctx.end()
	}
	return ctx.serialise()
}

func main() {
	err := simulateForProject(1)
	if err != nil {
		panic(err)
	}
}
