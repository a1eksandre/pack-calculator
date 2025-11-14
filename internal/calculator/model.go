package calculator

type PackResult struct {
	Packs      map[int]int // pack size -> count
	TotalItems int         // items shipped in total
	ExtraItems int         // TotalItems - Items requested
}
