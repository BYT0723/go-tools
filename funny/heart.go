package funny

import (
	"math"
)

func Heart() error {
	return Draw(func(variables ...float64) bool {
		if len(variables) < 2 {
			return false
		}
		xf := variables[0] / 20
		yf := variables[1] / 20
		return math.Pow(math.Pow(xf, 2)+math.Pow(yf, 2)-1, 3)-math.Pow(xf, 2)*math.Pow(yf, 3) <= 0
	})
}
