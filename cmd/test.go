package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"treereconstruction/io"

	"github.com/spf13/cobra"
)

type TestResult struct {
	InputFile         string
	OutputFile        string
	ExpectedFile      string
	Status            TestStatus
	Error             string
	Duration          time.Duration
	ComparisonDetails *CompareResult
}

type TestStatus int

const (
	TestPassed TestStatus = iota
	TestFailed
	TestSkipped
	TestError
)

func (s TestStatus) String() string {
	switch s {
	case TestPassed:
		return "PASS"
	case TestFailed:
		return "FAIL"
	case TestSkipped:
		return "SKIP"
	case TestError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test <directory>",
	Short: "Run batch tests on all '*.input.txt' files in a directory",
	Long:  `Run the reconstruct command on all '*.input.txt' files in the specified directory and compare results with corresponding '*.output.txt' files.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			fmt.Printf("Error: directory %s does not exist\n", directory)
			return
		}

		inputFiles, err := findInputFiles(directory)
		if err != nil {
			fmt.Printf("Error finding input files: %v\n", err)
			return
		}

		if len(inputFiles) == 0 {
			fmt.Printf("No '*.input.txt' files found in directory %s\n", directory)
			return
		}

		fmt.Printf("Running tests on %d input files in %s...\n\n", len(inputFiles), directory)

		var results []TestResult
		for _, inputFile := range inputFiles {
			result := runSingleTest(inputFile)
			results = append(results, result)
			printTestResult(result)
		}

		printTestSummary(results)
	},
}

func findInputFiles(directory string) ([]string, error) {
	var inputFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".input.txt") {
			inputFiles = append(inputFiles, path)
		}

		return nil
	})

	return inputFiles, err
}

func runSingleTest(inputFile string) TestResult {
	start := time.Now()

	result := TestResult{
		InputFile: inputFile,
	}

	expectedFile := strings.TrimSuffix(inputFile, ".input.txt") + ".output.txt"
	result.ExpectedFile = expectedFile

	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		result.Status = TestSkipped
		result.Error = "Expected output file not found"
		result.Duration = time.Since(start)
		return result
	}

	tmpDir := os.TempDir()
	outputFile := filepath.Join(tmpDir, fmt.Sprintf("test_output_%d.txt", time.Now().UnixNano()))
	result.OutputFile = outputFile

	reconstructResult := runReconstructCommand(inputFile, outputFile, io.SerializationTypeNeighborLists)

	if reconstructResult.Error != nil {
		result.Status = TestError
		result.Error = fmt.Sprintf("Reconstruction failed: %v", reconstructResult.Error)
		result.Duration = time.Since(start)
		// Clean up temporary file
		os.Remove(outputFile)
		return result
	}

	// Compare results using the extracted function
	compareResult := runCompareCommand(outputFile, expectedFile)
	if compareResult.Error != nil {
		result.Status = TestError
		result.Error = fmt.Sprintf("Comparison failed: %v", compareResult.Error)
		result.Duration = time.Since(start)
		// Clean up temporary file
		os.Remove(outputFile)
		return result
	}

	if compareResult.TopologiesMatch {
		result.Status = TestPassed
	} else {
		result.Status = TestFailed
		result.Error = "Tree topologies differ"
		result.ComparisonDetails = &compareResult
	}

	result.Duration = time.Since(start)

	// Clean up temporary file
	os.Remove(outputFile)

	return result
}

func printTestResult(result TestResult) {
	status := result.Status.String()
	inputName := filepath.Base(result.InputFile)

	switch result.Status {
	case TestPassed:
		fmt.Printf("✓ [%s] %s (%.2fs)\n", status, inputName, result.Duration.Seconds())
	case TestFailed:
		fmt.Printf("✗ [%s] %s (%.2fs) - %s\n", status, inputName, result.Duration.Seconds(), result.Error)
		if result.ComparisonDetails != nil {
			fmt.Printf("  Generated output:\n")
			for _, line := range result.ComparisonDetails.Tree1Summary {
				fmt.Printf("    %s\n", line)
			}
			fmt.Printf("  Expected output:\n")
			for _, line := range result.ComparisonDetails.Tree2Summary {
				fmt.Printf("    %s\n", line)
			}
		}
	case TestSkipped:
		fmt.Printf("- [%s] %s - %s\n", status, inputName, result.Error)
	case TestError:
		fmt.Printf("! [%s] %s (%.2fs) - %s\n", status, inputName, result.Duration.Seconds(), result.Error)
	}
}

func printTestSummary(results []TestResult) {
	passed := 0
	failed := 0
	skipped := 0
	errors := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		totalDuration += result.Duration
		switch result.Status {
		case TestPassed:
			passed++
		case TestFailed:
			failed++
		case TestSkipped:
			skipped++
		case TestError:
			errors++
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 50) + "\n")
	fmt.Printf("TEST SUMMARY\n")
	fmt.Printf(strings.Repeat("=", 50) + "\n")
	fmt.Printf("Total tests: %d\n", len(results))
	fmt.Printf("Passed:      %d\n", passed)
	fmt.Printf("Failed:      %d\n", failed)
	fmt.Printf("Skipped:     %d\n", skipped)
	fmt.Printf("Errors:      %d\n", errors)
	fmt.Printf("Duration:    %.2fs\n", totalDuration.Seconds())

	if failed > 0 || errors > 0 {
		fmt.Printf("\nFAILED/ERROR DETAILS:\n")
		for _, result := range results {
			if result.Status == TestFailed || result.Status == TestError {
				fmt.Printf("  %s: %s\n", filepath.Base(result.InputFile), result.Error)
			}
		}
	}

	if passed+failed > 0 {
		successRate := float64(passed) / float64(passed+failed) * 100
		fmt.Printf("\nSuccess rate: %.1f%% (%d/%d)\n", successRate, passed, passed+failed)
	}
}
