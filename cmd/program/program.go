// Package program provides the
// main functionality of Blueprint
package program

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/noxymon-mekari/mkrgen/cmd/flags"
	tpl "github.com/noxymon-mekari/mkrgen/cmd/template"
	"github.com/noxymon-mekari/mkrgen/cmd/template/advanced"
	"github.com/noxymon-mekari/mkrgen/cmd/template/dbdriver"
	"github.com/noxymon-mekari/mkrgen/cmd/template/docker"
	"github.com/noxymon-mekari/mkrgen/cmd/template/framework"
	"github.com/noxymon-mekari/mkrgen/cmd/utils"
)

// A Project contains the data for the project folder
// being created, and methods that help with that process
type Project struct {
	ProjectName       string
	Exit              bool
	AbsolutePath      string
	ProjectType       flags.Framework
	DBDriver          flags.Database
	Docker            flags.Database
	FrameworkMap      map[flags.Framework]Framework
	DBDriverMap       map[flags.Database]Driver
	DockerMap         map[flags.Database]Docker
	AdvancedOptions   map[string]bool
	AdvancedTemplates AdvancedTemplates
	GitOptions        flags.Git
	OSCheck           map[string]bool
}

type AdvancedTemplates struct {
	TemplateRoutes  string
	TemplateImports string
}

// A Framework contains the name and templater for a
// given Framework
type Framework struct {
	packageName []string
	templater   Templater
}

type Driver struct {
	packageName []string
	templater   DBDriverTemplater
}

type Docker struct {
	packageName []string
	templater   DockerTemplater
}

// A Templater has the methods that help build the files
// in the Project folder, and is specific to a Framework
type Templater interface {
	Main() []byte
	Server() []byte
	Routes() []byte
	TestHandler() []byte
	HtmxTemplRoutes() []byte
	HtmxTemplImports() []byte
	WebsocketImports() []byte
	SwaggerImports() []byte
	SwaggerRoutes() []byte
}

type DBDriverTemplater interface {
	Service() []byte
	Env() []byte
	Tests() []byte
	SqlcConfig() []byte
	SchemaExample() []byte
	QueryExample() []byte
	UsersSchema() []byte
	PostsSchema() []byte
	UsersQuery() []byte
	PostsQuery() []byte
}

type DockerTemplater interface {
	Docker() []byte
}

type WorkflowTemplater interface {
	Releaser() []byte
	Test() []byte
	ReleaserConfig() []byte
}

var (
	chiPackage     = []string{"github.com/go-chi/chi/v5"}
	gorillaPackage = []string{"github.com/gorilla/mux"}
	routerPackage  = []string{"github.com/julienschmidt/httprouter"}
	ginPackage     = []string{"github.com/gin-gonic/gin"}
	fiberPackage   = []string{"github.com/gofiber/fiber/v2"}
	echoPackage    = []string{"github.com/labstack/echo/v4", "github.com/labstack/echo/v4/middleware"}

	mysqlDriver    = []string{"github.com/go-sql-driver/mysql", "github.com/sqlc-dev/sqlc/cmd/sqlc"}
	postgresDriver = []string{"github.com/jackc/pgx/v5/stdlib", "github.com/sqlc-dev/sqlc/cmd/sqlc"}
	sqliteDriver   = []string{"github.com/mattn/go-sqlite3", "github.com/sqlc-dev/sqlc/cmd/sqlc"}
	mongoDriver    = []string{"go.mongodb.org/mongo-driver"}
	gocqlDriver    = []string{"github.com/gocql/gocql"}
	scyllaDriver   = "github.com/scylladb/gocql@v1.14.4" // Replacement for GoCQL

	godotenvPackage = []string{"github.com/joho/godotenv"}
	templPackage    = []string{"github.com/a-h/templ"}
	swaggerPackage  = []string{"github.com/swaggo/swag/cmd/swag", "github.com/swaggo/files", "github.com/swaggo/gin-swagger", "github.com/swaggo/echo-swagger", "github.com/swaggo/fiber-swagger", "github.com/swaggo/http-swagger"}
)

const (
	root                   = "/"
	cmdApiPath             = "cmd/api"
	cmdWebPath             = "cmd/web"
	internalServerPath     = "cmd/api/server"
	internalDatabasePath   = "pkg/database"
	gitHubActionPath       = ".github/workflows"
	dbSchemaPath           = "db/schema"
	dbQueryPath            = "db/query"
	pkgRepositoryPath      = "pkg/repository"
	dbSchemaUsersPath      = "db/schema/users"
	dbSchemaPostsPath      = "db/schema/posts"
	dbQueryUsersPath       = "db/query/users"
	dbQueryPostsPath       = "db/query/posts"
	pkgRepositoryUsersPath = "pkg/repository/users"
	pkgRepositoryPostsPath = "pkg/repository/posts"
)

// CheckOs checks Operation system and generates MakeFile and `go build` command
// Based on Project.Unixbase
func (p *Project) CheckOS() {
	p.OSCheck = make(map[string]bool)

	if runtime.GOOS != "windows" {
		p.OSCheck["UnixBased"] = true
	}
	if runtime.GOOS == "linux" {
		p.OSCheck["linux"] = true
	}
	if runtime.GOOS == "darwin" {
		p.OSCheck["darwin"] = true
	}
}

