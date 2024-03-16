package cmd

import (
	"fmt"
	"image"

	"github.com/ljurk/ebbe/render"
	"github.com/spf13/cobra"
)

var (
	imgPath  string
	x        int
	y        int
	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "Creates a command list based on an image",
		Run: func(cmd *cobra.Command, args []string) {
			img, err := render.ReadImage(imgPath)

			if err != nil {
				return
			}

			for _, i := range render.CommandsFromImage(img, render.NewOrder("l"), image.Point{x, y}).ToString() {
				fmt.Print(i)
			}
		},
	}
)

func init() {
	imageCmd.Flags().StringVarP(&imgPath, "image", "i", "", "path to image (required)")
	imageCmd.Flags().IntVarP(&x, "x", "x", 0, "x coordinate of image (default: 0)")
	imageCmd.Flags().IntVarP(&y, "y", "y", 0, "y coordinate of image (default: 0)")
	err := imageCmd.MarkFlagRequired("image")
	if err != nil {
		fmt.Println(err)
	}

	RootCmd.AddCommand(imageCmd)
}
