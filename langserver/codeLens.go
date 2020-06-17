// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.  // You may obtain a copy of the License at
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

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
)

// CodeLens is required by the protocol.Server interface.
func (s *server) CodeLens(_ context.Context, _ *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	// As of version 0.4.0 of gopls it is not possible to instruct the language
	// client to stop asking for Code Lenses and Document Links. To prevent
	// VS Code from showing error messages, this feature is implemented by
	// returning empty values.
	return nil, nil
}
