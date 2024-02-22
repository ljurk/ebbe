package cmd

import (
	"fmt"
	"image"
	"strings"

	"github.com/ljurk/ebbe/render"
	"github.com/spf13/cobra"
)

var (
	text    string
	textCmd = &cobra.Command{
		Use:   "text",
		Short: "Creates a command list based on text",
		Run: func(cmd *cobra.Command, args []string) {
			x, _ := cmd.Flags().GetInt("x")
			y, _ := cmd.Flags().GetInt("y")
			fontsize, _ := cmd.Flags().GetFloat64("fontsize")
			hexColor, _ := cmd.Flags().GetString("color")
			if !strings.HasPrefix(hexColor, "#") {
				hexColor = "#" + hexColor
			}

			// Parse the hex color string
			c, err := render.HexToRGBA(hexColor)

			if err != nil {
				fmt.Println("Error parsing hex color:", err)
				return
			}

			textColor := image.NewUniform(c)
			textImage := render.RenderText(text, fontsize, textColor, image.Transparent)
			comms := render.CommandsFromImage(textImage, render.NewOrder("l"), image.Point{x, y})
			for _, i := range comms.ToString() {
				fmt.Print(i)
			}

		},
	}
)

func init() {
	textCmd.Flags().StringVarP(&text, "text", "t", "", "text to print (required)")
	textCmd.Flags().StringP("color", "c", "000000", "color in hex")
	textCmd.Flags().Float64P("fontsize", "f", 10, "fontsize")
	textCmd.Flags().IntP("x", "x", 0, "x coordinate of image (default: 0)")
	textCmd.Flags().IntP("y", "y", 0, "y coordinate of image (default: 0)")
	err := textCmd.MarkFlagRequired("text")
	if err != nil {
		fmt.Println(err)
	}
	rootCmd.AddCommand(textCmd)
}
