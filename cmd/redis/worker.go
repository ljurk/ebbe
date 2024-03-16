package rediscmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ljurk/ebbe/helper"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("worker called")
		fmt.Println("redis")
		url, _ := cmd.Flags().GetString("url")
		host, _ := cmd.Flags().GetString("host")
		fmt.Println(url)
		fmt.Println(host)
		opt, err := redis.ParseURL(url) //"redis://<user>:<pass>@localhost:6379/<db>")
		if err != nil {
			fmt.Print(err)
			panic(err)
		}

		client := redis.NewClient(opt)
		defer client.Close()

		ctx := context.Background()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Perform cleanup actions in a goroutine
		go func() {
			<-sigChan // Wait for termination signal
			fmt.Println("Received SIGTERM. Performing cleanup...")
			unregisterNode(client, ctx)
			// Perform cleanup tasks here
			fmt.Println("Cleanup complete. Exiting.")
			os.Exit(0)
		}()
		nodeId := int64(registerNode(client, ctx))
		fmt.Printf("NodeId: %d\n", nodeId)
		//readFromRedis(client, ctx, host)
		dataChan := make(chan string)

		// Start the listener goroutine
		go listenForChanges(client, ctx, nodeId, dataChan)

		// Start the sender goroutine
		go sendDataOverSocket(dataChan, host)

		// Keep the main goroutine running
		select {}

	},
}

func unregisterNode(client *redis.Client, ctx context.Context) int {
	// get current number
	fmt.Println("UNREGISTER NODE")
	val, err := client.Get(ctx, "numberOfNodes").Result()
	fmt.Printf("current number:%s\n", val)
	if err != nil {
		if err != redis.Nil {
			fmt.Print(err)
			panic(err)
		}
		val = "0"
	}
	numberOfNodes, err := strconv.Atoi(val)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	if numberOfNodes == 0 {
		return 0
	}
	err = client.Set(ctx, "numberOfNodes", strconv.Itoa(numberOfNodes-1), 0).Err()
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	return numberOfNodes
}
func registerNode(client *redis.Client, ctx context.Context) int {
	// get current number
	val, err := client.Get(ctx, "numberOfNodes").Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Print(err)
			panic(err)
		}
		val = "0"
	}
	numberOfNodes, err := strconv.Atoi(val)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	err = client.Set(ctx, "numberOfNodes", strconv.Itoa(numberOfNodes+1), 0).Err()
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	return numberOfNodes
}

func listenForChanges(client *redis.Client, ctx context.Context, nodeId int64, dataChan chan<- string) {
	var oldElements []string
	var oldNumberOfNodes int64
	for {
		size, err := client.LLen(ctx, "commands").Result()
		if err != nil {
			panic(err)
		}
		snumberOfNodes, err := client.Get(ctx, "numberOfNodes").Result()
		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		numberOfNodes, err := strconv.ParseInt(snumberOfNodes, 10, 64)
		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		// calculate start and stop index
		start := size / numberOfNodes * nodeId
		// stop := size/numberOfNodes*nodeId + (size / numberOfNodes)
		stop := start + (size / numberOfNodes) - 1

		// Ensure the last node gets any remaining elements
		if nodeId == numberOfNodes-1 {
			stop = size - 1
		}

		elements, err := client.LRange(ctx, "commands", start, stop).Result()
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(elements, oldElements) || numberOfNodes != oldNumberOfNodes {
			fmt.Printf("update length:%d numberOfNodes:%d\n", len(elements), numberOfNodes)
			fmt.Printf("start:%d stop:%d total:%d\n", start, stop, size)
			oldNumberOfNodes = numberOfNodes
			// dd := strings.SplitAfter(val, "\n")
			oldElements = elements //strings.Join(helper.SplitList(dd, numberOfNodes)[nodeId], "")
			dataChan <- strings.Join(elements, "")
			// dataChan <- strings.Join(helper.SplitList(dd, int(numberOfNodes))[nodeId], "") // Send the updated data to the sender goroutine
			time.Sleep(time.Second) // Add a delay between successive reads
		}
	}
}

func sendDataOverSocket(dataChan <-chan string, host string) {
	socket, err := helper.OpenSockets(1, host)
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return
	}
	defer func() {
		if err := socket[0].Close(); err != nil {
			fmt.Println("Error closing socket:", err)
		}
	}()

	var cachedData []string // To cache data from dataChan

	for {
		select {
		case data := <-dataChan:
			// Cache the data
			cachedData = strings.SplitAfter(data, "\n")
			cachedData = helper.PacketizeList(cachedData, 1024)
		default:
			// If no new data, send the cached data without delay
			for _, val := range cachedData {
				//for _, val := range strings.SplitAfter(val, "\n") {
				fmt.Fprint(socket[0], val)
				//}
			}
			// Clear the cache
			//cachedData = nil
			// Add a delay before checking for new data
			time.Sleep(time.Millisecond) // Adjust the delay as needed
		}
	}
}
func init() {
	workerCmd.Flags().StringP("url", "u", "redis://localhost:6379/0", "Redis URL")
	workerCmd.Flags().StringP("host", "h", "", "pixelserver URL")
	redisCmd.AddCommand(workerCmd)
}
