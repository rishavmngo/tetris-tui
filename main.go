package main

import (
	"fmt"
	"os"
	"time"

	"rishavmngo/tetris-tui-v2/board"

	"golang.org/x/term"
)

func enableRawMode() (*term.State, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return oldState, nil
}

func GetInput(inputChan chan string) {
	go func() {
		for {
			buf := make([]byte, 3)

			n, err := os.Stdin.Read(buf)
			if err != nil {
				break
			}

			if n == 1 {
				byteVal := buf[0]
				if byteVal == 3 {
					inputChan <- "QUIT"
				}
				if byteVal == 'q' {
					inputChan <- "QUIT"
				}
				if byteVal == 32 {
					inputChan <- "SPACE"
				}
			}

			if n == 3 && buf[0] == 27 && buf[1] == 91 {
				switch buf[2] {
				case 'A':
					inputChan <- "UP"
				case 'B':
					inputChan <- "DOWN"
				case 'C':
					inputChan <- "RIGHT"
				case 'D':
					inputChan <- "LEFT"
				}
			}
		}
	}()
}

func main() {
	oldState, _ := enableRawMode()

	defer term.Restore(int(os.Stdin.Fd()), oldState)
	game := board.NewBoard()
	game.SpwanPiece()

	quitChan := make(chan int)
	inputChan := make(chan string)
	renderTicker := time.NewTicker(50 * time.Millisecond)
	gravityTicker := time.NewTicker(800 * time.Millisecond)

	GetInput(inputChan)
	for !game.GameOver {
		select {
		case <-quitChan:
			return
		case key := <-inputChan:
			handleInput(key, game)
		case <-gravityTicker.C:
			if hitted := game.MoveDown(); !hitted {
				game.LockPosition()
			}
		case <-renderTicker.C:
			game.Render()
		}
	}
	fmt.Printf("Total Score: %d", game.Score)
}

func handleInput(key string, game *board.Board) {
	switch key {
	case "UP":
		game.Rotate()
	case "DOWN":
		game.MoveDown()
	case "LEFT":
		game.MoveLeft()
	case "RIGHT":
		game.MoveRight()
	case "SPACE":
		for game.MoveDown() {
		}
		game.LockPosition()
	case "QUIT":
		game.GameOver = true
	}
}
