package main

import (
	// "fmt"
	// "io/ioutil"
	"os/exec"
	"testing"
)

type TestCases struct {
	Input    string
	Expected string
}

func Build(sourceFile string) error {
	// TODO: find gofmt dynamically.
	cmd := exec.Command("go", "build", sourceFile)
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// func TestNoDirectives(t *testing.T) {
// 	input := "test-cases/no-directives.go"
// 	actual := "test-cases/no-directives-actual.go"
// 	output := "test-cases/no-directives-expected.go"
// 	ProcessFile(input, actual)

// 	if !CompareFiles(actual, output) {
// 		t.Errorf("processing input: %s, result: %s differ from expectation: %s\n", input, actual, output)
// 	}
// }

// func TestPrimitives(t *testing.T) {
// 	input := "./test-cases/primitive-both.go"
// 	actual := "./test-cases/primitive-both-actual.go"
// 	output := "./test-cases/primitive-both-expected.go"
// 	ProcessFile(input, actual)

// 	if !CompareFiles(actual, output) {
// 		t.Errorf("processing input: %s, result: %s differ from expectation: %s\n", input, actual, output)
// 	}
// }

// func TestPrimitivesOnlyEntry(t *testing.T) {
// 	input := "./test-cases/primitive-entry.go"
// 	actual := "./test-cases/primitive-entry-actual.go"
// 	output := "./test-cases/primitive-entry-expected.go"
// 	ProcessFile(input, actual)

// 	if !CompareFiles(actual, output) {
// 		t.Errorf("processing input: %s, result: %s differ from expectation: %s\n", input, actual, output)
// 	}
// }

// func TestPrimitivesOnlyReturn(t *testing.T) {
// 	input := "./test-cases/primitive-return.go"
// 	actual := "./test-cases/primitive-return-actual.go"
// 	output := "./test-cases/primitive-return-expected.go"
// 	ProcessFile(input, actual)

// 	if !CompareFiles(actual, output) {
// 		t.Errorf("processing input: %s, result: %s differ from expectation: %s\n", input, actual, output)
// 	}
// }

func TestArray(t *testing.T) {
	input := "./test-cases/array.go"
	output := "./test-cases/array-output.go"
	ProcessFile(input, output)

	// Build output
	if err := Build(output); err != nil {
		t.Error(err)
	}

	// Execute.
}

// func TestPointerArray(t *testing.T) {
// 	input := "./test-cases/pointer-array.go"
// 	actual := "./test-cases/pointer-array-actual.go"
// 	output := "./test-cases/"
// 	ProcessFile(input, actual)

// 	if !CompareFiles(actual, output) {
// 		t.Errorf("processing input: %s, result: %s differ from expectation: %s\n", input, actual, output)
// 	}
// }

// func TestPointer(t *testing.T) {
// 	input := "./test-cases/pointer.go"
// 	actual := "./test-cases/pointer-actual.go"
// 	output := "./test-cases/"
// 	ProcessFile(input, actual)

// 	if !CompareFiles(actual, output) {
// 		t.Errorf("processing input: %s, result: %s differ from expectation: %s\n", input, actual, output)
// 	}
// }
