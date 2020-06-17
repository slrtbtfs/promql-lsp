// Copyright 2019 Tobias Guggenmos
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

package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
)

func TestCache(t *testing.T) { // nolint:funlen
	c := &DocumentCache{}

	c.Init()

	doc, err := c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{
			URI:        "test_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       "test_text",
		})
	if err != nil {
		panic("Failed to AddDocument() to cache")
	}

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "test_file",
			LanguageID: "yaml",
			Version:    1,
			Text:       "test_text",
		})
	if err == nil {
		panic("Should not be able to add same document twice")
	}

	doc1, err := c.GetDocument("test_file")
	if err != nil {
		panic("Failed to GetDocument() from cache")
	}

	if doc1.doc != doc.doc {
		panic("Cache returned wrong document")
	}

	err = c.RemoveDocument("test_file")
	if err != nil {
		panic("Failed to RemoveDocument() from cache")
	}

	err = c.RemoveDocument("test_file")
	if err == nil {
		panic("should have failed to RemoveDocument() twice")
	}

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "test_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       "test_text:",
		})
	if err != nil {
		panic("Should be able to readd document after removing it")
	}

	wrongYaml := "asdf["

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "wrong_yaml_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       wrongYaml,
		})
	if err != nil {
		panic("Should be able to handle yaml with syntax errors")
	}

	rulesFile := `
groups:
  - name: example
    rules:
    - record: job:http_inprogress_requests:sum
      expr: sum(http_inprogress_requests) by (job)
    - record: job:http_inprogress_requests:sum:wrong
      expr: 
          sum(http_inprogress_requests) by (job
    - record: job:http_inprogress_requests:sum:quoted
      expr: "sum(http_inprogress_requests) by (job)"
    - expr: sum(http_inprogress_requests) by (job)
      record: job:http_inprogress_requests:sum:2
    - expr: |
        sum(http_inprogress_requests) by (job)
      record: |
        job:http_inprogress_requests:sum:3
    - 1
    - a:
      - b :2
`

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "rules_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       rulesFile,
		})
	if err != nil {
		panic("adding rules file failed")
	}

	doc, err = c.GetDocument("rules_file")
	if err != nil {
		panic("failed to get rules file")
	}

	diagnostics, err := doc.GetDiagnostics()
	if err != nil {
		panic("failed to get diagnostics for rules file")
	}

	if len(diagnostics) != 2 {
		fmt.Println(diagnostics)
		panic("expected exactly 2 error messages for rules file got " + fmt.Sprint(len(diagnostics)))
	}

	queries, err := doc.GetQueries()
	if err != nil {
		panic("failed to get queries for rules file")
	}

	if len(queries) != 4 {
		fmt.Println(queries)
		panic("expected exactly 4 queries for rules file got " + fmt.Sprint(len(queries)))
	}

	_, err = doc.getQuery(97)
	if err == nil {
		panic("should have failed to get query")
	}

	_, err = doc.getQuery(98)
	if err != nil {
		panic("failed to get query: " + err.Error())
	}

	_, err = c.Find(&protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "rules_file",
		},
		Position: protocol.Position{
			Line:      5,
			Character: 25,
		},
	})
	if err != nil {
		panic("failed to find query: " + err.Error())
	}

	_, err = c.Find(&protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "rules_file",
		},
		Position: protocol.Position{
			Line:      4,
			Character: 25,
		},
	})
	if err == nil {
		panic("should have failed to find query")
	}

	_, err = c.Find(&protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "rules_file_nonexistent",
		},
		Position: protocol.Position{
			Line:      4,
			Character: 25,
		},
	})
	if err == nil {
		panic("should have failed to find query")
	}

	expectedNewContent := `
groups:
  - name: example
    rules:
    - record: job:http_inprogress_requests:sum
      expr: sum(http_inprogress_requests) by (job)
    - record: job:http_inprogress_requests:sum:wrong
      expr: 
          sum(http_inprogress_requests) by (job
    - expr: sum(http_inprogress_requests) by (job)
      record: job:http_inprogress_requests:sum:2
    - expr: |
        sum(http_inprogress_requests) by (job)
      record: |
        job:http_inprogress_requests:sum:3
    - 1
    - a:
      - b :2
`

	newContent, err := doc.ApplyIncrementalChanges([]protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{
					Line:      9.0,
					Character: 0.0,
				},
				End: EndOfLine(protocol.Position{
					Line:      10.0,
					Character: 0.1,
				}),
			},
			Text: "",
		},
	}, 2)

	if err != nil {
		panic("Failed to apply incremental changes: " + err.Error())
	}

	if newContent != expectedNewContent {
		panic("incremental update did not result in expected content")
	}

	_, err = doc.ApplyIncrementalChanges(nil, -1)
	if err == nil {
		panic("File update without version update should have failed")
	}

	err = doc.SetContent(context.Background(), "foo", 2, false)
	if err != nil {
		panic("file update failed")
	}

	err = doc.SetContent(context.Background(), "foo", 2, false)
	if err == nil {
		panic("File update without version update should have failed")
	}
}
