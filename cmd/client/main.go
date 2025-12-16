package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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
	res, err := client.Send(
		context.Background(),
		tcp.ListSchools, "")
	if err != nil {
		fmt.Printf("Error listing schools: %v\n", err)
		return
	}

	dataBytes, err := json.Marshal(res.Data)
	if err != nil {
		fmt.Printf("Error parsing schools data: %v\n", err)
		return
	}

	var schools []struct {
		Id   uint   `json:"Id"`
		Name string `json:"Name"`
	}

	if err := json.Unmarshal(dataBytes, &schools); err != nil {
		fmt.Printf("Error unmarshaling schools: %v\n", err)
		return
	}

	if len(schools) == 0 {
		fmt.Println("No schools found.")
		return
	}

	printSchoolsTable(schools)
}

func handleCreateClass(client *tcp.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the class name:")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Println("Class name cannot be empty")
		return
	}

	fmt.Println("Enter the school ID:")
	scanner.Scan()
	schoolIdStr := strings.TrimSpace(scanner.Text())
	schoolId, err := strconv.ParseUint(schoolIdStr, 10, 32)
	if err != nil {
		fmt.Printf("Invalid school ID: %v\n", err)
		return
	}

	fmt.Println("Enter the teacher ID:")
	scanner.Scan()
	teacherIdStr := strings.TrimSpace(scanner.Text())
	teacherId, err := strconv.ParseUint(teacherIdStr, 10, 32)
	if err != nil {
		fmt.Printf("Invalid teacher ID: %v\n", err)
		return
	}

	res, err := client.Send(
		context.Background(),
		tcp.CreateClass,
		dto.CreateClassReq{
			Name:      name,
			SchoolId:  uint(schoolId),
			TeacherId: uint(teacherId),
		},
	)
	if err != nil {
		fmt.Printf("Error creating class: %v\n", err)
		return
	}
	fmt.Printf("Class created successfully: %+v\n", res.Data)
}

func handleListClasses(client *tcp.Client) {
	res, err := client.Send(
		context.Background(),
		tcp.ListClasses, "")
	if err != nil {
		fmt.Printf("Error listing classes: %v\n", err)
		return
	}

	dataBytes, err := json.Marshal(res.Data)
	if err != nil {
		fmt.Printf("Error parsing classes data: %v\n", err)
		return
	}

	var classes []struct {
		Id       uint   `json:"Id"`
		Name     string `json:"Name"`
		SchoolId uint   `json:"SchoolId"`
	}

	if err := json.Unmarshal(dataBytes, &classes); err != nil {
		fmt.Printf("Error unmarshaling classes: %v\n", err)
		return
	}

	if len(classes) == 0 {
		fmt.Println("No classes found.")
		return
	}

	printClassesTable(classes)
}

func handleAddStudentToClass(client *tcp.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the class ID:")
	scanner.Scan()
	classIdStr := strings.TrimSpace(scanner.Text())
	classId, err := strconv.ParseUint(classIdStr, 10, 32)
	if err != nil {
		fmt.Printf("Invalid class ID: %v\n", err)
		return
	}

	fmt.Println("Enter the student ID:")
	scanner.Scan()
	studentIdStr := strings.TrimSpace(scanner.Text())
	studentId, err := strconv.ParseUint(studentIdStr, 10, 32)
	if err != nil {
		fmt.Printf("Invalid student ID: %v\n", err)
		return
	}

	res, err := client.Send(
		context.Background(),
		tcp.AddStudentToClass,
		dto.AddStudentToClassReq{
			ClassId:   uint(classId),
			StudentId: uint(studentId),
		},
	)
	if err != nil {
		fmt.Printf("Error adding student to class: %v\n", err)
		return
	}
	fmt.Printf("%v\n", res.Data)
}

func handleCreatePerson(client *tcp.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the person name:")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Println("Person name cannot be empty")
		return
	}

	fmt.Println("Enter the role (student/teacher):")
	scanner.Scan()
	role := strings.TrimSpace(scanner.Text())
	if role != "student" && role != "teacher" {
		fmt.Println("Role must be 'student' or 'teacher'")
		return
	}

	fmt.Println("Enter the school ID:")
	scanner.Scan()
	schoolIdStr := strings.TrimSpace(scanner.Text())
	schoolId, err := strconv.ParseUint(schoolIdStr, 10, 32)
	if err != nil {
		fmt.Printf("Invalid school ID: %v\n", err)
		return
	}

	res, err := client.Send(
		context.Background(),
		tcp.CreatePerson,
		dto.CreatePersonReq{
			Name:     name,
			Role:     role,
			SchoolId: uint(schoolId),
		},
	)
	if err != nil {
		fmt.Printf("Error creating person: %v\n", err)
		return
	}
	fmt.Printf("Person created successfully: %+v\n", res.Data)
}

