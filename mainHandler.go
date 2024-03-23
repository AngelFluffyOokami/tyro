package main

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func onLobbyJoin(host bool) {

	var currentLobby LobbyStruct

	currentLobby.Lobby.ID = uuid.New()
	LobbyHistory = append(LobbyHistory, currentLobby)

	currentLobby.Lobby.WinningTeam = "Invalid"

	for {
		currentMessage := <-Message
		currentLobby = updateLobbyCount(currentLobby)
		if !currentLobby.Lobby.PreLobby.LoadedIn {

			currentLobby = preLobbyHandler(currentMessage, currentLobby, host)

		}
		switch {
		case strings.Contains(currentMessage, "Setting up slot UI: "):
			currentLobby = onSlotUISetup(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "ilst_"):
			currentLobby = onIlstUpdate(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "Resetting objective "):
			currentLobby = onResetObjective(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "Setting up objective "):
			currentLobby = onBeginObjective(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "FlightLogger: Completed objective: "):
			currentLobby = onCompleteObjective(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "FlightLogger: Failed objective: "):
			currentLobby = onFailObjective(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "has disconnected.") && strings.Contains(currentMessage, "FlightLogger: "):
			currentLobby = onPlayerLeave(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "Lobby message from ") && strings.Contains(currentMessage, " killed ") && !strings.Contains(currentMessage, " was killed by "):
			currentLobby = onKill(currentMessage, currentLobby)
		case strings.Contains(currentMessage, " was killed by") && strings.Contains(currentMessage, "Lobby message from") && !strings.Contains(currentMessage, "environment") && !strings.Contains(currentMessage, "$log_EF-24G") && !strings.Contains(currentMessage, "$log_T-55") && !strings.Contains(currentMessage, "$log_AH-94"):
			currentLobby = onEnvDeath(currentMessage, currentLobby)
		case strings.Contains(currentMessage, " was killed by") && strings.Contains(currentMessage, "Lobby message from") && !strings.Contains(currentMessage, "environment") && (strings.Contains(currentMessage, "$log_EF-24G") || strings.Contains(currentMessage, "$log_T-55")) || strings.Contains(currentMessage, "$log_AH-94"):
			currentLobby = onEnvDeathMC(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "identity updated"):
			currentLobby = onIdentityUpdate(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "Endmission - Final Winner: "):
			if strings.Contains(currentMessage, "Allied") {
				currentLobby.Lobby.WinningTeam = "Allied"

			} else if strings.Contains(currentMessage, "Enemy") {
				currentLobby.Lobby.WinningTeam = "Enemy"
			}
		case strings.Contains(currentMessage, ").SetTeam("):
			currentLobby = onSetTeam(currentMessage, currentLobby)

		case currentMessage == "LeaveLobby()":
			idle()
			InLobby = false
			currentLobby.Lobby.LeaveTime = time.Now()

			for x, y := range LobbyHistory {
				if y.Lobby.ID == currentLobby.Lobby.ID {
					LobbyHistory[x] = currentLobby
					done <- true
					return
				}
			}
		}

		for _, x := range LobbyHistory {
			if x.Lobby.ID == currentLobby.Lobby.ID {
				x = currentLobby
			}
		}
		updateRichPresence(currentLobby)
		done <- true

	}
}

/*
*	preLobbyHandler handles the prelobby joining mechanism, and fills in information about the lobby,
*	into the array that will be needed later once joined and in game
 */
func preLobbyHandler(currentMessage string, currentLobby LobbyStruct, host bool) LobbyStruct {
	switch {
	case strings.Contains(currentMessage, "Launching Multiplayer game for ") && host:

		currentLobby.Lobby.PreLobby.JoinAttempted = true
		currentLobby.Lobby.PreLobby.LobbyInfo, _ = strings.CutPrefix(currentMessage, "Attempting to join lobby ")
		currentLobby.Lobby.ID64 = "host"
		currentLobby.Lobby.Name, _ = matchHostedMap(currentMessage)
		currentLobby.Lobby.PreLobby.JoinAttempted = true
		currentLobby.Lobby.PreLobby.JoinReqStatus = true

		currentLobby.Lobby.PreLobby.ScenarioInfo = currentLobby.Lobby.Name
	case strings.Contains(currentMessage, "Attempting to join lobby"):
		var name string
		currentLobby.Lobby.PreLobby.JoinAttempted = true
		currentLobby.Lobby.PreLobby.LobbyInfo, _ = strings.CutPrefix(currentMessage, "Attempting to join lobby ")
		currentLobby.Lobby.ID64, name, _ = strings.Cut(currentLobby.Lobby.PreLobby.LobbyInfo, " (")
		currentLobby.Lobby.ID64 = strings.TrimPrefix(currentLobby.Lobby.ID64, "VTMPMainMenu: Attempting to join lobby ")
		currentLobby.Lobby.Name, _ = strings.CutSuffix(name, ")")
	case strings.Contains(currentMessage, "Join request accepted!"):
		currentLobby.Lobby.PreLobby.JoinReqStatus = true
	case strings.Contains(currentMessage, "getting scenario"):
		currentLobby.Lobby.PreLobby.ScenarioInfo, _ = strings.CutPrefix(currentMessage, "getting scenario ")
	case strings.Contains(currentMessage, "Connecting to host: "):
		var trimmedMessage string
		var ID string
		_, trimmedMessage, _ = strings.Cut(currentMessage, "Creating socket client. Connecting to host: ")
		currentLobby.Lobby.HostName, ID, _ = strings.Cut(trimmedMessage, " (")
		currentLobby.Lobby.HostID64, _, _ = strings.Cut(ID, ")")
	case currentMessage == "Connected":
		currentLobby.Lobby.PreLobby.LoadedIn = true
		currentLobby.Lobby.JoinTime = time.Now()
	}
	return currentLobby
}