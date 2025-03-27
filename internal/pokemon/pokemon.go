package pokemon

type Pokemon struct {
	Name       string `json:"name"`
	Experience int    `json:"base_experience"`
}

type Encounter struct {
	Poke Pokemon `json:"pokemon"`
}

type EncountersResponse struct {
	Pokemon_Encounters []Encounter `json:"pokemon_encounters"`
}

type Pokedex struct {
	Pkdex map[string]Pokemon
}
