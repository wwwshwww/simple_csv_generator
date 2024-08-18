package dummy_value_generator

import (
	"math/rand"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

const (
	seed                = 999
	floatRoundDownScale = 2
)

var Selector = rand.New(rand.NewSource(seed))

func GetDummiesString(count int, length int) []string {
	return lo.Times(count, func(_ int) string {
		s := lo.Times(length, func(_ int) rune {
			return lo.LowerCaseLettersCharset[Selector.Intn(len(lo.LowerCaseLettersCharset))]
		})
		return string(s)
	})
}

func GetDummiesInt(count int, min, max int) []int {
	return lo.Times(count, func(_ int) int {
		return Selector.Intn(max-min) + min
	})
}

func GetDummiesFloat(count int, min, max float64) []float64 {
	return lo.Times(count, func(_ int) float64 {
		min := decimal.NewFromFloat(min)
		max := decimal.NewFromFloat(max)
		result, _ := decimal.NewFromFloat(Selector.Float64()).
			Mul(max.Sub(min)).
			Add(min).
			RoundDown(floatRoundDownScale).
			Float64()
		return result
	})
}

func GetDummiesBool(count int) []bool {
	return lo.Times(count, func(_ int) bool {
		return Selector.Intn(2) == 1
	})
}

func GetDummiesDatetime(count int, min, max time.Time) []time.Time {
	return lo.Times(count, func(_ int) time.Time {
		return time.Unix(Selector.Int63n(max.Unix()-min.Unix())+min.Unix(), 0)
	})
}

func GetDummiesMultilineString(count int, length int, lines int) []string {
	return lo.Times(count, func(_ int) string {
		ls := lo.Times(lines, func(_ int) string {
			s := lo.Times(length, func(_ int) rune {
				return lo.LettersCharset[Selector.Intn(len(lo.LettersCharset))]
			})
			return string(s)
		})
		return strings.Join(ls, "\n")
	})
}

func GetDummiesURL(count int, length int) []string {
	return lo.Times(count, func(_ int) string {
		s := lo.Times(length, func(_ int) rune {
			return lo.LowerCaseLettersCharset[Selector.Intn(len(lo.LowerCaseLettersCharset))]
		})
		return "https://" + string(s) + ".dummy.com"
	})
}

func Select[T any](choices []T) T {
	return choices[Selector.Intn(len(choices))]
}
