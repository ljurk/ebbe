package helper

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	//"io/ioutil"
	"net"
)

func OpenSockets(count int, host string) ([]net.Conn, error) {
	// Connect to the server
	var connections []net.Conn
	for i := 0; i < count; i++ {
		c, err := net.Dial("tcp", host)
		if err != nil {
			return nil, err
		}
		connections = append(connections, c)
	}
	return connections, nil
}

func GetCanvasSize(conn net.Conn) (int, int, error) {
	// Send the command "SIZE" to the server
	command := "SIZE\n"
	_, err := conn.Write([]byte(command))
	if err != nil {
		return 0, 0, fmt.Errorf("error sending command to server: %v", err)
	}

	// Wait for a response for up to 1 second
	err = conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		return 0, 0, fmt.Errorf("error setting timeout: %v", err)
	}
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return 0, 0, fmt.Errorf("error reading response from server: %v", err)
	}

	// Parse the response
	coordinates := strings.Fields(response)
	if len(coordinates) != 3 {
		return 0, 0, fmt.Errorf("unexpected response format")
	}

	// Extract X and Y coordinates
	x, err := strconv.Atoi(coordinates[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error converting X coordinate: %v", err)
	}
	y, err := strconv.Atoi(coordinates[2])
	if err != nil {
		return 0, 0, fmt.Errorf("error converting Y coordinate: %v", err)
	}

	return x, y, nil
}

func SendMessages(conn net.Conn, messages []string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	for {
		for _, val := range messages {
			fmt.Fprint(conn, val)
		}
	}
}
