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

package langserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/prometheus-community/promql-langserver/config"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/jsonrpc2"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
	"github.com/prometheus-community/promql-langserver/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNotImplemented checks whether unimplemented functions return the approbiate Error.
func TestNotImplemented(t *testing.T) { // nolint: gocognit, funlen, gocyclo
	s := &server{}

	err := s.DidChangeWorkspaceFolders(context.Background(), &protocol.DidChangeWorkspaceFoldersParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.DidSave(context.Background(), &protocol.DidSaveTextDocumentParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.WillSave(context.Background(), &protocol.WillSaveTextDocumentParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.DidChangeWatchedFiles(context.Background(), &protocol.DidChangeWatchedFilesParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.Progress(context.Background(), &protocol.ProgressParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.SelectionRange(context.Background(), &protocol.SelectionRangeParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.SetTraceNotification(context.Background(), &protocol.SetTraceParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.LogTraceNotification(context.Background(), &protocol.LogTraceParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Implementation(context.Background(), &protocol.ImplementationParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.TypeDefinition(context.Background(), &protocol.TypeDefinitionParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentColor(context.Background(), &protocol.DocumentColorParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ColorPresentation(context.Background(), &protocol.ColorPresentationParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.FoldingRange(context.Background(), &protocol.FoldingRangeParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.NonstandardRequest(context.Background(), "", nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Declaration(context.Background(), &protocol.DeclarationParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.WillSaveWaitUntil(context.Background(), &protocol.WillSaveTextDocumentParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Resolve(context.Background(), &protocol.CompletionItem{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Definition(context.Background(), &protocol.DefinitionParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.References(context.Background(), &protocol.ReferenceParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentHighlight(context.Background(), &protocol.DocumentHighlightParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentSymbol(context.Background(), &protocol.DocumentSymbolParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.CodeAction(context.Background(), &protocol.CodeActionParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Symbol(context.Background(), &protocol.WorkspaceSymbolParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.CodeLens(context.Background(), &protocol.CodeLensParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ResolveCodeLens(context.Background(), &protocol.CodeLens{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Formatting(context.Background(), &protocol.DocumentFormattingParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.RangeFormatting(context.Background(), &protocol.DocumentRangeFormattingParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.OnTypeFormatting(context.Background(), &protocol.DocumentOnTypeFormattingParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Rename(context.Background(), &protocol.RenameParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.PrepareRename(context.Background(), &protocol.PrepareRenameParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentLink(context.Background(), &protocol.DocumentLinkParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ResolveDocumentLink(context.Background(), &protocol.DocumentLink{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ExecuteCommand(context.Background(), &protocol.ExecuteCommandParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.IncomingCalls(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.OutgoingCalls(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.PrepareCallHierarchy(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.SemanticTokens(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.SemanticTokensEdits(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.SemanticTokensRange(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.WorkDoneProgressCancel(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.WorkDoneProgressCreate(context.Background(), nil)
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		require.Fail(t, "Expected a jsonrpc2 Error with CodeMethodNotFound")
	}
}

// dummyStream is a fake jsonrpc2.Stream for Test purposes.
type dummyStream struct {
	readQueue []byte
}

func (d *dummyStream) Read(_ context.Context) ([]byte, int64, error) {
	ret := d.readQueue
	d.readQueue = []byte{}

	return ret, int64(len(ret)), nil
}

func (d *dummyStream) Write(_ context.Context, text []byte) (int64, error) {
	return int64(len(text)), nil
}

// Push adds a text to the readQueue.
func (d *dummyStream) Push(text []byte) {
	d.readQueue = append(d.readQueue, text...)
}

type dummyWriter struct{}

func (d *dummyWriter) Write(text []byte) (int, error) {
	return len(text), nil
}

// TestServerState tries to emulate a full server lifetime.
func TestServer(t *testing.T) { //nolint:funlen, gocognit, gocyclo
	var stream jsonrpc2.Stream = &dummyStream{}
	stream = jSONLogStream(stream, &dummyWriter{})
	_, server := ServerFromStream(context.Background(), stream, &config.Config{LogFormat: config.TextFormat})
	s := mustServerTest(t, server.server)

	// Initialize Server
	_, err := s.server.Initialize(context.Background(), &protocol.ParamInitialize{})
	if err != nil {
		require.Fail(t, "Failed to initialize Server")
	}

	_, err = s.server.Initialize(context.Background(), &protocol.ParamInitialize{})
	if err == nil {
		require.Fail(t, "cannot initialize server twice")
	}
	// Confirm Initialisation
	err = s.server.Initialized(context.Background(), &protocol.InitializedParams{})
	if err != nil {
		require.Fail(t, "Failed to initialize Server")
	}

	err = s.server.Initialized(context.Background(), &protocol.InitializedParams{})
	if err == nil {
		require.Fail(t, "cannot confirm server initialisation twice")
	}

	// Add a document to the server
	err = s.server.DidOpen(context.Background(), &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        s.doc.DocumentURI(),
			Version:    s.doc.NextVersion(),
			LanguageID: "promql",
			Text:       "",
		},
	})
	if err != nil {
		require.Fail(t, "Failed to open document")
	}

	// Apply a Full Change to the document
	err = s.server.DidChange(context.Background(), &protocol.DidChangeTextDocumentParams{
		TextDocument: s.doc.NextVersionID(),
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Range:       nil,
				RangeLength: 0,
				Text:        "sum()",
			},
		},
	})
	if err != nil {
		require.Fail(t, "Failed to apply full change to document")
	}

	hover, err := s.server.Hover(context.Background(), &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: s.doc.ID(),
			Position: protocol.Position{
				Line:      0.0,
				Character: 1.0,
			},
		},
	})

	if err != nil {
		require.Fail(t, "Failed to get hovertext")
	}

	if hover == nil || strings.Contains("sum", hover.Contents.Value) {
		fmt.Println(hover)
		require.Fail(t, "unexpected or no hovertext")
	}

	// Apply a Full Change to the document
	err = s.server.DidChange(context.Background(), &protocol.DidChangeTextDocumentParams{
		TextDocument: s.doc.NextVersionID(),
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Range:       nil,
				RangeLength: 0,
				Text:        "metric_name",
			},
		},
	})
	if err != nil {
		require.Fail(t, "Failed to apply full change to document")
	}

	hover, err = s.server.Hover(context.Background(), &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: s.doc.ID(),
			Position: protocol.Position{
				Line:      0.0,
				Character: 1.0,
			},
		},
	})

	if err != nil {
		require.Fail(t, "Failed to get hovertext")
	}

	if hover == nil || strings.Contains("metric_name", hover.Contents.Value) {
		fmt.Println(hover)
		require.Fail(t, "unexpected or no hovertext")
	}
	// Apply a partial Change to the document
	err = s.server.DidChange(context.Background(), &protocol.DidChangeTextDocumentParams{
		TextDocument: s.doc.NextVersionID(),
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Range: &protocol.Range{
					Start: protocol.Position{
						Line:      0.0,
						Character: 0.0,
					},
					End: protocol.Position{
						Line:      0.0,
						Character: 0.0,
					},
				},
				RangeLength: 5,
				Text:        "rate(",
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to apply change to document: %s", err.Error()))
	}

	// Wait for diagnostics
	doc, err := s.server.cache.GetDocument("test.promql")
	if err != nil {
		require.Fail(t, "Failed to get document")
	}

	if diagnostics, err := doc.GetDiagnostics(); err != nil && len(diagnostics) != 0 {
		require.Fail(t, "expected nonempty diagnostics")
	}

	// Apply a partial Change to the document
	err = s.server.DidChange(context.Background(), &protocol.DidChangeTextDocumentParams{
		TextDocument: s.doc.NextVersionID(),
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Range: &protocol.Range{
					Start: protocol.Position{
						Line:      0.0,
						Character: 11.0,
					},
					End: protocol.Position{
						Line:      0.0,
						Character: 16.0,
					},
				},
				RangeLength: 5,
				Text:        ")",
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to apply change to document: %s", err.Error()))
	}

	// Wait for diagnostics
	doc, err = s.server.cache.GetDocument("test.promql")
	if err != nil {
		require.Fail(t, "Failed to get document")
	}

	if diagnostics, err := doc.GetDiagnostics(); err != nil && len(diagnostics) != 0 {
		require.Fail(t, "expected empty diagnostics")
	}

	var content string

	content, err = doc.GetContent()
	if err != nil {
		require.Fail(t, "failed to get document content")
	}

	if content != "rate(metric)" {
		panic(fmt.Sprintf("unexpected content, expected \"rate(metric)\", got %s", content))
	}

	// Apply a Full Change to the document
	err = s.server.DidChange(context.Background(), &protocol.DidChangeTextDocumentParams{
		TextDocument: s.doc.NextVersionID(),
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Range:       nil,
				RangeLength: 0,
				Text:        "rat",
			},
		},
	})
	if err != nil {
		require.Fail(t, "Failed to apply full change to document")
	}

	completion, err := s.server.Completion(context.Background(), &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: s.doc.ID(),
			Position: protocol.Position{
				Line:      0.0,
				Character: 1.0,
			},
		},
	})

	if err != nil || completion == nil || len(completion.Items) == 0 || completion.Items[0].Label != "rate" {
		fmt.Println(completion)
		require.Fail(t, "Failed to get completion")
	}

	// Apply a Full Change to the document
	err = s.server.DidChange(context.Background(), &protocol.DidChangeTextDocumentParams{
		TextDocument: s.doc.NextVersionID(),
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Range:       nil,
				RangeLength: 0,
				Text:        "rat()",
			},
		},
	})
	if err != nil {
		require.Fail(t, "Failed to apply full change to document")
	}

	completion, err = s.server.Completion(context.Background(), &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: s.doc.ID(),
			Position: protocol.Position{
				Line:      0.0,
				Character: 1.0,
			},
		},
	})

	if err != nil || completion == nil || len(completion.Items) == 0 || completion.Items[0].Label != "rate" {
		fmt.Println(completion)
		require.Fail(t, "Failed to get completion")
	}

	// Close a document
	err = s.server.DidClose(context.Background(), &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: s.doc.DocumentURI(),
		},
	})
	if err != nil {
		require.Fail(t, "Failed to close document")
	}

	_, err = s.server.cache.GetDocument("test.promql")
	if err == nil {
		require.Fail(t, "getting a closed document should have failed")
	}

	// Close a document twice
	err = s.server.DidClose(context.Background(), &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: s.doc.DocumentURI(),
		},
	})
	if err == nil {
		require.Fail(t, "should have failed to close document")
	}

	// Reopen a closed document
	s.doc.ResetVersion()
	err = s.server.DidOpen(context.Background(), &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        s.doc.DocumentURI(),
			Version:    s.doc.NextVersion(),
			LanguageID: "promql",
			Text:       "abs()",
		},
	})
	if err != nil {
		require.Fail(t, "Failed to reopen document")
	}

	signature, err := s.server.SignatureHelp(context.Background(), &protocol.SignatureHelpParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: s.doc.ID(),
			Position: protocol.Position{
				Line:      1.0,
				Character: 0.0,
			},
		},
	})

	if err != nil {
		require.Fail(t, "Failed to get signature")
	}

	if signature != nil && len(signature.Signatures) != 0 {
		fmt.Println(signature)
		require.Fail(t, "Wrong number of signatures returned")
	}

	signature, err = s.server.SignatureHelp(context.Background(), &protocol.SignatureHelpParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: s.doc.ID(),
			Position: protocol.Position{
				Line:      0,
				Character: 4,
			},
		},
	})

	if err != nil {
		require.Fail(t, "Failed to get signature")
	}

	if signature == nil || len(signature.Signatures) != 1 {
		fmt.Println(signature.Signatures)
		require.Fail(t, "Wrong number of signatures returned")
	}

	hover, err = s.server.Hover(context.Background(), &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: s.doc.ID(),
			Position: protocol.Position{
				Line:      0.0,
				Character: 1.0,
			},
		},
	})

	if err != nil {
		require.Fail(t, "Failed to get hovertext")
	}

	if hover == nil || strings.Contains("abs", hover.Contents.Value) {
		fmt.Println(hover)
		require.Fail(t, "unexpected or no hovertext")
	}

	// Run completion metadata tests.
	t.Run("completion label name: sum(metric_name{})", func(t *testing.T) {
		// Apply a Full Change to the document.
		err := s.server.DidChange(context.Background(),
			&protocol.DidChangeTextDocumentParams{
				TextDocument: s.doc.NextVersionID(),
				ContentChanges: []protocol.TextDocumentContentChangeEvent{
					{
						Range:       nil,
						RangeLength: 0,
						Text:        "sum(metric_name{})",
					},
				},
			})
		require.NoError(t, err, "Failed to apply full change to document")

		// Simulate completions for metric_name.
		metaServer := s.SetupTestMetaServer()
		defer metaServer.TearDown()

		metaServer.HandleFunc("/api/v1/series",
			func(w http.ResponseWriter, r *http.Request) {
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, []string{"metric_name"}, r.Form["match[]"])

				json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "success",
					"data": []interface{}{
						map[string]string{
							"foo_name": "foo_value",
							"bar_name": "bar_value",
						},
					},
				})
			})

		completion, err := s.server.Completion(context.Background(),
			&protocol.CompletionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: s.doc.ID(),
					Position: protocol.Position{
						Line:      0.0,
						Character: 16.0,
					},
				},
			})
		require.NoError(t, err, "Failed to get completion")
		expected := []string{
			"bar_name",
			"foo_name",
		}
		actual := completionValuesSorted(t, completion)
		require.Equal(t, expected, actual)
	})

	t.Run("completion label name: sum(metric_name{foo_name=\"foo_value\",})", func(t *testing.T) {
		// Apply a Full Change to the document.
		err := s.server.DidChange(context.Background(),
			&protocol.DidChangeTextDocumentParams{
				TextDocument: s.doc.NextVersionID(),
				ContentChanges: []protocol.TextDocumentContentChangeEvent{
					{
						Range:       nil,
						RangeLength: 0,
						Text:        "sum(metric_name{foo_name=\"foo_value\",})",
					},
				},
			})
		require.NoError(t, err, "Failed to apply full change to document")

		// Simulate completions for metric_name.
		metaServer := s.SetupTestMetaServer()
		defer metaServer.TearDown()

		metaServer.HandleFunc("/api/v1/series",
			func(w http.ResponseWriter, r *http.Request) {
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, []string{"metric_name{foo_name=\"foo_value\"}"}, r.Form["match[]"])

				json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "success",
					"data": []interface{}{
						map[string]string{
							"foo_name": "foo_value",
							"bar_name": "bar_value",
						},
						map[string]string{
							"foo_name": "foo_value",
							"baz_name": "baz_value",
						},
					},
				})
			})

		completion, err := s.server.Completion(context.Background(),
			&protocol.CompletionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: s.doc.ID(),
					Position: protocol.Position{
						Line:      0.0,
						Character: 37.0,
					},
				},
			})
		require.NoError(t, err, "Failed to get completion")
		expected := []string{
			"bar_name",
			"baz_name",
		}
		actual := completionValuesSorted(t, completion)
		require.Equal(t, expected, actual)
	})

	t.Run("completion label name: sum(metric_name{foo_name=~\"foo_value\",baz_name=\"\"})", func(t *testing.T) {
		// Apply a Full Change to the document.
		err := s.server.DidChange(context.Background(),
			&protocol.DidChangeTextDocumentParams{
				TextDocument: s.doc.NextVersionID(),
				ContentChanges: []protocol.TextDocumentContentChangeEvent{
					{
						Range:       nil,
						RangeLength: 0,
						Text:        "sum(metric_name{foo_name=~\"foo_value\",baz_name=\"\"})",
					},
				},
			})
		require.NoError(t, err, "Failed to apply full change to document")

		// Simulate completions for metric_name.
		metaServer := s.SetupTestMetaServer()
		defer metaServer.TearDown()

		metaServer.HandleFunc("/api/v1/series",
			func(w http.ResponseWriter, r *http.Request) {
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, []string{"metric_name{foo_name=~\"foo_value\"}"}, r.Form["match[]"])

				json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "success",
					"data": []interface{}{
						map[string]string{
							"foo_name": "foo_value",
							"bar_name": "bar_value",
						},
						map[string]string{
							"foo_name": "foo_value",
							"baz_name": "baz_value",
						},
						map[string]string{
							"foo_name": "foo_value",
							"baz_name": "baz_value2",
						},
					},
				})
			})

		completion, err := s.server.Completion(context.Background(),
			&protocol.CompletionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: s.doc.ID(),
					Position: protocol.Position{
						Line:      0.0,
						Character: 48.0,
					},
				},
			})
		require.NoError(t, err, "Failed to get completion")
		expected := []string{
			"\"baz_value\"",
			"\"baz_value2\"",
		}
		actual := completionValuesSorted(t, completion)
		require.Equal(t, expected, actual)
	})

	t.Run("completion label name: sum(metric_name{foo_name=~\"foo_value\"}) by ()", func(t *testing.T) {
		// Apply a Full Change to the document.
		err := s.server.DidChange(context.Background(),
			&protocol.DidChangeTextDocumentParams{
				TextDocument: s.doc.NextVersionID(),
				ContentChanges: []protocol.TextDocumentContentChangeEvent{
					{
						Range:       nil,
						RangeLength: 0,
						Text:        "sum(metric_name{foo_name=~\"foo_value\"}) by ()",
					},
				},
			})
		require.NoError(t, err, "Failed to apply full change to document")

		// Simulate completions for metric_name.
		metaServer := s.SetupTestMetaServer()
		defer metaServer.TearDown()

		metaServer.HandleFunc("/api/v1/series",
			func(w http.ResponseWriter, r *http.Request) {
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, []string{"metric_name{foo_name=~\"foo_value\"}"}, r.Form["match[]"])

				json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "success",
					"data": []interface{}{
						map[string]string{
							"foo_name": "foo_value",
							"bar_name": "bar_value",
						},
						map[string]string{
							"foo_name": "foo_value",
							"baz_name": "baz_value",
						},
						map[string]string{
							"foo_name": "foo_value",
							"baz_name": "baz_value2",
						},
					},
				})
			})

		completion, err := s.server.Completion(context.Background(),
			&protocol.CompletionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: s.doc.ID(),
					Position: protocol.Position{
						Line:      0.0,
						Character: 44.0,
					},
				},
			})
		require.NoError(t, err, "Failed to get completion")
		expected := []string{
			"bar_name",
			"baz_name",
		}
		actual := completionValuesSorted(t, completion)
		require.Equal(t, expected, actual)
	})

	// Shutdown Server
	err = s.server.Shutdown(context.Background())
	if err != nil {
		require.Fail(t, "Failed to initialize Server")
	}

	err = s.server.Shutdown(context.Background())
	if err == nil {
		require.Fail(t, "cannot shutdown server twice")
	}
	// Left out until it does something else than calling os.Exit()
	// Confirm Shutdown
	err = s.server.Exit(context.Background())
	if err != nil {
		require.Fail(t, "Failed to initialize Server")
	}
}

