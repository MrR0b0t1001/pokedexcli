package pokemon

type Pokemon struct {
	Name string `json:"name"`

	Height     int `json:"height"`
	Weight     int `json:"weight"`
	Experience int `json:"base_experience"`

	Stats []PokeStats `json:"stats"`
	Types []PokeTypes `json:"types"`
}

type Encounter struct {
	Poke Pokemon `json:"pokemon"`
}

type EncountersResponse struct {
	Pokemon_Encounters []Encounter `json:"pokemon_encounters"`
}

type PokeStats struct {
	BaseStat int      `json:"base_stat"`
	Stat     StatInfo `json:"stat"`
}

type StatInfo struct {
	Name string `json:"name"`
}

type PokeTypes struct {
	Slot int      `json:"slot"`
	Type TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name string `json:"name"`
}

type Pokedex struct {
	Pkdex map[string]Pokemon
}

func (pkdx *Pokedex) Add(name string, pkmn Pokemon) {
	pkdx.Pkdex[name] = pkmn
}

func (pkdx *Pokedex) Get(name string) (Pokemon, bool) {
	pokemon, ok := pkdx.Pkdex[name]
	if !ok {
		return Pokemon{}, false
	}

	return pokemon, true
}
