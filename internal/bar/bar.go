package bar

import "fmt"

const DEFAULT_GRAPH = "#"

type Bar struct {
	percent uint64
	cur     uint64
	total   uint64
	rate    string
	graph   string
}

func New(total uint64) *Bar {
	bar := Bar{
		cur: 0,
		total: total,
		graph: DEFAULT_GRAPH,
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph // initial progress position
	}

	return &bar
}

func (bar *Bar) getPercent() uint64 {
	return uint64((float32(bar.cur) / float32(bar.total)) * 100)
}

func (bar *Bar) Play(cur uint64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Finish() {
	fmt.Println()
}
