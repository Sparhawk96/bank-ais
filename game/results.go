package game

import "github.com/Sparhawk96/bank-ais/table"

type results struct {
	players            map[string]*playerNode
	firstPlayer        *playerNode
	largestName        int
	humanPlayers       int
	bankedHumanPlayers int
}

type playerNode struct {
	Player

	pts    uint
	banked bool
	ahead  *playerNode
	behind *playerNode
}

/**
 * Adds a Player to the results table
 *
 * @param player Player to be added to the results table
 */
func (r *results) addPlayer(player Player) {
	if r.players == nil {
		r.players = make(map[string]*playerNode)
	}

	if _, have := r.players[player.Name()]; !have {
		pn := &playerNode{Player: player}
		r.players[pn.Name()] = pn

		if !pn.AiAgent() {
			r.humanPlayers++
		}

		if nameLen := len(pn.Name()); r.largestName < nameLen {
			r.largestName = nameLen
		}

		if r.firstPlayer == nil {
			r.firstPlayer = pn
		} else {
			lastPlayer := r.firstPlayer
			for lastPlayer.behind != nil {
				lastPlayer = lastPlayer.behind
			}
			lastPlayer.behind = pn
		}
	}
}

/**
 * Marks a Player as Banked and updates their score
 *
 * @param player Player who is banking
 * @param pts Points they accrued from the round
 *
 * @return True if all human players have banked, otherwise False
 */
func (r *results) playerBanks(player Player, pts uint) bool {
	pn := r.players[player.Name()]
	pn.pts += pts
	pn.banked = true

	if !pn.AiAgent() {
		r.bankedHumanPlayers++
	}

	playerAhead := pn.ahead
	for playerAhead != nil && playerAhead.pts < pn.pts {
		playerAhead = playerAhead.ahead
	}

	// Did player get enough pts to move up in the ranking?
	if playerAhead != pn.ahead {
		// If moved it can't be in first place so
		// there won't be a nil pointer dereference
		pn.ahead.behind = pn.behind
		pn.behind.ahead = pn.ahead

		// Will the player be moved to first place?
		if playerAhead == nil {
			pn.ahead = nil
			pn.behind = r.firstPlayer
			r.firstPlayer.ahead = pn
			r.firstPlayer = pn
		} else {
			pn.ahead = playerAhead
			pn.behind = playerAhead.behind
			playerAhead.behind = pn
			playerAhead.behind.ahead = pn
		}
	}

	return r.humanPlayers == r.bankedHumanPlayers
}

/**
 * Dictates if a player has banked for the round
 *
 * @param player To check if they have banked
 *
 * @return True if the player has banked, otherwise false
 */
func (r *results) playerBanked(player Player) bool {
	return r.players[player.Name()].banked
}

/**
 * Unmarks all players as banked
 */
func (r *results) unbankAllPlayers() {
	r.bankedHumanPlayers = 0
	player := r.firstPlayer

	for player != nil {
		player.banked = false
		player = player.behind
	}
}

/**
 * Gets all of the players who haven't banked yet this round
 *
 * @return List of unbanked players
 */
func (r *results) getUnbankedPlayers() []Player {
	unbankedPlayers := make([]Player, 0)

	for player := r.firstPlayer; player != nil; player = player.behind {
		if !player.banked {
			unbankedPlayers = append(unbankedPlayers, player)
		}
	}

	return unbankedPlayers
}

/**
 * Gets the player data from the current results
 *
 * @param player Player to get data for
 *
 * @return The current data about the player
 */
func (r *results) getPlayerData(player Player) PlayerDataSnapshot {
	p := r.players[player.Name()]

	return PlayerDataSnapshot{
		Name:   p.Name(),
		Points: p.pts,
		Banked: p.banked,
	}
}

func (r *results) String() string {
	t := new(table.Table)

	// Setup Headers
	playerHdr := "Players"
	aiAgentHdr := "AI Agent"
	bankedHdr := "Banked"
	pointsHdr := "Points"

	t.CreateColumn(playerHdr, table.LEFT, 0)
	t.CreateColumn(aiAgentHdr, table.CENTER, 0)
	t.CreateColumn(bankedHdr, table.CENTER, 0)
	t.CreateColumn(pointsHdr, table.LEFT, 0)

	for player := r.firstPlayer; player != nil; player = player.behind {
		data := map[string]any{
			playerHdr: player.Name(),
			pointsHdr: player.pts,
		}

		if player.banked {
			data[bankedHdr] = "✔"
		}

		if player.AiAgent() {
			data[aiAgentHdr] = "✔"
		}

		t.AddEntry(data)
	}

	return t.String()
}
