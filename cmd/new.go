package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	projectName string
	moduleName  string
	outputDir   string
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new kratos project",
	Long:  `Generate a new kratos project based on the windo template.`,
	Args:  cobra.ExactArgs(1),
	Run:   runNew,
}

func init() {
	newCmd.Flags().StringVarP(&moduleName, "module", "m", "", "module name (e.g., git.code.tencent.com/windo-/bd/services)")
	newCmd.Flags().StringVarP(&outputDir, "path", "p", "git.code.tencent.com/windo-/bd/services", "output directory for the new project (supports relative and absolute paths)")
}

func runNew(cmd *cobra.Command, args []string) {
	projectName = args[0]

	// 转换 outputDir 为绝对路径，支持相对路径和绝对路径
	var absOutputDir string
	if filepath.IsAbs(outputDir) {
		absOutputDir = outputDir
	} else {
		var err error
		absOutputDir, err = filepath.Abs(outputDir)
		if err != nil {
			fmt.Printf("Error resolving output directory: %v\n", err)
			os.Exit(1)
		}
	}

	basePath := filepath.Join(absOutputDir, projectName)

	if moduleName == "" {
		// 从路径中提取 moduleName，查找 "src/" 后的路径
		if _, after, found := strings.Cut(absOutputDir, "src/"); found {
			moduleName = after
		} else {
			moduleName = absOutputDir
		}
		moduleName = filepath.Join(moduleName, projectName)
	}

	fmt.Printf("Creating new kratos project: %s\n", projectName)
	fmt.Printf("Module name: %s\n", moduleName)
	fmt.Printf("Base path: %s\n", basePath)

	if err := os.MkdirAll(basePath, 0755); err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		os.Exit(1)
	}
	
	files := map[string]string{
		"go.mod":                          goModTemplate,
		"Makefile":                        makefileTemplate,
		"Dockerfile":                      dockerfileTemplate,
		"README.md":                       readmeTemplate,
		".gitignore":                      gitignoreTemplate,
		".gitlab-ci.yml":                  gitlabCiTemplate,
		"cmd/" + projectName + "/main.go": cmdMainTemplate,
		"cmd/" + projectName + "/wire.go": wireTemplate,
		"internal/biz/biz.go":             bizTemplate,
		"internal/biz/greeter.go":         bizGreeterTemplate,
		"internal/data/data.go":           dataTemplate,
		"internal/data/data_repo.go":      dataRepoTemplate,
		"internal/server/server.go":       serverTemplate,
		"internal/server/grpc.go":         grpcTemplate,
		"internal/server/http.go":         httpTemplate,
		"internal/service/service.go":     serviceTemplate,
		"internal/service/greeter.go":     serviceGreeterTemplate,
		"configs/config.yaml":             configTemplate,
		"profile/dev.yaml":                profileDevTemplate,
		"profile/prod.yaml":               profileProdTemplate,
	}

	funcMap := template.FuncMap{
		"toLower":     strings.ToLower,
		"toUpper":     strings.ToUpper,
		"title":       strings.Title,
		"modulePath":  getModulePath,
		"projectName": func() string { return projectName },
	}

	for filePath, content := range files {
		fullPath := filepath.Join(basePath, filePath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}

		tmpl, err := template.New(filePath).Funcs(funcMap).Parse(content)
		if err != nil {
			fmt.Printf("Error parsing template for %s: %v\n", filePath, err)
			os.Exit(1)
		}

		file, err := os.Create(fullPath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", fullPath, err)
			os.Exit(1)
		}
		defer file.Close()

		if err := tmpl.Execute(file, nil); err != nil {
			fmt.Printf("Error executing template for %s: %v\n", filePath, err)
			os.Exit(1)
		}

		fmt.Printf("Created: %s\n", filePath)
	}

	fmt.Printf("\nProject %s created successfully!\n", projectName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  make init")
	fmt.Println("  make api")
	fmt.Println("  make generate")
	fmt.Println("  make build")
}

func getModulePath() string {
	return moduleName
}