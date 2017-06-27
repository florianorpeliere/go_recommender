package neighborhood_model

import (
	"github.com/shopspring/decimal"
	"testing"
	"math"
)

func TestSimilarity(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}

	userAVec := []float64{7, 6, 7, 4, 5, 4}
	userBVec := []float64{0, 3, 3, 1, 1, 0}

	assertEqual(math.Ceil(CosineSim(userBVec, userBVec)), float64(1))
	assertEqual(PearsonSim(userBVec, userBVec), float64(1))

	value, _ := decimal.NewFromFloat(CosineSim(userAVec, userBVec)).Round(3).Float64()
	assertEqual(value, float64(0.956))

	value, _ = decimal.NewFromFloat(PearsonSim(userAVec, userBVec)).Round(3).Float64()
	assertEqual(value, float64(0.894))
}

func TestPredictionUserBased(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}

	usersItems := [][]float64{
		{7, 6, 7, 4, 5, 4},
		{6, 7, 0, 4, 3, 4},
		{0, 3, 3, 1, 1, 0},
		{1, 2, 2, 3, 3, 4},
		{1, 0, 1 ,2, 3, 3},
	}

	pred, _ := getPredictionUserBased(usersItems, 2, 0, 2)
	value, _ := decimal.NewFromFloat(pred).Round(2).Float64()
	assertEqual(value, float64(3.35))

	pred, _ = getPredictionUserBased(usersItems, 2, 5, 2)
	value, _ = decimal.NewFromFloat(pred).Round(2).Float64()
	assertEqual(value, float64(0.86))
}

func TestPredictionItemBased(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}

	usersItems := [][]float64{
		{7, 6, 7, 4, 5, 4},
		{6, 7, 0, 4, 3, 4},
		{0, 3, 3, 1, 1, 0},
		{1, 2, 2, 3, 3, 4},
		{1, 0, 1 ,2, 3, 3},
	}

	pred, _ := getPredictionItemBased(usersItems, 2, 0, 2)
	value, _ := decimal.NewFromFloat(pred).Round(2).Float64()
	assertEqual(value, float64(3))

	pred, _ = getPredictionItemBased(usersItems, 2, 5, 2)
	value, _ = decimal.NewFromFloat(pred).Round(2).Float64()
	assertEqual(value, float64(1))
}