package cmd

import (
	"bufio"
	"bytes"
	goctx "context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/axone-protocol/axone-mcp/internal/mocks"
	"github.com/axone-protocol/axone-mcp/internal/version"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

const testTimeout = 5 * time.Second

func TestServeStdioCommand(t *testing.T) {
	Convey("Testing Serve Stdio command", t, func() {
		tests := []struct {
			name     string
			args     []string
			input    string
			expected string
		}{
			{
				name:     "Ping",
				args:     []string{"serve", "stdio"},
				input:    `{"jsonrpc": "2.0", "id": 42, "method": "ping", "params": {}}`,
				expected: `{"jsonrpc":"2.0","id":42,"result":{}}`,
			},
		}
		for _, tt := range tests {
			Convey(fmt.Sprintf("Given a new server executed by serve stdio command for %s", tt.name),
				withCommandArguments([]string{"serve", "stdio"},
					withPipedIOStreams(func(c C, stdinW io.Writer, stdoutR io.Reader, stderrR io.Reader) {
						go func() {
							ctrl := gomock.NewController(t)
							c.Reset(ctrl.Finish)

							cc := mocks.NewMockClientConnInterface(ctrl)
							ctx := WithGrpcClientConn(goctx.Background(), cc)
							Execute(ctx)
						}()

						Convey(fmt.Sprintf("When sending input: %s", tt.input),
							func(c C) {
								go func() {
									writer := bufio.NewWriter(stdinW)
									_, err := fmt.Fprintf(writer, "%s\n", tt.input)
									c.So(err, ShouldBeNil)
									err = writer.Flush()
									c.So(err, ShouldBeNil)
								}()

								Convey(fmt.Sprintf("Then the response should be: %s", tt.expected), func(c C) {
									var got string
									done := make(chan struct{})
									go func() {
										scanner := bufio.NewScanner(stdoutR)
										scanner.Scan()
										c.So(scanner.Err(), ShouldBeNil)
										got = scanner.Text()
										close(done)
									}()

									select {
									case <-done:
									case <-time.After(testTimeout):
										t.Fatalf("timeout")
									}

									So(got, shouldJSONEqual, tt.expected)
								})
							})
					})))
		}
	})
}

func TestInvalidServeStdioCommand(t *testing.T) {
	Convey("Testing Invalid Serve Stdio command", t,
		withCommandArguments([]string{"serve", "stdio", "--node-grpc", "%1"},
			withPipedIOStreams(func(c C, stdinW io.Writer, stdoutR io.Reader, stderrR io.Reader) {
				Convey("When launching the command with invalid arguments", func() {
					ctx, cancel := goctx.WithTimeout(goctx.Background(), testTimeout)
					Reset(func() {
						cancel()
					})
					got := serveStdioCmd.ExecuteContext(ctx)

					Convey("Then the command should return an error", func() {
						c.So(got, ShouldBeError, `parse "dns:///%1": invalid URL escape "%1"`)
					})
				})
			})))
}

func withCommandArguments(args []string, f func(c C)) func(c C) {
	return func(c C) {
		origArgs := os.Args

		Reset(func() {
			os.Args = origArgs
		})

		os.Args = append([]string{version.Name}, args...)

		f(c)
	}
}

func withPipedIOStreams(f func(c C, stdinW io.Writer, stdoutR io.Reader, stderr io.Reader)) func(c C) {
	return func(c C) {
		stdinReader, stdinWriter := io.Pipe()
		stdoutReader, stdoutWriter := io.Pipe()
		stderrReader, stderrWriter := io.Pipe()

		origStdin := MCPStdin
		origStdout := MCPStdout
		origStderr := MCPStderr

		Reset(func() {
			MCPStdin = origStdin
			MCPStdout = origStdout
			MCPStderr = origStderr

			So(stdinWriter.Close(), ShouldBeNil)
			So(stdoutWriter.Close(), ShouldBeNil)
			So(stderrWriter.Close(), ShouldBeNil)
		})

		MCPStdin = stdinReader
		MCPStdout = stdoutWriter
		MCPStderr = stderrWriter

		f(c, stdinWriter, stdoutReader, stderrReader)
	}
}

func shouldJSONEqual(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("This assertion requires exactly %d comparison values (you provided %d).", 1, len(expected))
	}

	left, leftIsString := actual.(string)
	right, rightIsString := expected[0].(string)

	if !leftIsString || !rightIsString {
		return fmt.Sprintf("Both arguments to this assertion must be strings (you provided %v and %v).", reflect.TypeOf(actual), reflect.TypeOf(expected[0]))
	}

	var leftNormalized, rightNormalized bytes.Buffer
	ShouldBeNil(json.Compact(&leftNormalized, []byte(left)))
	ShouldBeNil(json.Compact(&rightNormalized, []byte(right)))

	return ShouldEqual(leftNormalized.String(), rightNormalized.String())
}
