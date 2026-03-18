package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var makeModuleCmd = &cobra.Command{
	Use:   "make:module [name]",
	Short: "Create a new module structure",
	Args:  cobra.ExactArgs(1),
	Run:   runMakeModule,
}

func init() {
	rootCmd.AddCommand(makeModuleCmd)
}

func runMakeModule(cmd *cobra.Command, args []string) {
	moduleName := strings.ToLower(args[0])
	fmt.Printf("📦 Creating module %s...\n", moduleName)

	basePath := filepath.Join("src", "modules", moduleName)

	subDirs := []string{
		"dto",
		"entities",
		"repositories",
		"services",
	}

	for _, dir := range subDirs {
		path := filepath.Join(basePath, dir)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", path, err)
			return
		}
	}

	moduleFilePath := filepath.Join(basePath, moduleName+".module.go")
	moduleContent := fmt.Sprintf(`package %s

import (
	"github.com/gin-gonic/gin"
)

func InitModule(router *gin.RouterGroup) {
	// Initialize Repository, Service, and Controller here
	// Example:
	// repo := NewRepository()
	// service := NewService(repo)
	// ctrl := NewController(service)
	
	// router.GET("/%s", ctrl.FindAll)
}
`, moduleName, moduleName)

	if err := os.WriteFile(moduleFilePath, []byte(moduleContent), os.ModePerm); err != nil {
		fmt.Printf("Error creating module file: %v\n", err)
		return
	}

	fmt.Printf("✅ Module %s created at %s\n", moduleName, basePath)
	fmt.Println("⚠️  Don't forget to register it in src/app/app.module.go")
}
