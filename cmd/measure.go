package cmd

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	bytesPerPixel = 3 // Each pixel is represented by 3 bytes (RGB)
	measureCmd    = &cobra.Command{
		Use:   "measure",
		Short: "dummy pixelflut server that shows througput",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetString("port")

			fmt.Println("Pixelflut server is running on", port)

			// Listen for incoming connections
			ln, err := net.Listen("tcp", port)
			if err != nil {
				fmt.Println("Error listening:", err.Error())
				os.Exit(1)
			}
			defer ln.Close()

			var (
				totalBytes uint64 = 0
				startTime         = time.Now()
				mu         sync.Mutex
			)
			// Function to handle incoming connections
			handleConnection := func(conn net.Conn) {
				defer conn.Close()
				buf := make([]byte, 1024)
				for {
					// Read the message from the connection
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					// Update total byte count
					mu.Lock()
					totalBytes += uint64(n)
					mu.Unlock()
				}
			}
			// Accept connections and handle them in separate goroutines
			go func() {
				for {
					conn, err := ln.Accept()
					if err != nil {
						fmt.Println("Error accepting connection:", err.Error())
						continue
					}
					go handleConnection(conn)
				}
			}()

			// Periodically print throughput
			go func() {
				for {
					time.Sleep(1 * time.Second)
					mu.Lock()
					elapsedTime := time.Since(startTime).Seconds()
					throughput := float64(totalBytes) / elapsedTime
					//fmt.Printf("Throughput: %.2f bytes/sec\n", throughput)
					//fmt.Printf("Throughput: %.2f KB/sec\n", bytesToKB(throughput))
					fmt.Printf("Throughput: %.2f GBit/sec\n", bytesToGB(throughput)*8)
					totalBytes = 0 // Reset totalBytes
					mu.Unlock()
					startTime = time.Now() // Reset startTime
				}
			}()

			// Keep the main goroutine running
			select {}

		},
	}
)

func bytesToKB(bytesPerSec float64) float64 {
	return bytesPerSec / 1024
}

func bytesToMB(bytesPerSec float64) float64 {
	return bytesPerSec / (1024 * 1024)
}
func bytesToGB(bytesPerSec float64) float64 {
	return bytesPerSec / (1024 * 1024 * 1024)
}
func init() {
	measureCmd.Flags().String("port", ":1337", "port to serve")
	RootCmd.AddCommand(measureCmd)
}
