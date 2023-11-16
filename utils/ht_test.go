package utils

import "testing"

func TestSanitizeParticularPunc(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{"Your string with special characters?“”"},
			want: "your-string-with-special-characters",
		},
		{
			name: "",
			args: args{"How to implement Suffix-Tree with golang?"},
			want: "how-to-implement-suffix-tree-with-golang",
		},
		{
			name: "",
			args: args{"How to implement Suffix-Tree with golang ?"},
			want: "how-to-implement-suffix-tree-with-golang-",
		},
		{
			name: "",
			args: args{"能否自己实现一个 raft 算法？"},
			want: "能否自己实现一个-raft-算法",
		},
		{
			name: "",
			args: args{"RSA+AES 混合加密的流程是什么？"},
			want: "rsaaes-混合加密的流程是什么",
		},
		{
			name: "",
			args: args{"综合：服务端优化？（输入网址后，会发生什么？每一步的优化方案？）"},
			want: "综合服务端优化输入网址后会发生什么每一步的优化方案",
		},
		{
			name: "",
			args: args{"xxx-------x-x-x-x"},
			want: "xxx-------x-x-x-x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeParticularPunc(tt.args.str); got != tt.want {
				t.Errorf("SanitizeParticularPunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
