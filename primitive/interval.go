package primitive

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Float | constraints.Integer
}

type interval[N Number] struct {
	Min N
	Max N
}

func (i *interval[N]) Contains(n N) bool {
	return i.Min <= n && n <= i.Max
}

func (i *interval[N]) Excludes(n N) bool {
	return i.Min > n || n > i.Max
}

type IntervalF struct{ interval[float64] }

func NewIntervalF(min, max float64) IntervalF {
	return IntervalF{
		interval[float64]{
			Min: min,
			Max: max,
		},
	}
}

type IntervalI struct{ interval[int64] }

func NewIntervalI(min, max int64) IntervalI {
	return IntervalI{
		interval[int64]{
			Min: min,
			Max: max,
		},
	}
}
