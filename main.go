package main

import (
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/urfave/cli"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	flags []cli.Flag
	t     string
	n     int
)

const (
	RegDetails       = `<details>[\s\S]*?</details>`
	RegUnorderedList = `(?m)^-\s(.*)`
)

func init() {
	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "t, tag",
			Value:       "ms",
			Usage:       "tag files",
			Destination: &t,
		},
		cli.IntFlag{
			Name:        "n, num",
			Value:       30,
			Usage:       "num of questions",
			Destination: &n,
		},
		cli.StringSliceFlag{
			Name:  "d, dirs",
			Value: &cli.StringSlice{"/Users/lhgtqb7bll/docs/arch", "/Users/lhgtqb7bll/docs/database"},
			Usage: "dirs of files",
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Weekly-Training"
	app.Usage = "use to collocate with gh-ac, generate weekly training"
	app.HideVersion = true
	app.Flags = flags
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(c *cli.Context) error {

	dirs := c.StringSlice("dirs")
	num := c.Int("num")
	res := make([]string, 0)

	for _, dir := range dirs {
		isDir := gfile.IsDir(dir)
		if !isDir {
			return cli.NewExitError("not a dir", 2)
		}
		files, err := gfile.ScanDir(dir, "*.md", true)
		if err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
		for _, file := range files {
			isFile := gfile.IsFile(file)
			if !isFile {
				return cli.NewExitError("not a file", 2)
			}
			qs := ExtractQuestion(file)
			res = append(res, qs...)
		}
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
		if strings.Contains(ss[i], "http") || (!strings.Contains(ss[i], "？") && !strings.Contains(ss[i], "?")) {
			ss = append(ss[:i], ss[i+1:]...)
			i--
		}
	}
	return ss
}

// GenerateMD 生成最终的md文档
func GenerateMD(qs []string) (rt string) {
	for i := 0; i < len(qs); i++ {
		rt += "- [ ] " + strings.Replace(qs[i], "- ", "", -1) + "\n"
	}
	return
}
