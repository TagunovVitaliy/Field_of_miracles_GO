package main

import (
	"Field_of_miracles/funcs"
	"Field_of_miracles/gamblers"
	"Field_of_miracles/myStrings"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//GAME LOGIC CODE
func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("===============")
	fmt.Println("FIELD OF MIRACLES")
	fmt.Println("===============")

	//Create the all player instances
	num_human := funcs.GetNumberBetween("How many human players?", 0, 10)
	players := make([]gamblers.Players, 0)

	//Create the human player instances
	var newName string
	for i := 0; i < num_human; i++ {
		fmt.Print("Enter the name of the new human player: ")
		fmt.Scanln(&newName)
		players = append(players, &gamblers.WOFHumanPlayer{gamblers.WOFPlayer{Name: newName}})
	}

	num_computer := funcs.GetNumberBetween("How many computer players?", 1, 10)
	//If there are computer players, ask how difficult they should be
	difficulty := funcs.GetNumberBetween("What difficulty for the computers? (1-10)", 1, 10)

	//Create the computer player instances
	for i := 0; i < num_computer; i++ {
		newName = fmt.Sprintf("Computer %d", i+1)
		players = append(players, &gamblers.WOFComputerPlayer{gamblers.WOFPlayer{Name: newName}, difficulty})
	}

	//No players, no game :(
	if len(players) == 0 {
		fmt.Println("We need players to play!")
		panic("Not enough players")
	}

	// category and phrase are strings.
	category_phrase := funcs.GetRandomCategoryAndPhrase()
	category := category_phrase[0]
	phrase := category_phrase[1]
	//guessed is a list of the letters that have been guessed
	guessed := make([]string, 0)
	//playerIndex keeps track of the index (0 to len(players)-1) of the player whose turn it is
	playerIndex := 0
	//will be set to the player instance when/if someone wins
	var winner gamblers.Players

	for {
		player := players[playerIndex]
		wheelPrize := funcs.SpinWheel()

		fmt.Println("\n---------------")
		myStrings.ShowBoard(category, funcs.ObscurePhrase(phrase, guessed), guessed)
		fmt.Printf("%s spins...", player.GetName())
		time.Sleep(2) //pause for dramatic effect!
		fmt.Printf("%s! ", wheelPrize["text"])
		time.Sleep(1) //pause again for more dramatic effect!

		if wheelPrize["type"] == "bankrupt" {
			player.GoBankrupt()
		} else if wheelPrize["type"] == "loseturn" {
			//do nothing, just move on to the next player
		} else if wheelPrize["type"] == "cash" {
			move := funcs.RequestPlayerMove(player, category_phrase, guessed)
			if move == "EXIT" { //leave the game
				fmt.Println("Until next time!")
				break
			} else if move == "PASS" { //will just move on to next player
				fmt.Printf("%s passes\n", player.GetName())
			} else if len(move) == 1 { //they guessed a letter
				guessed = append(guessed, move)
				fmt.Printf("%s guesses '%s'\n", player.GetName(), move)
				if strings.Contains(myStrings.VOWELS, move) {
					player.AddMoney(-myStrings.VOWEL_COST)
				}
				count := strings.Count(phrase, move) //returns an integer with how many times this letter appears
				if count > 0 {
					if count == 1 {
						fmt.Printf("There is one %s\n", move)
					} else {
						fmt.Printf("There are %d %s's\n", count, move)
					}
					//Give them the money and the prizes
					t1 := int(wheelPrize["value"].(float64))
					player.AddMoney(count * t1)
					t2 := wheelPrize["prize"]
					if _, ok := t2.(string); ok {
						player.AddPrize(wheelPrize["prize"].(string))
					}
					//All the letters have been guessed
					if funcs.ObscurePhrase(phrase, guessed) == phrase {
						winner = player
						break
					}
					continue //this player gets to go again
				} else if count == 0 {
					fmt.Printf("There is no %s\n", move)
				}
			} else { //they guessed the whole phrase
				if move == phrase {
					winner = player
					//Give them the money and the prizes
					t := int(wheelPrize["value"].(float64))
					player.AddMoney(t)
					t2 := wheelPrize["prize"]
					if _, ok := t2.(string); ok {
						player.AddPrize(wheelPrize["prize"].(string))
					}
					break
				} else {
					fmt.Printf("%s was not the phrase", move)
				}
			}
		}
		//Move on to the next player (or go back to player[0] if we reached the end)
		playerIndex = (playerIndex + 1) % len(players)
	}

	if winner.GetName() != "" {
		//In your head, you should hear this as being announced by a game show host
		fmt.Printf("%s wins! The phrase was %s\n", winner.GetName(), phrase)
		fmt.Printf("%s won $%d\n", winner.GetName(), winner.GetPrizeMoney())
		if len(winner.GetPrizes()) > 0 {
			fmt.Printf("%s also won:\n", winner.GetName())
			for _, prize := range winner.GetPrizes() {
				fmt.Printf("    - %s\n", prize)
			}
		}
	} else {
		fmt.Printf("Nobody won. The phrase was %s\n", phrase)
	}
}