// ExitCLI checks if the Project has been exited, and closes
// out of the CLI if it has
func (p *Project) ExitCLI(tprogram *tea.Program) {
	if p.Exit {
		// logo render here
		if err := tprogram.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

// createFrameWorkMap adds the current supported
// Frameworks into a Project's FrameworkMap
func (p *Project) createFrameworkMap() {
	p.FrameworkMap[flags.Chi] = Framework{
		packageName: chiPackage,
		templater:   framework.ChiTemplates{},
	}

	p.FrameworkMap[flags.StandardLibrary] = Framework{
		packageName: []string{},
		templater:   framework.StandardLibTemplate{},
	}

	p.FrameworkMap[flags.Gin] = Framework{
		packageName: ginPackage,
		templater:   framework.GinTemplates{},
	}

	p.FrameworkMap[flags.Fiber] = Framework{
		packageName: fiberPackage,
		templater:   framework.FiberTemplates{},
	}

	p.FrameworkMap[flags.GorillaMux] = Framework{
		packageName: gorillaPackage,
		templater:   framework.GorillaTemplates{},
	}

	p.FrameworkMap[flags.HttpRouter] = Framework{
		packageName: routerPackage,
		templater:   framework.RouterTemplates{},
	}

	p.FrameworkMap[flags.Echo] = Framework{
		packageName: echoPackage,
		templater:   framework.EchoTemplates{},
	}
}

func (p *Project) createDBDriverMap() {
	p.DBDriverMap[flags.MySql] = Driver{
		packageName: mysqlDriver,
		templater:   dbdriver.MysqlTemplate{},
	}
	p.DBDriverMap[flags.Postgres] = Driver{
		packageName: postgresDriver,
		templater:   dbdriver.PostgresTemplate{},
	}
	p.DBDriverMap[flags.Sqlite] = Driver{
		packageName: sqliteDriver,
		templater:   dbdriver.SqliteTemplate{},
	}
	p.DBDriverMap[flags.Mongo] = Driver{
		packageName: mongoDriver,
		templater:   dbdriver.MongoTemplate{},
	}

	p.DBDriverMap[flags.Scylla] = Driver{
		packageName: gocqlDriver,
		templater:   dbdriver.ScyllaTemplate{},
	}
}

func (p *Project) createDockerMap() {
	p.DockerMap = make(map[flags.Database]Docker)

	p.DockerMap[flags.MySql] = Docker{
		packageName: []string{},
		templater:   docker.MysqlDockerTemplate{},
	}
	p.DockerMap[flags.Postgres] = Docker{
		packageName: []string{},
		templater:   docker.PostgresDockerTemplate{},
	}
	p.DockerMap[flags.Mongo] = Docker{
		packageName: []string{},
		templater:   docker.MongoDockerTemplate{},
	}
	p.DockerMap[flags.Scylla] = Docker{
		packageName: []string{},
		templater:   docker.ScyllaDockerTemplate{},
	}
}

// CreateMainFile creates the project folders and files,
// and writes to them depending on the selected options
func (p *Project) CreateMainFile() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0o754); err != nil {
			log.Printf("Could not create directory: %v", err)
			return err
		}
	}

	// Check if user.email is set.
	if p.GitOptions.String() != flags.Skip {

		emailSet, err := utils.CheckGitConfig("user.email")
		if err != nil {
			return err
		}
		if !emailSet {
			fmt.Println("user.email is not set in git config.")
			fmt.Println("Please set up git config before trying again.")
			panic("\nGIT CONFIG ISSUE: user.email is not set in git config.\n")
		}
	}

	p.ProjectName = strings.TrimSpace(p.ProjectName)

	// Create a new directory with the project name
	projectPath := filepath.Join(p.AbsolutePath, utils.GetRootDir(p.ProjectName))
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		err := os.MkdirAll(projectPath, 0o751)
		if err != nil {
			log.Printf("Error creating root project directory %v\n", err)
			return err
		}
	}

	// Define Operating system
	p.CheckOS()

	// Create the map for our program
	p.createFrameworkMap()

	// Create go.mod
	err := utils.InitGoMod(p.ProjectName, projectPath)
	if err != nil {
		log.Printf("Could not initialize go.mod in new project %v\n", err)
		return err
	}

	// Install the correct package for the selected framework
	if p.ProjectType != flags.StandardLibrary {
		err = utils.GoGetPackage(projectPath, p.FrameworkMap[p.ProjectType].packageName)
		if err != nil {
			log.Println("Could not install go dependency for the chosen framework")
			return err
		}
	}

	// Install the correct package for the selected driver
	if p.DBDriver != "none" {
		p.createDBDriverMap()
		err = utils.GoGetPackage(projectPath, p.DBDriverMap[p.DBDriver].packageName)
		if err != nil {
			log.Println("Could not install go dependency for chosen driver")
			return err
		}

		err = p.CreatePath(internalDatabasePath, projectPath)
		if err != nil {
			log.Printf("Error creating path: %s", internalDatabasePath)
			return err
		}

		err = p.CreateFileWithInjection(internalDatabasePath, projectPath, "database.go", "database")
		if err != nil {
			log.Printf("Error injecting database.go file: %v", err)
			return err
		}

		if p.DBDriver != "sqlite" {
			err = p.CreateFileWithInjection(internalDatabasePath, projectPath, "database_test.go", "integration-tests")
			if err != nil {
				log.Printf("Error injecting database_test.go file: %v", err)
				return err
			}
		}

		// Create sqlc files for SQL databases (MySQL, PostgreSQL, SQLite)
		if p.DBDriver == "mysql" || p.DBDriver == "postgres" || p.DBDriver == "sqlite" {
			// Create sqlc.yaml configuration file
			err = p.CreateFileWithInjection(root, projectPath, "sqlc.yaml", "sqlc-config")
			if err != nil {
				log.Printf("Error creating sqlc.yaml file: %v", err)
				return err
			}

			// Create db/schema directory and schema files
			err = p.CreatePath(dbSchemaPath, projectPath)
			if err != nil {
				log.Printf("Error creating path: %s", dbSchemaPath)
				return err
			}

			err = p.CreateFileWithInjection(dbSchemaPath, projectPath, "users.sql", "users-schema")
			if err != nil {
				log.Printf("Error creating users schema file: %v", err)
				return err
			}

			err = p.CreateFileWithInjection(dbSchemaPath, projectPath, "posts.sql", "posts-schema")
			if err != nil {
				log.Printf("Error creating posts schema file: %v", err)
				return err
			}

			// Create db/query/users directory and users queries
			err = p.CreatePath(dbQueryUsersPath, projectPath)
			if err != nil {
				log.Printf("Error creating path: %s", dbQueryUsersPath)
				return err
			}

			err = p.CreateFileWithInjection(dbQueryUsersPath, projectPath, "users.sql", "users-query")
			if err != nil {
				log.Printf("Error creating users query file: %v", err)
				return err
			}

			// Create db/query/posts directory and posts queries
			err = p.CreatePath(dbQueryPostsPath, projectPath)
			if err != nil {
				log.Printf("Error creating path: %s", dbQueryPostsPath)
				return err
			}

			err = p.CreateFileWithInjection(dbQueryPostsPath, projectPath, "posts.sql", "posts-query")
			if err != nil {
				log.Printf("Error creating posts query file: %v", err)
				return err
			}

			// Create pkg/repository/users directory (sqlc will generate files here)
			err = p.CreatePath(pkgRepositoryUsersPath, projectPath)
			if err != nil {
				log.Printf("Error creating path: %s", pkgRepositoryUsersPath)
				return err
			}

			// Create pkg/repository/posts directory (sqlc will generate files here)
			err = p.CreatePath(pkgRepositoryPostsPath, projectPath)
			if err != nil {
				log.Printf("Error creating path: %s", pkgRepositoryPostsPath)
				return err
			}
		}
	}

	// Create correct docker compose for the selected driver
	if p.DBDriver != "none" {
		if p.DBDriver != "sqlite" {
			p.createDockerMap()
			p.Docker = p.DBDriver

			err = p.CreateFileWithInjection(root, projectPath, "docker-compose.yml", "db-docker")
			if err != nil {
				log.Printf("Error injecting docker-compose.yml file: %v", err)
				return err
			}
		} else {
			fmt.Println(" SQLite doesn't support docker-compose.yml configuration")
		}
	}

	// Install the godotenv package
	err = utils.GoGetPackage(projectPath, godotenvPackage)
	if err != nil {
		log.Println("Could not install go dependency")

		return err
	}

	if p.DBDriver == flags.Scylla {
		replace := fmt.Sprintf("%s=%s", gocqlDriver[0], scyllaDriver)
		err = utils.GoModReplace(projectPath, replace)
		if err != nil {
			log.Printf("Could not replace go dependency %v\n", err)
			return err
		}
	}

	err = p.CreatePath(cmdApiPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		return err
	}

	err = p.CreateFileWithInjection(cmdApiPath, projectPath, "main.go", "main")
	if err != nil {
		return err
	}

	makeFile, err := os.Create(filepath.Join(projectPath, "Makefile"))
	if err != nil {
		return err
	}

	defer makeFile.Close()

	// inject makefile template
	makeFileTemplate := template.Must(template.New("makefile").Parse(string(framework.MakeTemplate())))
	err = makeFileTemplate.Execute(makeFile, p)
	if err != nil {
		return err
	}

	readmeFile, err := os.Create(filepath.Join(projectPath, "README.md"))
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	// inject readme template
	readmeFileTemplate := template.Must(template.New("readme").Parse(string(framework.ReadmeTemplate())))
	err = readmeFileTemplate.Execute(readmeFile, p)
	if err != nil {
		return err
	}

	err = p.CreatePath(internalServerPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", internalServerPath)
		return err
	}

	if p.AdvancedOptions[string(flags.React)] {
		// deselect htmx option automatically since react is selected
		p.AdvancedOptions[string(flags.Htmx)] = false
		if err := p.CreateViteReactProject(projectPath); err != nil {
			return fmt.Errorf("failed to set up React project: %w", err)
		}

		// if everything went smoothly, remove tailwing flag option
		p.AdvancedOptions[string(flags.Tailwind)] = false
	}

	if p.AdvancedOptions[string(flags.Tailwind)] {
		// select htmx option automatically since tailwind is selected
		p.AdvancedOptions[string(flags.Htmx)] = true

		err = os.MkdirAll(fmt.Sprintf("%s/%s/assets/css", projectPath, cmdWebPath), 0o755)
		if err != nil {
			return err
		}

		err = os.MkdirAll(fmt.Sprintf("%s/%s/styles", projectPath, cmdWebPath), 0o755)
		if err != nil {
			return fmt.Errorf("failed to create styles directory: %w", err)
		}

		inputCssFile, err := os.Create(fmt.Sprintf("%s/%s/styles/input.css", projectPath, cmdWebPath))
		if err != nil {
			return err
		}
		defer inputCssFile.Close()

		inputCssTemplate := advanced.InputCssTemplate()
		err = os.WriteFile(fmt.Sprintf("%s/%s/styles/input.css", projectPath, cmdWebPath), inputCssTemplate, 0o644)
		if err != nil {
			return err
		}

		outputCssFile, err := os.Create(fmt.Sprintf("%s/%s/assets/css/output.css", projectPath, cmdWebPath))
		if err != nil {
			return err
		}
		defer outputCssFile.Close()

		outputCssTemplate := advanced.OutputCssTemplate()
		err = os.WriteFile(fmt.Sprintf("%s/%s/assets/css/output.css", projectPath, cmdWebPath), outputCssTemplate, 0o644)
		if err != nil {
			return err
		}
	}

	if p.AdvancedOptions[string(flags.Htmx)] {
		// create folders and hello world file
		err = p.CreatePath(cmdWebPath, projectPath)
		if err != nil {
			return err
		}
		helloTemplFile, err := os.Create(fmt.Sprintf("%s/%s/hello.templ", projectPath, cmdWebPath))
		if err != nil {
			return err
		}
		defer helloTemplFile.Close()

		// inject hello.templ template
		helloTemplTemplate := template.Must(template.New("hellotempl").Parse((string(advanced.HelloTemplTemplate()))))
		err = helloTemplTemplate.Execute(helloTemplFile, p)
		if err != nil {
			return err
		}

		baseTemplFile, err := os.Create(fmt.Sprintf("%s/%s/base.templ", projectPath, cmdWebPath))
		if err != nil {
			return err
		}
		defer baseTemplFile.Close()

		baseTemplTemplate := template.Must(template.New("basetempl").Parse((string(advanced.BaseTemplTemplate()))))
		err = baseTemplTemplate.Execute(baseTemplFile, p)
		if err != nil {
			return err
		}

		err = os.MkdirAll(fmt.Sprintf("%s/%s/assets/js", projectPath, cmdWebPath), 0o755)
		if err != nil {
			return err
		}

		htmxMinJsFile, err := os.Create(fmt.Sprintf("%s/%s/assets/js/htmx.min.js", projectPath, cmdWebPath))
		if err != nil {
			return err
		}
		defer htmxMinJsFile.Close()

		htmxMinJsTemplate := advanced.HtmxJSTemplate()
		err = os.WriteFile(fmt.Sprintf("%s/%s/assets/js/htmx.min.js", projectPath, cmdWebPath), htmxMinJsTemplate, 0o644)
		if err != nil {
			return err
		}

		efsFile, err := os.Create(fmt.Sprintf("%s/%s/efs.go", projectPath, cmdWebPath))
		if err != nil {
			return err
		}
		defer efsFile.Close()

		efsTemplate := template.Must(template.New("efs").Parse((string(advanced.EfsTemplate()))))
		err = efsTemplate.Execute(efsFile, p)
		if err != nil {
			return err
		}
		err = utils.GoGetPackage(projectPath, templPackage)
		if err != nil {
			log.Println("Could not install go dependency")
			return err
		}

		helloGoFile, err := os.Create(fmt.Sprintf("%s/%s/hello.go", projectPath, cmdWebPath))
		if err != nil {
			return err
		}
		defer efsFile.Close()

		if p.ProjectType == "fiber" {
			helloGoTemplate := template.Must(template.New("efs").Parse((string(advanced.HelloFiberGoTemplate()))))
			err = helloGoTemplate.Execute(helloGoFile, p)
			if err != nil {
				return err
			}
			err = utils.GoGetPackage(projectPath, []string{"github.com/gofiber/fiber/v2/middleware/adaptor"})
			if err != nil {
				log.Println("Could not install go dependency")
				return err
			}
		} else {
			helloGoTemplate := template.Must(template.New("efs").Parse((string(advanced.HelloGoTemplate()))))
			err = helloGoTemplate.Execute(helloGoFile, p)
			if err != nil {
				return err
			}
		}

	}

	// Create .github/workflows folder and inject release.yml and go-test.yml
	if p.AdvancedOptions[string(flags.GoProjectWorkflow)] {
		err = p.CreatePath(gitHubActionPath, projectPath)
		if err != nil {
			log.Printf("Error creating path: %s", gitHubActionPath)
			return err
		}

		err = p.CreateFileWithInjection(gitHubActionPath, projectPath, "release.yml", "releaser")
		if err != nil {
			log.Printf("Error injecting release.yml file: %v", err)
			return err
		}

		err = p.CreateFileWithInjection(gitHubActionPath, projectPath, "go-test.yml", "go-test")
		if err != nil {
			log.Printf("Error injecting go-test.yml file: %v", err)
			return err
		}

		err = p.CreateFileWithInjection(root, projectPath, ".goreleaser.yml", "releaser-config")
		if err != nil {
			log.Printf("Error injecting .goreleaser.yml file: %v", err)
			return err
		}
	}

	// if the websocket option is checked, a websocket dependency needs to
	// be added to the routes depending on the framework choosen.
	// Only fiber uses a different websocket library, the other frameworks
	// all work with the same one
	if p.AdvancedOptions[string(flags.Websocket)] {
		p.CreateWebsocketImports(projectPath)
	}

	if p.AdvancedOptions[string(flags.Docker)] {
		Dockerfile, err := os.Create(filepath.Join(projectPath, "Dockerfile"))
		if err != nil {
			return err
		}
		defer Dockerfile.Close()

		// inject Docker template
		dockerfileTemplate := template.Must(template.New("Dockerfile").Parse(string(advanced.Dockerfile())))
		err = dockerfileTemplate.Execute(Dockerfile, p)
		if err != nil {
			return err
		}

		if p.DBDriver == "none" || p.DBDriver == "sqlite" {

			Dockercompose, err := os.Create(filepath.Join(projectPath, "docker-compose.yml"))
			if err != nil {
				return err
			}
			defer Dockercompose.Close()

			// inject DockerCompose template
			dockerComposeTemplate := template.Must(template.New("docker-compose.yml").Parse(string(advanced.DockerCompose())))
			err = dockerComposeTemplate.Execute(Dockercompose, p)
			if err != nil {
				return err
			}
		}
	}

	if p.AdvancedOptions[string(flags.Worker)] {
		err := p.CreateWorkerFiles(projectPath)
		if err != nil {
			return err
		}
	}

	if p.AdvancedOptions[string(flags.Kafka)] {
		err := p.CreateKafkaFiles(projectPath)
		if err != nil {
			return err
		}
	}

	if p.AdvancedOptions[string(flags.Redis)] {
		err := p.CreateRedisFiles(projectPath)
		if err != nil {
			return err
		}
	}

	if p.AdvancedOptions[string(flags.Swagger)] {
		err := utils.GoGetPackage(projectPath, swaggerPackage)
		if err != nil {
			log.Println("Could not install go dependency for Swagger")
			return err
		}

		err = p.CreateSwaggerFiles(projectPath)
		if err != nil {
			return err
		}
	}

	// Call CreateAdvancedTemplates if any advanced features that need template injection are enabled
	if p.AdvancedOptions[string(flags.Htmx)] || p.AdvancedOptions[string(flags.Swagger)] {
		p.CreateAdvancedTemplates()
	}

	// Using the embedded static files for tailwind and htmx
	err = p.CreateFileWithInjection(internalServerPath, projectPath, "routes.go", "routes")
	if err != nil {
		log.Printf("Error injecting routes.go file: %v", err)
		return err
	}

	err = p.CreateFileWithInjection(internalServerPath, projectPath, "routes_test.go", "tests")
	if err != nil {
		return err
	}

	err = p.CreateFileWithInjection(internalServerPath, projectPath, "server.go", "server")
	if err != nil {
		log.Printf("Error injecting server.go file: %v", err)
		return err
	}

	err = p.CreateFileWithInjection(root, projectPath, ".env", "env")
	if err != nil {
		log.Printf("Error injecting .env file: %v", err)
		return err
	}

	// Create gitignore
	gitignoreFile, err := os.Create(filepath.Join(projectPath, ".gitignore"))
	if err != nil {
		return err
	}
	defer gitignoreFile.Close()

	// inject gitignore template
	gitignoreTemplate := template.Must(template.New(".gitignore").Parse(string(framework.GitIgnoreTemplate())))
	err = gitignoreTemplate.Execute(gitignoreFile, p)
	if err != nil {
		return err
	}

	// Create .air.toml file
	airTomlFile, err := os.Create(filepath.Join(projectPath, ".air.toml"))
	if err != nil {
		return err
	}

	defer airTomlFile.Close()

	// inject air.toml template
	airTomlTemplate := template.Must(template.New("airtoml").Parse(string(framework.AirTomlTemplate())))
	err = airTomlTemplate.Execute(airTomlFile, p)
	if err != nil {
		return err
	}

	err = utils.GoTidy(projectPath)
	if err != nil {
		log.Printf("Could not go tidy in new project %v\n", err)
		return err
	}

	err = utils.GoFmt(projectPath)
	if err != nil {
		log.Printf("Could not gofmt in new project %v\n", err)
		return err
	}

	if p.GitOptions != flags.Skip {
		nameSet, err := utils.CheckGitConfig("user.name")
		if err != nil {
			return err
		}

		if !nameSet {
			fmt.Println("user.name is not set in git config.")
			fmt.Println("Please set up git config before trying again.")
			panic("\nGIT CONFIG ISSUE: user.name is not set in git config.\n")
		}
		// Initialize git repo
		err = utils.ExecuteCmd("git", []string{"init"}, projectPath)
		if err != nil {
			log.Printf("Error initializing git repo: %v", err)
			return err
		}

		// Git add files
		err = utils.ExecuteCmd("git", []string{"add", "."}, projectPath)
		if err != nil {
			log.Printf("Error adding files to git repo: %v", err)
			return err
		}

		if p.GitOptions == flags.Commit {
			// Git commit files
			err = utils.ExecuteCmd("git", []string{"commit", "-m", "Initial commit"}, projectPath)
			if err != nil {
				log.Printf("Error committing files to git repo: %v", err)
				return err
			}
		}
	}
	return nil
}

