package cmd

import (
	"cobra/internal/word"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

const (
	ModeUpper                      = iota + 1 // 全部单词转为大写
	ModeLower                                 // 全部单词转为小写
	ModeUnderscoreToUpperCamelCase            // 下画线单词转为大写驼峰单词
	ModeUnderscoreToLowerCamelCase            // 下画线单词转为小写驼峰单词
	ModeCamelCaseToUnderscore                 // 驼峰单词转为下画线单词
)

var str string
var mode int8

func init() {
	WordCmd.Flags().StringVarP(&str, "str", "s", "", "请输入单词内容")
	WordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "请输入单词转换的模式")
}

var WordCmd = &cobra.Command{
	Use:   "word",
	Short: "单词格式转换",
	Long: strings.Join([]string{
		"该子命令支持各种单词格式转换，模式如下：",
		"1：全部单词转为大写",
		"2：全部单词转为小写",
		"3：下画线单词转为大写驼峰单词",
		"4：下画线单词转为小写驼峰单词",
		"5：驼峰单词转为下画线单词",
	}, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case ModeUpper:
			content = word.ToUpper(str)
		case ModeLower:
			content = word.ToLower(str)
		case ModeUnderscoreToUpperCamelCase:
			content = word.UnderscoreToUpperCamelCase(str)
		case ModeUnderscoreToLowerCamelCase:
			content = word.UnderscoreToLowerCamelCase(str)
		case ModeCamelCaseToUnderscore:
			content = word.CamelCaseToUnderscore(str)
		default:
			log.Fatalf("暂不支持该转换模式，请执行 help word 查看帮助文档")
		}
		log.Printf("输出结果：%s", content)
	},
}