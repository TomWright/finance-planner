package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/service"
	"github.com/tomwright/finance-planner/internal/errs"
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

			printTransactionsFn := func(t *domain.Transaction) error {
				tagStr := strings.Join(t.Tags, ", ")
				fmt.Printf("\t%s - %d - [%s]\n", t.Label, t.Amount, tagStr)
				return nil
			}

			incomingTransactions := profile.Transactions.Subset(func(t *domain.Transaction) bool {
				return t.Amount > 0
			})
			outgoingTransactions := profile.Transactions.Subset(func(t *domain.Transaction) bool {
				return t.Amount < 0
			})

			fmt.Printf("Profile: %s\n", profile.Name)
			fmt.Printf("Incoming Transactions (%d):\n", incomingTransactions.Sum())
			if err := incomingTransactions.Range(nil, printTransactionsFn); err != nil {
				return errs.FromErr(err)
			}
			fmt.Printf("Outgoing Transactions (%d):\n", outgoingTransactions.Sum())
			if err := outgoingTransactions.Range(nil, printTransactionsFn); err != nil {
				return errs.FromErr(err)
			}

			fmt.Printf("End balance: %d\n", profile.Transactions.Sum())

			return nil
		},
	}

	cmd.Flags().String("profile", "", "Profile to interact with")

	_ = cmd.MarkFlagRequired("profile")

	return cmd
}
