/*
* Mini max algorythm for Tictactoe Bot (https://botsArena.tinad.fr or bollosseum)
* By Gnieark https://blog-du-grouik.tinad.fr 2018-06
* I am learning golang it's my first script, don't take it seriously
*/

package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"./tictactoe"
)

//Plate somewhere to put the json encoded map
type Plate struct{
	Cell00 string `json:"0-0,omitempty"`
	Cell01 string `json:"0-1,omitempty"`
	Cell02 string `json:"0-2,omitempty"`
	Cell10 string `json:"1-0,omitempty"`
	Cell11 string `json:"1-1,omitempty"`
	Cell12 string `json:"1-2,omitempty"`
	Cell20 string `json:"2-0,omitempty"`
	Cell21 string `json:"2-1,omitempty"`
	Cell22 string `json:"2-2,omitempty"`
}

//QuestionMessage Somewhere to put the whole JSON message from botsarena
type QuestionMessage struct {
	GameID      string  `json:"game-id,omitempty"`
	Action      string  `json:"action,omitempty"`   
	Game        string  `json:"game,omitempty"`   
	Players     int     `json:"players,omitempty"`  
	Board       Plate	`json:"board"` 
	You         string 	`json:"you,omitempty"`
	PlayerIndex int 	`json:"player-index,omitempty"`
}


func parseQuery(w http.ResponseWriter, r *http.Request){

	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

	decoder := json.NewDecoder(r.Body)
	var questionMessage QuestionMessage
	err := decoder.Decode(&questionMessage)

	if err != nil {
		panic(err)
	}

	switch questionMessage.Action{
		case "init":
			w.Write([]byte("{\"name\":\"GniearkGolangTictactoe\"}"))
		case "play-turn":

			//Convert the map to an exploitable array
			var tmap [3][3]int
			tmap[0][0] = tictactoeSymbolsToInt(questionMessage.Board.Cell00,questionMessage.You)
			tmap[0][1] = tictactoeSymbolsToInt(questionMessage.Board.Cell01,questionMessage.You)
			tmap[0][2] = tictactoeSymbolsToInt(questionMessage.Board.Cell02,questionMessage.You)
			tmap[1][0] = tictactoeSymbolsToInt(questionMessage.Board.Cell10,questionMessage.You)
			tmap[1][1] = tictactoeSymbolsToInt(questionMessage.Board.Cell11,questionMessage.You)
			tmap[1][2] = tictactoeSymbolsToInt(questionMessage.Board.Cell12,questionMessage.You)
			tmap[2][0] = tictactoeSymbolsToInt(questionMessage.Board.Cell20,questionMessage.You)
			tmap[2][1] = tictactoeSymbolsToInt(questionMessage.Board.Cell21,questionMessage.You)
			tmap[2][2] = tictactoeSymbolsToInt(questionMessage.Board.Cell22,questionMessage.You)
		
			target, score := tictactoe.PlayOn(tmap, 1)
			fmt.Fprintf(w, "{\"play\":\"%d-%d\",\"Comment\":\"score %d\"}", target.X, target.Y,score)

		default:
			w.Write([]byte("Error " + questionMessage.Action ))
			
	}

}
func tictactoeSymbolsToInt (symbolToTest string,mySymbol string) int{
	switch symbolToTest {
		case mySymbol:
			return 1
		case "": 
			return 0
		case " ": 
			return 0
		default:
			return 2
	}
}
func arena(w http.ResponseWriter, r *http.Request){
	data, err := ioutil.ReadFile("./arena.html")
	if err == nil {
        w.Write(data)
    } else {
        w.WriteHeader(404)
        w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
	}

}
func main() {
	http.HandleFunc("/arena", arena)
    http.HandleFunc("/", parseQuery)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

