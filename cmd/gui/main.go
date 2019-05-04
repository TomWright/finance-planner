package main

import (
	"fmt"
	"fyne.io/fyne"
	fyneApp "fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/service"
	"github.com/tomwright/finance-planner/internal/repository"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	app                fyne.App
	profileInputWindow fyne.Window
	profileWindow      fyne.Window

	profileService service.Profile
)

func main() {
	loadDeps()

	app = fyneApp.New()

	ProfileInputWindow()

	profileInputWindow.ShowAndRun()
}

func loadDeps() {
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

	profileService = service.NewProfileService(profileRepo, transactionRepo)
}

func refreshProfileWindowTransactionList(profile *domain.Profile, infoLabel *widget.Label, transactionList *widget.Entry) {
	infoLabel.SetText("Refreshing...")
	transactionList.SetText("")
	err := profile.Transactions.Range(nil, func(t domain.Transaction) error {
		txt := transactionList.Text
		if txt != "" {
			txt += "\n"
		}
		txt += fmt.Sprintf("%s: %d [%s]", t.Label, t.Amount, strings.Join(t.Tags, ", "))
		transactionList.SetText(txt)
		return nil
	})
	if err != nil {
		infoLabel.SetText("Error: " + err.Error())
	} else {
		infoLabel.SetText("Refreshed")
	}
}

func ProfileWindow(profileName string) (func(), error) {
	profile, err := profileService.LoadOrCreateProfile(profileName)
	if err != nil {
		return func() {}, err
	}

	var w fyne.Window
	if profileWindow != nil {
		w = profileWindow
		w.SetTitle("Finance - " + profile.Name)
	} else {
		w = app.NewWindow("Finance - " + profile.Name)
	}

	infoLabel := widget.NewLabel("")
	infoLabel.Show()

	transactionList := widget.NewMultiLineEntry()
	transactionList.Show()

	newTransactionLabel := widget.NewEntry()
	newTransactionAmount := widget.NewEntry()
	newTransactionSubmitBtn := widget.NewButton("Add Transaction", func() {
		amount, err := strconv.ParseInt(newTransactionAmount.Text, 10, 64)
		if err != nil {
			infoLabel.SetText("Error: Invalid transaction amount: " + err.Error())
			return
		}
		t := domain.NewTransaction().
			WithLabel(newTransactionLabel.Text).
			WithAmount(amount)

		if err := t.Validate(); err != nil {
			infoLabel.SetText("Error: " + err.Error())
			return
		}

		profile.Transactions.Add(t)
		if err := profileService.SaveProfile(profile); err != nil {
			infoLabel.SetText("Error: " + err.Error())
			return
		}

		newTransactionLabel.SetText("")
		newTransactionAmount.SetText("")

		refreshProfileWindowTransactionList(profile, infoLabel, transactionList)
	})

	newTransactionBox := widget.NewHBox(newTransactionLabel, newTransactionAmount, newTransactionSubmitBtn)

	content := widget.NewVBox(
		newTransactionBox,
		transactionList,
		infoLabel,
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 300))

	profileWindow = w

	return func() {
		refreshProfileWindowTransactionList(profile, infoLabel, transactionList)
	}, nil
}

func ProfileInputWindow() {
	w := app.NewWindow("Finance - Enter profile")

	errorLabel := widget.NewLabel("")
	errorLabel.Hidden = true

	profileEntry := widget.NewEntry()
	profileSubmit := widget.NewButton("Submit", func() {
		errorLabel.Hidden = false
		errorLabel.SetText("Loading profile...")
		initFn, err := ProfileWindow(profileEntry.Text)
		if err != nil {
			errorLabel.SetText("Error: " + err.Error())
		} else {
			errorLabel.Hidden = true
			profileInputWindow.Hide()
			profileWindow.Show()
			if initFn != nil {
				initFn()
			}
		}
	})

	content := widget.NewVBox(
		profileEntry,
		profileSubmit,
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
		errorLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(480, 130))

	profileInputWindow = w
}
