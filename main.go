package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/net/context"
	kirk "qiniupkg.com/kirk/kirksdk"
)

var args struct {
	// Verbose bool `short:"v" long:"verbose" description:"Show verbose information"`
	AccessKey string `long:"access-key" description:"Kirk AccessKey" required:"true"`
	SecretKey string `long:"secret-key" description:"Kirk SecretKey" required:"true"`
	Host      string `long:"host" description:"Kirk API Host" required:"true"`
	IP        string `long:"ip" description:"Kirk Container IP" required:"true"`
}

func main() {
	var (
		commands []string
		err      error
		ret      kirk.ExecContainerRet
	)
	if commands, err = flags.ParseArgs(&args, os.Args); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	commands = commands[1:] // Remove first args, which is kirk-ssh itself
	if len(commands) == 0 {
		commands = []string{"/bin/bash"}
	}

	kirkConfig := kirk.QcosConfig{
		AccessKey: args.AccessKey,
		SecretKey: args.SecretKey,
		Host:      args.Host,
	}

	client := kirk.NewQcosClient(kirkConfig)
	execArgs := kirk.ExecContainerArgs{
		Command: commands,
	}

	if ret, err = client.ExecContainer(context.TODO(), args.IP, execArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute container: %s\n", err)
		os.Exit(1)
	}

	errChan := make(chan error)

	go func() {
		if err := <-errChan; err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start container: %s\n", err)
			os.Exit(1)
		}
	}()

	startArgs := kirk.StartContainerExecArgs{Mode: "auto"}
	startOpts := kirk.StartContainerExecOpts{
		InStream:  os.Stdin,
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
		ErrorCh:   errChan,
	}
	if err = client.StartContainerExec(context.TODO(), args.IP, ret.ExecID, startArgs, startOpts); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start container: %s\n", err)
		os.Exit(1)
	}
}
