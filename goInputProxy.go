package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	logs, err := os.OpenFile("/dev/ttyS0", os.O_RDWR, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer logs.Close()
	fmt.Fprintf(logs, "inputProxy: started!\n")

	env := os.Environ()
	args, envvars := strings.Join(os.Args, " "), strings.Join(env, "\n")

	logfile, err := os.OpenFile("/tmp/LOG_CGIBIN", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(logs, "inputProxy: opening /tmp/LOG_CGIBIN failed, %s\n", err)
		os.Exit(1)
	}
	defer logfile.Close()

	bufLen, stdin, l := 0, []byte{}, os.Getenv("CONTENT_LENGTH")
	if l != "" {
		if bufLen, err = strconv.Atoi(l); err != nil {
			fmt.Fprintf(logs, "inputProxy: couldn't get CONTENT_LENGTH, %s\n", err)
			os.Exit(1)
		}
		reader := bufio.NewReader(os.Stdin)
		stdin, err = reader.Peek(bufLen)
		if err != nil {
			fmt.Fprintf(logs, "inputProxy: opening Stdin failed, %s\n", err)
			os.Exit(1)
		}
		os.Stdin.Close()
	}

	fmt.Fprintf(logs, "inputProxy: calling %s: stdin=%s, args=%s, env=\n%s\n", os.Args[0], stdin, args, envvars)

	var bStdout, bStderr bytes.Buffer
	cmd := exec.Command("/htdocs/cgibin_")
	cmd.Stdout, cmd.Stderr, cmd.Stdin, cmd.Env, cmd.Args = &bStdout, &bStderr, strings.NewReader(fmt.Sprintf("%s", stdin)), env, os.Args
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(logs, "inputProxy: error executing /htdocs/cgibin_, %s\n", err)
		os.Exit(254)
	}
	cmd.Wait()
	stdout, stderr := string(bStdout.Bytes()), string(bStderr.Bytes())
	fmt.Print(stdout)
	_, err = logfile.WriteString(args + "\nSTDIN: ____\n" + string(stdin) + "\nENV: ____\n" + envvars + "\nSTDOUT: ____\n" + stdout + "\nSTDERR: ____\n" + stderr + "\n________\n")
	if err != nil {
		fmt.Fprintf(logs, "inputProxy: could not write to /tmp/LOG_CGIBIN, %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(logs, "inputProxy: Finished execution. stdout=%s, stderr=%s\n", stdout, stderr)
}
