package main

import (
	"fmt"
	"github.com/tomwright/finance-planner/internal/application/service"
	"github.com/tomwright/finance-planner/internal/command"
	"github.com/tomwright/finance-planner/internal/repository"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("could not get users working dir: %s", err)
		os.Exit(1)
	}
	storageDir := filepath.Join(u.HomeDir, "finance_planner")
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		fmt.Printf("could not create storage dir: %s", err)
		os.Exit(1)
	}

	profileRepo := repository.NewProfile(storageDir)
	transactionRepo := repository.NewTransaction(storageDir)

	profileService := service.NewProfileService(profileRepo, transactionRepo)

	rootCmd := command.Load(profileService)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
