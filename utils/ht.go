package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"os"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"

	"github.com/olekukonko/tablewriter"
)

const (
	RegDetails       = `<details>[\s\S]*?</details>`
	RegUnorderedList = `(?m)^-\s(.*)`
	// RegHeaders       = `^#{1,6}\s+(.*)`
	// RegHeaders = `(?m)^\*{1,3}(.+?)\*{1,3}(.*)$`
	// RegHeaders     = `(?m)(?<=# \*\*|## \*\*|### \*\*\*)[^*]+`
	RegHeaders     = `(?m)(#+ \*\*|#+ \*\*\*|#+ \*\*\*\*)\s*([^*]+)`
	RegMD          = `*.md`
	MarkMD         = `.md`
	MarkDel        = "~~"
	MarkURL        = "http"
	MarkQuestionEN = "?"
	MarkQuestionCN = "？"
)

// ExtractQuestion 从md中提取问题
// func ExtractQuestion(file string) []string {
// 	fb := gfile.GetContents(file)
// 	reg := regexp.MustCompile(RegDetails)
// 	ff := reg.ReplaceAllString(fb, "")
//
// 	reg = regexp.MustCompile(RegUnorderedList)
// 	ss := reg.FindAllString(ff, -1)
//
// 	// 剔除所有有url以及没有？的
// 	for i := 0; i < len(ss); i++ {
// 		if strings.Contains(ss[i], MarkURL) || strings.Contains(ss[i], MarkDel) || (!strings.Contains(ss[i], MarkQuestionCN) && !strings.Contains(ss[i], MarkQuestionEN)) {
// 			ss = append(ss[:i], ss[i+1:]...)
// 			i--
// 		}
// 	}
// 	return ss
// }

// ExtractQuestion 从md中提取问题
func ExtractQuestion(file string) (extractedHeaders []Question) {
	fb := gfile.GetContents(file)
	reg := regexp.MustCompile(RegHeaders)
	headers := reg.FindAllStringSubmatch(fb, -1)

	// reg = regexp.MustCompile(RegUnorderedList)
	// ss := reg.FindAllString(ff, -1)

	// 剔除所有有url以及没有？的
	// for i := 0; i < len(ss); i++ {
	// 	if strings.Contains(ss[i], MarkURL) || strings.Contains(ss[i], MarkDel) || (!strings.Contains(ss[i], MarkQuestionCN) && !strings.Contains(ss[i], MarkQuestionEN)) {
	// 		ss = append(ss[:i], ss[i+1:]...)
	// 		i--
	// 	}
	// }

	for _, header := range headers {
		headerCts := header[2]
		fs := strings.ReplaceAll(file, ".md", "")
		// determine whether duplicate
		fx := garray.NewStrArrayFrom(strings.Split(fs, "/")).Unique().Join("/")
		extractedHeaders = append(extractedHeaders, Question{
			text: headerCts,
			url:  fmt.Sprintf("%s%s#%s", os.Getenv("BaseURL"), fx, headerCts),
		})
	}

	return extractedHeaders
}

// GenerateMD 生成最终的md文档
func GenerateMD(qs []string) (rt string) {
	for i := 0; i < len(qs); i++ {
		rt += ReplaceUnorderedListWithTask(qs[i])
	}
	return
}

// ReplaceUnorderedListWithTask 将无序列表替换为任务列表
func ReplaceUnorderedListWithTask(str string) string {
	return "- [ ] " + strings.Replace(str, "- ", "", -1) + "\n\n"
}

// GenerateMDTable 生成md表格
func GenerateMDTable(res [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Dir", "File", "Count", "Total"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoMergeCells(true)
	table.AppendBulk(res)
	table.Render()
}

// SanitizeParticularPunc
// determine whether same name
// remove particular punctuations
func SanitizeParticularPunc(str string) string {

	str = strings.ReplaceAll(str, " ", "-")
	reg := regexp.MustCompile(`[\"?？“”【】]+`)
	return reg.ReplaceAllString(str, "")
}
