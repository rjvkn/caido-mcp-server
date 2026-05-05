package tools_test

import (
	"testing"

	"github.com/c0tton-fluff/caido-mcp-server/internal/testutil"
	"github.com/c0tton-fluff/caido-mcp-server/internal/tools"
	caido "github.com/caido-community/sdk-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestBatchSend(t *testing.T) {
	tests := []struct {
		name      string
		args      map[string]any
		setup     func(*testutil.MockHandler)
		wantError bool
	}{
		{
			name:      "rejects empty requests array",
			args:      map[string]any{"requests": []map[string]any{}},
			setup:     func(m *testutil.MockHandler) {},
			wantError: true,
		},
		{
			name: "rejects more than 50 requests",
			args: func() map[string]any {
				reqs := make([]map[string]any, 51)
				for i := range reqs {
					reqs[i] = map[string]any{"label": "x", "raw": "GET / HTTP/1.1\r\nHost: x\r\n\r\n"}
				}
				return map[string]any{"requests": reqs}
			}(),
			setup:     func(m *testutil.MockHandler) {},
			wantError: true,
		},
		{
			name: "rejects request with empty raw",
			args: map[string]any{
				"requests": []map[string]any{
					{"label": "test", "raw": ""},
				},
			},
			setup:     func(m *testutil.MockHandler) {},
			wantError: true,
		},
		{
			name: "rejects request with raw over 1MB",
			args: map[string]any{
				"requests": []map[string]any{
					{"label": "big", "raw": string(make([]byte, 1048577))},
				},
			},
			setup:     func(m *testutil.MockHandler) {},
			wantError: true,
		},
		{
			name: "successful batch sends requests",
			args: map[string]any{
				"requests": []map[string]any{
					{"label": "req-a", "raw": "GET /a HTTP/1.1\r\nHost: example.com\r\n\r\n"},
					{"label": "req-b", "raw": "GET /b HTTP/1.1\r\nHost: example.com\r\n\r\n"},
				},
				"concurrency": 2,
			},
			setup: func(m *testutil.MockHandler) {
				m.On("CreateReplaySession", testutil.CreateReplaySessionResponse("batch-s1"))
				m.On("StartReplayTask", testutil.StartReplayTaskResponse())
				m.On("GetReplaySession", testutil.GetReplaySessionResponse("batch-s1", "be-1"))
				m.On("GetReplayEntry", testutil.GetReplayEntryResponse("be-1", "br-1", 200, "ok"))
				m.On("DeleteReplaySessions", map[string]any{
					"deleteReplaySessions": map[string]any{
						"deletedIds": []string{},
					},
				})
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := testutil.NewMCPTestEnv(t, func(s *mcp.Server, c *caido.Client) {
				tools.RegisterBatchSendTool(s, c)
			})
			tt.setup(env.Mock)

			result := env.CallTool(t, "caido_batch_send", tt.args)

			if tt.wantError {
				if !result.IsError {
					t.Fatal("expected error result")
				}
				return
			}

			output := testutil.UnmarshalToolResult[tools.BatchSendOutput](t, result)
			if len(output.Results) != 2 {
				t.Fatalf("want 2 results, got %d", len(output.Results))
			}
		})
	}
}
