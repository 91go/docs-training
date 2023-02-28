package main

import (
	"fmt"
	"log"
	"os"

	"github.com/91go/weekly-training/dir"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/urfave/cli"
)

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
			Name:   "count",
			Usage:  "count the number of questions in each file",
			Action: Count,
			Flags:  flags,
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Weekly-Training"
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

func Action(c *cli.Context) error {
	wf := c.StringSlice("wf")
	ex := c.StringSlice("exclude")
	num := c.Int("num")

	var zk []string
	for _, w := range wf {
		var qs []string
		if gfile.IsDir(w) {
			qs = dir.NewDir(w).Xz().Exclude(ex).GetQuestions()
		}
		if gfile.IsFile(w) {
			qs = dir.NewFile(w).Xz().GetQuestions()
		}
		zk = append(zk, qs...)
	}

	lzk := len(zk)
	if lzk < num {
		num = lzk
		log.Printf("%v, the number of questions is less than %d, so use %d", wf, num, lzk)
	}

	uzk := garray.NewStrArrayFrom(zk).Unique()
	fmt.Println("unique questions: ", len(uzk.Slice()))
	// 随机打乱，再取前n个
	rands := garray.NewStrArrayFrom(zk).Shuffle().SubSlice(0, num)
	rt := dir.GenerateMD(rands)
	fmt.Println(rt)
	return nil
}
