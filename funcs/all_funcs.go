package funcs

import (
	"Field_of_miracles/gamblers"
	"Field_of_miracles/myStrings"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetNumberBetween(prompt string, min, max int) int {
	var s string
	for {
		fmt.Println(prompt)
		fmt.Scanln(&s)
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("%s is not a number\n", s)
			continue
		} else {
			if n < min {
				fmt.Printf("Must be at least %d\n", min)
			} else if n > max {
				fmt.Printf("Must be at most %d\n", max)
			} else {
				return n
			}
		}
	}
}

/*
 Spins the wheel of fortune wheel to give a random prize
 Examples:
    { "type": "cash", "text": "$950", "value": 950, "prize": "A trip to Ann Arbor!" },
    { "type": "bankrupt", "text": "Bankrupt", "prize": false },
    { "type": "loseturn", "text": "Lose a turn", "prize": false }
*/

func SpinWheel() map[string]interface{} {
	jsonStr, err := os.ReadFile("funcs/wheel1.json")
	if err != nil {
		fmt.Print(err)
	}
	wheel := make([]map[string]interface{}, 0)
	json.Unmarshal(jsonStr, &wheel)
	randomIndex := rand.Intn(len(wheel))
	return wheel[randomIndex]
}

/* Returns a category & phrase (as a tuple) to guess
Example:
    ("Artist & Song", "Whitney Houston's I Will Always Love You")*/
func GetRandomCategoryAndPhrase() []string {
	jsonStr, err := os.ReadFile("funcs/phrases.json")
	if err != nil {
		fmt.Print(err)
	}
	phrases := make(map[string][]string, 0)
	cat_and_phrase := make([]string, 0)
	count := 0
	json.Unmarshal(jsonStr, &phrases)
	keys := make([]string, len(phrases))
	randomIndex1 := rand.Intn(len(phrases))
	for k, _ := range phrases {
		keys[count] = k
		count += 1
	}
	category := keys[randomIndex1]
	randomIndex2 := rand.Intn(len(phrases[category]))
	phrase := strings.ToUpper(phrases[category][randomIndex2])
	cat_and_phrase = append(cat_and_phrase, category, phrase)
	return cat_and_phrase
}

/* Given a phrase and a list of guessed letters, returns an obscured version
Example:
    guessed: ['L', 'B', 'E', 'R', 'N', 'P', 'K', 'X', 'Z']
    phrase:  "GLACIER NATIONAL PARK"
    returns> "_L___ER N____N_L P_RK"
*/
func ObscurePhrase(phrase string, guessed []string) string {
	rv := ""
	for _, s := range phrase {
		let := strings.ToUpper(string(s))
		if strings.Contains(myStrings.Letters, let) && !myStrings.IsInList(guessed, let) {
			rv = rv + "_"
		} else {
			rv = rv + string(s)
		}
	}
	return rv
}

func RequestPlayerMove(player gamblers.Players, category_phrase, guessed []string) string {
	for { //we're going to keep asking the player for a move until they give a valid one
		time.Sleep(1) //added so that any feedback is printed out before the next prompt
		move := player.GetMove(category_phrase[0], ObscurePhrase(category_phrase[1], guessed), guessed)
		move = strings.ToUpper(move) //convert whatever the player entered to UPPERCASE

		if move == "EXIT" || move == "PASS" {
			return move
		} else if len(move) == 1 { //they guessed a character
			if !strings.Contains(myStrings.Letters, move) { //the user entered an invalid letter (such as @, #, or $)
				fmt.Println("Guesses should be letters. Try again.")
				continue
			} else if myStrings.IsInList(guessed, move) { //this letter has already been guessed
				fmt.Printf("%s has already been guessed. Try again.", move)
				continue
			} else if strings.Contains(myStrings.VOWELS, move) && player.GetPrizeMoney() < myStrings.VOWEL_COST { //if it's a vowel, we need to be sure the player has enough
				fmt.Printf("Need $%d to guess a vowel. Try again.", myStrings.VOWEL_COST)
				continue
			} else {
				return move
			}
		} else { //they guessed the phrase
			return move
		}
	}
}
