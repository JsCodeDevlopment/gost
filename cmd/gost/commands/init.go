package commands

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gost"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	projectName string
	template    string
	modules     []string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Gost project",
	Run:   runInit,
}

func init() {
	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	initCmd.Flags().StringVarP(&template, "template", "t", "", "Template type (Full or Basic)")
	initCmd.Flags().StringSliceVarP(&modules, "modules", "m", []string{}, "Modules to include (comma-separated: auth,messaging,i18n)")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Welcome to the Gost CLI Project Initializer!")

	var answers struct {
		ProjectName  string
		TemplateType string
		Modules      []string
	}

	if projectName != "" {
		answers.ProjectName = projectName
		if template == "Full" || template == "foundational" {
			answers.TemplateType = "Full (Includes all modules)"
			answers.Modules = []string{"Authentication (JWT/Redis)", "Messaging (RabbitMQ)", "Internationalization (i18n)"}
		} else {
			answers.TemplateType = "Basic (Pick modules)"
			for _, m := range modules {
				switch strings.ToLower(m) {
				case "auth", "authentication":
					answers.Modules = append(answers.Modules, "Authentication (JWT/Redis)")
				case "messaging", "mq", "rabbitmq":
					answers.Modules = append(answers.Modules, "Messaging (RabbitMQ)")
				case "i18n", "translation":
					answers.Modules = append(answers.Modules, "Internationalization (i18n)")
				}
			}
		}
	} else {
		questions := []*survey.Question{
			{
				Name: "ProjectName",
				Prompt: &survey.Input{
					Message: "What is your project's name?",
					Default: "gost-app",
				},
				Validate: survey.Required,
			},
			{
				Name: "TemplateType",
				Prompt: &survey.Select{
					Message: "Choose the template type:",
					Options: []string{"Full (Includes all modules)", "Basic (Pick modules)"},
					Default: "Full (Includes all modules)",
				},
			},
		}

		err := survey.Ask(questions, &answers)
		if err != nil {
			fmt.Printf("Error during prompts: %v\n", err)
			return
		}

		if answers.TemplateType == "Basic (Pick modules)" {
			modulePrompt := &survey.MultiSelect{
				Message: "Select optional modules to include:",
				Options: []string{"Authentication (JWT/Redis)", "Messaging (RabbitMQ)", "Internationalization (i18n)"},
			}
			err = survey.AskOne(modulePrompt, &answers.Modules)
			if err != nil {
				fmt.Printf("Error selecting modules: %v\n", err)
				return
			}
		} else {
			answers.Modules = []string{"Authentication (JWT/Redis)", "Messaging (RabbitMQ)", "Internationalization (i18n)"}
		}
	}

	fmt.Printf("\n🏗️ Scaffolding project %s...\n", answers.ProjectName)

	destDir := answers.ProjectName

	if err := os.MkdirAll(destDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create directory: %v\n", err)
		return
	}
	err := extractEmbedFS(gost.TemplateFS, ".", destDir, []string{".git", ".gemini", "cmd", "test", answers.ProjectName})
	if err != nil {
		fmt.Printf("❌ Failed to scaffold project: %v\n", err)
		return
	}

	pruneModules(destDir, answers.Modules)

	patchNewProject(destDir, answers.ProjectName, answers.Modules)

	fmt.Printf("\n✅ Project %s initialized successfully!\n", answers.ProjectName)
	fmt.Printf("👉 cd %s\n", answers.ProjectName)
	fmt.Println("👉 go mod tidy")
	fmt.Println("👉 go run main.go")
}

func extractEmbedFS(efs embed.FS, srcDir string, dstDir string, ignore []string) error {
	return fs.WalkDir(efs, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(srcDir, path)
		if relPath == "." {
			return nil
		}
		for _, skip := range ignore {
			if strings.HasPrefix(relPath, skip) || relPath == skip {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		destPath := filepath.Join(dstDir, relPath)
		if relPath == "main.go.tpl" {
			destPath = filepath.Join(dstDir, "main.go")
		}

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0644)
	})
}

func hasModule(selected []string, module string) bool {
	for _, s := range selected {
		if strings.Contains(s, module) {
			return true
		}
	}
	return false
}

func pruneModules(destDir string, selectedModules []string) {
	if !hasModule(selectedModules, "Authentication") {
		os.RemoveAll(filepath.Join(destDir, "src", "modules", "auth"))
		os.RemoveAll(filepath.Join(destDir, "src", "common", "guards"))
		os.RemoveAll(filepath.Join(destDir, "src", "common", "security"))
	}
	if !hasModule(selectedModules, "Messaging") {
		os.RemoveAll(filepath.Join(destDir, "src", "common", "messaging"))
		os.Remove(filepath.Join(destDir, "src", "config", "rabbitmq.go"))
	}
	if !hasModule(selectedModules, "i18n") {
		os.RemoveAll(filepath.Join(destDir, "src", "common", "i18n"))
		os.RemoveAll(filepath.Join(destDir, "locales"))
		os.RemoveAll(filepath.Join(destDir, "docs", "11-internationalization-i18n.md"))
	}
}

func patchNewProject(destDir, projectName string, selectedModules []string) {
	goModPath := filepath.Join(destDir, "go.mod")
	replaceInFile(goModPath, "module gost", "module "+projectName)

	appModulePath := filepath.Join(destDir, "src", "app", "app.module.go")
	content, _ := os.ReadFile(appModulePath)
	lines := strings.Split(string(content), "\n")
	var newLines []string

	isAuth := hasModule(selectedModules, "Authentication")
	isMQ := hasModule(selectedModules, "Messaging")
	isI18n := hasModule(selectedModules, "i18n")

	for _, line := range lines {
		if !isMQ && (strings.Contains(line, "config.ConnectRabbitMQ") || strings.Contains(line, "src/config") && !isAuth) {
			if strings.Contains(line, "config.ConnectRabbitMQ") {
				continue
			}
		}

		if !isMQ && strings.Contains(line, "config.ConnectRabbitMQ") {
			continue
		}
		if !isMQ && strings.Contains(line, "\"gost/src/common/messaging\"") {
			continue
		}

		if !isI18n && (strings.Contains(line, "i18n.Initialize") ||
			strings.Contains(line, "i18n.InitValidator") ||
			strings.Contains(line, "i18n.Middleware") ||
			strings.Contains(line, "\"gost/src/common/i18n\"") ||
			strings.Contains(line, "\"golang.org/x/text/language\"")) {
			continue
		}

		if !isAuth && (strings.Contains(line, "auth.InitModule(api)") ||
			strings.Contains(line, "\"gost/src/modules/auth\"")) {
			continue
		}

		line = strings.ReplaceAll(line, "\"gost/", "\""+projectName+"/")

		newLines = append(newLines, line)
	}

	os.WriteFile(appModulePath, []byte(strings.Join(newLines, "\n")), os.ModePerm)

	filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			replaceInFile(path, "\"gost/", "\""+projectName+"/")
		}
		return nil
	})
}

func replaceInFile(path, oldStr, newStr string) {
	content, _ := os.ReadFile(path)
	newContent := strings.ReplaceAll(string(content), oldStr, newStr)
	os.WriteFile(path, []byte(newContent), os.ModePerm)
}
