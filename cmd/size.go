package cmd

import (
	"fmt"

	"github.com/ljurk/ebbe/helper"
	"github.com/spf13/cobra"
)

var (
	sizeCmd = &cobra.Command{
		Use:   "size",
		Short: "outputs canvas size",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			connections, err := helper.OpenSockets(1, host)
			if err != nil {
				return
			}
			x, y, err := helper.GetCanvasSize(connections[0])
			if err != nil {
				return
			}
			fmt.Printf("size %d %d\n", x, y)
		},
	}
)

func init() {
	sizeCmd.Flags().StringVar(&host, "host", ":1337", "address of pixelflut server")
	rootCmd.AddCommand(sizeCmd)
}
