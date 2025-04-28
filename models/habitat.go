package models

type Habitat struct {
	ID      int    `json:"id"`
	Adresse string `json:"adresse"`
	User    int    `json:"user"`
}

type HabitatGestion struct {
	ID         int  `json:"id"`
	User       int  `json:"user"`
	Habitat    int  `json:"habitat"`
	Validation bool `json:"validation"`
}

type RoomsTable struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Habitat int    `json:"habitat"`
}
