package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ebbe",
	Short: "A Tool to create Pixelflut commands and send them to a server",
	Long: `ebbe is a modular application, there are commands to create pixelflut commands and there is one to send data to a pixelflut server.
to combine these commands, you can either pipe them together:

  ebbe image --image enton.png | ebbe send --host :1337 --input -

or write them to file, merge them with other commands and then send it:

  ebbe image --image enton.png > data.txt
  ebbe pattern --color 000000,ffffff >> data.txt
  ebbe send --input data.txt

In the above example the image will fight against the color. To remove all image pixels from color pixels you can run:

  ebbe merge --input colors.txt --input image.txt | ebbe send -i -`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	RootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")
}
