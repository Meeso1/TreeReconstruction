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

type TimeResult struct {
	InputFile string
	Duration  time.Duration
	Error     error
}

var (
	timeOutputFile              string
	timeSerializationTypeString string
)

func init() {
	timeCmd.Flags().StringVarP(&timeOutputFile, "output", "o", "", "Output file to save reconstruction times (required)")
	timeCmd.Flags().StringVarP(&timeSerializationTypeString, "serialization", "s", "neighbor-lists", "Serialization type (brackets, brackets-shortened, neighbor-lists)")
	timeCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(timeCmd)
}

var timeCmd = &cobra.Command{
	Use:   "time <directory>",
	Short: "Time reconstruction on all '*.input.txt' files in a directory",
	Long:  `Run the reconstruction algorithm on all '*.input.txt' files in the specified directory and save timing results to a file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			fmt.Printf("Error: directory %s does not exist\n", directory)
			return
		}

		// Parse serialization type
		var serializationType io.SerializationType
		switch timeSerializationTypeString {
		case "brackets":
			serializationType = io.SerializationTypeBrackets
		case "brackets-shortened":
			serializationType = io.SerializationTypeBracketsShortened
		case "neighbor-lists":
			serializationType = io.SerializationTypeNeighborLists
		default:
			fmt.Printf("Invalid serialization type: %s\n", timeSerializationTypeString)
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

		fmt.Printf("Timing reconstruction on %d input files in %s...\n\n", len(inputFiles), directory)

		var results []TimeResult
		for _, inputFile := range inputFiles {
			result := runTimingTest(inputFile, serializationType)
			results = append(results, result)
			printTimingResult(result)
		}

		err = saveTimesToFile(timeOutputFile, results)
		if err != nil {
			fmt.Printf("Error: Failed to save times to file %s: %v\n", timeOutputFile, err)
			return
		}

		fmt.Printf("\nTiming completed. Results saved to %s\n", timeOutputFile)
		printTimingSummary(results)
	},
}

func runTimingTest(inputFile string, serializationType io.SerializationType) TimeResult {
	result := TimeResult{
		InputFile: inputFile,
	}

	tmpDir := os.TempDir()
	outputFile := filepath.Join(tmpDir, fmt.Sprintf("time_output_%d.txt", time.Now().UnixNano()))

	start := time.Now()
	reconstructResult := runReconstructCommand(inputFile, outputFile, serializationType)
	result.Duration = time.Since(start)

	if reconstructResult.Error != nil {
		result.Error = reconstructResult.Error
	}

	// Clean up temporary file
	os.Remove(outputFile)

	return result
}

func printTimingResult(result TimeResult) {
	inputName := filepath.Base(result.InputFile)

	if result.Error != nil {
		fmt.Printf("✗ %s - ERROR: %v\n", inputName, result.Error)
	} else {
		fmt.Printf("✓ %s - %.6fs\n", inputName, result.Duration.Seconds())
	}
}

func saveTimesToFile(filename string, results []TimeResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, result := range results {
		// Extract input name without '.input.txt' suffix
		inputName := filepath.Base(result.InputFile)
		inputName = strings.TrimSuffix(inputName, ".input.txt")

		// Only save successful results
		if result.Error == nil {
			timeSeconds := result.Duration.Seconds()
			_, err := fmt.Fprintf(file, "%s;%.6f\n", inputName, timeSeconds)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func printTimingSummary(results []TimeResult) {
	successful := 0
	failed := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		if result.Error == nil {
			successful++
			totalDuration += result.Duration
		} else {
			failed++
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 40) + "\n")
	fmt.Printf("TIMING SUMMARY\n")
	fmt.Printf(strings.Repeat("=", 40) + "\n")
	fmt.Printf("Total files:     %d\n", len(results))
	fmt.Printf("Successful:      %d\n", successful)
	fmt.Printf("Failed:          %d\n", failed)
	if successful > 0 {
		fmt.Printf("Total time:      %.6fs\n", totalDuration.Seconds())
		fmt.Printf("Average time:    %.6fs\n", totalDuration.Seconds()/float64(successful))
	}

	if failed > 0 {
		fmt.Printf("\nFAILED FILES:\n")
		for _, result := range results {
			if result.Error != nil {
				fmt.Printf("  %s: %v\n", filepath.Base(result.InputFile), result.Error)
			}
		}
	}
}
