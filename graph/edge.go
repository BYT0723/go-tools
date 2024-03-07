package graph

type Edge struct {
	// 原点字符
	Zero string
	// X轴
	X string
	// Y轴
	Y string
	// 符合函数的字符
	FC string
	// 刻度线
	Scale string
}

var (
	DefaultEdge = &Edge{Zero: " ", X: " ", Y: " ", FC: "*", Scale: " "}
	StandbyEdge = &Edge{Zero: "+", X: "-", Y: "|", FC: "*", Scale: "+"}
)
