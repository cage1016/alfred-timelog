/*
Copyright © 2022 KAI CHU CHUNG
*/
package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/cage1016/alfred-timelog/cmd"
)

func main() {
	cmd.Execute()
}