package tools

import (
	"context"
	"fmt"

	gql "github.com/Khan/genqlient/graphql"
	caido "github.com/caido-community/sdk-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// The Go SDK v0.5.0 does not wrap WebSocket (stream) queries, so we issue
// raw GraphQL via client.GraphQL.MakeRequest -- the same approach used by
// create_tamper_rule.go. Field names and the WS protocol enum are taken
// from the bundled graphql/schema.graphql (type Stream, enum StreamProtocol).

type listWsStreamsVars struct {
	First    int     `json:"first"`
	After    *string `json:"after,omitempty"`
	Protocol string  `json:"protocol"`
	ScopeID  *string `json:"scopeId,omitempty"`
}

type listWsStreamsResp struct {
	Streams struct {
		Edges []struct {
			Cursor string `json:"cursor"`
			Node   struct {
				ID        string `json:"id"`
				Host      string `json:"host"`
				Port      int    `json:"port"`
				Path      string `json:"path"`
				IsTls     bool   `json:"isTls"`
				Direction string `json:"direction"`
				Source    string `json:"source"`
				CreatedAt int64  `json:"createdAt"`
			} `json:"node"`
		} `json:"edges"`
		PageInfo struct {
			HasNextPage bool    `json:"hasNextPage"`
			EndCursor   *string `json:"endCursor"`
		} `json:"pageInfo"`
	} `json:"streams"`
}

const listWsStreamsQuery = `
query ListWsStreams($first: Int!, $after: String, $protocol: StreamProtocol!, $scopeId: ID) {
	streams(first: $first, after: $after, protocol: $protocol, scopeId: $scopeId) {
		edges {
			cursor
			node { id host port path isTls direction source createdAt }
		}
		pageInfo { hasNextPage endCursor }
	}
}`

// ListWsStreamsInput is the input for the list_ws_streams tool
type ListWsStreamsInput struct {
	Limit   int    `json:"limit,omitempty" jsonschema:"Maximum number of streams to return (default 20, max 100)"`
	After   string `json:"after,omitempty" jsonschema:"Cursor for pagination from previous response nextCursor"`
	ScopeID string `json:"scope_id,omitempty" jsonschema:"Optional scope ID to filter streams"`
}

// WsStreamSummary is a minimal representation of a WebSocket stream
type WsStreamSummary struct {
	ID        string `json:"id"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Path      string `json:"path"`
	IsTLS     bool   `json:"isTls"`
	Direction string `json:"direction"`
	Source    string `json:"source"`
	CreatedAt int64  `json:"createdAt"`
}

// ListWsStreamsOutput is the output of the list_ws_streams tool
type ListWsStreamsOutput struct {
	Streams    []WsStreamSummary `json:"streams"`
	HasMore    bool              `json:"hasMore"`
	NextCursor string            `json:"nextCursor,omitempty"`
}

func listWsStreamsHandler(
	client *caido.Client,
) func(context.Context, *mcp.CallToolRequest, ListWsStreamsInput) (*mcp.CallToolResult, ListWsStreamsOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input ListWsStreamsInput,
	) (*mcp.CallToolResult, ListWsStreamsOutput, error) {
		limit := input.Limit
		if limit <= 0 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}

		vars := &listWsStreamsVars{First: limit, Protocol: "WS"}
		if input.After != "" {
			vars.After = &input.After
		}
		if input.ScopeID != "" {
			vars.ScopeID = &input.ScopeID
		}

		gqlReq := &gql.Request{
			OpName:    "ListWsStreams",
			Query:     listWsStreamsQuery,
			Variables: vars,
		}
		data := &listWsStreamsResp{}
		if err := client.GraphQL.MakeRequest(
			ctx, gqlReq, &gql.Response{Data: data},
		); err != nil {
			return nil, ListWsStreamsOutput{}, fmt.Errorf(
				"failed to list ws streams: %w", err,
			)
		}

		conn := data.Streams
		output := ListWsStreamsOutput{
			Streams: make([]WsStreamSummary, 0, len(conn.Edges)),
		}
		for _, edge := range conn.Edges {
			n := edge.Node
			output.Streams = append(output.Streams, WsStreamSummary{
				ID:        n.ID,
				Host:      n.Host,
				Port:      n.Port,
				Path:      n.Path,
				IsTLS:     n.IsTls,
				Direction: n.Direction,
				Source:    n.Source,
				CreatedAt: n.CreatedAt,
			})
		}
		if conn.PageInfo.HasNextPage {
			output.HasMore = true
			if conn.PageInfo.EndCursor != nil {
				output.NextCursor = *conn.PageInfo.EndCursor
			}
		}

		return nil, output, nil
	}
}

// RegisterListWsStreamsTool registers the tool with the MCP server
func RegisterListWsStreamsTool(server *mcp.Server, client *caido.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "caido_list_ws_streams",
		Description: `List WebSocket streams (connections) from the WebSocket tab. ` +
			`Returns id/host/port/path/direction/source. ` +
			`Use the stream id with caido_list_ws_messages to read frames. ` +
			`Params: limit (default 20, max 100), after (cursor), scope_id (optional).`,
	}, listWsStreamsHandler(client))
}
