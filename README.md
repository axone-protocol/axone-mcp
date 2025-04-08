# axone-mcp

> 🤖 [Axone](https://axone.xyz)’s [MCP](https://modelcontextprotocol.io/introduction) server – gateway to the dataverse for AI-powered tools

![MCP Protocol](https://img.shields.io/badge/MCP-Compatible-green?style=for-the-badge)
[![version](https://img.shields.io/github/v/release/axone-protocol/axone-mcp?style=for-the-badge&logo=github)](https://github.com/axone-protocol/axone-mcp/releases)
[![lint](https://img.shields.io/github/actions/workflow/status/axone-protocol/axone-mcp/lint.yml?branch=main&label=lint&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axone-mcp/actions/workflows/lint.yml)
[![build](https://img.shields.io/github/actions/workflow/status/axone-protocol/axone-mcp/build.yml?branch=main&label=build&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axone-mcp/actions/workflows/build.yml)
[![test](https://img.shields.io/github/actions/workflow/status/axone-protocol/axone-mcp/test.yml?branch=main&label=test&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axone-mcp/actions/workflows/test.yml)
[![codecov](https://img.shields.io/codecov/c/github/axone-protocol/axone-mcp?style=for-the-badge&token=6NL9ICGZQS&logo=codecov)](https://codecov.io/gh/axone-protocol/axone-mcp)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge&logo=conventionalcommits)](https://conventionalcommits.org)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg?style=for-the-badge)](https://github.com/semantic-release/semantic-release)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg?style=for-the-badge)](https://github.com/axone-protocol/.github/blob/main/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg?style=for-the-badge)](https://opensource.org/licenses/BSD-3-Clause)

## Axone’s MCP server

[Axone](https://axone.xyz)’s [MCP](https://modelcontextprotocol.io/introduction) server is a lightweight implementation that
exposes Axone’s capabilities through the standardized Model-Context Protocol.

```mermaid
flowchart LR
    classDef actor stroke:#808
    classDef system stroke:#0ff
    classDef resource stroke:#f00

    actor:::actor@{ shape: stadium, label: "Host with MCP Client<br>(Claude, IDEs, Tools)" }
    mcpServer:::system@{ shape: rounded, label: "Axone<br>MCP server" }
    axone:::system@{ shape: das, label: "🔗 Axone chain" }


    actor -- query --> mcpServer

    mcpServer -. query .-> axone
```

## Available tools

### `get_resource_governance_code`

Get the governance code attached to the given resource (if any).

#### Input schema

```json
{
  "resource": {
    "type": "string",
    "description": "The resource DID to get the governance code for."
  }
}
```

## Usage

Install the MCP server:

```sh
go install github.com/axone-protocol/axone-mcp@latest
```

### Usage with [Claude Desktop](https://claude.ai/download)

Add this to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "axone-mcp",
      "args": [
        "serve",
        "stdio"
      ]
    }
  }
}
```

### Run with SSE transport

```sh
axone-mcp serve sse --listen-addr localhost:8080
```

### Run with STDIO transport

```sh
axone-mcp serve stdio
```

## Build

- Be sure you have [Golang](https://go.dev/doc/install) installed.
- [Docker](https://docs.docker.com/engine/install/) as well if you want to use the Makefile.

```sh
make build
```
