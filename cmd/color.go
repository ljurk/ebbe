package cmd

import (
	"fmt"

	"github.com/ljurk/ebbe/render"
	"github.com/spf13/cobra"
)

var (
	colors   []string
	width    int
	height   int
	vertical bool
	colorCmd = &cobra.Command{
		Use:   "color",
		Short: "Creates color commands",
		Run: func(cmd *cobra.Command, args []string) {
			x, _ := cmd.Flags().GetInt("x")
			y, _ := cmd.Flags().GetInt("y")
			var order render.RenderOrder
			if vertical {
				order = render.NewOrder("t")
			} else {
				order = render.NewOrder("l")
			}

			commands, err := render.OnlyColor(x, y, width, height, colors, order)
			if err != nil {
				fmt.Printf("Error creating colors: %v", err)
				return
			}

			for _, i := range commands {
				fmt.Print(i)
			}
		},
	}
)

func init() {
	colorCmd.Flags().StringSliceVarP(&colors, "color", "c", []string{}, "specify multiple colors")
	colorCmd.Flags().BoolVarP(&vertical, "vertical", "v", false, "draw color vertical")
	colorCmd.Flags().IntVarP(&width, "width", "w", 161, "canvas width")
	colorCmd.Flags().IntVarP(&height, "height", "h", 161, "canvas height")
	colorCmd.Flags().IntP("x", "x", 0, "starting position")
	colorCmd.Flags().IntP("y", "y", 0, "startin position")
	err := colorCmd.MarkFlagRequired("color")
	if err != nil {
		fmt.Println(err)
	}
	rootCmd.AddCommand(colorCmd)
}
