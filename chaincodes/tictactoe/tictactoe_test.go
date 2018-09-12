package main

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestCheckForWinner(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+0].Value = symbolX
	game.Cells[gridWidth*0+1].Value = symbolX
	game.Cells[gridWidth*0+2].Value = symbolX
	winner:= CheckForWinner(game)
	if winner!= "org2" {
		t.Fatal("row win validation failed")
	}
}

func TestCheckForWinnerDiagnol(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+0].Value = symbolO
	game.Cells[gridWidth*1+1].Value = symbolO
	game.Cells[gridWidth*2+2].Value = symbolO

	winner:= CheckForWinner(game)
	if winner!= "org1" {
		t.Fatal("diagnol test failed")
	}
}

func TestCheckForWinnerDiagnol2(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+2].Value = symbolO
	game.Cells[gridWidth*1+1].Value = symbolO
	game.Cells[gridWidth*2+0].Value = symbolO

	winner:= CheckForWinner(game)
	if winner!= "org1" {
		t.Fatal("diagnol test failed")
	}
}

func TestCheckForWinnerSecondRow(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*1+0].Value = symbolO
	game.Cells[gridWidth*1+1].Value = symbolO
	game.Cells[gridWidth*1+2].Value = symbolO

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "org1" {
		t.Fatal("2nd row test failed")
	}
}

func TestCheckForWinner3rdRow(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*2+0].Value = symbolX
	game.Cells[gridWidth*2+1].Value = symbolX
	game.Cells[gridWidth*2+2].Value = symbolX

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "org2" {
		t.Fatal("3rd row test failed")
	}
}

func TestCheckForWinner1stCol(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+0].Value = symbolX
	game.Cells[gridWidth*1+0].Value = symbolX
	game.Cells[gridWidth*2+0].Value = symbolX

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "org2" {
		t.Fatal("3rd row test failed")
	}
}

func TestCheckForWinner2ndCol(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+1].Value = symbolX
	game.Cells[gridWidth*1+1].Value = symbolX
	game.Cells[gridWidth*2+1].Value = symbolX

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "org2" {
		t.Fatal("3rd row test failed")
	}
}

func TestCheckForWinner3rdCol(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+2].Value = symbolX
	game.Cells[gridWidth*1+2].Value = symbolX
	game.Cells[gridWidth*2+2].Value = symbolX

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "org2" {
		t.Fatal("3rd row test failed")
	}
}

func TestCheckForWinnerNoWinner(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+2].Value = symbolX
	game.Cells[gridWidth*1+2].Value = symbolO
	game.Cells[gridWidth*2+2].Value = symbolX

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "" {
		t.Fatal("no winner test failed")
	}
}

func TestCheckForWinnerNoWinner2(t *testing.T) {
	game:= createGame()
	game.Cells[gridWidth*0+0].Value = symbolX
	game.Cells[gridWidth*0+1].Value = symbolO
	game.Cells[gridWidth*0+2].Value = symbolX

	game.Cells[gridWidth*1+0].Value = symbolO
	game.Cells[gridWidth*1+1].Value = symbolO
	game.Cells[gridWidth*1+2].Value = symbolX

	game.Cells[gridWidth*2+0].Value = symbolO
	game.Cells[gridWidth*2+1].Value = symbolX
	game.Cells[gridWidth*2+2].Value = symbolO

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "" {
		t.Fatal("no winner test failed")
	}
}

func TestCheckForWinnerOnEmptyGame(t *testing.T) {
	game:= createGame()

	winner:= CheckForWinner(game)
	fmt.Println(winner)
	if winner!= "" {
		t.Fatal("empty game test failed")
	}
}

func TestTictactoeGame_CreateNewGame(t *testing.T) {
	gameBoard:= TictactoeGame{}
	game:= gameBoard.CreateNewGame("org1", 1)
	rowMap:= make(map[int]int)
	colMap:= make(map[int]int)

	for _,c:= range game.Cells {
		rowMap[c.Row] = rowMap[c.Row]+1;
		colMap[c.Column] = colMap[c.Column]+1;
	}

	fmt.Println(rowMap)
	fmt.Println(colMap)
	response,err:= GenerateResponse("", []Game{game})
	if err != nil {
		t.Error(err)
	}

	payload:= TictactoeGameResponse{}
	err = json.Unmarshal(response, &payload)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(payload)

	//for _,c:= range game.Cells {
	//	fmt.Println(c)
	//}
}

func createGame() Game {
	game:= Game{
		Id:1,
		Cells: [9]Cell{
			{Row:0, Column:0, Value:""},
			{Row:0, Column:1, Value:""},
			{Row:0, Column:2, Value:""},
			{Row:1, Column:0, Value:""},
			{Row:1, Column:1, Value:""},
			{Row:1, Column:2, Value:""},
			{Row:2, Column:0, Value:""},
			{Row:2, Column:1, Value:""},
			{Row:2, Column:2, Value:""},
		},
		Players: [2]Player{
			{Name:"org1", Symbol:symbolO},
			{Name:"org2", Symbol:symbolX},
		},
		PlayerToPlayIndex: 0,
		IsCompleted:false,
		Winner:"",
	}

	return game
}
