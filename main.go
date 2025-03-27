package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	pc "github.com/MrR0b0t1001/pokedexcli/internal/pokecache"
)

var cnfg = &config{}

type configResponse struct {
	Results  []locations `json:"results"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
}

type config struct {
	Next     *string
	Previous *string
}

type locations struct {
	Name string `json:"name"`
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
	cnfg        *config
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func commandMap() error {
	fmt.Println("This is the map!")
	// We need to make a request to the Poke API to retrieve locations

	url := "https://pokeapi.co/api/v2/location-area/"

	data, ok := pc.NewCache().Get(url)

	if !ok {

		if cnfg.Next != nil {
			url = *cnfg.Next
		}

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		cnfgResponse := configResponse{}

		if err := json.NewDecoder(res.Body).Decode(&cnfgResponse); err != nil {
			return err
		}

		for _, location := range cnfgResponse.Results {
			fmt.Println(location.Name)
		}

		cnfg.Next = cnfgResponse.Next
		cnfg.Previous = cnfgResponse.Previous

		newData, err := json.Marshal(cnfgResponse.Results)
		if err != nil {
			return err
		}

		pc.NewCache().Add(url, newData)
	} else {
		locationData := []locations{}
		if err := json.Unmarshal(data, &locationData); err != nil {
			return err
		}

		for _, location := range locationData {
			fmt.Println(location.Name)
		}
	}

	return nil
}

func commandMapB() error {
	// We need to make a request to the Poke API to retrieve locations
	if cnfg.Previous == nil {
		fmt.Println("This is the first page")
		return nil
	}
	url := *cnfg.Previous

	data, ok := pc.NewCache().Get(url)

	if !ok {

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		cnfgResponse := configResponse{}

		if err := json.NewDecoder(res.Body).Decode(&cnfgResponse); err != nil {
			return err
		}

		for _, location := range cnfgResponse.Results {
			fmt.Println(location.Name)
		}

		cnfg.Next = cnfgResponse.Next
		cnfg.Previous = cnfgResponse.Previous

		newData, err := json.Marshal(cnfgResponse.Results)
		if err != nil {
			return err
		}

		pc.NewCache().Add(url, newData)

	} else {

		locationData := []locations{}
		if err := json.Unmarshal(data, &locationData); err != nil {
			return err
		}

		for _, location := range locationData {
			fmt.Println(location.Name)
		}

	}

	return nil
}

func main() {
	supportedCommands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			cnfg:        cnfg,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			cnfg:        cnfg,
		},
		"map": {
			name:        "map",
			description: "Give you the map",
			callback:    commandMap,
			cnfg:        cnfg,
		},
		"mapb": {
			name:        "mapb",
			description: "shows you the previous page",
			callback:    commandMapB,
			cnfg:        cnfg,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			cleaned := cleanInput(text)

			cmd, ok := supportedCommands[cleaned[0]]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}

			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Print("Usage:\n")
			for _, val := range supportedCommands {
				fmt.Printf("%v: %v\n", val.name, val.description)
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

	}
}
