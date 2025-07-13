package game

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/Sparhawk96/bank-ais/table"
)

const MAX_ROUNDS = 20

const GAME_HAS_STARTED_ERR_MSG = "game has already started"

type Game struct {
	started      bool
	currentRound uint8 // 1 - 20
	rounds       [20]round
	players      map[string]Player
	results      *results
	seed         int64
	r            *rand.Rand

	// True if all of the players are AI Agents,
	// otherwise at least one human is playing
	onlyAI bool
}

type round struct {
	points uint
	rolls  []Dice
}

/**
 * Creates a new game of Bank
 *
 * @param players List of all Players who are going to be part of the game
 *
 * @return The new game of Bank
 */
func NewGame() *Game {
	seed := int64(rand.Uint64()) // Allow +/- numbers
	return &Game{
		currentRound: 0,
		players:      make(map[string]Player, 0),
		results:      new(results),
		seed:         seed,
		r:            rand.New(rand.NewSource(seed)),
		onlyAI:       true,
	}
}

/**
 * Sets the seed for the game.
 *
 * @note This allows the game to be ran the exact same way with dice rolls
 */
func (g *Game) SetSeed(seed int64) error {
	if g.started {
		return errors.New(GAME_HAS_STARTED_ERR_MSG)
	}

	g.seed = seed
	g.r = rand.New(rand.NewSource(seed))

	return nil
}

/**
 * Adds a Player to the game
 *
 * @param player Player to add to the game
 *
 * @return An error if the player is nil, another player with the same name exists, or the game has started
 */
func (g *Game) AddPlayer(player Player) error {
	if g.started {
		return errors.New(GAME_HAS_STARTED_ERR_MSG)
	} else if player == nil {
		return fmt.Errorf("added nil player")
	} else if _, exists := g.players[player.Name()]; exists {
		return fmt.Errorf("player already exists with that name: '%s'", player.Name())
	}

	g.players[player.Name()] = player
	g.results.addPlayer(player)
	g.onlyAI = g.onlyAI && player.AiAgent()
	return nil
}

/**
 * Starts the game of Bank
 *
 * @return An error if the game is already started
 */
func (g *Game) StartGame() error {
	if g.started {
		return errors.New(GAME_HAS_STARTED_ERR_MSG)
	}
	g.started = true

	// Notify the players who is playing
	playerHdr := "Player"
	aiAgentHdr := "AI Agent"

	players := new(table.Table)
	players.CreateColumn(playerHdr, table.LEFT, 0)
	players.CreateColumn(aiAgentHdr, table.CENTER, '\u2716') // ✖

	for _, player := range g.players {
		data := map[string]any{playerHdr: player.Name()}
		if player.AiAgent() {
			data[aiAgentHdr] = '\u2714' // ✔
		}
		players.AddEntry(data)
	}

	fmt.Println("Starting Game ...")
	fmt.Println(players)

	// Start the game
	for ; g.currentRound < MAX_ROUNDS; g.currentRound++ {
		g.startRound()
		g.results.unbankAllPlayers()
		fmt.Println(g.results)
	}

	// TODO: Ties?
	fmt.Printf("Player '%s' won!\n\r", g.results.firstPlayer.Name())
	return nil
}

/**
 * Starts a round of Bank
 */
func (g *Game) startRound() {
	round := g.rounds[g.currentRound]

	fmt.Printf("\n\r### Starting Round %d of 20 ###\n\r", g.currentRound+1)

	dice, keepRolling := g.roll(&round)
	bankedRound := false
	for keepRolling {
		fmt.Println()
		fmt.Println(dice)
		fmt.Printf("Current Points: %d\n\r", round.points)
		fmt.Printf("Roll Number: %d\n\r\n\r", len(round.rolls))

		g.askAiAgentsToBank()
		for keepPrompting := true; keepPrompting; {
			switch prompt() {
			case PRINT_POINTS:
				fmt.Println(g.results)
			case PLAYERS_BANK:
				bankingPlayers := getBankingPlayers(g.results.getUnbankedPlayers())
				for _, player := range bankingPlayers {
					keepPrompting = !g.results.playerBanks(g.players[player], round.points)
				}

				// This allows the agents and humans to have the
				// same advantage of banking after someone else banks.
				if len(bankingPlayers) > 0 {
					g.askAiAgentsToBank()
				}

			case ROLL_DICE:
				keepPrompting = false
			}
		}

		// Round is over when all players bank or a 7 is rolled
		if len(g.results.getUnbankedPlayers()) == 0 {
			keepRolling = false
			bankedRound = true
		} else {
			dice, keepRolling = g.roll(&round)
		}
	}

	if !bankedRound {
		fmt.Println()
		fmt.Println(dice)
		fmt.Printf("Roll Number: %d\n\r", len(round.rolls))
	}
	fmt.Printf("Round %d done!\n\r\n\r", g.currentRound+1)
}

/**
 * Asks all AI Agents if they want to bank.
 *
 * @note Agents will be asked to bank every time regardless if they have banked or not.
 *       This is only to pass the game data to the AI Agent easily and won't change when
 *       they initially stated they wanted to bank.
 */
func (g *Game) askAiAgentsToBank() {
	for _, player := range g.players {
		if player.AiAgent() {
			response := player.Bank(g)
			if !g.results.playerBanked(player) && response {
				g.results.playerBanks(player, g.rounds[g.currentRound].points)
				fmt.Printf("AI Agent '%s' banked!\n\r", player.Name())
			}
		}
	}
}

/**
 * Rolls the dice for a given round
 */
func (g *Game) roll(r *round) (Dice, bool) {
	roll := new(Dice).roll(g.r)
	r.rolls = append(r.rolls, roll)
	newPts, cont := roll.Points(len(r.rolls), r.points)
	if cont {
		r.points = newPts
	}
	return roll, cont
}

//////////////////////////////////////////////
//                                          //
// Copy of the game data so that players    //
// such as AI can't manipulate the the Game //
//                                          //
//////////////////////////////////////////////

type BankDataSnapshot struct {
	CurrentRound uint8
	RoundPoints  uint
	Roll         Dice // Last Rolled Dice
	Players      []PlayerDataSnapshot
	Seed         int64
}

type PlayerDataSnapshot struct {
	Name   string
	Points uint // Total Points thus far
	Banked bool // True if banked this round, otherwise false
}

/**
 * Gets a snapshot of the game data
 *
 * @param requestor Player who is requesting the data. This is so as to
 *                  exclude them from the snapshot data that is returned.
 *
 * @return Snapshot of game data
 */
func (g *Game) GetData(requestor Player) BankDataSnapshot {
	round := g.rounds[g.currentRound]
	playerData := make([]PlayerDataSnapshot, len(g.players)-1)

	for name, player := range g.players {
		if name != requestor.Name() {
			playerData = append(playerData, g.results.getPlayerData(player))
		}
	}

	return BankDataSnapshot{
		CurrentRound: g.currentRound,
		RoundPoints:  round.points,
		Roll:         round.rolls[len(round.rolls)-1],
		Players:      playerData,
		Seed:         g.seed,
	}
}
