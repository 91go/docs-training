/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package utils

import (
	"math/rand"
	"strconv"

	"github.com/gogf/gf/v2/container/garray"
)

func NewFile(name string) *File {
	return &File{Name: name}
}

func (f *File) Xz() *File {
	f.Questions = ExtractQuestion(f.Name)
	f.Num = len(f.Questions)
	return f
}

func (f *File) GetQuestions() []string {
	return f.Questions
}

// GetTableData 组装tablewriter需要的数据
func (f *File) GetTableData(dirname string, total int) (data [][]string) {
	if total == 0 {
		total = f.Num
	}
	data = append(data, []string{dirname, f.Name, strconv.Itoa(f.Num), strconv.Itoa(total)})
	return
}

func (f *File) ConvertToMarkdown() (res string) {
	az := garray.NewStrArrayFrom(f.Questions)
	// 随机数
	if az.Len() == 0 {
		return
	}
	azi := rand.Intn(az.Len())
	azz, _ := az.Get(azi)
	_ = az.Set(azi, ReplaceUnorderedListWithTask(azz))

	res += "## " + f.Name + "\n\n"
	res += az.Join("\n") + "\n"
	res += "\n\n"
	return
}