// CreatePath creates the given directory in the projectPath
func (p *Project) CreatePath(pathToCreate string, projectPath string) error {
	path := filepath.Join(projectPath, pathToCreate)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0o751)
		if err != nil {
			log.Printf("Error creating directory %v\n", err)
			return err
		}
	}

	return nil
}

// CreateFileWithInjection creates the given file at the
// project path, and injects the appropriate template
func (p *Project) CreateFileWithInjection(pathToCreate string, projectPath string, fileName string, methodName string) error {
	createdFile, err := os.Create(filepath.Join(projectPath, pathToCreate, fileName))
	if err != nil {
		return err
	}

	defer createdFile.Close()

	switch methodName {
	case "main":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.Main())))
		err = createdTemplate.Execute(createdFile, p)
	case "server":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.Server())))
		err = createdTemplate.Execute(createdFile, p)
	case "routes":
		routeFileBytes := p.FrameworkMap[p.ProjectType].templater.Routes()
		createdTemplate := template.Must(template.New(fileName).Parse(string(routeFileBytes)))
		err = createdTemplate.Execute(createdFile, p)
	case "releaser":
		createdTemplate := template.Must(template.New(fileName).Parse(string(advanced.Releaser())))
		err = createdTemplate.Execute(createdFile, p)
	case "go-test":
		createdTemplate := template.Must(template.New(fileName).Parse(string(advanced.Test())))
		err = createdTemplate.Execute(createdFile, p)
	case "releaser-config":
		createdTemplate := template.Must(template.New(fileName).Parse(string(advanced.ReleaserConfig())))
		err = createdTemplate.Execute(createdFile, p)
	case "database":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.Service())))
		err = createdTemplate.Execute(createdFile, p)
	case "db-docker":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DockerMap[p.Docker].templater.Docker())))
		err = createdTemplate.Execute(createdFile, p)
	case "integration-tests":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.Tests())))
		err = createdTemplate.Execute(createdFile, p)
	case "tests":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.TestHandler())))
		err = createdTemplate.Execute(createdFile, p)
	case "env":
		if p.DBDriver != "none" {

			envBytes := [][]byte{
				tpl.GlobalEnvTemplate(),
				p.DBDriverMap[p.DBDriver].templater.Env(),
			}

			if p.AdvancedOptions[string(flags.Worker)] {
				envBytes = append(envBytes, advanced.WorkerEnvTemplate())
			}

			if p.AdvancedOptions[string(flags.Kafka)] {
				envBytes = append(envBytes, advanced.KafkaEnvTemplate())
			}

			if p.AdvancedOptions[string(flags.Redis)] {
				envBytes = append(envBytes, advanced.RedisEnv())
			}

			createdTemplate := template.Must(template.New(fileName).Parse(string(bytes.Join(envBytes, []byte("\n")))))
			err = createdTemplate.Execute(createdFile, p)

		} else {
			envBytes := [][]byte{
				tpl.GlobalEnvTemplate(),
			}

			if p.AdvancedOptions[string(flags.Worker)] {
				envBytes = append(envBytes, advanced.WorkerEnvTemplate())
			}

			if p.AdvancedOptions[string(flags.Kafka)] {
				envBytes = append(envBytes, advanced.KafkaEnvTemplate())
			}

			if p.AdvancedOptions[string(flags.Redis)] {
				envBytes = append(envBytes, advanced.RedisEnv())
			}

			createdTemplate := template.Must(template.New(fileName).Parse(string(bytes.Join(envBytes, []byte("\n")))))
			err = createdTemplate.Execute(createdFile, p)
		}
	case "sqlc-config":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.SqlcConfig())))
		err = createdTemplate.Execute(createdFile, p)
	case "schema-example":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.SchemaExample())))
		err = createdTemplate.Execute(createdFile, p)
	case "query-example":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.QueryExample())))
		err = createdTemplate.Execute(createdFile, p)
	case "users-schema":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.UsersSchema())))
		err = createdTemplate.Execute(createdFile, p)
	case "posts-schema":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.PostsSchema())))
		err = createdTemplate.Execute(createdFile, p)
	case "users-query":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.UsersQuery())))
		err = createdTemplate.Execute(createdFile, p)
	case "posts-query":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBDriverMap[p.DBDriver].templater.PostsQuery())))
		err = createdTemplate.Execute(createdFile, p)
	}

	if err != nil {
		return err
	}

	return nil
}

