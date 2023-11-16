/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/91go/docs-training/internal"
	"github.com/91go/docs-training/utils"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "count the number of questions in each file",
	Run: func(cmd *cobra.Command, args []string) {
		wf, _ := cmd.Flags().GetStringSlice("wf")
		ex, _ := cmd.Flags().GetStringSlice("exclude")
		res := make([][]string, 0)

		for _, w := range wf {
			var qs [][]string
			if gfile.IsDir(w) {
				qs = internal.NewDir(w).Xz(internal.ExtractQuestion).Exclude(ex).GetTableData()
			}
			if gfile.IsFile(w) {
				qs = internal.NewFile(w).Xz().GetTableData(w, 0)
			}
			res = append(res, qs...)
		}

		utils.GenerateMDTable(res)
	},
}

func init() {
	rootCmd.AddCommand(countCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// countCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// countCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
