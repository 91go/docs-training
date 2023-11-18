package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/91go/docs-training/utils"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/gogf/gf/v2/text/gregex"
)

type Dir struct {
	Name  string
	Files []File
	// SubDirs []*Dir
}

type File struct {
	Name      string
	Questions []Question
	Num       int
}

type Question struct {
	Text string
	URL  string
}

func NewDir(name string) *Dir {
	return &Dir{Name: name}
}

func (d *Dir) Xz(fn func(string) []Question) *Dir {
	err := filepath.Walk(d.Name, func(path string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if fi.IsDir() {
			return nil
		}
		// 过滤指定格式的文件
		if strings.HasSuffix(fi.Name(), utils.MarkMD) {
			// qs := ExtractQuestion(path)
			qs := fn(path)
			d.Files = append(d.Files, File{
				Name:      fi.Name(),
				Questions: qs,
				Num:       len(qs),
			})
		}
		return nil
	})
	if err != nil {
		return nil
	}

	return d
}

// Exclude 根据文件名，排除指定文件
// 直接写文件名，不需要带路径。比如devops.md、mysql.md等，否则无法匹配。
func (d *Dir) Exclude(names []string) *Dir {
	for _, name := range names {
		for i, file := range d.Files {
			if file.Name == name {
				d.Files = append(d.Files[:i], d.Files[i+1:]...)
			}
		}
	}
	return d
}

func (d *Dir) AddFile(name string, questions []Question) {
	d.Files = append(d.Files, File{Name: name, Questions: questions, Num: len(questions)})
}

func (d *Dir) AddFiles(files []File) {
	d.Files = append(d.Files, files...)
}

func (d *Dir) GetFiles() []File {
	return d.Files
}

func (d *Dir) GetFile(name string) *File {
	for _, file := range d.Files {
		if file.Name == name {
			return &file
		}
	}
	return nil
}

func (d *Dir) GetFileNum() int {
	return len(d.Files)
}

// GetQuestionNum 获取所有题目数量
func (d *Dir) GetQuestionNum() int {
	return len(d.GetQuestions())
}

// GetQuestionNumByFile 获取指定文件的题目数量
func (d *Dir) GetQuestionNumByFile(name string) int {
	for _, file := range d.Files {
		if file.Name == name {
			return file.Num
		}
	}
	return 0
}

// GetQuestionNumByFiles 获取指定文件的题目数量
func (d *Dir) GetQuestionNumByFiles(names []string) int {
	var num int
	for _, name := range names {
		num += d.GetQuestionNumByFile(name)
	}
	return num
}

// GetQuestionNumByFileReg 获取指定文件的题目数量
func (d *Dir) GetQuestionNumByFileReg(reg string) int {
	var num int
	for _, file := range d.Files {
		if gregex.IsMatchString(reg, file.Name) {
			num += file.Num
		}
	}
	return num
}

func (d *Dir) GetQuestionNumByFileRegEx(reg string) int {
	var num int
	for _, file := range d.Files {
		if gregex.IsMatchString(reg, file.Name) {
			num += file.Num
		}
	}
	return num
}

func (d *Dir) GetQuestionNumByFileRegExs(regs []string) int {
	var num int
	for _, file := range d.Files {
		for _, reg := range regs {
			if gregex.IsMatchString(reg, file.Name) {
				num += file.Num
			}
		}
	}
	return num
}

// GetQuestions 获取所有Questions
func (d *Dir) GetQuestions() (qs []string) {
	for _, file := range d.Files {
		var res []string
		// flatten Question struct
		for _, q := range file.Questions {
			res = append(res, fmt.Sprintf("%s [%s](%s)", q.Text, q.Text, q.URL))
		}

		qs = append(qs, res...)
	}
	return
}

func (d *Dir) GetInterviews() (qs []string) {
	for _, file := range d.Files {
		var res []string
		// flatten Question struct
		for _, q := range file.Questions {
			res = append(res, q.Text)
		}

		qs = append(qs, res...)
	}
	return
}

// GetTableData 组装tablewriter需要的数据
func (d *Dir) GetTableData() (data [][]string) {
	for _, file := range d.Files {
		data = append(data, file.GetTableData(d.Name, d.GetQuestionNum())...)
	}
	return
}

// InterviewsToMarkdown Convert Dir to Markdown string
func (d *Dir) InterviewsToMarkdown(count int) (res string) {
	for _, file := range d.Files {
		rsk := make([]string, count)
		qLen := grand.Perm(len(file.Questions))[:count]
		for i, index := range qLen {
			rsk[i] = file.Questions[index].Text
		}

		res += fmt.Sprintf("## %s \n\n", file.Name)
		res += fmt.Sprintf("%s \n\n", garray.NewStrArrayFrom(rsk).Join("\n"))
	}
	return
}

const (
	RegHeaders       = `(?m)(#+ \*|#+ \*\*|#+ \*\*\*|#+ \*\*\*\*)\s*([^*]+)`
	RegUnorderedList = `(?m)^-\s(.*)`
)

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
		fds := strings.ReplaceAll(file, ".md", "")
		// determine whether duplicate
		fx := garray.NewStrArrayFrom(strings.Split(fds, "/")).Unique().Join("/")
		qURL := fmt.Sprintf("%s#%s", fx, headerCts)
		extractedHeaders = append(extractedHeaders, Question{
			Text: headerCts,
			URL:  fmt.Sprintf("%s%s", os.Getenv("BaseURL"), utils.SanitizeParticularPunc(qURL)),
		})
	}

	return
}

func ExtractInterviews(file string) (extractedHeaders []Question) {
	fb := gfile.GetContents(file)
	reg := regexp.MustCompile(RegUnorderedList)
	headers := reg.FindAllString(fb, -1)

	for _, header := range headers {
		extractedHeaders = append(extractedHeaders, Question{
			Text: header,
			URL:  "",
		})
	}

	return
}
