package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"user-access-management/handlers"
	"user-access-management/jwtparser"
	"user-access-management/middleware"
	"user-access-management/routes"
	"user-access-management/service"

	"github.com/gin-gonic/gin"
)

var (
	// Derived from ldflags -X
	buildRevision string
	buildVersion  string
	buildTime     string

	// general options
	versionFlag bool
	helpFlag    bool

	// server port
	port string

	// program controller
	done = make(chan struct{})
	errc = make(chan error)
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "show current version and exit")
	flag.BoolVar(&helpFlag, "help", false, "show usage and exit")
	flag.StringVar(&port, "port", ":3540", "server port")
}

func setBuildVariables() {
	if buildRevision == "" {
		buildRevision = "dev"
	}
	if buildVersion == "" {
		buildVersion = "dev"
	}
	if buildTime == "" {
		buildTime = time.Now().UTC().Format(time.RFC3339)
	}
}

func parseFlags() {
	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Printf("%s %s %s\n", buildRevision, buildVersion, buildTime)
		os.Exit(0)
	}
}

func openDB() (*sql.DB, error) {
	var (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "user-access-management"
	)

	psqlInfo := os.Getenv("POSTGRESQL_CONN_STRING")
	if len(psqlInfo) == 0 {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func handleInterrupts() {
	log.Println("start handle interrupts")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	sig := <-interrupt
	log.Printf("caught sig: %v", sig)
	// close resource here
	done <- struct{}{}
}

func main() {
	setBuildVariables()
	parseFlags()

	//read environment variables
	tokenServiceBaseURL := os.Getenv("TOKEN_SERVICE_URL")
	if len(tokenServiceBaseURL) == 0 {
		log.Println("error: could not read TOKEN_SERVICE_URL from environment")
		return
	}

	go handleInterrupts()

	server := gin.Default()

	psqlInfo, err := openDB()
	if err != nil {
		log.Printf("error connecting DB: %v", err)
		return
	}
	log.Println("DB connection is successful")
	defer psqlInfo.Close()

	userManagementService := service.NewUserManagementService(psqlInfo)
	tokenManager := jwtparser.NewJWTTokenManager()
	userAccessManagementHandler := handlers.NewUserAccessManagementHandler(userManagementService, tokenManager)

	apiRoutes := routes.NewRoutes(userAccessManagementHandler)
	authMiddleware := middleware.NewAuthMiddleware(tokenServiceBaseURL)
	routes.AttachRoutes(server, apiRoutes, authMiddleware)

	go func() {
		errc <- server.Run(port)
	}()

	select {
	case err := <-errc:
		log.Printf("ListenAndServe error: %v", err)
	case <-done:
		log.Println("shutting down the server ...")
	}
	time.AfterFunc(1*time.Second, func() {
		close(done)
		close(errc)
	})
}
