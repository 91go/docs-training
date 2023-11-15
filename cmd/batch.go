package cmd

import (
	"fmt"

	"github.com/91go/docs-training/utils"
	"github.com/spf13/cobra"
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		wf, _ := cmd.Flags().GetStringSlice("wf")
		num, _ := cmd.Flags().GetInt("num")
		var rs string
		for _, w := range wf {
			rs += utils.NewDir(w).Xz(utils.ExtractInterviews).InterviewsToMarkdown(num)
		}
		fmt.Println(rs)
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// batchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// batchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
