/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// lcExtractCmd represents the lcExtract command
var lcExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		// 定义URL和初始页码
		url := "https://codetop.cc/api/questions/?page=%d&search=&ordering=-frequency&rate=0"
		page := 1

		// 创建JSON文件
		file, err := os.Create("lc.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 创建JSON编码器
		encoder := json.NewEncoder(file)

		// 获取第一页数据
		resp, err := http.Get(fmt.Sprintf(url, page))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// 解析响应
		var response CodeTop
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			panic(err)
		}
		var res []string

		// 输出第一页数据并写入JSON文件
		for _, question := range response.List {
			if question.Value >= Freq {
				fmt.Printf("ID: %d, Title: %s, Difficulty: %d, SlugTitle: %s\n",
					question.Leetcode.QuestionID, question.Leetcode.Title, question.Leetcode.Level, question.Leetcode.SlugTitle)
				res = append(res, fmt.Sprintf("- [ ] [%d. %s](https://leetcode.cn/problems/%s)", question.Leetcode.QuestionID, question.Leetcode.Title, question.Leetcode.SlugTitle))
			}
		}

		// 获取后续页码数据，直到没有下一页为止
		for resp.StatusCode == http.StatusOK {
			page++
			resp, err = http.Get(fmt.Sprintf(url, page))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				panic(err)
			}

			for _, question := range response.List {
				if question.Value >= Freq {
					fmt.Printf("ID: %d, Title: %s, Difficulty: %d, SlugTitle: %s\n",
						question.Leetcode.QuestionID, question.Leetcode.Title, question.Leetcode.Level, question.Leetcode.SlugTitle)
					res = append(res, fmt.Sprintf("- [ ] [%s](https://leetcode.cn/problems/%s)", question.Leetcode.Title, question.Leetcode.SlugTitle))
				}
			}
		}

		err = encoder.Encode(res)
		if err != nil {
			return
		}
	},
}

func init() {
	lcCmd.AddCommand(lcExtractCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lcExtractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lcExtractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var Freq = 30

type CodeTop struct {
	List     []List `json:"list"`
	Finished []int  `json:"finished"`
	Count    int    `json:"count"`
}

type Leetcode struct {
	FrontendQuestionID string `json:"frontend_question_id"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	SlugTitle          string `json:"slug_title"`
	ID                 int    `json:"id"`
	QuestionID         int    `json:"question_id"`
	Level              int    `json:"level"`
	Expand             bool   `json:"expand"`
}

type List struct {
	Time         time.Time `json:"time"`
	Leetcode     Leetcode  `json:"leetcode"`
	ID           int       `json:"id"`
	Value        int       `json:"value"`
	Rate         int       `json:"rate"`
	CommentCount int       `json:"comment_count"`
	Status       bool      `json:"status"`
	NoteStatus   bool      `json:"note_status"`
}
