// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The linter doesn't see the uses since the consumers are conditionally compiled tests.
//
//nolint:unused,deadcode,varcheck
package ints

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	ptesting "github.com/pulumi/pulumi/sdk/v3/go/common/testing"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/fsutil"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const WindowsOS = "windows"

// assertPerfBenchmark implements the integration.TestStatsReporter interface, and reports test
// failures when a scenario exceeds the provided threshold.
type assertPerfBenchmark struct {
	T                  *testing.T
	MaxPreviewDuration time.Duration
	MaxUpdateDuration  time.Duration
}

func (t assertPerfBenchmark) ReportCommand(stats integration.TestCommandStats) {
	var maxDuration *time.Duration
	if strings.HasPrefix(stats.StepName, "pulumi-preview") {
		maxDuration = &t.MaxPreviewDuration
	}
	if strings.HasPrefix(stats.StepName, "pulumi-update") {
		maxDuration = &t.MaxUpdateDuration
	}

	if maxDuration != nil && *maxDuration != 0 {
		if stats.ElapsedSeconds < maxDuration.Seconds() {
			t.T.Logf(
				"Test step %q was under threshold. %.2fs (max %.2fs)",
				stats.StepName, stats.ElapsedSeconds, maxDuration.Seconds())
		} else {
			t.T.Errorf(
				"Test step %q took longer than expected. %.2fs vs. max %.2fs",
				stats.StepName, stats.ElapsedSeconds, maxDuration.Seconds())
		}
	}
}

func testComponentSlowLocalProvider(t *testing.T) integration.LocalDependency {
	return integration.LocalDependency{
		Package: "testcomponent",
		Path:    filepath.Join("construct_component_slow", "testcomponent"),
	}
}

func testComponentProviderSchema(t *testing.T, path string) {
	t.Parallel()

	runComponentSetup(t, "component_provider_schema")

	tests := []struct {
		name          string
		env           []string
		version       int32
		expected      string
		expectedError string
	}{
		{
			name:     "Default",
			expected: "{}",
		},
		{
			name:     "Schema",
			env:      []string{"INCLUDE_SCHEMA=true"},
			expected: `{"hello": "world"}`,
		},
		{
			name:          "Invalid Version",
			version:       15,
			expectedError: "unsupported schema version 15",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Start the plugin binary.
			cmd := exec.Command(path, "ignored")
			cmd.Env = append(os.Environ(), test.env...)
			stdout, err := cmd.StdoutPipe()
			assert.NoError(t, err)
			err = cmd.Start()
			assert.NoError(t, err)
			defer func() {
				// Ignore the error as it may fail with access denied on Windows.
				cmd.Process.Kill() //nolint:errcheck
			}()

			// Read the port from standard output.
			reader := bufio.NewReader(stdout)
			bytes, err := reader.ReadBytes('\n')
			assert.NoError(t, err)
			port := strings.TrimSpace(string(bytes))

			// Create a connection to the server.
			conn, err := grpc.Dial(
				"127.0.0.1:"+port,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				rpcutil.GrpcChannelOptions(),
			)
			assert.NoError(t, err)
			client := pulumirpc.NewResourceProviderClient(conn)

			// Call GetSchema and verify the results.
			resp, err := client.GetSchema(context.Background(), &pulumirpc.GetSchemaRequest{Version: test.version})
			if test.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedError)
			} else {
				assert.Equal(t, test.expected, resp.GetSchema())
			}
		})
	}
}

