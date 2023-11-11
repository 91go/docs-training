/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/91go/docs-training/utils"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/spf13/cobra"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "A brief description of your command",
	Long: `提取markdown文件中的header和无序列表
只提取一级header，如果一级header下没有无序列表，则不提取`,
	Run: func(cmd *cobra.Command, args []string) {
		wf, _ := cmd.Flags().GetStringSlice("wf")
		ex, _ := cmd.Flags().GetStringSlice("exclude")

		files := make([]utils.File, 0)
		for _, w := range wf {
			if gfile.IsDir(w) {
				files = utils.NewDir(w).Xz().Exclude(ex).GetFiles()
			}
			if gfile.IsFile(w) {
				files = append(files, *utils.NewFile(w).Xz())
			}
		}
		// 提取files的Name为header，Questions为无序列表
		res := ""
		for _, file := range files {
			res += file.ConvertToMarkdown()
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(allCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
