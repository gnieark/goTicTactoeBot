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
)

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
type QuestionMessage struct {
	GameId      string  `json:"game-id,omitempty"`
	Action      string  `json:"action,omitempty"`   
	Game        string  `json:"game,omitempty"`   
	Players     int     `json:"players,omitempty"`  
	Board       Plate	`json:"board"` 
	You         string 	`json:"you,omitempty"`
	PlayerIndex int 	`json:"player-index,omitempty"`
}

type Coords struct{
	X int
	Y int
}

/**
* Give a score to a cell where to play
* @param tmap the grid Values 0 are empty cells
* 		 target. The coords to test
*		 currentPlayer int.  His digit 1 or 2
* @return int the score
*/
func scoreTarget (tmap [3][3] int, target Coords, currentPlayer int) int{
	
	
	tmap[target.X][target.Y] = currentPlayer
	//count the depth
	depth :=0
	for i := 0; i<3 ; i++{
		for j := 0; j<3 ; j++{
			if tmap[i][j] > 0 {
				depth++
			}
		}
	}

	/*
	* 0-0 | 0-1 | 0-2
	* 1-0 | 1-1 | 1-2
	* 2-0 | 2-1 | 2-2
	*/

	alignments := [8][3]Coords{
		{Coords{0,0},Coords{0,1},Coords{0,2}},
		{Coords{1,0},Coords{1,1},Coords{1,2}},
		{Coords{2,0},Coords{2,1},Coords{2,2}},
		{Coords{0,0},Coords{1,0},Coords{2,0}},
		{Coords{0,1},Coords{1,1},Coords{2,1}},
		{Coords{0,2},Coords{1,2},Coords{2,2}},
		{Coords{0,0},Coords{1,1},Coords{2,2}},
		{Coords{0,2},Coords{1,1},Coords{2,0}},
	}

	win:=false
	for i:=0; i < len(alignments) ; i++ {
		win=true
		for j:=0; j < 3 ; j++ {
			if tmap[alignments[i][j].X][alignments[i][j].Y] != currentPlayer{
				win = false
			}
		}
		if win  {	
			return 100 - depth
		}

	}

	//if it was the last cell
	if depth == 9 { return 0}

	//test if this target prevent to loose
	var newPlayer int
	if currentPlayer == 1 {
		newPlayer = 2
	}else{
		newPlayer = 1
	}
	_ ,nextScore := playOn(tmap,newPlayer)
	return -nextScore
	
}
/**
* return the better cell, and his score where to play
* @param tmap the grid Values 0 are empty cells
*		 currentPlayer int.  His digit 1 or 2
* @return beastCoord,beastScore
*/
func playOn (tmap [3][3]int, currentPlayer int) (Coords,int){
	
	beastScore := -999
	beastCoord := Coords{-1,-1} 
	//scorer les emplacements libres
	for i := 0; i < 3; i++ {
		for j:= 0; j < 3; j++ {
			if tmap[i][j] == 0 {
				sc:=scoreTarget(tmap,Coords{i,j},currentPlayer)

				if sc > beastScore {
					beastScore = sc
					beastCoord = Coords{i,j}
				}
			}
		}
	}
	
	return beastCoord,beastScore
}

//******** http, and parsing functions *******

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
		
			target, score := playOn(tmap, 1)
			fmt.Fprintf(w, "{\"play\":\"%d-%d\",\"Comment\":\"score %d\"}", target.X, target.Y,score)

		default:
			w.Write([]byte("Error " + questionMessage.Action ))
			
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