//  Copyright (c) 2023 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build vectors
// +build vectors

package searcher

import (
	"context"

	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/scorer"
	index "github.com/blevesearch/bleve_index_api"
)

type KNNSearcher struct {
	field        string
	vector       []float32
	k            int64
	indexReader  index.IndexReader
	vectorReader index.VectorReader
	scorer       *scorer.KNNQueryScorer
	count        uint64
	vd           index.VectorDoc
}

func NewKNNSearcher(ctx context.Context, i index.IndexReader, m mapping.IndexMapping,
	options search.SearcherOptions, field string, vector []float32, k int64,
	boost float64, similarityMetric string) (search.Searcher, error) {
	if vr, ok := i.(index.VectorIndexReader); ok {
		vectorReader, _ := vr.VectorReader(ctx, vector, field, k)

		count, err := i.DocCount()
		if err != nil {
			_ = vectorReader.Close()
			return nil, err
		}

		knnScorer := scorer.NewKNNQueryScorer(vector, field, boost,
			vectorReader.Count(), count, options, similarityMetric)
		return &KNNSearcher{
			indexReader:  i,
			vectorReader: vectorReader,
			field:        field,
			vector:       vector,
			k:            k,
			scorer:       knnScorer,
		}, nil
	}
	return nil, nil
}

func (s *KNNSearcher) Advance(ctx *search.SearchContext, ID index.IndexInternalID) (
	*search.DocumentMatch, error) {
	knnMatch, err := s.vectorReader.Next(s.vd.Reset())
	if err != nil {
		return nil, err
	}

	if knnMatch == nil {
		return nil, nil
	}

	docMatch := s.scorer.Score(ctx, knnMatch)

	return docMatch, nil
}

func (s *KNNSearcher) Close() error {
	return s.vectorReader.Close()
}

func (s *KNNSearcher) Count() uint64 {
	return s.vectorReader.Count()
}

func (s *KNNSearcher) DocumentMatchPoolSize() int {
	return 1
}

func (s *KNNSearcher) Min() int {
	return 0
}

func (s *KNNSearcher) Next(ctx *search.SearchContext) (*search.DocumentMatch, error) {
	knnMatch, err := s.vectorReader.Next(s.vd.Reset())
	if err != nil {
		return nil, err
	}

	if knnMatch == nil {
		return nil, nil
	}

	docMatch := s.scorer.Score(ctx, knnMatch)

	return docMatch, nil
}

func (s *KNNSearcher) SetQueryNorm(qnorm float64) {
	s.scorer.SetQueryNorm(qnorm)
}

func (s *KNNSearcher) Size() int {
	return 0
}

func (s *KNNSearcher) Weight() float64 {
	return s.scorer.Weight()
}
