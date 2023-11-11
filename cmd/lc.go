/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/spf13/cobra"
)

//go:embed lc.json
var questions string

// lcCmd represents the lc command
var lcCmd = &cobra.Command{
	Use:   "lc",
	Short: "根据codetop Top100的题目，生成leetcode的题目",
	Run: func(cmd *cobra.Command, args []string) {
		num, _ := cmd.Flags().GetInt("num")
		qj, _ := gjson.LoadJson(questions)

		rands := garray.NewArrayFrom(qj.Array()).Shuffle().SubSlice(0, num)
		rt := ""
		for _, rand := range rands {
			rt += gconv.String(rand) + "\n"
		}
		fmt.Println(rt)
	},
}

func init() {
	rootCmd.AddCommand(lcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
