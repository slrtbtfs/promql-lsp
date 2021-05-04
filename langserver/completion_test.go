// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package langserver

import (
	"context"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	promql "github.com/prometheus/prometheus/promql/parser"

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
	"github.com/prometheus-community/promql-langserver/langserver/cache"
)

// This is a mock object which serves as a test helper
type MockMetadataService struct {
	metadata    map[string][]v1.Metadata
	labelNames  []string
	labelValues []model.LabelValue
}

func (m *MockMetadataService) MetricMetadata(ctx context.Context, metric string) (v1.Metadata, error) {
	return m.metadata[metric][0], nil
}

func (m *MockMetadataService) AllMetricMetadata(ctx context.Context) (map[string][]v1.Metadata, error) {
	return m.metadata, nil
}

func (m *MockMetadataService) LabelNames(ctx context.Context, metricName string, startTime time.Time, endTime time.Time) ([]string, error) {
	return m.labelNames, nil
}

func (m *MockMetadataService) LabelValues(ctx context.Context, label string, startTime time.Time, endTime time.Time) ([]model.LabelValue, error) {
	return m.labelValues, nil
}

func (m *MockMetadataService) ChangeDataSource(prometheusURL string) error {
	return nil
}

func (m *MockMetadataService) GetURL() string {
	return "testhost:9090"
}

func TestMetricNameCompletion(t *testing.T) {
	s := &server{
		metadataService: &MockMetadataService{
			metadata: map[string][]v1.Metadata{
				"a1": {
					v1.Metadata{
						Type: v1.MetricTypeCounter,
						Help: "For completion test",
						Unit: "bytes",
					},
				},
				"b2": {
					v1.Metadata{
						Type: v1.MetricTypeGauge,
						Help: "Not selected",
						Unit: "seconds",
					},
				},
				"c3": {
					v1.Metadata{
						Type: v1.MetricTypeUnknown,
						Help: "Not selected",
						Unit: "ratio",
					},
				},
			},
		},
	}

	dc := &cache.DocumentCache{}
	dc.Init()
	doc, err := dc.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{
			URI:        "test.promql",
			LanguageID: "promql",
			Version:    0,
			Text:       "",
		})

	if err != nil {
		t.Fatalf("Error occurred when adding document: %s", err)
	}

	//TODO add compiled quries via dc.compileQuery

	l := &cache.Location{
		Doc: doc,
		Query: &cache.CompiledQuery{
			Pos: 0,
		},
		Node: &promql.Item{},
	}
	items := new([]protocol.CompletionItem)
	err = s.completeMetricName(context.Background(), items, l, "a")

	if err != nil {
		t.Errorf("Error occurred when calling completeMetricName: %s.\n", err)
	}

	if len(*items) != 1 {
		t.Errorf("Expected to have %d items, got: %d.\n", 1, len(*items))
	}
}

func TestLabelNameCompletion(t *testing.T) {
	s := &server{
		metadataService: &MockMetadataService{
			labelNames: []string{
				"a1",
				"a2",
				"b3",
				"c4",
			},
		},
	}

	dc := &cache.DocumentCache{}
	dc.Init()
	doc, err := dc.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{
			URI:        "test.promql",
			LanguageID: "promql",
			Version:    0,
			Text:       "",
		})

	if err != nil {
		t.Fatalf("Error occurred when adding document: %s", err)
	}

	l := &cache.Location{
		Doc: doc,
		Query: &cache.CompiledQuery{
			Pos: 0,
		},
		Node: &promql.Item{
			Val: "a",
		},
	}
	items := new([]protocol.CompletionItem)
	err = s.completeLabel(context.Background(), items, l, nil)

	if err != nil {
		t.Errorf("Error occurred when calling completeLabel: %s.\n", err)
	}

	if len(*items) != 2 {
		t.Errorf("Expected to have %d items, got: %d.\n", 2, len(*items))
	}
}

func TestFunctionNameCompletion(t *testing.T) {
	s := &server{
		metadataService: &MockMetadataService{},
	}

	dc := &cache.DocumentCache{}
	dc.Init()
	doc, err := dc.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{
			URI:        "test.promql",
			LanguageID: "promql",
			Version:    0,
			Text:       "",
		})

	if err != nil {
		t.Fatalf("Error occurred when adding document: %s", err)
	}

	l := &cache.Location{
		Doc: doc,
		Query: &cache.CompiledQuery{
			Pos: 0,
		},
		Node: &promql.Item{},
	}
	items := new([]protocol.CompletionItem)
	name := "std"
	// expect to match 5 functions:
	// - sort_desc
	// - stddev_over_time
	// - stdvar_over_time
	// - stddev
	// - stdvar
	err = s.completeFunctionName(items, l, name)

	if err != nil {
		t.Errorf("Error occurred when calling completeFunctionName: %s.\n", err)
	}

	if len(*items) != 5 {
		results := make([]string, len(*items))
		for i, it := range *items {
			results[i] = it.Label
		}
		t.Errorf("Expected to have %d items, got: %d, are: %v.\n", 5, len(*items), results)
	}
}