// Test remote component inputs properly handle unknowns.
func testConstructUnknown(t *testing.T, lang string, dependencies ...string) {
	t.Parallel()

	const testDir = "construct_component_unknown"
	runComponentSetup(t, testDir)

	tests := []struct {
		componentDir string
	}{
		{
			componentDir: "testcomponent",
		},
		{
			componentDir: "testcomponent-python",
		},
		{
			componentDir: "testcomponent-go",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.componentDir, func(t *testing.T) {
			localProviders := []integration.LocalDependency{
				{Package: "testcomponent", Path: filepath.Join(testDir, test.componentDir)},
			}
			integration.ProgramTest(t, &integration.ProgramTestOptions{
				Dir:                    filepath.Join(testDir, lang),
				Dependencies:           dependencies,
				LocalProviders:         localProviders,
				SkipRefresh:            true,
				SkipPreview:            false,
				SkipUpdate:             true,
				SkipExportImport:       true,
				SkipEmptyPreviewUpdate: true,
				Quick:                  false,
			})
		})
	}
}

// Test methods properly handle unknowns.
func testConstructMethodsUnknown(t *testing.T, lang string, dependencies ...string) {
	t.Parallel()

	const testDir = "construct_component_methods_unknown"
	runComponentSetup(t, testDir)
	tests := []struct {
		componentDir string
	}{
		{
			componentDir: "testcomponent",
		},
		{
			componentDir: "testcomponent-python",
		},
		{
			componentDir: "testcomponent-go",
		},
	}
	for _, test := range tests {
		test := test

		t.Run(test.componentDir, func(t *testing.T) {
			localProviders := []integration.LocalDependency{
				{Package: "testcomponent", Path: filepath.Join(testDir, test.componentDir)},
			}
			integration.ProgramTest(t, &integration.ProgramTestOptions{
				Dir:                    filepath.Join(testDir, lang),
				Dependencies:           dependencies,
				LocalProviders:         localProviders,
				SkipRefresh:            true,
				SkipPreview:            false,
				SkipUpdate:             true,
				SkipExportImport:       true,
				SkipEmptyPreviewUpdate: true,
				Quick:                  false,
			})
		})
	}
}

func runComponentSetup(t *testing.T, testDir string) {
	ptesting.YarnInstallMutex.Lock()
	defer ptesting.YarnInstallMutex.Unlock()

	setupFilename, err := filepath.Abs("component_setup.sh")
	require.NoError(t, err, "could not determine absolute path")
	// Even for Windows, we want forward slashes as bash treats backslashes as escape sequences.
	setupFilename = filepath.ToSlash(setupFilename)

	synchronouslyDo(t, filepath.Join(testDir, ".lock"), 10*time.Minute, func() {
		out := newTestLogWriter(t)

		cmd := exec.Command("bash", "-x", setupFilename)
		cmd.Dir = testDir
		cmd.Stdout = out
		cmd.Stderr = out
		err := cmd.Run()

		// This runs in a separate goroutine, so don't use 'require'.
		assert.NoError(t, err, "failed to run setup script")
	})

	// The function above runs in a separate goroutine
	// so it can't halt test execution.
	// Verify that it didn't fail separately
	// and halt execution if it did.
	require.False(t, t.Failed(), "component setup failed")
}

func synchronouslyDo(t testing.TB, lockfile string, timeout time.Duration, fn func()) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	lockWait := make(chan struct{})
	go func() {
		mutex := fsutil.NewFileMutex(lockfile)

		// ctx.Err will be non-nil when the context finishes
		// either because it timed out or because it got canceled.
		for ctx.Err() == nil {
			if err := mutex.Lock(); err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				defer func() {
					assert.NoError(t, mutex.Unlock())
				}()
				break
			}
		}

		// Context may hav expired
		// by the time we acquired the lock.
		if ctx.Err() == nil {
			fn()
			close(lockWait)
		}
	}()

	select {
	case <-ctx.Done():
		t.Fatalf("timed out waiting for lock on %s", lockfile)
	case <-lockWait:
		// waited for fn, success.
	}
}

// Verifies that if a file lock is already acquired,
// synchronouslyDo is able to time out properly.
func TestSynchronouslyDo_timeout(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "foo")
	mu := fsutil.NewFileMutex(path)
	require.NoError(t, mu.Lock())
	defer func() {
		assert.NoError(t, mu.Unlock())
	}()

	fakeT := nonfatalT{T: t}
	synchronouslyDo(&fakeT, path, 10*time.Millisecond, func() {
		t.Errorf("timed-out operation should not be called")
	})

	assert.True(t, fakeT.fatal, "must have a fatal failure")
	if assert.Len(t, fakeT.messages, 1) {
		assert.Contains(t, fakeT.messages[0], "timed out waiting")
	}
}