func handleListPersons(client *tcp.Client) {
	res, err := client.Send(
		context.Background(),
		tcp.ListPersons, "")
	if err != nil {
		fmt.Printf("Error listing persons: %v\n", err)
		return
	}

	dataBytes, err := json.Marshal(res.Data)
	if err != nil {
		fmt.Printf("Error parsing persons data: %v\n", err)
		return
	}

	var persons []struct {
		Id     uint   `json:"Id"`
		Name   string `json:"Name"`
		Role   string `json:"Role"`
		School struct {
			Id   uint   `json:"Id"`
			Name string `json:"Name"`
		} `json:"School"`
	}

	if err := json.Unmarshal(dataBytes, &persons); err != nil {
		fmt.Printf("Error unmarshaling persons: %v\n", err)
		return
	}

	if len(persons) == 0 {
		fmt.Println("No persons found.")
		return
	}

	printPersonsTable(persons)
}

func handleWhoAmI(client *tcp.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the person ID:")
	scanner.Scan()
	personIdStr := strings.TrimSpace(scanner.Text())
	personId, err := strconv.ParseUint(personIdStr, 10, 32)
	if err != nil {
		fmt.Printf("Invalid person ID: %v\n", err)
		return
	}

	res, err := client.Send(
		context.Background(),
		tcp.WhoAmI,
		dto.WhoAmIReq{
			PersonId: uint(personId),
		},
	)
	if err != nil {
		fmt.Printf("Error getting person info: %v\n", err)
		return
	}

	dataBytes, err := json.Marshal(res.Data)
	if err != nil {
		fmt.Printf("Error parsing person data: %v\n", err)
		return
	}

	var person struct {
		Id     uint   `json:"Id"`
		Name   string `json:"Name"`
		Role   string `json:"Role"`
		School struct {
			Id   uint   `json:"Id"`
			Name string `json:"Name"`
		} `json:"School"`
	}

	if err := json.Unmarshal(dataBytes, &person); err != nil {
		fmt.Printf("Error unmarshaling person: %v\n", err)
		return
	}

	printPersonDetails(person)
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

func printSchoolsTable(schools []struct {
	Id   uint   `json:"Id"`
	Name string `json:"Name"`
}) {
	fmt.Println("\n┌─────┬────────────────────────────────────────┐")
	fmt.Printf("│ %-3s │ %-38s │\n", "ID", "Name")
	fmt.Println("├─────┼────────────────────────────────────────┤")

	for _, school := range schools {
		name := school.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		fmt.Printf("│ %-3d │ %-38s │\n", school.Id, name)
	}

	fmt.Println("└─────┴────────────────────────────────────────┘")
	fmt.Printf("\nTotal: %d school(s)\n\n", len(schools))
}

func printClassesTable(classes []struct {
	Id       uint   `json:"Id"`
	Name     string `json:"Name"`
	SchoolId uint   `json:"SchoolId"`
}) {
	fmt.Println("\n┌─────┬────────────────────────────────────────┬──────────┐")
	fmt.Printf("│ %-3s │ %-38s │ %-8s │\n", "ID", "Name", "SchoolID")
	fmt.Println("├─────┼────────────────────────────────────────┼──────────┤")

	for _, class := range classes {
		name := class.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		fmt.Printf("│ %-3d │ %-38s │ %-8d │\n", class.Id, name, class.SchoolId)
	}

	fmt.Println("└─────┴────────────────────────────────────────┴──────────┘")
	fmt.Printf("\nTotal: %d class(es)\n\n", len(classes))
}

func printPersonsTable(persons []struct {
	Id     uint   `json:"Id"`
	Name   string `json:"Name"`
	Role   string `json:"Role"`
	School struct {
		Id   uint   `json:"Id"`
		Name string `json:"Name"`
	} `json:"School"`
}) {
	fmt.Println("\n┌─────┬────────────────────────────────────────┬──────────┬────────────────────────────────────────┐")
	fmt.Printf("│ %-3s │ %-38s │ %-8s │ %-38s │\n", "ID", "Name", "Role", "School")
	fmt.Println("├─────┼────────────────────────────────────────┼──────────┼────────────────────────────────────────┤")

	for _, person := range persons {
		name := person.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		schoolName := person.School.Name
		if len(schoolName) > 38 {
			schoolName = schoolName[:35] + "..."
		}
		fmt.Printf("│ %-3d │ %-38s │ %-8s │ %-38s │\n", person.Id, name, person.Role, schoolName)
	}

	fmt.Println("└─────┴────────────────────────────────────────┴──────────┴────────────────────────────────────────┘")
	fmt.Printf("\nTotal: %d person(s)\n\n", len(persons))
}

func printPersonDetails(person struct {
	Id     uint   `json:"Id"`
	Name   string `json:"Name"`
	Role   string `json:"Role"`
	School struct {
		Id   uint   `json:"Id"`
		Name string `json:"Name"`
	} `json:"School"`
}) {
	fmt.Println("\n┌──────────────────────────────────────────────────────────┐")
	fmt.Printf("│ ID:     %-50d │\n", person.Id)
	fmt.Printf("│ Name:   %-50s │\n", person.Name)
	fmt.Printf("│ Role:   %-50s │\n", person.Role)
	fmt.Printf("│ School: %-50s │\n", person.School.Name)
	fmt.Println("└──────────────────────────────────────────────────────────┘")
	fmt.Println()
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
