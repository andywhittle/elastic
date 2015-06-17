// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// The function_score allows you to modify the score of documents that
// are retrieved by a query. This can be useful if, for example,
// a score function is computationally expensive and it is sufficient
// to compute the score on a filtered set of documents.
// For more details, see
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-function-score-query.html
type FunctionScoreQuery struct {
	query      Query
	filter     Filter
	boost      *float32
	maxBoost   *float32
	scoreMode  string
	boostMode  string
	functions  []FunctionScoreQueryFunction
	minScore   *float32
}

// NewFunctionScoreQuery creates a new function score query.
func NewFunctionScoreQuery() FunctionScoreQuery {
	return FunctionScoreQuery{
		functions: make([]FunctionScoreQueryFunction, 0),
	}
}

func (q FunctionScoreQuery) Query(query Query) FunctionScoreQuery {
	q.query = query
	q.filter = nil
	return q
}

func (q FunctionScoreQuery) Filter(filter Filter) FunctionScoreQuery {
	q.query = nil
	q.filter = filter
	return q
}

func (q FunctionScoreQuery) Add(filter Filter, scoreFunc ScoreFunction) FunctionScoreQuery {
	q.functions = append(
		q.functions,
		FunctionScoreQueryFunction{
			filter: filter,
			scoreFunction: scoreFunc,
		},
	)
	return q
}

func (q FunctionScoreQuery) AddScoreFunc(scoreFunc ScoreFunction) FunctionScoreQuery {
	q.functions = append(
		q.functions,
		FunctionScoreQueryFunction{
			scoreFunction: scoreFunc,
		},
	)
	return q
}

func (q FunctionScoreQuery) AddScoreFuncWithWeight(scoreFunc ScoreFunction, weight Weight) FunctionScoreQuery {
	q.functions = append(
		q.functions,
		FunctionScoreQueryFunction{
			scoreFunction: scoreFunc,
			weight: weight,
		},
	)
	return q
}

func (q FunctionScoreQuery) ScoreMode(scoreMode string) FunctionScoreQuery {
	q.scoreMode = scoreMode
	return q
}

func (q FunctionScoreQuery) BoostMode(boostMode string) FunctionScoreQuery {
	q.boostMode = boostMode
	return q
}

func (q FunctionScoreQuery) MaxBoost(maxBoost float32) FunctionScoreQuery {
	q.maxBoost = &maxBoost
	return q
}

func (q FunctionScoreQuery) Boost(boost float32) FunctionScoreQuery {
	q.boost = &boost
	return q
}

func (q FunctionScoreQuery) MinScore(minScore float32) FunctionScoreQuery {
	q.minScore = &minScore
	return q
}

// Source returns JSON for the function score query.
func (q FunctionScoreQuery) Source() interface{} {
	source := make(map[string]interface{})
	query := make(map[string]interface{})
	source["function_score"] = query

	if q.query != nil {
		query["query"] = q.query.Source()
	} else if q.filter != nil {
		query["filter"] = q.filter.Source()
	}

	if len(q.functions) == 1 && q.functions[0].filter == nil {
		scoreFunc := q.functions[0].scoreFunction
		query[scoreFunc.Name()] = scoreFunc.Source()
	} else {
		funcs := make([]interface{}, len(q.functions))
		for i, function := range q.functions {
			funcs[i] = function.Source()
		}
		query["functions"] = funcs
	}

	if q.scoreMode != "" {
		query["score_mode"] = q.scoreMode
	}
	if q.boostMode != "" {
		query["boost_mode"] = q.boostMode
	}
	if q.maxBoost != nil {
		query["max_boost"] = *q.maxBoost
	}
	if q.boost != nil {
		query["boost"] = *q.boost
	}
	if q.minScore != nil {
		query["min_score"] = *q.minScore
	}

	return source
}
