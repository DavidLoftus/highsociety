package main

import (
	"flag"
	"fmt"
	"github.com/DavidLoftus/highsociety"
	"log"
)

var flNumUsers = flag.Int("users", 0, "number of users playing")

func main() {
	gameState := highsociety.NewGame(*flNumUsers)

	for !gameState.GameOver() {
		player := gameState.CurrentPlayer()

		fmt.Printf("Player %d, it's your turn.", player+1)

		_, err := fmt.Scanln()
		if err != nil {
			log.Fatal(err)
		}
	}

}
