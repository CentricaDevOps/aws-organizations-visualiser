package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSetupLoggingTrue tests the setupLogging function when the logging
// parameter is true.
func TestSetupLoggingTrue(t *testing.T) {
	setupLogging("true")
	require.NotNil(t, logs, "Expected logs to be initialised")
	require.Equal(t, true, logs.Enabled, "Expected logs to be enabled")
}

// TestSetupLoggingFalse tests the setupLogging function when the logging
// parameter is false.
func TestSetupLoggingFalse(t *testing.T) {
	setupLogging("false")
	require.NotNil(t, logs, "Expected logs to be initialised")
	require.Equal(t, false, logs.Enabled, "Expected logs to be disabled")
}

// TestSetupLoggingInvalid tests the setupLogging function when the logging
// parameter is invalid.
func TestSetupLoggingInvalid(t *testing.T) {
	setupLogging("invalid")
	require.NotNil(t, logs, "Expected logs to be initialised")
	require.Equal(t, false, logs.Enabled, "Expected logs to be disabled")
}

// TestSetupLoggingEmpty tests the setupLogging function when the logging
// parameter is empty.
func TestSetupLoggingEmpty(t *testing.T) {
	setupLogging("")
	require.NotNil(t, logs, "Expected logs to be initialised")
	require.Equal(t, false, logs.Enabled, "Expected logs to be disabled")
}

// TestPrintlnLogging tests the Println function of the Logs struct.
func TestPrintlnLogging(t *testing.T) {
	setupLogging("true")
	output1 := captureOutput(func() {
		logs.Println("Test")
	})
	require.Equal(t, "Test\n", output1, "Expected output to be 'Test\n'")

	setupLogging("false")
	output2 := captureOutput(func() {
		logs.Println("Test")
	})
	require.Equal(t, "", output2, "Expected output to be ''")
}

// captureOutput is a helper function to capture the output of a function.
// This is used to test the output of the display functions.
func captureOutput(f func()) string {
	// Store the old stdout and replace it with a pipe.
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	// Run the function.
	f()
	// Close the pipe and restore stdout.
	w.Close()
	os.Stdout = original

	// Read the output from the pipe.
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}
	// Return the output.
	return buf.String()
}
