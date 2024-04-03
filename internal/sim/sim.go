package sim

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/model"
)

// type Event struct {
// 	Params    model.Params `json:"params,omitempty"`
// 	Event     string       `json:"event"`
// 	UserId    string       `json:"user_id"`
// 	Timestamp int64        `json:"timestamp"`
// }

type choose_tier struct {
	data   any
	weight int
}

func chooseRand(choices []choose_tier) any {
	return choose(rand.IntN(100), choices)
}

func choose(weight int, choices []choose_tier) any {
	total := 0
	for _, v := range choices {
		total += v.weight
	}

	mod_weight := weight % total
	run := 0
	for _, v := range choices {
		run += v.weight
		if mod_weight <= run {
			return v.data
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
	{"Hero Mode", 10},
}

type challenge struct {
	guid         string
	series       string
	chapter      string
	challenge    string
	winCondition string
	scoreLow     int
	scoreHigh    int
}

var challenges = []choose_tier{
	{challenge{"1", "The Problem With AI", "It Started Here", "Hey! Stop That", "descending", 0, 10}, 10},
	{challenge{"2", "The Problem With AI", "It Started Here", "Hold On Clem", "ascending", 20, 200}, 10},
	{challenge{"3", "The Problem With AI", "It Started Here", "Big Mac's Mill", "ascending", 25, 200}, 10},
	{challenge{"11", "The Problem With AI", "Robots Have Rights", "Returning Home", "ascending", 30, 200}, 10},
	{challenge{"12", "The Problem With AI", "Robots Have Rights", "Under Warranty", "ascending", 30, 200}, 10},
	{challenge{"13", "The Problem With AI", "Robots Have Rights", "How Many?", "ascending", 30, 200}, 10},
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
	{"Pirate Docks", 10},
	{"Pirate Ship", 10},
	{"Top 1", 10},
	{"Top 2", 10},
	{"Top 3", 10},
	{"Island", 10},
	{"At Sea", 10},
	{"Local_Store", 10},
	{"Barn", 10},
	{"Radio Array", 10},
}

type session_context struct {
	startTime time.Time
	userId    string
	events    []model.Event
	simTime   float32
	bucket    int
}

func (s *session_context) addEvent(eventType string, params model.Params) {
	s.addEventT(5+rand.Float32()*30.0, eventType, params)
}

func (s *session_context) addEventT(duration float32, eventType string, params model.Params) {
	if eventType == "$uiScreen" {
		return
	}
	s.simTime += duration
	e := model.Event{
		Id:        int64(len(s.events) + 1),
		Event:     eventType,
		UserId:    s.userId,
		Timestamp: s.startTime.Add(time.Duration(s.simTime) * time.Second).UnixMilli(),
		Params:    params,
	}

	s.events = append(s.events, e)
}

func (s *session_context) sim_identify() {
	paymentTier := choose(s.bucket, payment_tiers)
	s.addEvent("$identify", model.Params{
		"name":         "kneehat",
		"payment_tier": paymentTier,
	})
}

func (s *session_context) sim_settings() {
	if rand.IntN(100) < 30 {
		s.addEvent("$uiScreen", model.Params{
			"name": "settings",
			"tab":  "game config",
		})
		s.addEvent("$uiScreen", model.Params{
			"name": "settings",
			"tab":  "controls",
		})
		s.addEvent("$uiScreen", model.Params{
			"name": "welcome",
		})
	}
}

func (s *session_context) sim_tutorial() {
	if rand.IntN(100) < 10 {
		s.addEvent("$tutorialBegin", nil)
		s.addEvent("$tutorialStep", model.Params{
			"step": 1,
		})
		s.addEvent("$tutorialStep", model.Params{
			"step": 2,
		})
		s.addEvent("$tutorialStep", model.Params{
			"step": 3,
		})
		s.addEvent("$tutorialEnd", nil)
		s.addEvent("$uiScreen", model.Params{
			"name": "welcome",
		})
	}
}

func (s *session_context) sim_customise(homeScreen string) {
	if rand.IntN(100) < 30 {
		// sim the character select screen
		s.addEvent("$uiScreen", model.Params{
			"name": "customise",
		})
		s.addEvent("$uiScreen", model.Params{
			"name": "change character",
			"tab":  "regular characters",
		})
		s.addEvent("$uiScreen", model.Params{
			"name": homeScreen,
		})
	}
}

func (s *session_context) sim_game_battle_arena(homeScreen string) {
	s.sim_customise(homeScreen)

	_type := chooseRand(battle_arena)
	length := rand.IntN(4)*2 + 1

	s.addEvent("$uiScreen", model.Params{
		"name": "game mode",
	})
	s.addEvent("$uiScreen", model.Params{
		"name": homeScreen,
	})
	s.addEvent("$gameBegin", model.Params{
		"type":   _type,
		"length": length,
	})

	numPlayers := rand.IntN(2) + 2

	for i := 0; i < length; i++ {
		location := chooseRand(locations)
		s.addEventT(3, "levelBegin", model.Params{
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

			score = append(score, rand.IntN(600)*50)

			//    numDeath := rand.IntN(5)
			// numBanana := 0
			// if numDeath > 0 {
			// 	numBanana = rand.IntN(numDeath)
			// }
			// death = append(death, numDeath)
			// banana = append(banana, numBanana)
		}

		s.addEventT(30+rand.Float32()*30, "levelEnd", model.Params{
			"player": player,
			"score":  score,
			// "death":  death,
			// "banana": banana,
		})
	}

	s.addEventT(1, "$gameEnd", nil)
	s.addEvent("$uiScreen", model.Params{
		"name": homeScreen,
	})
}

func (s *session_context) sim_challenge(homeScreen string) {
	chal := chooseRand(challenges).(challenge)

	score := rand.IntN(chal.scoreHigh-chal.scoreLow) + chal.scoreLow

	s.addEventT(1, "challengeBegin", model.Params{})

	s.addEventT(30, "challengeEnd", model.Params{
		"guid":  chal.guid,
		"name":  chal.challenge,
		"sort":  chal.winCondition,
		"score": score,
	})

	s.addEvent("$uiScreen", model.Params{
		"name": homeScreen,
	})

	// s.addEvent("challengeBegin", eventPayload{
	// 	"type":   _type,
	// })

	// for i := 0; i < length; i++ {
	// 	location := chooseRand(locations)
	// 	s.addEventT(3, "levelBegin", eventPayload{
	// 		"name": location,
	// 	})
	//
	// 	player := []string{}
	// 	score := []int{}
	//
	// 	for i := 0; i < numPlayers; i++ {
	// 		if i == 0 {
	// 			player = append(player, s.userId)
	// 		} else {
	// 			player = append(player, fmt.Sprintf("player%d", i+1))
	// 		}
	//
	// 		score = append(score, rand.Intn(600)*50)
	// 	}
	//
	// 	s.addEventT(30+rand.Float32()*30, "levelEnd", eventPayload{
	// 		"player": player,
	// 		"score":  score,
	// 		// "death":  death,
	// 		// "banana": banana,
	// 	})
	// }
}

func (s *session_context) sim_home() {
	s.addEvent("$uiScreen", model.Params{
		"name": "home",
	})

	switch val := rand.IntN(100); {
	case val < 100:
		numChallenges := rand.IntN(8) + 1
		for i := 0; i < numChallenges; i++ {
			s.sim_challenge("home")
		}
		// case val < 50:
		// numGames := rand.Intn(4) + 1
		// for i := 0; i < numGames; i++ {
		// 	s.sim_game_battle_arena("home")
		// }

	default:
		s.addEvent("$uiScreen", model.Params{
			"name": "play online",
		})
		s.addEvent("$uiScreen", model.Params{
			"name": "create lobby",
		})
		s.addEvent("$uiScreen", model.Params{
			"name": "al's imports",
		})
		numGames := rand.IntN(4) + 1
		for i := 0; i < numGames; i++ {
			s.sim_game_battle_arena("al's imports")
		}
		s.addEvent("$uiScreen", model.Params{
			"name": "home",
		})
	}
}

func newSession(userId int, startTime time.Time) session_context {
	s := session_context{}
	s.bucket = userId
	s.startTime = startTime
	s.userId = fmt.Sprintf("STEAM#%d", s.bucket)
	s.simTime = 0
	return s
}

func (s *session_context) begin() {
	s.sim_identify()
	s.addEvent("$sessionBegin", model.Params{
		"branch":  "developer",
		"vendor":  "steam",
		"version": "0.1.6669",
	})
	s.addEvent("$uiScreen", model.Params{
		"name": "welcome",
	})
}

func (s *session_context) end() {
	s.addEvent("$sessionEnd", nil)
}

// func (s *session_context) serialise() error {
// 	r, err := json.MarshalIndent(s.events, "", "  ")
// 	// r, err := json.Marshal(s.events)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("%s", string(r))
// 	return nil
// }

func Serialise(events []model.Event) (string, error) {
	r, err := json.MarshalIndent(events, "", "  ")
	// r, err := json.Marshal(s.events)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func SimulateSessionForUser(userId int, startTime time.Time) []model.Event {
	ctx := newSession(userId, startTime)
	ctx.begin()
	ctx.sim_settings()
	// ctx.sim_tutorial()
	ctx.sim_home()
	ctx.end()
	return ctx.events
}
