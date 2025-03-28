package main

import (
	"bufio"
	"fmt"
	"os"

	cli "github.com/MrR0b0t1001/pokedexcli/internal/cliCommand"
	cf "github.com/MrR0b0t1001/pokedexcli/internal/config"
	pk "github.com/MrR0b0t1001/pokedexcli/internal/pokemon"
)

func printError(err error) {
	fmt.Println("Error: %v occurred", err)
}

func main() {
	supportedCommands := map[string]cli.CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    cli.CommandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    cli.CommandExit,
		},
		"map": {
			Name:        "map",
			Description: "Give you the map",
			Callback:    cli.CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "shows you the previous page",
			Callback:    cli.CommandMapB,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore an area",
			Callback:    cli.CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a pokemon",
			Callback:    cli.CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a caught pokemon's info",
			Callback:    cli.CommandInspect,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	cnfg := &cf.Config{
		Next:     nil,
		Previous: nil,
	}
	pokedex := &pk.Pokedex{
		Pkdex: map[string]pk.Pokemon{},
	}
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			cleaned := cli.CleanInput(text)

			cmd, ok := supportedCommands[cleaned[0]]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}

			switch cmd.Name {
			case "explore":
				err := cmd.Callback(cnfg, cleaned[1], pokedex)
				if err != nil {
					printError(err)
				}
			case "catch":
				err := cmd.Callback(cnfg, cleaned[1], pokedex)
				if err != nil {
					printError(err)
				}
			case "inspect":
				err := cmd.Callback(cnfg, cleaned[1], pokedex)
				if err != nil {
					printError(err)
				}
			default:
				err := cmd.Callback(cnfg, "", pokedex)
				if err != nil {
					printError(err)
				}
			}

			fmt.Print("\n\nUsage:\n")
			for _, val := range supportedCommands {
				fmt.Printf("\n%v: %v\n", val.Name, val.Description)
			}
		}
		if err := scanner.Err(); err != nil {
			printError(err)
			os.Exit(-1)
		}
	}
}
