package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculatePacks_BasicCases(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}

	// 1 item -> must send 1 x 250
	res, err := CalculatePacks(1, sizes)
	require.NoError(t, err)
	assert.Equal(t, 250, res.TotalItems)
	assert.Equal(t, 249, res.ExtraItems)
	assert.Equal(t, map[int]int{250: 1}, res.Packs)

	// 250 items -> exactly 1 x 250
	res, err = CalculatePacks(250, sizes)
	require.NoError(t, err)
	assert.Equal(t, 250, res.TotalItems)
	assert.Equal(t, 0, res.ExtraItems)
	assert.Equal(t, map[int]int{250: 1}, res.Packs)

	// 251 items -> 1 x 500 (not 2 x 250)
	res, err = CalculatePacks(251, sizes)
	require.NoError(t, err)
	assert.Equal(t, 500, res.TotalItems)
	assert.Equal(t, 249, res.ExtraItems)
	assert.Equal(t, map[int]int{500: 1}, res.Packs)

	// 501 items -> 1 x 500 + 1 x 250 (750 total, 2 packs)
	res, err = CalculatePacks(501, sizes)
	require.NoError(t, err)
	assert.Equal(t, 750, res.TotalItems)
	assert.Equal(t, 249, res.ExtraItems)
	assert.Equal(t, map[int]int{250: 1, 500: 1}, res.Packs)

	// 780 items -> 1 x 1000
	res, err = CalculatePacks(780, sizes)
	require.NoError(t, err)
	assert.Equal(t, 1000, res.TotalItems)
	assert.Equal(t, map[int]int{1000: 1}, res.Packs)

	// 12001 items -> 2 x 5000 + 1 x 2000 + 1 x 250
	res, err = CalculatePacks(12001, sizes)
	require.NoError(t, err)
	assert.Equal(t, 12250, res.TotalItems)
	assert.Equal(t, map[int]int{5000: 2, 2000: 1, 250: 1}, res.Packs)
}

func TestCalculatePacks_EdgeCase(t *testing.T) {
	sizes := []int{23, 31, 53}

	res, err := CalculatePacks(500000, sizes)
	require.NoError(t, err)

	expected := map[int]int{23: 2, 31: 7, 53: 9429}
	assert.Equal(t, expected, res.Packs)

	total := 0
	for size, qty := range res.Packs {
		total += size * qty
	}
	assert.Equal(t, total, res.TotalItems)
	assert.GreaterOrEqual(t, res.TotalItems, 500000)
	assert.Equal(t, res.TotalItems-500000, res.ExtraItems)
}
