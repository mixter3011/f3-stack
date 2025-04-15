package cmd

import (
	"f3-stack/internal/generator"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "f3-stack",
	Short: "F3 App Generator - Flutter, Firebase, Freezed",
	Long: `F3 App Generator is a CLI tool to quickly bootstrap Flutter projects with
Firebase integration and Freezed for code generation, following clean architecture principles.
Think of it as T3 stack but for Flutter.`,
}

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new F3 app",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		generator := generator.Project_generator(projectName)

		fmt.Printf("Creating F3 app: %s\n", projectName)
		if err := generator.Generate(); err != nil {
			fmt.Printf("Error generating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n🎉 Successfully created project: %s\n", projectName)
		fmt.Println("\n📝 Next steps:")
		fmt.Println("1. Connect your project to Firebase:")
		fmt.Println("   - Go to Firebase Console (https://console.firebase.google.com/)")
		fmt.Println("   - Create a new project with the same name")
		fmt.Println("   - Follow the instructions to add your app to Firebase")
		fmt.Println("2. Run your app:")
		fmt.Printf("   - cd %s\n", projectName)
		fmt.Println("   - flutter run")
		fmt.Println("\n🚀 Happy coding!")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
}
