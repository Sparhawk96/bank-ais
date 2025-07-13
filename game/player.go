package game

type Player interface {
	/**
	 * Gets the name of the player
	 *
	 * @return The player name
	 */
	Name() string

	/**
	 * Dictates if the player will bank or not.
	 *
	 * @note This is only used for AI agents. Will always be called regardless if the
	 *       agent banked or not so as to always provide the game data to the agents.
	 *       If the agent has already banked the game will preform a noop.
	 *
	 * @return True if the player wants to bank, otherwise false
	 */
	Bank(*Game) bool

	/**
	 * Dictates if the player is an AI Agent
	 *
	 * @return True if player is an AI Agent, otherwise false
	 */
	AiAgent() bool
}

/**
 * Creates a new Human Player
 *
 * @param name Human player's name
 *
 * @return The created player
 */
func NewHumanPlayer(name string) Player {
	return HumanPlayer{name}
}

type HumanPlayer struct {
	name string
}

func (r HumanPlayer) Name() string {
	return r.name
}

func (r HumanPlayer) Bank(g *Game) bool {
	return false // Real Player will do so via prompts
}

func (r HumanPlayer) AiAgent() bool {
	return false
}
