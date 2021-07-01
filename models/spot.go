package models

type Spot struct {
	SpotId    int64
	GameId    int64
	Value     string
	X         int
	Y         int
	NearSpots map[string]*Spot
	Status    string
}
