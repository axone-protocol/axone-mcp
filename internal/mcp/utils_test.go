package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	success = ""
)

var requestId = mcp.NewRequestId("42")

// ShouldBeJSONRPCResponseSuccessWithText validates the result of a tool call in a
// JSON-RPC response.
func ShouldBeJSONRPCResponseSuccessWithText(actual any, expected ...any) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	expectedText, ok := expected[0].(string)
	if !ok {
		return fmt.Sprintf("Expected text must be a string, got: %T", expected[0])
	}

	if fail := ShouldNotBeNil(actual); fail != "" {
		return fail
	}

	if fail := ShouldHaveSameTypeAs(actual, mcp.JSONRPCResponse{}); fail != "" {
		return fail
	}

	response := actual.(mcp.JSONRPCResponse)

	return shouldBeToolResultText(response, false, expectedText)
}

// ShouldBeJSONRPCResponseErrorWithText validates a JSON-RPC response object marked as error.
func ShouldBeJSONRPCResponseErrorWithText(actual any, expected ...any) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	expectedText, ok := expected[0].(string)
	if !ok {
		return fmt.Sprintf("Expected text must be a string, got: %T", expected[0])
	}

	if fail := ShouldNotBeNil(actual); fail != "" {
		return fail
	}

	if fail := ShouldHaveSameTypeAs(actual, mcp.JSONRPCResponse{}); fail != "" {
		return fail
	}

	response := actual.(mcp.JSONRPCResponse)

	return shouldBeToolResultText(response, true, expectedText)
}

// ShouldBeJSONRPCErrorWithText validates a JSON-RPC error object.
func ShouldBeJSONRPCErrorWithText(actual any, expected ...any) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	expectedText, ok := expected[0].(string)
	if !ok {
		return fmt.Sprintf("Expected text must be a string, got: %T", expected[0])
	}

	if fail := ShouldNotBeNil(actual); fail != "" {
		return fail
	}

	if fail := ShouldHaveSameTypeAs(actual, mcp.JSONRPCError{}); fail != "" {
		return fail
	}

	response := actual.(mcp.JSONRPCError)
	if fail := ShouldResemble(response.ID, requestId); fail != "" {
		return fmt.Sprintf("ID: %s", fail)
	}

	if fail := ShouldEqual(response.JSONRPC, mcp.JSONRPC_VERSION); fail != "" {
		return fmt.Sprintf("JSONRPC: %s", fail)
	}

	if fail := ShouldEqual(response.Error.Message, expectedText); fail != "" {
		return fmt.Sprintf("Message: %s", fail)
	}

	return success
}

func shouldBeToolResultText(response mcp.JSONRPCResponse, isError bool, expectedContentText string) string {
	if fail := ShouldResemble(response.ID, requestId); fail != "" {
		return fmt.Sprintf("ID: %s", fail)
	}

	if fail := ShouldEqual(response.JSONRPC, mcp.JSONRPC_VERSION); fail != "" {
		return fmt.Sprintf("JSONRPC: %s", fail)
	}

	if fail := ShouldHaveSameTypeAs(response.Result, mcp.CallToolResult{}); fail != "" {
		return fmt.Sprintf("Result: %s", fail)
	}

	ctr := response.Result.(mcp.CallToolResult)
	if fail := ShouldEqual(ctr.IsError, isError); fail != "" {
		return fmt.Sprintf("IsError: %s", fail)
	}

	if fail := ShouldHaveLength(ctr.Content, 1); fail != "" {
		return fmt.Sprintf("Content length: %s", fail)
	}

	if fail := ShouldHaveSameTypeAs(ctr.Content[0], mcp.TextContent{}); fail != "" {
		return fmt.Sprintf("Content type: %s", fail)
	}

	content := ctr.Content[0].(mcp.TextContent)
	if fail := ShouldEqual(content.Type, "text"); fail != "" {
		return fmt.Sprintf("Content type: %s", fail)
	}

	if fail := ShouldEqual(content.Text, expectedContentText); fail != "" {
		return fmt.Sprintf("Content text: %s", fail)
	}

	return success
}

func need(needed int, expected []any) string {
	if len(expected) != needed {
		return fmt.Sprintf("This assertion requires exactly %d comparison values (you provided %d).", needed, len(expected))
	}

	return success
}
