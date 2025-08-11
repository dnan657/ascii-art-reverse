package main

import (
	"testing"
)

func TestReverseAsciiArt(t *testing.T) {
	testCases := []struct {
		name           string
		inputFile      string
		banner         string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "Test Case 00: Hello World",
			inputFile:      "example00.txt",
			banner:         "standard",
			expectedOutput: "Hello World",
			expectError:    false,
		},
		{
			name:           "Test Case 01: Numbers",
			inputFile:      "example01.txt",
			banner:         "standard",
			expectedOutput: "123",
			expectError:    false,
		},
		{
			name:           "Test Case 02: Special Characters 1",
			inputFile:      "example02.txt",
			banner:         "standard",
			expectedOutput: "#=\\[",
			expectError:    false,
		},
		{
			name:           "Test Case 03: Mixed",
			inputFile:      "example03.txt",
			banner:         "standard",
			expectedOutput: "something&234",
			expectError:    false,
		},
		{
			name:           "Test Case 04: Lowercase Alphabet",
			inputFile:      "example04.txt",
			banner:         "standard",
			expectedOutput: "abcdefghijklmnopqrstuvwxyz",
			expectError:    false,
		},
		{
			name:           "Test Case 05: Special Characters 2",
			inputFile:      "example05.txt",
			banner:         "standard",
			expectedOutput: `\!" #$%&'()*+,-./`,
			expectError:    false,
		},
		{
			name:           "Test Case 06: Special Characters 3",
			inputFile:      "example06.txt",
			banner:         "standard",
			expectedOutput: `:;{=}?@`,
			expectError:    false,
		},
		{
			name:           "Test Case 07: Uppercase Alphabet",
			inputFile:      "example07.txt",
			banner:         "standard",
			expectedOutput: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			expectError:    false,
		},
		{
			name:           "Test Case 08: Non-existent file",
			inputFile:      "nonexistent.txt",
			banner:         "standard",
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reverseAsciiArt(tc.inputFile, tc.banner)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
				if got != tc.expectedOutput {
					t.Errorf("expected [%s], got [%s]", tc.expectedOutput, got)
				}
			}
		})
	}
}