// checkNpmInstalled checks if npm is installed on the system
func checkNpmInstalled() error {
	_, err := exec.LookPath("npm")
	if err != nil {
		return fmt.Errorf("npm is not installed or not in PATH: %w", err)
	}
	return nil
}

func (p *Project) CreateViteReactProject(projectPath string) error {
	if err := checkNpmInstalled(); err != nil {
		return err
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			fmt.Fprintf(os.Stderr, "failed to change back to original directory: %v\n", err)
		}
	}()

	// change into the project directory to run vite command
	err = os.Chdir(projectPath)
	if err != nil {
		fmt.Println("failed to change into project directory: %w", err)
	}

	// the interactive vite command will not work as we can't interact with it
	fmt.Println("Installing create-vite (using cache if available)...")
	cmd := exec.Command("npm", "create", "vite@latest", "frontend", "--",
		"--template", "react-ts",
		"--prefer-offline",
		"--no-fund")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to use create-vite: %w", err)
	}

	frontendPath := filepath.Join(projectPath, "frontend")
	if err := os.MkdirAll(frontendPath, 0755); err != nil {
		return fmt.Errorf("failed to create frontend directory: %w", err)
	}

	if err := os.Chdir(frontendPath); err != nil {
		return fmt.Errorf("failed to change to frontend directory: %w", err)
	}

	srcDir := filepath.Join(frontendPath, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return fmt.Errorf("failed to create src directory: %w", err)
	}

	if err := os.WriteFile(filepath.Join(srcDir, "App.tsx"), advanced.ReactAppfile(), 0644); err != nil {
		return fmt.Errorf("failed to write App.tsx template: %w", err)
	}

	// Create the global `.env` file from the template
	err = p.CreateFileWithInjection("", projectPath, ".env", "env")
	if err != nil {
		return fmt.Errorf("failed to create global .env file: %w", err)
	}

	// Read from the global `.env` file and create the frontend-specific `.env`
	globalEnvPath := filepath.Join(projectPath, ".env")
	vitePort := "8080" // Default fallback

	// Read the global .env file
	if data, err := os.ReadFile(globalEnvPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PORT=") {
				vitePort = strings.SplitN(line, "=", 2)[1] // Get the backend port value
				break
			}
		}
	}

	// Use a template to generate the frontend .env file
	frontendEnvContent := fmt.Sprintf("VITE_PORT=%s\n", vitePort)
	if err := os.WriteFile(filepath.Join(frontendPath, ".env"), []byte(frontendEnvContent), 0644); err != nil {
		return fmt.Errorf("failed to create frontend .env file: %w", err)
	}

	// Handle Tailwind configuration if selected
	if p.AdvancedOptions[string(flags.Tailwind)] {
		fmt.Println("Installing Tailwind dependencies (using cache if available)...")
		cmd := exec.Command("npm", "install",
			"--prefer-offline",
			"--no-fund",
			"tailwindcss@^4", "@tailwindcss/vite")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Tailwind: %w", err)
		}

		// Create the vite + react + Tailwind v4 configuration
		if err := os.WriteFile(filepath.Join(frontendPath, "vite.config.ts"), advanced.ViteTailwindConfigFile(), 0644); err != nil {
			return fmt.Errorf("failed to write vite.config.ts: %w", err)
		}

		srcDir := filepath.Join(frontendPath, "src")
		if err := os.MkdirAll(srcDir, 0755); err != nil {
			return fmt.Errorf("failed to create src directory: %w", err)
		}

		err = os.WriteFile(filepath.Join(srcDir, "index.css"), advanced.InputCssTemplateReact(), 0644)
		if err != nil {
			return fmt.Errorf("failed to update index.css: %w", err)
		}

		if err := os.WriteFile(filepath.Join(srcDir, "App.tsx"), advanced.ReactTailwindAppfile(), 0644); err != nil {
			return fmt.Errorf("failed to write App.tsx template: %w", err)
		}

		if err := os.Remove(filepath.Join(srcDir, "App.css")); err != nil {
			// Don't return error if file doesn't exist
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove App.css: %w", err)
			}
		}

		// set to false to not re-do in next step
		p.AdvancedOptions[string(flags.Tailwind)] = false
	}

	return nil
}

