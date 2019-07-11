package main

import (
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

func save(file *os.File) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-quit:
			file.Write([]byte("exit\n"))
			return
		}
	}
}

func run() *os.File {
	command := exec.Command("./run.sh", "-config", "/config/server-config.txt")
	file, err := pty.Start(command)
	check(err)
	return file
}

func main() {
	file := run()
	defer file.Close()

	go save(file)

	io.Copy(os.Stdout, file)
}
