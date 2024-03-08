package functions

import (
	"math"
)

type BinaryFunction func(x, y float64) bool

// 爱心
func Love() BinaryFunction {
	return func(x, y float64) bool {
		return math.Pow(math.Pow(x, 2)+math.Pow(y, 2)-1, 3)-math.Pow(x, 2)*math.Pow(y, 3) <= 0
	}
}

// 循环爱心
func CircularLove() BinaryFunction {
	return func(x, y float64) bool {
		return y <= math.Sqrt(math.Abs(math.Sin(x)))+math.Sqrt(math.Abs(math.Cos(x))) && y >= math.Sqrt(math.Abs(math.Sin(x)))-math.Sqrt(math.Abs(math.Cos(x)))
	}
}

func RoseLine(leaf, scale float64) BinaryFunction {
	return func(x, y float64) bool {
		var (
			theta = math.Atan2(y, x)
			r     = scale * math.Cos(leaf*theta)
			xp    = r * math.Cos(theta)
			yp    = r * math.Sin(theta)
		)
		return (x >= 0 && y >= 0 && x <= xp && y <= yp) ||
			(x <= 0 && y >= 0 && x >= xp && y <= yp) ||
			(x <= 0 && y <= 0 && x >= xp && y >= yp) ||
			(x >= 0 && y <= 0 && x <= xp && y >= yp)
	}
}
