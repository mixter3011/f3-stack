package generator

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Generator struct {
	ProjectName string
	TemplateDir string
}

func Project_generator(projectName string) *Generator {
	return &Generator{
		ProjectName: projectName,
	}
}

func (pg *Generator) Generate() error {
	if err := pg.Create_project(); err != nil {
		return fmt.Errorf("failed to create Flutter project: %w", err)
	}

	if err := pg.Add_packages(); err != nil {
		return fmt.Errorf("failed to add packages: %w", err)
	}

	if err := pg.Create_structure(); err != nil {
		return fmt.Errorf("failed to create folder structure: %w", err)
	}

	if err := pg.Add_assets(); err != nil {
		return fmt.Errorf("failed to copy assets: %w", err)
	}

	if err := pg.Update_yaml(); err != nil {
		return fmt.Errorf("failed to update pubspec.yaml: %w", err)
	}

	if err := pg.Generate_files(); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	if err := pg.Update_iOS(); err != nil {
		return fmt.Errorf("failed to update iOS files: %w", err)
	}

	if err := pg.Runcmd(); err != nil {
		return fmt.Errorf("failed to run build_runner: %w", err)
	}

	return nil
}

func (pg *Generator) Create_project() error {
	cmd := exec.Command("flutter", "create", pg.ProjectName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (pg *Generator) Add_packages() error {
	packages := []string{
		"flutter_bloc:^9.1.0",
		"freezed_annotation:^3.0.0",
		"firebase_core:^3.13.0",
		"firebase_auth:^5.5.2",
		"json_annotation:^4.9.0",
		"equatable:^2.0.7",
		"shadcn_ui:^0.24.0",
		"google_fonts:^6.2.1",
		"url_launcher:^6.3.1",
		"cached_network_image:^3.4.1",
		"url_launcher_ios:^6.3.3",
		"google_sign_in:^6.3.0",
		"build_runner:^2.4.8",
		"freezed:^3.0.1",
		"json_serializable:^6.7.1",
	}

	projectDir := filepath.Join(".", pg.ProjectName)

	for _, pkg := range packages {
		cmd := exec.Command("flutter", "pub", "add", pkg)
		cmd.Dir = projectDir
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Println(string(output))
			return err
		}
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (pg *Generator) Create_structure() error {
	libDir := filepath.Join(pg.ProjectName, "lib")

	folders := []string{
		"core/constants",
		"core/error",
		"core/services",
		"core/utils",
		"core/widgets",

		"features/auth/data/datasources",
		"features/auth/data/models",
		"features/auth/data/repositories",

		"features/auth/domain/entities",
		"features/auth/domain/repositories",
		"features/auth/domain/usecases",

		"features/auth/presentation/bloc",
		"features/auth/presentation/pages",
		"features/auth/presentation/widgets",

		"features/home/data/datasources",
		"features/home/data/models",
		"features/home/data/repositories",

		"features/home/domain/entities",
		"features/home/domain/repositories",
		"features/home/domain/usecases",

		"features/home/presentation/bloc",
		"features/home/presentation/pages",
		"features/home/presentation/widgets",
	}

	for _, folder := range folders {
		path := filepath.Join(libDir, folder)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	assetsDir := filepath.Join(pg.ProjectName, "assets", "images")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return err
	}

	logoPath := filepath.Join(assetsDir, "logo.png")
	if err := os.WriteFile(logoPath, []byte("placeholder image content"), 0644); err != nil {
		return err
	}
	time.Sleep(300 * time.Millisecond)
	return nil
}

func (pg *Generator) Add_assets() error {
	destImagesDir := filepath.Join(pg.ProjectName, "assets", "images")
	if err := os.MkdirAll(destImagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets/images directory: %w", err)
	}

	rootPath := "assets"

	return fs.WalkDir(embeddedAssets, rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == rootPath {
			return nil
		}

		_, filename := filepath.Split(path)
		destPath := filepath.Join(destImagesDir, filename)

		if !d.IsDir() {
			content, err := embeddedAssets.ReadFile(path)
			if err != nil {
				return err
			}

			if err := os.WriteFile(destPath, content, 0644); err != nil {
				return err
			}
		}

		return nil
	})
}

func (pg *Generator) Update_yaml() error {
	content, exists := templateData["pubspec.yaml"]
	if !exists {
		return fmt.Errorf("pubspec.yaml template not found in embedded data")
	}

	updatedContent := strings.ReplaceAll(content, "name: f3stack", fmt.Sprintf("name: %s", pg.ProjectName))

	pubspecPath := filepath.Join(pg.ProjectName, "pubspec.yaml")
	time.Sleep(300 * time.Millisecond)
	return os.WriteFile(pubspecPath, []byte(updatedContent), 0644)
}

func (pg *Generator) Generate_files() error {
	files := map[string]string{
		"lib/core/constants/routes.dart": "routes.dart",

		"lib/features/auth/data/models/user_model.dart":                 "user_model.dart",
		"lib/features/auth/data/repositories/auth_repository_impl.dart": "auth_repository_impl.dart",

		"lib/features/auth/domain/entities/user_entity.dart":           "user_entity.dart",
		"lib/features/auth/domain/repositories/auth_repository.dart":   "auth_repository.dart",
		"lib/features/auth/domain/usecases/google_signin_usecase.dart": "google_signin_usecase.dart",
		"lib/features/auth/domain/usecases/signin_usecase.dart":        "signin_usecase.dart",
		"lib/features/auth/domain/usecases/signup_usecase.dart":        "signup_usecase.dart",
		"lib/features/auth/domain/usecases/siginout_usecase.dart":      "signout_usecase.dart",

		"lib/features/auth/presentation/bloc/auth_bloc.dart":     "auth_bloc.dart",
		"lib/features/auth/presentation/bloc/auth_event.dart":    "auth_event.dart",
		"lib/features/auth/presentation/bloc/auth_state.dart":    "auth_state.dart",
		"lib/features/auth/presentation/pages/signin_page.dart":  "signin_page.dart",
		"lib/features/auth/presentation/pages/signup_page.dart":  "signup_page.dart",
		"lib/features/auth/presentation/pages/auth_wrapper.dart": "auth_wrapper.dart",

		"lib/features/home/presentation/pages/home_page.dart":       "home_page.dart",
		"lib/features/home/presentation/widgets/action_button.dart": "action_button.dart",
		"lib/features/home/presentation/widgets/bottom_bar.dart":    "bottom_bar.dart",
		"lib/features/home/presentation/widgets/content.dart":       "content.dart",
		"lib/features/home/presentation/widgets/feature_grid.dart":  "feature_grid.dart",
		"lib/features/home/presentation/widgets/feature_card.dart":  "feature_card.dart",
		"lib/features/home/presentation/widgets/features.dart":      "features.dart",
		"lib/features/home/presentation/widgets/hero.dart":          "hero.dart",
		"lib/features/home/presentation/widgets/started.dart":       "started.dart",

		"lib/main.dart":         "main.dart",
		"test/widget_test.dart": "widget_test.dart",
		"ios/Runner/Info.plist": "Info.plist",
	}

	for filePath, templateKey := range files {
		content, exists := templateData[templateKey]
		if !exists {
			content = "// Implement " + templateKey + "\n"
		}

		content = strings.ReplaceAll(content, "f3stack", pg.ProjectName)

		fullPath := filepath.Join(pg.ProjectName, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (pg *Generator) Update_iOS() error {
	infoPlistPath := filepath.Join(pg.ProjectName, "ios", "Runner", "Info.plist")
	plistContent, exists := templateData["Info.plist"]
	if !exists {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(infoPlistPath), 0755); err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	return os.WriteFile(infoPlistPath, []byte(plistContent), 0644)
}

func (pg *Generator) Runcmd() error {
	cmd := exec.Command("dart", "run", "build_runner", "build", "--delete-conflicting-outputs")
	cmd.Dir = pg.ProjectName
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}
