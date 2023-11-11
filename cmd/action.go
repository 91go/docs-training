/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/91go/docs-training/utils"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/spf13/cobra"
)

// actionCmd represents the action command
var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		wf, _ := cmd.Flags().GetStringSlice("wf")
		ex, _ := cmd.Flags().GetStringSlice("exclude")
		num, _ := cmd.Flags().GetInt("num")

		zk := make([]string, 0)
		for _, w := range wf {
			// var qs []string
			if gfile.IsDir(w) {
				zk = append(zk, utils.NewDir(w).Xz().Exclude(ex).GetQuestions()...)
			}
			if gfile.IsFile(w) {
				zk = append(zk, utils.NewFile(w).Xz().GetQuestions()...)
			}
		}

		lzk := len(zk)
		if lzk < num {
			num = lzk
			log.Printf("%v, the number of questions is less than %d, so use %d", wf, num, lzk)
		}

		// 随机打乱，再取前n个
		rands := garray.NewStrArrayFrom(zk).Shuffle().SubSlice(0, num)
		rt := utils.GenerateMD(rands)
		fmt.Println(rt)
	},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// actionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// actionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
