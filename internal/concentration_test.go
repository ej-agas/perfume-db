package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConcentration_String(t *testing.T) {
	assert.Equal(t, "Eau Fraiche", EauFraiche.String())
	assert.Equal(t, "Eau De Cologne", EauDeCologne.String())
	assert.Equal(t, "Eau De Toilette", EauDeToilette.String())
	assert.Equal(t, "Eau De Parfum", EauDeParfum.String())
	assert.Equal(t, "Parfum", Parfum.String())
	assert.Equal(t, "Unknown", Concentration(-1).String())
}

func TestConcentrationFromString(t *testing.T) {
	ef, err := ConcentrationFromString("Eau Fraiche")
	assert.Nil(t, err)
	assert.Equal(t, EauFraiche, ef)

	edt, err := ConcentrationFromString("Eau De Toilette")
	assert.Nil(t, err)
	assert.Equal(t, EauDeToilette, edt)

	edp, err := ConcentrationFromString("Eau De Parfum")
	assert.Nil(t, err)
	assert.Equal(t, EauDeParfum, edp)

	p, err := ConcentrationFromString("Parfum")
	assert.Nil(t, err)
	assert.Equal(t, Parfum, p)

	unknown, err := ConcentrationFromString("foo")
	assert.Error(t, err, "foo")
	assert.Equal(t, Concentration(-1), unknown)
}
