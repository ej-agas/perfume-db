package internal

import "fmt"

type Concentration int

const (
	Unknown Concentration = iota
	EauFraiche
	EauDeCologne
	EauDeToilette
	EauDeParfum
	Parfum
	ExtraitDeParfum
)

var ConcentrationMap = map[string]Concentration{
	"Unknown":           Unknown,
	"Eau Fraiche":       EauFraiche,
	"Eau De Cologne":    EauDeCologne,
	"Eau De Toilette":   EauDeToilette,
	"Eau De Parfum":     EauDeParfum,
	"Parfum":            Parfum,
	"Extrait De Parfum": ExtraitDeParfum,
}

func ConcentrationFromString(s string) (Concentration, error) {
	concentration, ok := ConcentrationMap[s]
	if !ok {
		return -1, fmt.Errorf("unknown concentration: %s", s)
	}

	return concentration, nil
}

func (c Concentration) String() string {
	switch c {
	case Unknown:
		return "Unknown"
	case EauFraiche:
		return "Eau Fraiche"
	case EauDeToilette:
		return "Eau De Toilette"
	case EauDeCologne:
		return "Eau De Cologne"
	case EauDeParfum:
		return "Eau De Parfum"
	case Parfum:
		return "Parfum"
	default:
		return "Unknown"
	}
}
