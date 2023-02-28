package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/urfave/cli"
)

var (
	cmds  []cli.Command
	flags []cli.Flag
	n     int
)

const (
	RegDetails       = `<details>[\s\S]*?</details>`
	RegUnorderedList = `(?m)^-\s(.*)`
	RegMD            = `*.md`
	MarkDel          = "~~"
	MarkURL          = "http"
	MarkQuestionEN   = "?"
	MarkQuestionCN   = "？"
)

type Dir struct {
	Name  string
	Files []File
}

type File struct {
	Name      string
	Questions []string
	Num       int
}

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
		isDir := gfile.IsDir(w)
		if isDir {
			files, err := gfile.ScanDir(w, RegMD, true)
			if err != nil {
				return cli.NewExitError(err.Error(), 2)
			}
			for _, file := range files {
				isFile := gfile.IsFile(file)
				if !isFile {
					return cli.NewExitError("not a file", 2)
				}
				fArr := strings.Split(file, "/")
				if !garray.NewStrArrayFrom(ex).Contains(fmt.Sprintf("%s/%s", fArr[len(fArr)-2], fArr[len(fArr)-1])) {
					qs := ExtractQuestion(file)
					res = append(res, [][]string{{w, fArr[len(fArr)-1], fmt.Sprintf("%d", len(qs))}}...)
				}
			}
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Dir", "File", "Count"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoMergeCells(true)
	table.AppendBulk(res)
	table.Render()

	return nil
}

func Action(c *cli.Context) error {
	wf := c.StringSlice("wf")
	ex := c.StringSlice("exclude")
	num := c.Int("num")
	res := make([]string, 0)

	for _, w := range wf {
		isDir := gfile.IsDir(w)
		if isDir {
			files, err := gfile.ScanDir(w, RegMD, true)
			if err != nil {
				return cli.NewExitError(err.Error(), 2)
			}
			for _, file := range files {
				isFile := gfile.IsFile(file)
				if !isFile {
					return cli.NewExitError("not a file", 2)
				}
				fArr := strings.Split(file, "/")
				if !garray.NewStrArrayFrom(ex).Contains(fmt.Sprintf("%s/%s", fArr[len(fArr)-2], fArr[len(fArr)-1])) {
					qs := ExtractQuestion(file)
					res = append(res, qs...)
				}
			}
		} else {
			isFile := gfile.IsFile(w)
			if !isFile {
				return cli.NewExitError("not a file", 2)
			}
			fArr := strings.Split(w, "/")
			if !garray.NewStrArrayFrom(ex).Contains(fmt.Sprintf("%s/%s", fArr[len(fArr)-2], fArr[len(fArr)-1])) {
				qs := ExtractQuestion(w)
				res = append(res, qs...)
			}
		}
	}

	count := len(res)
	if count < num {
		num = count
		log.Printf("%v, the number of questions is less than %d, so use %d", wf, num, count)
	}

	// 随机打乱，再取前n个
	rands := garray.NewStrArrayFrom(res).Shuffle().Rands(num)
	rt := GenerateMD(rands)
	fmt.Println(rt)

	return nil
}

// ExtractQuestion 从md中提取问题
func ExtractQuestion(file string) []string {
	fb := gfile.GetContents(file)
	reg := regexp.MustCompile(RegDetails)
	ff := reg.ReplaceAllString(fb, "")

	reg = regexp.MustCompile(RegUnorderedList)
	ss := reg.FindAllString(ff, -1)

	// 剔除所有有url以及没有？的
	for i := 0; i < len(ss); i++ {
		if strings.Contains(ss[i], MarkURL) || strings.Contains(ss[i], MarkDel) || (!strings.Contains(ss[i], MarkQuestionCN) && !strings.Contains(ss[i], MarkQuestionEN)) {
			ss = append(ss[:i], ss[i+1:]...)
			i--
		}
	}
	return ss
}

// GenerateMD 生成最终的md文档
func GenerateMD(qs []string) (rt string) {
	for i := 0; i < len(qs); i++ {
		rt += "- [ ] " + strings.Replace(qs[i], "- ", "", -1) + "\n\n"
	}
	return
}
