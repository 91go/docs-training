package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/91go/docs-training/dir"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/urfave/cli"
)

//go:embed lc.json
var questions string

var (
	cmds  []cli.Command
	flags []cli.Flag
	n     int
)

func init() {
	flags = []cli.Flag{
		cli.IntFlag{
			Name:        "n, num",
			Value:       30,
			Usage:       "num of questions",
			Destination: &n,
			Required:    false,
		},
		cli.StringSliceFlag{
			Name:     "w, wf",
			Value:    &cli.StringSlice{},
			Usage:    "files",
			Required: false,
		},
		cli.StringSliceFlag{
			Name:     "e, exclude",
			Value:    &cli.StringSlice{},
			Usage:    "exclude specified files",
			Required: false,
		},
	}
	cmds = []cli.Command{
		{
			Name:   "all",
			Usage:  "all the questions in each file, 生成json文件",
			Action: All,
			Flags:  flags,
		},
		{
			Name:   "count",
			Usage:  "count the number of questions in each file",
			Action: Count,
			Flags:  flags,
		},
		{
			Name:   "lc",
			Usage:  "根据codetop Top100的题目，生成leetcode的题目",
			Action: Leetcode,
			Flags:  flags,
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Docs-Training"
	app.Usage = "use to collocate with gh-ac, generate weekly training items"
	app.HideVersion = true
	app.Flags = flags
	app.Action = Action
	app.Commands = cmds

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Count(c *cli.Context) error {
	wf := c.StringSlice("wf")
	ex := c.StringSlice("exclude")
	res := make([][]string, 0)

	for _, w := range wf {
		var qs [][]string
		if gfile.IsDir(w) {
			qs = dir.NewDir(w).Xz().Exclude(ex).GetTableData()
		}
		if gfile.IsFile(w) {
			qs = dir.NewFile(w).Xz().GetTableData(w, 0)
		}
		res = append(res, qs...)
	}

	dir.GenerateMDTable(res)

	return nil
}

// All 提取markdown文件中的header和无序列表
// 只提取一级header，如果一级header下没有无序列表，则不提取
func All(c *cli.Context) error {
	wf := c.StringSlice("wf")
	ex := c.StringSlice("exclude")

	files := make([]dir.File, 0)
	for _, w := range wf {
		if gfile.IsDir(w) {
			files = dir.NewDir(w).Xz().Exclude(ex).GetFiles()
		}
		if gfile.IsFile(w) {
			files = append(files, *dir.NewFile(w).Xz())
		}
	}
	// 提取files的Name为header，Questions为无序列表
	res := ""
	for _, file := range files {
		res += file.ConvertToMarkdown()
	}
	fmt.Println(res)
	return nil
}

// Leetcode 生成leetcode的题目
func Leetcode(c *cli.Context) error {
	num := c.Int("num")
	qj, err := gjson.LoadJson(questions)
	if err != nil {
		return err
	}
	rands := garray.NewArrayFrom(qj.Array()).Shuffle().SubSlice(0, num)
	rt := ""
	for _, rand := range rands {
		rt += gconv.String(rand) + "\n"
	}
	fmt.Println(rt)

	return nil
}

func Action(c *cli.Context) error {
	wf := c.StringSlice("wf")
	ex := c.StringSlice("exclude")
	num := c.Int("num")

	zk := make([]string, 0)
	for _, w := range wf {
		// var qs []string
		if gfile.IsDir(w) {
			zk = append(zk, dir.NewDir(w).Xz().Exclude(ex).GetQuestions()...)
		}
		if gfile.IsFile(w) {
			zk = append(zk, dir.NewFile(w).Xz().GetQuestions()...)
		}
	}

	lzk := len(zk)
	if lzk < num {
		num = lzk
		log.Printf("%v, the number of questions is less than %d, so use %d", wf, num, lzk)
	}

	// 随机打乱，再取前n个
	rands := garray.NewStrArrayFrom(zk).Shuffle().SubSlice(0, num)
	rt := dir.GenerateMD(rands)
	fmt.Println(rt)
	return nil
}
