package myStrings

import "fmt"

const Letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const VOWELS = "AEIOU"
const VOWEL_COST = 250

func IsInList(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//Returns a string representing the current state of the game
func ShowBoard(category, obscuredPhrase string, guessed []string) {
	fmt.Printf("Category: %s\nPhrase: %s\nGuessed: %s\n", category, obscuredPhrase, guessed)
}
