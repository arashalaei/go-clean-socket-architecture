package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp"
	"github.com/arashalaei/go-clean-socket-architecture/pkg/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func client(cfg *config.Config) {
	printBanner()
	client := tcp.NewClient(
		tcp.WithClientCfg(mapToClientCfg(&cfg.Client)),
	)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	prompt := promptui.Select{
		Label: "Select An Action",
		Items: []string{
			"1. Create School",
			"2. Create Class",
			"3. Add Student To Class",
			"4. Who Am I ?",
		},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	<-stop
	log.Println("Closing signal received")
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func mapToClientCfg(cfg *config.ClientConfig) tcp.ClientConfig {
	return tcp.ClientConfig{
		Network:         cfg.Network,
		Address:         cfg.Address,
		ConnectTimeout:  cfg.Timeouts.Connect,
		ReadTimeout:     cfg.Timeouts.Read,
		WriteTimeout:    cfg.Timeouts.Write,
		KeepAlivePeriod: cfg.Timeouts.KeepAlivePeriod,
		KeepAlive:       cfg.Limits.KeepAlive,
		MaxRetries:      cfg.Limits.MaxRetries,
		RetryDelay:      cfg.Limits.RetryDelay,
	}
}

func printBanner() {
	fmt.Println(`
	_____ _ _            _   
	/ ____| (_)          | |  
	| |    | |_  ___ _ __ | |_ 
	| |    | | |/ _ \ '_ \| __|
	| |____| | |  __/ | | | |_ 
	\_____|_|_|\___|_| |_|\__|

	━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	Author : Arash Alaei
	GitHub : github.com/arashalaei
	Version: 1.0.0
	━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	`)
}

func Register(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "run tcp client",
		Run: func(cmd *cobra.Command, args []string) {
			path, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatal(err)
			}

			config, err := config.Load(path)
			if err != nil {
				log.Fatal(err)
			}

			client(config)
		},
	}

	cmd.Flags().StringP("config", "c", "", "client config path")
	cmd.MarkFlagRequired("config")
	rootCmd.AddCommand(cmd)
}
