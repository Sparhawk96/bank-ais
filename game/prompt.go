package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Sparhawk96/bank-ais/table"
)

const (
	PROMPT      = "> "
	MAIN_PROMPT = "Enter '?' for help\n\r" + PROMPT
	BANK_PROMPT = "Enter 'd' or 'done' to submit\n\r" + PROMPT
)

type PromptRequest int

const (
	PRINT_POINTS PromptRequest = iota
	PLAYERS_BANK
	ROLL_DICE
)

/**
 * Prompts real players for some kind of action
 *
 * @return The request from the real players
 */
func prompt() PromptRequest {
	var request PromptRequest

	for keepPrompting := true; keepPrompting; {
		input := GetInput(MAIN_PROMPT, true)

		// Determine Players Action
		switch input {
		case "?", "help", "h":
			printPromptMenu()
		case "p", "points", "pp", "print points", "player points":
			keepPrompting = false
			request = PRINT_POINTS
		case "b", "bank", "pb", "player bank", "players bank", "bp", "bank points":
			keepPrompting = false
			request = PLAYERS_BANK
		case "", "r", "roll", "rd", "roll dice":
			keepPrompting = false
			request = ROLL_DICE
		default:
			fmt.Printf("Invalid Input: '%s'\n\r", input)
		}
	}

	return request
}

/**
 * Prints the Prompt Menu
 */
func printPromptMenu() {
	menu := new(table.Table)

	actHdr := "Action"
	cmdsHdr := "Commands"
	descHdr := "Description"

	menu.CreateColumn(actHdr, table.LEFT, 0)
	menu.CreateColumn(cmdsHdr, table.LEFT, 0)
	menu.CreateColumn(descHdr, table.LEFT, 0)

	menu.AddEntry(map[string]any{
		actHdr:  "Player Points",
		cmdsHdr: "[p, points, print points]",
		descHdr: "Prints the current player points",
	})

	menu.AddEntry(map[string]any{
		actHdr:  "Bank Points",
		cmdsHdr: "[b, bank, bank points]",
		descHdr: "Enables players to bank the current points",
	})

	menu.AddEntry(map[string]any{
		actHdr:  "Roll Dice",
		cmdsHdr: "[r, roll, roll dice]",
		descHdr: "Keep going and roll the dice",
	})

	fmt.Println(menu)
}

/**
 * Gets the list of all players who are banking
 *
 * @param players List of players who haven't banked
 *
 * @return List of players who are banking
 */
func getBankingPlayers(players []Player) []string {
	playersBanking := make([]string, 0)
	playerMap := make(map[string]bool)
	posPlayerMap := make(map[string]string)

	for idx, player := range players {
		// Human Players can't bank for AI Agents
		if !player.AiAgent() {
			fmt.Printf("%d) %s\n\r", idx+1, player.Name())
			playerMap[player.Name()] = true
			posPlayerMap[fmt.Sprint(idx+1)] = player.Name()
		}
	}
	fmt.Println("")

	prompt := BANK_PROMPT
	for keepPrompting := true; keepPrompting; prompt = PROMPT {
		input := GetInput(prompt, false)

		if playerMap[input] {
			playersBanking = append(playersBanking, input)
		} else if playerMap[posPlayerMap[input]] {
			playersBanking = append(playersBanking, posPlayerMap[input])
		} else if strings.ToLower(input) == "d" ||
			strings.ToLower(input) == "done" ||
			strings.ToLower(input) == "submit" {

			keepPrompting = false
		} else {
			fmt.Printf("Invalid Player Name or Number: %s\n\r", input)
		}
	}

	return playersBanking
}

var inputReader *bufio.Reader

/**
 * Gets input from reader (stdin) and trims spaces\
 *
 * @param prompt Sends a prompt to stdout requesting input from the user
 * @param lowerCase If True the input is lower cased
 *
 * @return The requested input
 */
func GetInput(prompt string, lowerCase bool) string {
	fmt.Print(prompt)

	// TODO: how can i suppress all non-visible characters from being seen
	if inputReader == nil {
		inputReader = bufio.NewReader(os.Stdin)
	}

	input, _ := inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	if lowerCase {
		input = strings.ToLower(input)
	}
	return input
}
