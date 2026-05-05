package resources

import (
	"context"
	"fmt"
	"strings"

	caido "github.com/caido-community/sdk-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerReplaySessionResource(server *mcp.Server, client *caido.Client) {
	server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "caido://replay-sessions/{id}",
			Name:        "caido-replay-session",
			Description: "Replay session with entry list and active entry info",
			MIMEType:    "text/plain",
		},
		replaySessionHandler(client),
	)
}

func replaySessionHandler(client *caido.Client) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		id := extractIDFromURI(req.Params.URI, "caido://replay-sessions/")
		if id == "" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		resp, err := client.Replay.GetSession(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get replay session %s: %w", id, err)
		}
		if resp.ReplaySession == nil {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		s := resp.ReplaySession
		var b strings.Builder
		fmt.Fprintf(&b, "# Replay Session: %s\n", s.Name)
		fmt.Fprintf(&b, "ID: %s\n", s.Id)

		if s.ActiveEntry != nil {
			fmt.Fprintf(&b, "Active Entry: %s\n", s.ActiveEntry.Id)
		}

		fmt.Fprintf(&b, "Collection: %s (%s)\n", s.Collection.Name, s.Collection.Id)

		if len(s.Entries.Edges) > 0 {
			fmt.Fprintf(&b, "\n## Entries (%d)\n", len(s.Entries.Edges))
			for _, edge := range s.Entries.Edges {
				e := edge.Node
				fmt.Fprintf(&b, "- %s | %s:%d (tls=%t)\n",
					e.Id, e.Connection.Host, e.Connection.Port, e.Connection.IsTLS,
				)
			}
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:  req.Params.URI,
				Text: b.String(),
			}},
		}, nil
	}
}
