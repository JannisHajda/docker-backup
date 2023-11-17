package main

import (
	"github.com/JannisHajda/docker-backup/internal/cli"
	"github.com/JannisHajda/docker-backup/internal/utils"
)

func main() {
	err := utils.PrepareEnv()

	if err != nil {
		panic(err)
	}

	cli.Start()
}