func (p *Project) CreateAdvancedTemplates() {
	routesPlaceHolder := ""
	importsPlaceHolder := ""
	if p.AdvancedOptions[string(flags.Htmx)] {
		routesPlaceHolder += string(p.FrameworkMap[p.ProjectType].templater.HtmxTemplRoutes())
		importsPlaceHolder += string(p.FrameworkMap[p.ProjectType].templater.HtmxTemplImports())
	}

	if p.AdvancedOptions[string(flags.Swagger)] {
		routesPlaceHolder += "\n\t" + string(p.FrameworkMap[p.ProjectType].templater.SwaggerRoutes())
		importsPlaceHolder += "\n" + string(p.FrameworkMap[p.ProjectType].templater.SwaggerImports())
	}

	routeTmpl, err := template.New("routes").Parse(routesPlaceHolder)
	if err != nil {
		log.Fatal(err)
	}
	importTmpl, err := template.New("imports").Parse(importsPlaceHolder)
	if err != nil {
		log.Fatal(err)
	}
	var routeBuffer bytes.Buffer
	var importBuffer bytes.Buffer
	err = routeTmpl.Execute(&routeBuffer, p)
	if err != nil {
		log.Fatal(err)
	}
	err = importTmpl.Execute(&importBuffer, p)
	if err != nil {
		log.Fatal(err)
	}
	p.AdvancedTemplates.TemplateRoutes = routeBuffer.String()
	p.AdvancedTemplates.TemplateImports = importBuffer.String()
}

