package utils

import (
	"log"
	"os"
	"strings"

	"github.com/gogf/gf/v2/text/gregex"
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

func NewDir(name string) *Dir {
	return &Dir{Name: name}
}

func (d *Dir) Xz() *Dir {
	dir := d.Name
	dirPath, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("read dir %s error: %v", dir, err)
	}
	sep := string(os.PathSeparator)
	for _, fi := range dirPath {
		// 过滤指定格式的文件
		ok := strings.HasSuffix(fi.Name(), MarkMD)
		if ok {
			filename := dir + sep + fi.Name()
			qs := ExtractQuestion(filename)
			d.Files = append(d.Files, File{
				Name:      fi.Name(),
				Questions: qs,
				Num:       len(qs),
			})
		}
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

func (d *Dir) AddFile(name string, questions []string) {
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
		qs = append(qs, file.Questions...)
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
