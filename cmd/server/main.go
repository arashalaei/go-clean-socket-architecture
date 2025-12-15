package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp"
	"github.com/arashalaei/go-clean-socket-architecture/pkg/config"
	"github.com/spf13/cobra"
)

func main(cfg *config.Config) {
	printBanner()
	server := tcp.NewServer(
		tcp.WithCfg(mapToSrvCfg(&cfg.Server)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Register handlers
	server.RegisterHandler(tcp.CreateClass, server.CreateClassHandler)
	server.RegisterHandler(tcp.CreatePerson, server.CreatePersonHandler)
	server.RegisterHandler(tcp.CreateSchool, server.CreateSchoolHandler)
	server.RegisterHandler(tcp.WhoAmI, server.WhoAmIHandler)

	<-stop
	log.Println("Shutdown signal received")

	err = server.Shutdown()
	if err != nil {
		log.Fatal(err)
	}
}

func mapToSrvCfg(cfg *config.ServerConfig) tcp.SrvCfg {
	return tcp.SrvCfg{
		Network:         cfg.Network,
		Address:         cfg.Address,
		ReadTimeout:     cfg.Timeouts.Read,
		WriteTimeout:    cfg.Timeouts.Write,
		IdleTimeout:     cfg.Timeouts.Idle,
		ShutdownTimeout: cfg.Timeouts.Shutdown,
		MaxConnections:  cfg.Limits.MaxConnectionsSize,
		MaxMessageSize:  cfg.Limits.MaxMessageSize,
	}
}

func printBanner() {
	fmt.Println(`
	_____                          
	/ ____|                         
	| (___   ___ _ ____   _____ _ __ 
	\___ \ / _ \ '__\ \ / / _ \ '__|
	____) |  __/ |   \ V /  __/ |   
	|_____/ \___|_|    \_/ \___|_|   

	━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	Author : Arash Alaei
	GitHub : github.com/arashalaei
	Version: 1.0.0
	━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	`)
}

func Register(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "run tcp server",
		Run: func(cmd *cobra.Command, args []string) {
			path, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatal(err)
			}

			cfg, err := config.Load(path)
			if err != nil {
				log.Fatal(err)
			}
			main(cfg)
		},
	}

	cmd.Flags().StringP("config", "c", "", "The config path")
	cmd.MarkFlagRequired("config")
	root.AddCommand(cmd)
}
