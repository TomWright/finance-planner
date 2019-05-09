package command

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/service"
	"math"
	"os"
	"strings"
)

func ListTransactions(profileService service.Profile) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-transactions",
		Short: "List all transactions for the profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName, _ := cmd.Flags().GetString("profile")

			profile, err := profileService.LoadProfileByName(profileName)
			if err != nil {
				return err
			}

			fmt.Printf("Profile: %s\n", profile.Name)

			outputTransactions("Incoming Transactions", profile.Transactions.Subset(func(t *domain.Transaction) bool {
				return t.Amount > 0
			}))
			outputTransactions("Outgoing Transactions", profile.Transactions.Subset(func(t *domain.Transaction) bool {
				return t.Amount < 0
			}))

			fmt.Printf("End balance: £%v\n", float64(profile.Transactions.Sum())/100)

			return nil
		},
	}

	cmd.Flags().String("profile", "", "Profile to interact with")

	_ = cmd.MarkFlagRequired("profile")

	return cmd
}

func outputTransactions(title string, collection *domain.TransactionCollection) {
	outputTable := tablewriter.NewWriter(os.Stdout)
	outputTable.SetAutoFormatHeaders(false)
	outputTable.SetHeader([]string{"ID", "Label", "Amount", "Tags"})

	fmt.Printf("%s:\n", title)
	_ = collection.Range(nil, func(t *domain.Transaction) error {
		outputTable.Append([]string{t.ID, t.Label, "£" + fmt.Sprint(math.Abs(float64(t.Amount)/100)), strings.Join(t.Tags, ", ")})
		return nil
	})
	outputTable.SetFooter([]string{"", "", "£" + fmt.Sprint(math.Abs(float64(collection.Sum())/100)), ""})
	outputTable.Render()
}
