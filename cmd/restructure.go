/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "os"
	"sort"
    "log"
	"strconv"
	"strings"
	"github.com/spf13/cobra"
)

// restructureCmd represents the restructure command
var restructureCmd = &cobra.Command{
	Use:   "restructure",
	Short: "restructuring issue numbers to be from 1 - n",
	Long: `restructuring issue numbers in case there are gaps, so that they match from 1-n`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restructuring Medium issues (M)")
        restructure("M");
		fmt.Println("restructuring High issues (H)")
        restructure("H");

	},
}

func init() {
	rootCmd.AddCommand(restructureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restructureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restructureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}




func restructure(prefix string) {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	existingFolders := []string{}
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), prefix+"-") {
			existingFolders = append(existingFolders, file.Name())
		}
	}
	sort.Strings(existingFolders)
	nextNum := getNextIssueNumber(prefix)
	fmt.Println("Nextnum would be", nextNum)
	issueCount := len(existingFolders)
	if  (nextNum - len(existingFolders)) > 1 {
		fmt.Println("Mismatch in numbers for folders")
		for i := 1; i < issueCount; i++ {
			folderName := fmt.Sprintf("%s-%03d", prefix, i)
			if !contains(existingFolders, folderName) {
				fmt.Println(folderName, "missing")
				lastFolder := existingFolders[len(existingFolders)-1]
				fmt.Println("Renaming", lastFolder, "to", folderName)
					err := os.Rename(lastFolder, folderName)
					if err != nil {
						log.Fatal(err)
					}
				existingFolders = existingFolders[:len(existingFolders)-1]
			}
		}
	} else {
		fmt.Println("No restructure needed")
	}
}

func getNextIssueNumber(prefix string) int {
	highestNum := 0

	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), prefix+"-") {
			parts := strings.Split(file.Name(), "-")
			if len(parts) != 2 {
				continue
			}

			num, err := strconv.Atoi(parts[1])
			if err == nil && num > highestNum {
				highestNum = num
			}
		}
	}

	return highestNum + 1
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

