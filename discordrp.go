package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hugolgst/rich-go/client"
)

func sendHandshake() {
	err := client.Login("1220960048704913448")
	if err != nil {
		log.Fatal(err)
	}
}

func startRP() {
	sendHandshake()
	go updateRichPresence()

	for {

		updateRP <- true
		time.Sleep(16 * time.Second)
	}
}

var idling bool

var updateRP = make(chan bool)

func updateRichPresence() {

	for {
		<-updateRP
		if idling {
			return
		}
		if len(LobbyHistory) == 0 {
			continue
		}
		currentLobby := LobbyHistory[len(LobbyHistory)-1]
		var player LobbyPlayerStruct
		for _, y := range currentLobby.Players {
			if y.Name == currentPilot && y.Active {
				player = y
			}
		}
		playerkdr := getKDR(player)

		state := playerkdr + "-" + "Objectives: " + countObjectives(currentLobby)

		details := currentLobby.Lobby.PreLobby.ScenarioInfo

		var aircraft string
		var smalltext string
		largetext := "Currently flying: "
		switch player.Aircraft {
		case "EF-24G":
			largetext += "EF-24G Mischief"
			aircraft = "ef24g"
			smalltext = "Split the throttles soulja boy!"
		case "F-45A":
			largetext += "F-45A Ghost"
			aircraft = "f45a"
			smalltext = "The mind of an f45 main cannot comprehend a 26 chaffing."
		case "F/A-26B":
			largetext += "F/A-26B Wasp"
			aircraft = "fa26b"
			smalltext = "Carrying literally the entire weight of a Shipping Container in bombs."
		case "T-55":
			largetext += "T-55 Tyro"
			aircraft = "t55"
			smalltext = "Courtesy of dubyaaa"
		case "AH-94":
			largetext += "AH-94 Dragonfly"
			aircraft = "ah94"
			smalltext = "Courtesy of kentuckyfrieda10wallsimper"
		case "AV-42C":
			largetext += "A/V-42C Kestrel"
			aircraft = "av42c"
			smalltext = "Sacred meeting of the kestrel kouncil."

		}

		if len(details) >= 20 {
			details = details[:16]
			details = details + "..."
		}

		err := client.SetActivity(client.Activity{
			State:      state,
			Details:    details,
			LargeImage: aircraft,
			LargeText:  largetext,
			SmallImage: "vtolvr",
			SmallText:  smalltext,

			Timestamps: &client.Timestamps{
				Start: &currentLobby.Lobby.JoinTime,
			},
		})
		if err != nil {
			fmt.Println(err)
		}
	}

}

func idle() {
	idling = true
	now := time.Now()
	err := client.SetActivity(client.Activity{
		State:      "Idling in game...",
		Details:    "Currently not in any match",
		LargeImage: "vtolvr",
		LargeText:  "VTOL VR",
		SmallImage: "vtolvr",
		SmallText:  "Plugin made by @f45a",

		Timestamps: &client.Timestamps{
			Start: &now,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func countObjectives(currentLobby LobbyStruct) (status string) {

	completed := 0
	total := 0
	for _, y := range currentLobby.Lobby.Objectives {
		if y.Result == "Completed" {
			completed += 1
		}
		total += 1
	}
	status = fmt.Sprint(completed) + " of " + fmt.Sprint(total) + " completed"
	return
}

func getKDR(player LobbyPlayerStruct) (playerkdr string) {

	playerk := fmt.Sprint(len(player.Kills))
	playerd := fmt.Sprint(len(player.Deaths))

	var playerdint int
	if len(player.Deaths) == 0 {
		playerdint = 1
	} else {
		playerdint = len(player.Deaths)
	}

	intpr := len(player.Kills) / playerdint

	playerr := fmt.Sprint(intpr)

	playerkdr = playerk + "K/" + playerd + "D/" + playerr + "R"

	return
}
