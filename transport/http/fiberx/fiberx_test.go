package fiberx

import (
	"strings"
	"testing"
)

func TestTrimSuffix(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "/",
			expectedOutput: "",
		},
		{
			input:          "/basePath",
			expectedOutput: "/basePath",
		},
		{
			input:          "/basePath/",
			expectedOutput: "/basePath",
		},
		{
			input:          "/basePath/abc",
			expectedOutput: "/basePath/abc",
		},
		{
			input:          "/basePath/abc/",
			expectedOutput: "/basePath/abc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := strings.TrimSuffix(tc.input, "/")
			t.Logf("Input: %s, Result: %s", tc.input, result)
			if result != tc.expectedOutput {
				t.Errorf("Expected: %s, Got: %s", tc.expectedOutput, result)
			}
		})
	}
}

func TestGetPath(t *testing.T) {
	fiberApp := NewFiberApp(WithBasePath("/basePath"))

	testCasesIncludeBasePath := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "/",
			expectedOutput: "/basePath/",
		},
		{
			input:          "/abc",
			expectedOutput: "/basePath/abc",
		},
	}

	for _, tc := range testCasesIncludeBasePath {
		t.Run(tc.input, func(t *testing.T) {
			result := fiberApp.getPath(tc.input, true)
			t.Logf("Input: %s, Result: %s", tc.input, result)
			if result != tc.expectedOutput {
				t.Errorf("Expected: %s, Got: %s", tc.expectedOutput, result)
			}
		})
	}

	testCasesNotIncludeBasePath := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "/",
			expectedOutput: "/",
		},
		{
			input:          "/abc",
			expectedOutput: "/abc",
		},
	}

	for _, tc := range testCasesNotIncludeBasePath {
		t.Run(tc.input, func(t *testing.T) {
			result := fiberApp.getPath(tc.input, false)
			t.Logf("Input: %s, Result: %s", tc.input, result)
			if result != tc.expectedOutput {
				t.Errorf("Expected: %s, Got: %s", tc.expectedOutput, result)
			}
		})
	}
}
