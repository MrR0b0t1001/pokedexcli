package main

import (
	"bufio"
	"fmt"
	"os"

	cli "github.com/MrR0b0t1001/pokedexcli/internal/cliCommand"
	cf "github.com/MrR0b0t1001/pokedexcli/internal/config"
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
	}
	scanner := bufio.NewScanner(os.Stdin)
	cnfg := &cf.Config{
		Next:     nil,
		Previous: nil,
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

			if cmd.Name == "explore" {
				err := cmd.Callback(cnfg, cleaned[1])
				if err != nil {
					printError(err)
				}
			} else if cmd.Name == "catch" {
				err := cmd.Callback(cnfg, cleaned[1])
				if err != nil {
					printError(err)
				}

			} else {
				err := cmd.Callback(cnfg, "")
				if err != nil {
					printError(err)
				}
			}

			fmt.Print("\n\nUsage:\n")
			for _, val := range supportedCommands {
				fmt.Printf("%v: %v\n", val.Name, val.Description)
			}
		}
		if err := scanner.Err(); err != nil {
			printError(err)
			os.Exit(-1)
		}
	}
}