func completionValuesSorted(
	t *testing.T,
	completions *protocol.CompletionList,
) []string {
	require.NotNil(t, completions)
	results := make([]string, 0, len(completions.Items))
	for _, item := range completions.Items {
		results = append(results, item.Label)
	}
	sort.Strings(results)
	return results
}

type serverTest struct {
	t      *testing.T
	server *server
	doc    *testDocument
}

func mustServerTest(t *testing.T, server *server) *serverTest {
	s := &serverTest{
		t:      t,
		server: server,
		doc: &testDocument{
			name:    "test.promql",
			version: 0,
		},
	}
	s.resetMetadataService()
	return s
}

func (t *serverTest) resetMetadataService() {
	// Default meta service.
	svc, err := prometheus.NewClient("", 5*time.Minute)
	require.NoError(t.t, err, "Failed to initialize metadata service")

	// Reset the metadata service on the server.
	t.server.metadataService = svc
}

func (t *serverTest) SetupTestMetaServer() *testMetadataServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/status/buildinfo",
		func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(&prometheus.BuildInfoResponse{
				Status: "up",
				Data: prometheus.BuildInfoData{
					Version:   prometheus.RequiredVersion.String(),
					Revision:  "HEAD",
					Branch:    "master",
					BuildUser: "prometheus",
					BuildDate: time.Now().String(),
					GoVersion: runtime.Version(),
				},
			})
		})

	server := httptest.NewServer(mux)

	// Create new metadata service.
	svc, err := prometheus.NewClient(server.URL, 5*time.Minute)
	require.NoError(t.t, err, "Failed to initialize fake metadata service")

	// Reset the metadata service on the server.
	t.server.metadataService = svc

	return &testMetadataServer{
		serverTest: t,
		mux:        mux,
		server:     server,
	}
}

type testDocument struct {
	name    string
	version int
}

func (d *testDocument) ID() protocol.TextDocumentIdentifier {
	return protocol.TextDocumentIdentifier{
		URI: d.DocumentURI(),
	}
}

func (d *testDocument) NextVersionID() protocol.VersionedTextDocumentIdentifier {
	return protocol.VersionedTextDocumentIdentifier{
		Version:                d.NextVersion(),
		TextDocumentIdentifier: d.ID(),
	}
}

func (d *testDocument) DocumentURI() protocol.DocumentURI {
	return protocol.DocumentURI(d.name)
}

func (d *testDocument) NextVersion() float64 {
	v := d.version
	d.version++
	return float64(v)
}

func (d *testDocument) ResetVersion() {
	d.version = 0
}

type testMetadataServer struct {
	serverTest *serverTest
	mux        *http.ServeMux
	server     *httptest.Server
}

func (s *testMetadataServer) HandleFunc(
	pattern string,
	handler func(http.ResponseWriter, *http.Request),
) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *testMetadataServer) TearDown() {
	// Close server.
	s.server.Close()
	// Reset metadata service.
	s.serverTest.resetMetadataService()
}
