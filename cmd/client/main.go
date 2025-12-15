package client

import "github.com/spf13/cobra"

//  Application entry point, dependency injection

func client() {

}

func Register(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "run tcp client",
		Run: func(cmd *cobra.Command, args []string) {
			client()
		},
	}

	rootCmd.AddCommand(cmd)
}
