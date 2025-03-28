package clicommand

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	cf "github.com/MrR0b0t1001/pokedexcli/internal/config"
	pc "github.com/MrR0b0t1001/pokedexcli/internal/pokecache"
	pk "github.com/MrR0b0t1001/pokedexcli/internal/pokemon"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error
}

func CleanInput(text string) []string {
	return strings.Fields(text)
}

func CommandExit(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func CommandMap(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error {
	fmt.Println("This is the map!")
	// We need to make a request to the Poke API to retrieve locations
	var expiration time.Duration = 20 * time.Second
	cache := pc.NewCache(expiration)

	url := "https://pokeapi.co/api/v2/location-area/"

	if cnfg.Next != nil {
		url = *cnfg.Next
	}

	data, ok := cache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		cnfgResponse := cf.ConfigResponse{}

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

		cache.Add(url, newData)
	} else {
		locationData := []cf.Locations{}
		if err := json.Unmarshal(data, &locationData); err != nil {
			return err
		}

		for _, location := range locationData {
			fmt.Println(location.Name)
		}
	}

	return nil
}

func CommandMapB(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error {
	// We need to make a request to the Poke API to retrieve locations
	if cnfg.Previous == nil {
		fmt.Println("This is the first page")
		return nil
	}

	var expiration time.Duration = 20 * time.Second
	cache := pc.NewCache(expiration)

	url := *cnfg.Previous
	data, ok := cache.Get(url)

	if !ok {

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		cnfgResponse := cf.ConfigResponse{}

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

		cache.Add(url, newData)

	} else {

		locationData := []cf.Locations{}
		if err := json.Unmarshal(data, &locationData); err != nil {
			return err
		}

		for _, location := range locationData {
			fmt.Println(location.Name)
		}

	}

	return nil
}

func CommandExplore(cfng *cf.Config, name string, pokedex *pk.Pokedex) error {
	url := "https://pokeapi.co/api/v2/location-area/" + name

	var expiration time.Duration = 20 * time.Second
	cache := pc.NewCache(expiration)

	data, ok := cache.Get(url)

	if !ok {

		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		areaEnc := pk.EncountersResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&areaEnc); err != nil {
			return err
		}

		fmt.Println("Found Pokemon:")
		for _, pokemons := range areaEnc.Pokemon_Encounters {
			fmt.Printf("- %v\n", pokemons.Poke.Name)
		}

		newData, err := json.Marshal(areaEnc.Pokemon_Encounters)
		if err != nil {
			return err
		}

		cache.Add(name, newData)

	} else {

		pokemonData := []pk.Encounter{}
		if err := json.Unmarshal(data, &pokemonData); err != nil {
			return err
		}

		for _, pokemons := range pokemonData {
			fmt.Println(pokemons.Poke.Name)
		}

	}

	return nil
}

func catchPokemon(exp int) bool {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	// the highest exp given from one pokemon is 666
	// We extract the first digit of the exp
	// and generate a number up to that digit i.e from 1 to 6
	// if that number matches our secretNum then pokemon caught
	// scales based on the Experience
	expStr := strconv.Itoa(exp)
	firstDigit := 1

	if len(expStr) >= 3 {
		firstDigit, _ = strconv.Atoi(string(expStr[0]))
	}

	secretNum := r.Intn(firstDigit) + 1
	randomNum := r.Intn(firstDigit) + 1
	if randomNum == secretNum {
		return true
	}

	return false
}

func CommandCatch(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	pokemonInfo := pk.Pokemon{}

	if err := json.NewDecoder(res.Body).Decode(&pokemonInfo); err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", name)
	time.Sleep(2 * time.Second)

	if ok := catchPokemon(pokemonInfo.Experience); !ok {
		fmt.Printf("%v escaped!\n", name)
	} else {
		fmt.Printf("%v was caught!\n", name)
		pokedex.Pkdex[name] = pokemonInfo
	}
	fmt.Println("You may now inspect it with the inspect command.")
	time.Sleep(1 * time.Second)

	return nil
}

func CommandInspect(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error {
	pokemon, ok := pokedex.Get(name)
	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, statValue := range pokemon.Stats {
		fmt.Printf(" -%v: %v\n", statValue.Stat.Name, statValue.BaseStat)
	}

	fmt.Println("Types:")
	for _, tp := range pokemon.Types {
		fmt.Printf(" - %v\n", tp.Type.Name)
	}

	return nil
}

func CommandPokedex(cnfg *cf.Config, name string, pokedex *pk.Pokedex) error {
	pokedex.Print()
	return nil
}