// nonfatalT wraps a testing.T to capture fatal errors.
type nonfatalT struct {
	*testing.T

	mu       sync.Mutex
	fatal    bool
	messages []string
}

func (t *nonfatalT) Fatalf(msg string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.fatal = true
	t.messages = append(t.messages, fmt.Sprintf(msg, args...))
}

// Test methods that create resources.
func testConstructMethodsResources(t *testing.T, lang string, dependencies ...string) {
	t.Parallel()

	const testDir = "construct_component_methods_resources"
	runComponentSetup(t, testDir)

	tests := []struct {
		componentDir string
	}{
		{
			componentDir: "testcomponent",
		},
		{
			componentDir: "testcomponent-python",
		},
		{
			componentDir: "testcomponent-go",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.componentDir, func(t *testing.T) {
			localProviders := []integration.LocalDependency{
				{Package: "testcomponent", Path: filepath.Join(testDir, test.componentDir)},
			}
			integration.ProgramTest(t, &integration.ProgramTestOptions{
				Dir:            filepath.Join(testDir, lang),
				Dependencies:   dependencies,
				LocalProviders: localProviders,
				Quick:          true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 6, len(stackInfo.Deployment.Resources))
					var hasExpectedResource bool
					var result string
					for _, res := range stackInfo.Deployment.Resources {
						if res.URN.Name().String() == "myrandom" {
							hasExpectedResource = true
							result = res.Outputs["result"].(string)
							assert.Equal(t, float64(10), res.Inputs["length"])
							assert.Equal(t, 10, len(result))
						}
					}
					assert.True(t, hasExpectedResource)
					assert.Equal(t, result, stackInfo.Outputs["result"])
				},
			})
		})
	}
}

// Test failures returned from methods are observed.
func testConstructMethodsErrors(t *testing.T, lang string, dependencies ...string) {
	t.Parallel()

	const testDir = "construct_component_methods_errors"
	runComponentSetup(t, testDir)

	tests := []struct {
		componentDir string
	}{
		{
			componentDir: "testcomponent",
		},
		{
			componentDir: "testcomponent-python",
		},
		{
			componentDir: "testcomponent-go",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.componentDir, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			expectedError := "the failure reason (the failure property)"

			localProvider := integration.LocalDependency{
				Package: "testcomponent", Path: filepath.Join(testDir, test.componentDir),
			}
			integration.ProgramTest(t, &integration.ProgramTestOptions{
				Dir:            filepath.Join(testDir, lang),
				Dependencies:   dependencies,
				LocalProviders: []integration.LocalDependency{localProvider},
				Quick:          true,
				Stderr:         stderr,
				ExpectFailure:  true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					output := stderr.String()
					assert.Contains(t, output, expectedError)
				},
			})
		})
	}
}

func testConstructOutputValues(t *testing.T, lang string, dependencies ...string) {
	t.Parallel()

	const testDir = "construct_component_output_values"
	runComponentSetup(t, testDir)

	tests := []struct {
		componentDir string
	}{
		{
			componentDir: "testcomponent",
		},
		{
			componentDir: "testcomponent-python",
		},
		{
			componentDir: "testcomponent-go",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.componentDir, func(t *testing.T) {
			localProviders := []integration.LocalDependency{
				{Package: "testcomponent", Path: filepath.Join(testDir, test.componentDir)},
			}
			integration.ProgramTest(t, &integration.ProgramTestOptions{
				Dir:            filepath.Join(testDir, lang),
				Dependencies:   dependencies,
				LocalProviders: localProviders,
				Quick:          true,
			})
		})
	}
}