func (p *Project) CreateWebsocketImports(appDir string) {
	websocketDependency := []string{"github.com/coder/websocket"}
	if p.ProjectType == flags.Fiber {
		websocketDependency = []string{"github.com/gofiber/contrib/websocket"}
	}

	// Websockets require a different package depending on what framework is
	// choosen. The application calls go mod tidy at the end so we don't
	// have to here
	err := utils.GoGetPackage(appDir, websocketDependency)
	if err != nil {
		log.Fatal(err)
	}

	importsPlaceHolder := string(p.FrameworkMap[p.ProjectType].templater.WebsocketImports())

	importTmpl, err := template.New("imports").Parse(importsPlaceHolder)
	if err != nil {
		log.Fatalf("CreateWebsocketImports failed to create template: %v", err)
	}
	var importBuffer bytes.Buffer
	err = importTmpl.Execute(&importBuffer, p)
	if err != nil {
		log.Fatalf("CreateWebsocketImports failed write template: %v", err)
	}
	newImports := strings.Join([]string{string(p.AdvancedTemplates.TemplateImports), importBuffer.String()}, "\n")
	p.AdvancedTemplates.TemplateImports = newImports
}

func (p *Project) CreateWorkerFiles(appDir string) error {
	// Install Worker dependencies
	workerDependencies := []string{
		"github.com/hibiken/asynq",
		"github.com/joho/godotenv",
	}
	err := utils.GoGetPackage(appDir, workerDependencies)
	if err != nil {
		return err
	}

	// Create cmd/worker directory
	workerDir := filepath.Join(appDir, "cmd", "worker")
	err = os.MkdirAll(workerDir, 0755)
	if err != nil {
		return err
	}

	// Create cmd/worker/main.go file
	workerMainFile, err := os.Create(filepath.Join(workerDir, "main.go"))
	if err != nil {
		return err
	}
	defer workerMainFile.Close()

	workerMainTemplate := template.Must(template.New("main.go").Parse(string(advanced.WorkerMainTemplate())))
	err = workerMainTemplate.Execute(workerMainFile, p)
	if err != nil {
		return err
	}

	// Create cmd/worker/tasks directory
	tasksDir := filepath.Join(workerDir, "tasks")
	err = os.MkdirAll(tasksDir, 0755)
	if err != nil {
		return err
	}

	// Create cmd/worker/tasks/hello_world_task.go file
	helloWorldTaskFile, err := os.Create(filepath.Join(tasksDir, "hello_world_task.go"))
	if err != nil {
		return err
	}
	defer helloWorldTaskFile.Close()

	helloWorldTaskTemplate := template.Must(template.New("hello_world_task.go").Parse(string(advanced.WorkerHelloWorldTaskTemplate())))
	err = helloWorldTaskTemplate.Execute(helloWorldTaskFile, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) CreateKafkaFiles(appDir string) error {
	// Install Kafka dependency
	kafkaDependency := []string{"github.com/segmentio/kafka-go"}
	err := utils.GoGetPackage(appDir, kafkaDependency)
	if err != nil {
		return err
	}

	// Create pkg/kafka/segmentio directory
	kafkaDir := filepath.Join(appDir, "pkg", "kafka", "segmentio")
	err = os.MkdirAll(kafkaDir, 0755)
	if err != nil {
		return err
	}

	// Create consumer.go file
	consumerFile, err := os.Create(filepath.Join(kafkaDir, "consumer.go"))
	if err != nil {
		return err
	}
	defer consumerFile.Close()

	consumerTemplate := template.Must(template.New("consumer.go").Parse(string(advanced.KafkaConsumerTemplate())))
	err = consumerTemplate.Execute(consumerFile, p)
	if err != nil {
		return err
	}

	// Create consumer_test.go file
	testFile, err := os.Create(filepath.Join(kafkaDir, "consumer_test.go"))
	if err != nil {
		return err
	}
	defer testFile.Close()

	testTemplate := template.Must(template.New("consumer_test.go").Parse(string(advanced.KafkaConsumerTestTemplate())))
	err = testTemplate.Execute(testFile, p)
	if err != nil {
		return err
	}

	// Create cmd/consumer directory for the consumer binary
	consumerCmdDir := filepath.Join(appDir, "cmd", "consumer")
	err = os.MkdirAll(consumerCmdDir, 0755)
	if err != nil {
		return err
	}

	// Create cmd/consumer/main.go file
	consumerMainFile, err := os.Create(filepath.Join(consumerCmdDir, "main.go"))
	if err != nil {
		return err
	}
	defer consumerMainFile.Close()

	consumerMainTemplate := template.Must(template.New("main.go").Parse(string(advanced.KafkaConsumerMainTemplate())))
	err = consumerMainTemplate.Execute(consumerMainFile, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) CreateRedisFiles(appDir string) error {
	// Install Redis dependency
	redisDependency := []string{"github.com/redis/go-redis/v9"}
	err := utils.GoGetPackage(appDir, redisDependency)
	if err != nil {
		return err
	}

	// Create pkg/redis directory
	redisDir := filepath.Join(appDir, "pkg", "redis")
	err = os.MkdirAll(redisDir, 0755)
	if err != nil {
		return err
	}

	// Create interface.go file
	interfaceFile, err := os.Create(filepath.Join(redisDir, "interface.go"))
	if err != nil {
		return err
	}
	defer interfaceFile.Close()

	interfaceTemplate := template.Must(template.New("interface.go").Parse(string(advanced.RedisInterface())))
	err = interfaceTemplate.Execute(interfaceFile, p)
	if err != nil {
		return err
	}

	// Create client.go file
	clientFile, err := os.Create(filepath.Join(redisDir, "client.go"))
	if err != nil {
		return err
	}
	defer clientFile.Close()

	clientTemplate := template.Must(template.New("client.go").Parse(string(advanced.RedisImplementation())))
	err = clientTemplate.Execute(clientFile, p)
	if err != nil {
		return err
	}

	// Create config.go file
	configFile, err := os.Create(filepath.Join(redisDir, "config.go"))
	if err != nil {
		return err
	}
	defer configFile.Close()

	configTemplate := template.Must(template.New("config.go").Parse(string(advanced.RedisConfig())))
	err = configTemplate.Execute(configFile, p)
	if err != nil {
		return err
	}

	// Create internal/example_service directory
	exampleServiceDir := filepath.Join(appDir, "internal", "example_service")
	err = os.MkdirAll(exampleServiceDir, 0755)
	if err != nil {
		return err
	}

	// Create example_service.go file
	exampleServiceFile, err := os.Create(filepath.Join(exampleServiceDir, "redis_example.go"))
	if err != nil {
		return err
	}
	defer exampleServiceFile.Close()

	exampleServiceTemplate := template.Must(template.New("redis_example.go").Parse(string(advanced.RedisExampleService())))
	err = exampleServiceTemplate.Execute(exampleServiceFile, p)
	if err != nil {
		return err
	}

	return nil
}

// CreateSwaggerFiles creates necessary docs directory and documentation files for Swagger
func (p *Project) CreateSwaggerFiles(appDir string) error {
	// Create docs directory
	docsDir := filepath.Join(appDir, "docs")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		return err
	}

	// Create docs.go file that includes swagger documentation
	docsFile, err := os.Create(filepath.Join(docsDir, "docs.go"))
	if err != nil {
		return err
	}
	defer docsFile.Close()

	// Write the docs.go content directly
	docsContent := fmt.Sprintf(`package docs

import "github.com/swaggo/swag"

const docTemplate = `+"`"+`{
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {}
}`+"`"+`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "%s API",
	Description:      "This is a sample server for %s.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
`, p.ProjectName, p.ProjectName)

	_, err = docsFile.WriteString(docsContent)
	if err != nil {
		return err
	}

	// Create swagger.json and swagger.yaml placeholder files
	swaggerJSONFile, err := os.Create(filepath.Join(docsDir, "swagger.json"))
	if err != nil {
		return err
	}
	defer swaggerJSONFile.Close()

	swaggerYAMLFile, err := os.Create(filepath.Join(docsDir, "swagger.yaml"))
	if err != nil {
		return err
	}
	defer swaggerYAMLFile.Close()

	// Create README file for Swagger documentation
	readmeFile, err := os.Create(filepath.Join(appDir, "README_SWAGGER.md"))
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	_, err = readmeFile.Write(advanced.SwaggerReadmeTemplate())
	if err != nil {
		return err
	}

	return nil
}
