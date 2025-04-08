# Axone MCP

## 1.0.0 (2025-04-08)


### ⚠ BREAKING CHANGES

* **mcp:** drop hello_world tool
* **sse:** refactor cmd hierarchy to group transports under serve

### Features

* **cli:** enable environment variable support for all CLI flags ([55f6a18](https://github.com/axone-protocol/axone-mcp/commit/55f6a182cfdd8a899766e06ead1f6c35cc09054e))
* **cli:** support console/json log output and level control ([8dcc595](https://github.com/axone-protocol/axone-mcp/commit/8dcc595245b710bfcd9a578155da069f385ac90a))
* **cli:** wire up minimal MCP server using SSE transport ([96e32a5](https://github.com/axone-protocol/axone-mcp/commit/96e32a5e0092e378a4020b0dfcdcfd839c189276))
* **mcp:** drop hello_world tool ([0a8c6e5](https://github.com/axone-protocol/axone-mcp/commit/0a8c6e5b9694914bba6eaca2966facc055704387))
* **mcp:** enables logging capabilities for the server ([41a9f1f](https://github.com/axone-protocol/axone-mcp/commit/41a9f1fa05a67da1c7dfe2e4e6c1526da24d6665))
* **mcp:** implement get_resource_governance_code tool ([7e2ac1a](https://github.com/axone-protocol/axone-mcp/commit/7e2ac1ab9d6594fc6a0356b6e4e7cbb85d6d4b0d))
* **mcp:** log session registration events ([28ed13e](https://github.com/axone-protocol/axone-mcp/commit/28ed13e9965280aa2647d4389f00a70218536f77))
* **sse:** refactor cmd hierarchy to group transports under serve ([1c7d13b](https://github.com/axone-protocol/axone-mcp/commit/1c7d13bbfafa7de7ec80784c40e9752aa1fb173b))
* **stdio:** allow injection of stdin/stdout/stderr for testing ([276f1c4](https://github.com/axone-protocol/axone-mcp/commit/276f1c40cc259f1166dc80a8e1a904ec5504f0b4))
* **stdio:** implement stdio transport for MCP server ([b6be665](https://github.com/axone-protocol/axone-mcp/commit/b6be665245c4d707e09e2bbce98a709f56237c65))


### Bug Fixes

* **cli:** correctly decode --log-format auto flag ([b1a5fef](https://github.com/axone-protocol/axone-mcp/commit/b1a5fefa029c51bb960e7464beccaaab16c37392))
* **security:** enforce TLS 1.2 minimum when skipping verification ([6e14030](https://github.com/axone-protocol/axone-mcp/commit/6e140308691390dead8584def5ae672cee7f9e66))
