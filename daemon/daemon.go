package daemon

import (
	"os"
	"time"
)

func Daemonize() {
	// make sure parent process died before child process
	time.Sleep(time.Second)
	if os.Getppid() != 1 {
		args := append([]string{os.Args[0]}, os.Args[1:]...)
		os.StartProcess(os.Args[0], args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
		os.Exit(0)
	}
}
