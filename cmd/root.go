package cmd

import (
	"bufio"
	"f3-stack/internal/generator"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "f3-stack",
	Short: "F3 App Generator - Flutter, Firebase, Freezed",
	Long: `F3 App Generator is a CLI tool to quickly bootstrap Flutter projects with
Firebase integration and Freezed for code generation, following clean architecture principles.
Think of it as T3 stack but for Flutter.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

const logoart = `
        %%%%                                                                                                
      %%%%%%   %%%%%%%%             %%%%%%%%%  %%%%%%%%%%%%%      %%%%%         %%%%%%%%%    %%%%    %%%%%% 
    %%%%%     %%%%%%%%%%%         %%%%%%%%%%%% %%%%%%%%%%%%%     %%%%%%%      %%%%%%%%%%%%%  %%%%  %%%%%%   
  %%%%%  %%%%%       %%%%         %%%%     %%%%     %%%%        %%%%%%%%     %%%%      %%%%% %%%%%%%%%%     
 %%%%%  %%%%     %%%%%%%           %%%%%%%%         %%%%        %%%  %%%%    %%%             %%%%%%%%       
 %%%%%%%%%%      %%%%%%%%               %%%%%%      %%%%       %%%%   %%%%   %%%%            %%%%%%%%%      
    %%%%%%   %%%%    %%%%         %%%%     %%%%     %%%%      %%%%    %%%%%  %%%%%     %%%%% %%%%  %%%%%    
     %%%%%%%  %%%%%%%%%%%         %%%%%%%%%%%%      %%%%     %%%%      %%%%   %%%%%%%%%%%%%  %%%%   %%%%%%  
        %%%%%  %%%%%%%%              %%%%%%%        %%%%     %%%%       %%%%     %%%%%%%     %%%      %%%%% 
          %%%%%                                                                                                 
`

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new F3 app",
	Run: func(cmd *cobra.Command, args []string) {
		logo := color.New(color.FgHiBlue).Add(color.Bold)
		title := color.New(color.FgHiMagenta)
		display(logoart, logo, title)

		var projectName string

		if len(args) >= 1 {
			projectName = args[0]
		} else {
			green := color.New(color.FgGreen)
			fmt.Print("|→")
			green.Print(" ? ")
			fmt.Print("What will your project be called? ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			projectName = strings.TrimSpace(input)
		}
		if projectName == "" {
			red := color.New(color.FgRed).Add(color.Bold)
			red.Println("Well a modern stack can't be empty init !")
			os.Exit(1)
		}
		main(projectName)
	},
}

func main(projectName string) {
	cyan := color.New(color.FgHiCyan)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue).Add(color.Bold)

	blue.Printf("Using: ")
	cyan.Printf("flutter create %s\n", projectName)

	generator := generator.Project_generator(projectName)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Create_project(); err != nil {
		s.Stop()
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error generating project: %v\n", err)
		os.Exit(1)
	}
	s.Stop()

	fmt.Print("\n")
	green.Print("✓ ")
	blue.Print("f3-stack ")
	cyan.Print(projectName)
	green.Println(" scaffolded successfully!")
	fmt.Println()

	fmt.Print("|→")
	blue.Print(" 〄 Installing dependencies...\n")

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Add_packages(); err != nil {
		s.Stop()
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error adding packages: %v\n", err)
		os.Exit(1)
	}
	s.Stop()
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully installed dependencies!")
	fmt.Println()

	fmt.Print("|→")
	blue.Print(" 〄 Generating folder structure...\n")

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Create_structure(); err != nil {
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error creating folder structure: %v\n", err)
		os.Exit(1)
	}

	s.Stop()
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully created folder structure!")
	fmt.Println()

	fmt.Print("|→")
	blue.Print(" 〄 Adding Assets...\n")

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Add_assets(); err != nil {
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error adding assets: %v\n", err)
		os.Exit(1)
	}

	s.Stop()
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully added assets!")
	fmt.Println()

	fmt.Print("|→")
	blue.Print(" 〄 Updating logs...\n")

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Update_yaml(); err != nil {
		s.Stop()
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error updating pubspec.yaml: %v\n", err)
		os.Exit(1)
	}
	s.Stop()
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully updated logs!")
	fmt.Println()

	fmt.Print("|→")
	blue.Print(" 〄 Adding boilerplate...\n")

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Generate_files(); err != nil {
		s.Stop()
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error generating files: %v\n", err)
		os.Exit(1)
	}

	s.Stop()
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully added boilerplate codes for auth!")
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully added boilerplate codes for layout!")
	fmt.Println()

	fmt.Print("|→")
	blue.Print(" 〄 Updating Info.plist file...\n")

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Update_iOS(); err != nil {
		s.Stop()
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error updating iOS files: %v\n", err)
		os.Exit(1)
	}

	s.Stop()
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully updated Info.plist!")
	fmt.Println()

	fmt.Print("|→")
	blue.Print(" 〄 Running build_runner...\n")

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if err := generator.Runcmd(); err != nil {
		s.Stop()
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("Error running build_runner: %v\n", err)
		os.Exit(1)
	}

	s.Stop()
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully ran build_runner!")
	fmt.Print("      [")
	green.Print("✓ ")
	green.Println("Successfully generated freezed files!")
	fmt.Println()

	fmt.Print("[")
	green.Println("✓ Project scaffolded successfully!")

	fmt.Println()
	blue.Println("Next steps:")
	blue.Println("  cd", projectName)
	blue.Println("  flutter run")
	blue.Println("  Connect to Firebase: https://console.firebase.google.com/")
}

func display(art string, logo, title *color.Color) {
	lines := strings.Split(art, "\n")

	for _, line := range lines {
		if len(line) == 0 {
			fmt.Println()
			continue
		}

		divisionPoint := int(float64(len(line)) * 0.3)

		logo.Print(line[:divisionPoint])

		title.Print(line[divisionPoint:])

		fmt.Println()
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
}
