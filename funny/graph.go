package funny

import "fmt"

type function func(variables ...float64) bool

func Draw(f function) error {
	size, err := GetTermSize()
	if err != nil {
		return err
	}
	var (
		maxRow = size.Row
		maxCol = size.Col
	)

	// 打印爱心形状
	for row := maxRow / 2; row > 0-maxRow/2; row-- {
		for col := 0 - maxCol/2; col < maxCol/2; col++ {
			if f(float64(col), float64(row)) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
	}
	return nil
}
