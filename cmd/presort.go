/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// presortCmd represents the presort command
var presortCmd = &cobra.Command{
	Use:   "presort",
	Short: "Presorting the current judging project",
	Long:  `This is running the interactive presorting of the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		presort(".")
	},
}

func init() {
	rootCmd.AddCommand(presortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// presortCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// presortCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func processIssue(issuePath string, issueType string) {
	switch issueType {
	case "i":
		err := os.Rename(issuePath, "invalid/"+filepath.Base(issuePath))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Issue moved to 'invalid' folder.")
	case "m", "h":
		prefix := strings.ToUpper(issueType)
		nextNumber := getNextIssueNumber(prefix)
		newFolder := fmt.Sprintf("%s-%03d", prefix, nextNumber)
		err := os.Mkdir(newFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Rename(issuePath, newFolder+"/"+filepath.Base(issuePath))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Issue moved to '%s' folder.\n", newFolder)
	case "d":
		existingFolders := []string{}
		files, err := os.ReadDir(".") //TODO: Make this relative to the current issue
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() {
				name := file.Name()
				if strings.HasPrefix(name, "M-") || strings.HasPrefix(name, "H-") {
					existingFolders = append(existingFolders, name)
				}
			}
		}

		sort.Strings(existingFolders)

		fmt.Println("Existing folders:")
		for i, folder := range existingFolders {
			_, title := getFirstIssueInfo(folder)
			fmt.Printf("%d. %s - %s\n", i+1, folder, title)
		}

		var choice int
		fmt.Print("Enter the number of the existing folder: ")
		_, err = fmt.Scan(&choice)
		if err != nil || choice < 1 || choice > len(existingFolders) {
			fmt.Println("Invalid choice. Issue not moved.")
		} else {
			existingFolder := existingFolders[choice-1]
			err = os.Rename(issuePath, existingFolder+"/"+filepath.Base(issuePath))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Issue moved to '%s' folder.\n", existingFolder)
		}
	default:
		fmt.Println("Invalid choice. Issue not moved.")
	}
}

func getFirstIssueInfo(folder string) (string, string) {
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.ToLower(file.Name()) != "comment.md" {
			filePath := filepath.Join(folder, file.Name())
			lines, err := readLines(filePath)
			if err != nil || len(lines) < 5 {
				continue
			}
			title := strings.TrimSpace(lines[4])
			return file.Name(), title
		}
	}

	return "", ""
}

func readLines(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func extractSummary(lines []string) string {
	summaryStart := 7
	summaryEnd := len(lines)
	for i, line := range lines[summaryStart:] {
		if strings.HasPrefix(line, "##") {
			summaryEnd = i + summaryStart
			break
		}
	}
	summary := strings.Join(lines[summaryStart:summaryEnd], "\n")
	return strings.TrimSpace(summary)
}

func presort(dir string) {
	// Create 'invalid' folder if it doesn't exist
	os.Mkdir("invalid", 0755)

	// Get all markdown files in the current directory
	fsys := os.DirFS(dir)
	issues, err := fs.Glob(fsys,"[0-9]*.md")
	if err != nil {
		log.Fatal(err)
	}


	if len(issues) == 0 {
		fmt.Println("No issues found.")
		return
	}

	sortIssues(issues)

	for _, issue := range issues {
		lines, err := readLines(issue)
		if err != nil {
			log.Fatal(err)
		}

		// auditor := strings.TrimSpace(strings.TrimPrefix(lines[0], "#"))
		// severity := strings.TrimSpace(lines[2])
		title := strings.TrimSpace(lines[4])
		summary := extractSummary(lines)

		fmt.Printf("\n\nIssue: %s\n", issue)
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Summary: %s\n", summary)

		var issueType string
		fmt.Print("Select the issue type ((i)nvalid/(m)edium/(h)igh/(d)uplicate/(s)kip): ")
		fmt.Scan(&issueType) //TODO: We can replace issueType with a char to be more efficient
		issueType = strings.TrimSpace(issueType)

		if issueType == "s" {
			continue
		} else if issueType == "q" {
			os.Exit(0)
		}
		processIssue(issue, issueType)
		fmt.Println()
	}
}

func sortIssues(issues []string) {
	sort.Slice(issues, func(i, j int) bool {
		return issues[i] < issues[j]
	})
}
