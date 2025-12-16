package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp"
	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
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

	// Run menu in a goroutine so we can handle shutdown signals
	go func() {
		runMainMenu(client)
	}()

	<-stop
	log.Println("Closing signal received")
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func runMainMenu(client *tcp.Client) {
	for {
		prompt := promptui.Select{
			Label: "Main Menu - Select Category",
			Items: []string{
				"1. School",
				"2. Class",
				"3. Person",
				"4. Exit",
			},
		}

		selected, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch selected {
		case 0:
			runSchoolMenu(client)
		case 1:
			runClassMenu(client)
		case 2:
			runPersonMenu(client)
		case 3:
			fmt.Println("Exiting...")
			return
		default:
			return
		}
	}
}

func runSchoolMenu(client *tcp.Client) {
	for {
		prompt := promptui.Select{
			Label: "School Menu - Select Action",
			Items: []string{
				"1. Add New School",
				"2. List All Schools",
				"3. Back to Main Menu",
			},
		}

		selected, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch selected {
		case 0:
			handleCreateSchool(client)
		case 1:
			handleListSchools(client)
		case 2:
			return
		default:
			return
		}
	}
}

func runClassMenu(client *tcp.Client) {
	for {
		prompt := promptui.Select{
			Label: "Class Menu - Select Action",
			Items: []string{
				"1. Add New Class",
				"2. List All Classes",
				"3. Add Student To Class",
				"4. Back to Main Menu",
			},
		}

		selected, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch selected {
		case 0:
			handleCreateClass(client)
		case 1:
			handleListClasses(client)
		case 2:
			handleAddStudentToClass(client)
		case 3:
			return
		default:
			return
		}
	}
}

func runPersonMenu(client *tcp.Client) {
	for {
		prompt := promptui.Select{
			Label: "Person Menu - Select Action",
			Items: []string{
				"1. Add New Person",
				"2. List All Persons",
				"3. Who Am I?",
				"4. Back to Main Menu",
			},
		}

		selected, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch selected {
		case 0:
			handleCreatePerson(client)
		case 1:
			handleListPersons(client)
		case 2:
			handleWhoAmI(client)
		case 3:
			return
		default:
			return
		}
	}
}

func handleCreateSchool(client *tcp.Client) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the school name:")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Println("School name cannot be empty")
		return
	}

	res, err := client.Send(
		context.Background(),
		tcp.CreateSchool,
		dto.CreateSchoolReq{Name: name},
	)
	if err != nil {
		fmt.Printf("Error creating school: %v\n", err)
		return
	}
	fmt.Printf("School created successfully: %+v\n", res.Data)
}

func handleListSchools(client *tcp.Client) {
	fmt.Println("not yet implemented")
}

func handleCreateClass(client *tcp.Client) {
	fmt.Println("not yet implemented")
}

func handleListClasses(client *tcp.Client) {
	fmt.Println("not yet implemented")
}

func handleAddStudentToClass(client *tcp.Client) {
	fmt.Println("not yet implemented")
}

func handleCreatePerson(client *tcp.Client) {
	fmt.Println("not yet implemented")
}

func handleListPersons(client *tcp.Client) {
	fmt.Println("not yet implemented")
}

func handleWhoAmI(client *tcp.Client) {
	fmt.Println("not yet implemented")
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
