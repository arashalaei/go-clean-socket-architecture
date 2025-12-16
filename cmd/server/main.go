package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp"
	store "github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite"
	"github.com/arashalaei/go-clean-socket-architecture/internal/usecase/class"
	"github.com/arashalaei/go-clean-socket-architecture/internal/usecase/person"
	"github.com/arashalaei/go-clean-socket-architecture/internal/usecase/school"
	"github.com/arashalaei/go-clean-socket-architecture/pkg/config"
	"github.com/spf13/cobra"
)

func main(cfg *config.Config) {
	printBanner()

	// set up repos
	db, err := store.NewSqlite(cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}
	schoolUsecases := school.NewSchoolUseCases(
		school.NewCreateSchoolUseCase(db),
		school.NewListSchoolsUseCase(db),
	)

	classUsecases := class.NewClassUseCases(
		class.NewCreateClassUseCase(db),
		class.NewListClassesUseCase(db),
		class.NewAddStudentToClassUseCase(db),
	)

	personUsecases := person.NewPersonUseCases(
		person.NewCreatePersonUseCase(db),
		person.NewListPersonsUseCase(db),
		person.NewWhoAmIUseCase(db),
		person.NewEnrollInSchoolStudentUseCase(db, db),
	)

	server := tcp.NewServer(
		tcp.WithCfg(mapToSrvCfg(&cfg.Server)),
		tcp.WithSchoolUsecases(*schoolUsecases),
		tcp.WithClassUsecases(*classUsecases),
		tcp.WithPersonUsecases(*personUsecases),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = server.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server.RegisterHandler(tcp.CreateSchool, server.CreateSchoolHandler)
	server.RegisterHandler(tcp.ListSchools, server.ListSchoolsHandler)
	server.RegisterHandler(tcp.CreatePerson, server.CreatePersonHandler)
	server.RegisterHandler(tcp.ListPersons, server.ListPersonsHandler)
	server.RegisterHandler(tcp.CreateClass, server.CreateClassHandler)
	server.RegisterHandler(tcp.ListClasses, server.ListClassesHandler)
	server.RegisterHandler(tcp.AddStudentToClass, server.AddStudentToClassHandler)
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
