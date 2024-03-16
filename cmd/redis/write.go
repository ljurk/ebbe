package rediscmd

import (
	"context"
	"fmt"

	"github.com/ljurk/ebbe/helper"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "write pixels to redis",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("write called")
		url, _ := cmd.Flags().GetString("url")
		input, _ := cmd.Flags().GetString("input")
		opt, err := redis.ParseURL(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		client := redis.NewClient(opt)
		defer client.Close()

		ctx := context.Background()
		var inputData []string
		if input == "-" {
			inputData, err = helper.ReadFromStdin()
		} else {
			inputData, err = helper.ReadFromFile(input)
		}

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		writeToRedis(client, ctx, inputData)
		helper.PrintByteSize(inputData)
		fmt.Println("Finished")
	},
}

func writeToRedis(client *redis.Client, ctx context.Context, input []string) {

	//clear list
	err := client.Del(ctx, "commands").Err()
	if err != nil {
		panic(err)
	}
	// Writing list to Redis
	err = client.LPush(ctx, "commands", input).Err()
	if err != nil {
		panic(err)
	}
	// err := client.Set(ctx, "foo", strings.Join(input, ""), 0).Err()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

func init() {
	writeCmd.Flags().StringP("input", "i", "", "write to redis")
	writeCmd.Flags().StringP("url", "u", "redis://localhost:6379/0", "redis url")
	err := writeCmd.MarkFlagRequired("input")
	if err != nil {
		fmt.Println(err)
	}
	redisCmd.AddCommand(writeCmd)
}
