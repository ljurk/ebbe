package cmd

import (
	"fmt"
	"sync"

	"github.com/ljurk/ebbe/helper"
	"github.com/spf13/cobra"
)

var (
	host        string
	connections int
	packetSize  int
	input       string
	sendCmd     = &cobra.Command{
		Use:   "send",
		Short: "Sends data to pixelflut server",
		Run: func(cmd *cobra.Command, args []string) {
			var inputData []string
			var err error
			oneshot, _ := cmd.Flags().GetBool("oneshot")
			if input == "-" {
				inputData, err = helper.ReadFromStdin()
			} else {
				inputData, err = helper.ReadFromFile(input)
			}
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			send(inputData, oneshot)
		},
	}
)

func send(messages []string, oneshot bool) {
	var messageParts [][]string
	sockets, err := helper.OpenSockets(connections, host)
	if err != nil {
		fmt.Println("Error opening Sockets :", err)
		return
	}
	fmt.Printf("%d Connections opened\n", len(sockets))
	fmt.Printf("Total messages %d\n", len(messages))
	messages = helper.PacketizeList(messages, packetSize)
	fmt.Printf("Total packetized messages %d\n", len(messages))

	messageParts = helper.SplitList(messages, connections)
	//fmt.Print(messageParts)
	var wg sync.WaitGroup
	wg.Add(connections)
	// Distribute the workload among connections using goroutines
	for i, conn := range sockets {
		go helper.SendMessages(conn, messageParts[i], oneshot, &wg)
	}
	// Wait for all goroutines to finish
	wg.Wait()
}

func init() {
	sendCmd.Flags().StringVarP(&host, "host", "h", ":1337", "address of pixelflut server")
	sendCmd.Flags().IntVarP(&connections, "connections", "c", 1, "number of sockets to open")
	sendCmd.Flags().IntVarP(&packetSize, "packetsize", "p", 1024, "size of packets to send in bytes")
	sendCmd.Flags().BoolP("oneshot", "o", false, "send all messages once")
	sendCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path, use '-' for stdin")
	err := sendCmd.MarkFlagRequired("input")
	if err != nil {
		fmt.Println(err)
	}
	rootCmd.AddCommand(sendCmd)
}
