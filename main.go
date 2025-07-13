package main

import (
	"fmt"

	"github.com/Sparhawk96/bank-ais/game"
)

func main() {
	bankGame := game.NewGame()

	fmt.Println("Enter 'd' or 'done' to stop adding players.")

	for keepPrompting := true; keepPrompting; {
		input := game.GetInput("Enter Player Name "+game.PROMPT, false)
		if input == "d" || input == "D" || input == "done" || input == "Done" || input == "DONE" {
			keepPrompting = false
		} else if len(input) == 0 {
			// NO-OP
		} else if bankGame.AddPlayer(game.NewHumanPlayer(input)) != nil {
			fmt.Println("Player already exists with that name.")
		}
	}

	for keepPrompting := true; keepPrompting; {
		input := game.GetInput("Add AI Agents (y/n) "+game.PROMPT, true)
		switch input {
		case "no", "n":
			keepPrompting = false
		case "yes", "y":
			keepPrompting = false
			addAiAgents(bankGame)
		}
	}

	fmt.Println()
	bankGame.StartGame()
}

func addAiAgents(game *game.Game) {
	fmt.Println("Adding in AI Agents ...")

	// TODO
}
