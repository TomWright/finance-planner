package command

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/service"
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

			in, _ := cmd.Flags().GetBool("in")
			out, _ := cmd.Flags().GetBool("out")
			if in && out {
				in = false
				out = false
			}

			transactions := profile.Transactions
			title := "All transactions"
			if in {
				title = "Incoming transactions"
				transactions = profile.Transactions.Subset(func(t *domain.Transaction) bool {
					return t.Amount > 0
				})
			}
			if out {
				title = "Outgoing transactions"
				transactions = profile.Transactions.Subset(func(t *domain.Transaction) bool {
					return t.Amount < 0
				})
			}

			outputTransactions(title, transactions)

			return nil
		},
	}

	cmd.Flags().String("profile", "", "Profile to interact with")
	cmd.Flags().Bool("in", false, "Only list incoming transactions")
	cmd.Flags().Bool("out", false, "Only list outgoing transactions")

	_ = cmd.MarkFlagRequired("profile")

	return cmd
}

func outputTransactions(title string, collection *domain.TransactionCollection) {
	outputTable := tablewriter.NewWriter(os.Stdout)
	outputTable.SetAutoFormatHeaders(false)
	outputTable.SetHeader([]string{"ID", "Label", "Tags", "Amount"})
	outputTable.SetAutoWrapText(false)
	outputTable.SetCaption(true, title)

	_ = collection.Range(nil, func(t *domain.Transaction) error {
		outputTable.Append([]string{t.ID, t.Label, strings.Join(t.Tags, ", "), "£" + fmt.Sprint(float64(t.Amount)/100)})
		return nil
	})
	outputTable.SetFooter([]string{"", "", "Total", "£" + fmt.Sprint(float64(collection.Sum())/100)})
	outputTable.Render()
}
