# Caido MCP Server TODO

- [ ] [2026-04-09] Add race condition support (research Caido API for low-level socket access)
- [ ] [2026-04-09] Fix oneof workarounds in delete_findings.go and export_findings.go when genqlient adds omitempty for nullable list fields
- [x] [2026-06-02] Implement WebSocket history read tools (list_ws_streams, list_ws_messages) via raw GraphQL for next release (done in-house; PR #25's 33-tool dump declined as scope creep)
- [ ] [2026-05-21] Add unit tests for edit_request and export_curl tools
- [ ] [2026-05-21] Add CI release workflow (auto-build + attach binaries on tag push)
- [ ] [2026-05-21] Body format converter tool (JSON/form-urlencoded/XML/multipart) - Caido v0.54.0 feature
- [ ] [2026-05-21] Increase test coverage from 14% to 60%+ (Findings, Intercept, Workflows, Tamper all untested)
- [ ] [2026-05-21] Extract hardcoded limits (50 batch, 20 concurrency, 1MB body) to named constants
