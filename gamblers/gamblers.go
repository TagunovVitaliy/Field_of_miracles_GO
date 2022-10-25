package gamblers

import (
	"Field_of_miracles/myStrings"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Players interface {
	GetName() string
	GetPrizeMoney() int
	GetPrizes() []string
	AddMoney(int)
	AddPrize(string)
	GoBankrupt()
	GetMove(category, obscuredPhrase string, guessed []string) string
}

type WOFPlayer struct {
	Name       string
	PrizeMoney int
	Prizes     []string
}

func (player *WOFPlayer) GetName() string {
	return player.Name
}

func (player *WOFPlayer) GetPrizeMoney() int {
	return player.PrizeMoney
}

func (player *WOFPlayer) GetPrizes() []string {
	return player.Prizes
}

func (player *WOFPlayer) GoBankrupt() {
	player.PrizeMoney = 0
}

func (player *WOFPlayer) AddPrize(prize string) {
	player.Prizes = append(player.Prizes, prize)
}

func (player *WOFPlayer) AddMoney(amt int) {
	player.PrizeMoney += amt
}

func (player *WOFPlayer) GetPlayer() string {
	prizeM := strconv.Itoa(player.PrizeMoney)
	state := player.Name + " ($" + prizeM + ")"
	return state
}

type WOFHumanPlayer struct{ WOFPlayer }

func (b WOFHumanPlayer) GetMove(category, obscuredPhrase string, guessed []string) string {
	//var scanned string
	myStrings.ShowBoard(category, obscuredPhrase, guessed)
	fmt.Print("Guess a letter, phrase, or type 'exit' or 'pass': ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	line := scanner.Text()
	return line
}

type WOFComputerPlayer struct {
	WOFPlayer
	Difficulty int
}

func (a WOFComputerPlayer) SmartCoinFlip() bool {
	if rand.Intn(9)+1 > a.Difficulty {
		return true
	} else {
		return false
	}
}

func (a WOFComputerPlayer) GetPossibleLetters(guessed []string) []string {
	list := make([]string, 0)
	if a.WOFPlayer.PrizeMoney >= 250 {
		for _, l := range myStrings.Letters {
			if !(myStrings.IsInList(guessed, string(l))) {
				list = append(list, string(l))
			}
		}
	} else {
		for _, l := range myStrings.Letters {
			if !(myStrings.IsInList(guessed, string(l))) && !(strings.Contains(myStrings.VOWELS, string(l))) {
				list = append(list, string(l))
			}
		}
	}
	return list
}

func (a WOFComputerPlayer) GetMove(category, obscuredPhrase string, guessed []string) string {
	SORTED_FREQUENCIES := "ZQXJKVBPYGFWMUCLDRHSNIOATE"
	list := a.GetPossibleLetters(guessed)
	FlipResult := a.SmartCoinFlip()
	if len(list) == 0 {
		return "pass"
	} else {
		if FlipResult {
			for _, l := range SORTED_FREQUENCIES {
				if myStrings.IsInList(list, string(l)) {
					return string(l)
				}
			}
		} else if !FlipResult {
			randomIndex := rand.Intn(len(list))
			return list[randomIndex]
		}
	}
	return ""
}
