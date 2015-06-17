// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

type FunctionScoreQueryFunction struct {
	filter        Filter
	scoreFunction ScoreFunction
	weight        Weighting
}

func (fn FunctionScoreQueryFunction) Source() interface{} {
	source := make(map[string]interface{})
	if fn.filter != nil {
		source["filter"] = fn.filter.Source()
	}
	if fn.scoreFunction != nil {
		source[fn.scoreFunction.Name()] = fn.scoreFunction.Source()
	}
	if fn.weight != nil {
		source["weight"] = fn.weight.Source()
	}
	return source
}

// -- Function weighting --

type Weighting interface {
	Source() interface{}
}

type Weight struct {
	value *float64
}

// NewWeight creates a new function score weighting.
func NewWeight(weight float64) Weight {
	return Weight{value: &weight}
}

func (w Weight) Value(weight float64) Weight {
	w.value = &weight
	return w
}

func (w Weight) Source() interface{} {
	return *w.value
}
