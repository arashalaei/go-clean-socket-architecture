package cmd

import (
	"fmt"
	"os"

	"github.com/arashalaei/go-clean-socket-architecture/cmd/client"
	"github.com/arashalaei/go-clean-socket-architecture/cmd/server"
	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:     "socket",
		Short:   "School Management System",
		Version: "0.1",
	}

	// Register cmds
	server.Register(rootCmd)
	client.Register(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
