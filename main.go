package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"strings"
	"sync"

	"github.com/ljurk/ebbe/helper"

	"github.com/ljurk/ebbe/render"
)

func main() {
	// define cli flags
	fmt.Print("ok")
	host := flag.String("host", "", "Target server address")
	imgPath := flag.String("image", "", "Filepath of an image to flut")
	imgX := flag.Int("ix", 0, "x coordinate of image")
	imgY := flag.Int("iy", 0, "y coordiante of image")
	color := flag.String("color", "", "if set just prints the whole screen in this color")
	text := flag.String("text", "", "text to print")
	textX := flag.Int("tx", 0, "x coordinate of text")
	textY := flag.Int("ty", 0, "y coordiante of text")
	con := flag.Int("con", 4, "Number of simultaneous connections and goroutines. Each routine send a part of the data over its own connection")
	//x := flag.Int("x", 0, "Offset of posted image from left border")
	//y := flag.Int("y", 0, "Offset of posted image from top border")
	random := flag.Bool("random", false, "if true randomize messageList")
	maxPacketSize := flag.Int("packetsize", 1024, "Size of packets to send in bytes")
	flag.Parse()

	if *host == "" {
		flag.Usage()
		os.Exit(1)
	}

	var comms render.Commands
	var messages []string
	var messageParts [][]string
	connections, err := helper.OpenSockets(*con, *host)
	if err != nil {
		fmt.Println("Error opening Sockets :", err)
		return
	}
	width, height, err := helper.GetCanvasSize(connections[0])
	if err != nil {
		fmt.Println("Error getting coordinates:", err)
		return
	}
	fmt.Printf("Canvas Size %d x %d", width, height)

	if *imgPath != "" {
		img, err := render.ReadImage(*imgPath)

		if err != nil {
			return
		}

		comms = render.CommandsFromImage(img, render.NewOrder("l"), image.Point{*imgX, *imgY})
		messages = helper.MergeLists(messages, comms.ToString())
	}

	if *text != "" {
		bg, _ := render.ReadImage("./blank-white.jpg")
		tex, _ := render.ReadImage("./zebra.jpg")
		textImage := render.RenderText(*text, 10, tex, bg)
		comms := render.CommandsFromImage(textImage, render.NewOrder("l"), image.Point{*textX, *textY})
		messages = helper.MergeLists(messages, helper.RemoveColorFromList(comms.ToString(), "ffffff"))
	}
	if *color != "" {
		tmpmessages, _ := render.OnlyColor(width, height, strings.Split(*color, ","))
		messages = helper.MergeLists(messages, tmpmessages)
	}
	if *random {
		messages = helper.ShuffleStrings(messages)
	}

	fmt.Printf("Total messages %d\n", len(messages))
	messages = helper.PacketizeList(messages, *maxPacketSize)
	fmt.Printf("Total packetized messages %d\n", len(messages))
	fmt.Printf("%d Connections opened\n", *con)
	// Split messages into equal parts for each goroutine
	messageParts = helper.SplitList(messages, *con)
	fmt.Printf("Splited messages %d\n", len(messageParts))

	var wg sync.WaitGroup
	wg.Add(*con)
	// Distribute the workload among connections using goroutines
	for i, conn := range connections {
		go helper.SendMessages(conn, messageParts[i], &wg)
	}
	// Wait for all goroutines to finish
	wg.Wait()
}
