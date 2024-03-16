package rediscmd

import (
	"github.com/ljurk/ebbe/cmd"
	"github.com/spf13/cobra"
)

var (
	redisCmd = &cobra.Command{
		Use:   "redis",
		Short: "Get Pixels from redis database",
	}
)

func init() {
	cmd.RootCmd.AddCommand(redisCmd)
}