var previewSummaryRegex = regexp.MustCompile(
	`{\s+"steps": \[[\s\S]+],\s+"duration": \d+,\s+"changeSummary": {[\s\S]+}\s+}`)

func assertOutputContainsEvent(t *testing.T, evt apitype.EngineEvent, output string) {
	evtJSON := bytes.Buffer{}
	encoder := json.NewEncoder(&evtJSON)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(evt)
	assert.NoError(t, err)
	assert.Contains(t, output, evtJSON.String())
}

// printfTestValidation is used by the TestPrintfXYZ test cases in the language-specific test
// files. It validates that there are a precise count of expected stdout/stderr lines in the test output.
func printfTestValidation(t *testing.T, stack integration.RuntimeValidationStackInfo) {
	var foundStdout int
	var foundStderr int
	for _, ev := range stack.Events {
		if de := ev.DiagnosticEvent; de != nil {
			if strings.HasPrefix(de.Message, fmt.Sprintf("Line %d", foundStdout)) {
				foundStdout++
			} else if strings.HasPrefix(de.Message, fmt.Sprintf("Errln %d", foundStderr+10)) {
				foundStderr++
			}
		}
	}
	assert.Equal(t, 11, foundStdout)
	assert.Equal(t, 11, foundStderr)
}

// testLogWriter is an io.Writer
// that writes to the provided testing.T.
//
// Use newTestLogWriter to ensure it is flushed
// of any input when the test finishes.
type testLogWriter struct {
	logf func(string, ...interface{})

	// Holds buffered text for the next write or flush
	// if we haven't yet seen a newline.
	buff bytes.Buffer
}

var _ io.Writer = (*testLogWriter)(nil)

func newTestLogWriter(t testing.TB) *testLogWriter {
	w := testLogWriter{logf: t.Logf}
	t.Cleanup(w.Flush)
	return &w
}

func (w *testLogWriter) Write(bs []byte) (int, error) {
	// t.Logf adds a newline so we should not write bs as-is.
	// Instead, we'll call t.Log one line at a time.
	//
	// To handle the case when Write is called with a partial line,
	// we use a buffer.
	total := len(bs)
	for len(bs) > 0 {
		idx := bytes.IndexByte(bs, '\n')
		if idx < 0 {
			// No newline. Buffer it for later.
			w.buff.Write(bs)
			break
		}

		var line []byte
		line, bs = bs[:idx], bs[idx+1:]

		if w.buff.Len() == 0 {
			// Nothing buffered from a prior partial write.
			// This is the majority case.
			w.logf("%s", line)
			continue
		}

		// There's a prior partial write. Join and flush.
		w.buff.Write(line)
		w.logf("%s", w.buff.String())
		w.buff.Reset()
	}
	return total, nil
}

func (w *testLogWriter) Flush() {
	if w.buff.Len() > 0 {
		w.logf("%s", w.buff.String())
	}
}

func TestTestLogWriter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string

		writes []string // individual write calls
		want   []string // expected log output
	}{
		{
			desc:   "empty strings",
			writes: []string{"", "", ""},
		},
		{
			desc:   "no newline",
			writes: []string{"foo", "bar", "baz"},
			want:   []string{"foobarbaz"},
		},
		{
			desc: "newline separated",
			writes: []string{
				"foo\n",
				"bar\n",
				"baz\n\n",
				"qux",
			},
			want: []string{
				"foo",
				"bar",
				"baz",
				"",
				"qux",
			},
		},
		{
			desc:   "partial line",
			writes: []string{"foo", "bar\nbazqux"},
			want: []string{
				"foobar",
				"bazqux",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			var got []string
			w := testLogWriter{
				logf: func(msg string, args ...interface{}) {
					got = append(got, fmt.Sprintf(msg, args...))
				},
			}

			for _, input := range tt.writes {
				n, err := w.Write([]byte(input))
				assert.NoError(t, err)
				assert.Equal(t, len(input), n)
			}

			w.Flush()

			assert.Equal(t, tt.want, got)
		})
	}
}
