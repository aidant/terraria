package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func startTerraria() *os.File {
	arguments := append([]string{"-config", "/config/server-config.txt"}, os.Args[1:]...)
	command := exec.Command("./run.sh", arguments...)
	stdio, err := pty.Start(command)
	check(err)
	return stdio
}

func pipeStdin(stdio *os.File) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		_, err := stdio.Write([]byte(scanner.Text() + "\n"))
		check(err)
	}

	err := scanner.Err()
	check(err)
}

func saveOnExit(stdio *os.File) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-quit:
			_, err := stdio.Write([]byte("exit\n"))
			check(err)
		}
	}

}

func main() {
	stdio := startTerraria()
	defer stdio.Close()

	go pipeStdin(stdio)
	go saveOnExit(stdio)

	io.Copy(os.Stdout, stdio)
}
