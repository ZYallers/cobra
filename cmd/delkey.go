package cmd

import (
	"cobra/internal/delkey"
	"github.com/spf13/cobra"
	"strings"
)

var host string
var port string
var auth string
var match string
var count int64
var ttl int8
var logDir string

func init() {
	DelKeyCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1", "请输入Redis服务地址")
	DelKeyCmd.Flags().StringVarP(&port, "port", "p", "6379", "请输入Redis服务端口")
	DelKeyCmd.Flags().StringVarP(&auth, "auth", "a", "", "请输入Redis服务密码")
	DelKeyCmd.Flags().StringVarP(&match, "match", "m", "", "请输入要删除keys的指定前缀")
	DelKeyCmd.Flags().Int64VarP(&count, "count", "c", 100, "请输入每次扫描的返回数量")
	DelKeyCmd.Flags().Int8VarP(&ttl, "ttl", "t", -1, "请输入要清除缓存的类型(-1永久缓存|1所有)")
	DelKeyCmd.Flags().StringVarP(&logDir, "logDir", "l", ".", "请输入日志目录")
}

var DelKeyCmd = &cobra.Command{
	Use:   "delKey",
	Short: "优雅删除Redis指定前缀的keys",
	Long: strings.Join([]string{
		"该子命令支持优雅删除Redis指定前缀的keys，优点如下：",
		"1：通过SCAN命令遍历符合要求的keys，不会造成阻塞",
		"2：可以指定是否只删除永久缓存的keys",
		"3：操作有日志记录(后悔药)",
	}, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		delkey.Processor(host, port, auth, match, count, ttl, logDir)
	},
}
