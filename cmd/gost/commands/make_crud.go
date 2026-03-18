package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var makeCrudCmd = &cobra.Command{
	Use:   "make:crud [name]",
	Short: "Generate a complete CRUD for a new module",
	Args:  cobra.ExactArgs(1),
	Run:   runMakeCrud,
}

func init() {
	rootCmd.AddCommand(makeCrudCmd)
}

func runMakeCrud(cmd *cobra.Command, args []string) {
	name := strings.ToLower(args[0])
	capitalName := strings.ToUpper(name[0:1]) + name[1:]
	pluralName := name + "s"
	modulePath := filepath.Join("src", "modules", pluralName)

	fmt.Printf("🛠️ Generating CRUD for module %s...\n", pluralName)

	if err := os.MkdirAll(filepath.Join(modulePath, "dto"), os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}
	if err := os.MkdirAll(filepath.Join(modulePath, "entities"), os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	projectName := getProjectName()

	generateEntity(modulePath, name, capitalName)
	generateDto(modulePath, name, capitalName)
	generateRepository(modulePath, name, capitalName, pluralName, projectName)
	generateService(modulePath, name, capitalName, pluralName, projectName)
	generateController(modulePath, name, capitalName, pluralName, projectName)
	generateModule(modulePath, name, capitalName, pluralName, projectName)

	registerModule(pluralName, projectName)

	fmt.Printf("✅ CRUD for %s successfully generated at %s\n", capitalName, modulePath)
	fmt.Printf("👉 Check http://localhost:3000/api/v1/%s after restarting the server.\n", pluralName)
}

func generateEntity(path, name, capName string) {
	content := fmt.Sprintf(`package entities

import "gorm.io/gorm"

type %s struct {
	gorm.Model
	Name  string `+"`json:\"name\"`"+`
}
`, capName)
	os.WriteFile(filepath.Join(path, "entities", name+".entity.go"), []byte(content), os.ModePerm)
}

func generateDto(path, name, capName string) {
	content := fmt.Sprintf(`package dto

type Create%sDto struct {
	Name string `+"`json:\"name\" binding:\"required,min=3\"`"+`
}

type Update%sDto struct {
	Name string `+"`json:\"name\"`"+`
}
`, capName, capName)
	os.WriteFile(filepath.Join(path, "dto", name+".dto.go"), []byte(content), os.ModePerm)
}

func generateRepository(path, name, capName, plural, proj string) {
	content := fmt.Sprintf(`package %s

import (
	"%s/src/modules/%s/entities"
	"gorm.io/gorm"
)

type %sRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *%sRepository {
	return &%sRepository{db: db}
}

func (r *%sRepository) Create(entity *entities.%s) error {
	return r.db.Create(entity).Error
}

func (r *%sRepository) FindAll() ([]entities.%s, error) {
	var results []entities.%s
	err := r.db.Find(&results).Error
	return results, err
}

func (r *%sRepository) FindOne(id uint) (*entities.%s, error) {
	var result entities.%s
	err := r.db.First(&result, id).Error
	return &result, err
}

func (r *%sRepository) Update(entity *entities.%s) error {
	return r.db.Save(entity).Error
}

func (r *%sRepository) Delete(id uint) error {
	return r.db.Delete(&entities.%s{}, id).Error
}
`, plural, proj, plural, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName)
	os.WriteFile(filepath.Join(path, plural+".repository.go"), []byte(content), os.ModePerm)
}

func generateService(path, name, capName, plural, proj string) {
	content := fmt.Sprintf(`package %s

import (
	"%s/src/modules/%s/dto"
	"%s/src/modules/%s/entities"
)

type %sService struct {
	repo *%sRepository
}

func NewService(repo *%sRepository) *%sService {
	return &%sService{repo: repo}
}

func (s *%sService) Create(data dto.Create%sDto) (*entities.%s, error) {
	entity := &entities.%s{
		Name: data.Name,
	}
	err := s.repo.Create(entity)
	return entity, err
}

func (s *%sService) FindAll() ([]entities.%s, error) {
	return s.repo.FindAll()
}

func (s *%sService) FindOne(id uint) (*entities.%s, error) {
	return s.repo.FindOne(id)
}

func (s *%sService) Update(id uint, data dto.Update%sDto) (*entities.%s, error) {
	entity, err := s.repo.FindOne(id)
	if err != nil {
		return nil, err
	}
	if data.Name != "" {
		entity.Name = data.Name
	}
	err = s.repo.Update(entity)
	return entity, err
}

func (s *%sService) Delete(id uint) error {
	return s.repo.Delete(id)
}
`, plural, proj, plural, proj, plural, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName)
	os.WriteFile(filepath.Join(path, plural+".service.go"), []byte(content), os.ModePerm)
}

func generateController(path, name, capName, plural, proj string) {
	content := fmt.Sprintf(`package %s

import (
	"%s/src/common/pipes"
	"%s/src/modules/%s/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type %sController struct {
	service *%sService
}

func NewController(service *%sService) *%sController {
	return &%sController{service: service}
}

func (ctrl *%sController) Create(c *gin.Context) {
	body, err := pipes.ValidateBody[dto.Create%sDto](c)
	if err != nil {
		return
	}

	result, err := ctrl.service.Create(*body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (ctrl *%sController) FindAll(c *gin.Context) {
	results, err := ctrl.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (ctrl *%sController) FindOne(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result, err := ctrl.service.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "%%s not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctrl *%sController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	body, err := pipes.ValidateBody[dto.Update%sDto](c)
	if err != nil {
		return
	}

	result, err := ctrl.service.Update(uint(id), *body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctrl *%sController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
`, plural, proj, proj, plural, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName, capName)

	content = strings.Replace(content, "%%s", capName, 1)

	os.WriteFile(filepath.Join(path, plural+".controller.go"), []byte(content), os.ModePerm)
}

func generateModule(path, name, capName, plural, proj string) {
	content := fmt.Sprintf(`package %s

import (
	"%s/src/config"
	"github.com/gin-gonic/gin"
)

func InitModule(router *gin.RouterGroup) {
	db := config.GetDB()
	repo := NewRepository(db)
	service := NewService(repo)
	ctrl := NewController(service)

	group := router.Group("/%s")
	{
		group.POST("/", ctrl.Create)
		group.GET("/", ctrl.FindAll)
		group.GET("/:id", ctrl.FindOne)
		group.PATCH("/:id", ctrl.Update)
		group.DELETE("/:id", ctrl.Delete)
	}
}
`, plural, proj, plural)
	os.WriteFile(filepath.Join(path, plural+".module.go"), []byte(content), os.ModePerm)
}

func registerModule(plural, proj string) {
	path := filepath.Join("src", "app", "app.module.go")
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	importStr := fmt.Sprintf("\t\"%s/src/modules/%s\"", proj, plural)
	newContent := strings.Replace(string(content), "import (", "import (\n"+importStr, 1)
	initStr := fmt.Sprintf("\t%s.InitModule(api)", plural)
	newContent = strings.Replace(newContent, "ws.InitModule(api)", "ws.InitModule(api)\n"+initStr, 1)

	os.WriteFile(path, []byte(newContent), os.ModePerm)
}

func getProjectName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "gost"
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 && parts[0] == "module" {
			return parts[1]
		}
	}
	return "gost"
}
