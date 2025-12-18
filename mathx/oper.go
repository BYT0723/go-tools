package mathx

import "golang.org/x/exp/constraints"

func DivCeil[T constraints.Integer](a, b T) T {
	if b == 0 {
		panic("division by zero")
	}

	q := a / b

	if a%b != 0 {
		// 如果 a 和 b 同号，向上加 1
		if (a < 0) == (b < 0) {
			q++
		}
	}
	return q
}

func DivFloor[T constraints.Integer](a, b T) T {
	if b == 0 {
		panic("division by zero")
	}

	q := a / b

	// 如果余数不为 0，且 a 和 b 异号，需要向下取整（比向零小 1）
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}
