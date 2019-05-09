package main

import (
	"fmt"
	"github.com/tomwright/finance-planner/internal/application/service"
	"github.com/tomwright/finance-planner/internal/application/validate"
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

	db, err := repository.ConnectSQLite(storageDir)
	if err != nil {
		fmt.Printf("could not create storage dir: %s", err)
		os.Exit(1)
	}

	profileRepo := repository.NewSQLiteProfile(db)
	if err := profileRepo.Init(); err != nil {
		fmt.Printf("could not init profile repo: %s", err)
		os.Exit(1)
	}
	transactionRepo := repository.NewSQLiteTransaction(db)
	if err := transactionRepo.Init(); err != nil {
		fmt.Printf("could not init transaction repo: %s", err)
		os.Exit(1)
	}

	validator := validate.NewValidator(profileRepo, transactionRepo)

	profileService := service.NewProfileService(profileRepo, transactionRepo, validator)

	rootCmd := command.Load(profileService)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
