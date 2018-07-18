/*
* Mini max algorythm for Tictactoe Bot (https://botsArena.tinad.fr or bollosseum)
* By Gnieark https://blog-du-grouik.tinad.fr 2018-06
* I am learning golang it's my first script, don't take it seriously
*/

package tictactoe

//Coords 2D coords
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

	var newPlayer int
	if currentPlayer == 1 {
		newPlayer = 2
	}else{
		newPlayer = 1
	}
	//recursion there
	_ ,nextScore := PlayOn(tmap,newPlayer)
	return -nextScore
	
}


// PlayOn return the better cell, and his score where to play
// @param tmap the grid Values 0 are empty cells
//		 currentPlayer int.  His digit 1 or 2
// @return beastCoord,beastScore
func PlayOn (tmap [3][3]int, currentPlayer int) (Coords,int){
	
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
