package resources

import (
	caido "github.com/caido-community/sdk-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterAll(server *mcp.Server, client *caido.Client) {
	registerRequestResource(server, client)
	registerReplaySessionResource(server, client)
	registerSitemapResource(server, client)
	registerFindingsResource(server, client)
}
