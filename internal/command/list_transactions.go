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

			profile, err := profileService.LoadOrCreateProfile(profileName)
			if err != nil {
				return err
			}

			fmt.Printf("Profile: %s\nTransactions:\n", profile.Name)
			if err := profile.Transactions.Range(nil, func(t domain.Transaction) error {
				tagStr := strings.Join(t.Tags, ", ")
				fmt.Printf("\t%s - %d - [%s]\n", t.Label, t.Amount, tagStr)
				return nil
			}); err != nil {
				return errs.FromErr(err)
			}

			{
				sum, err := profile.Transactions.Sum()
				if err != nil {
					return errs.FromErr(err)
				}
				fmt.Printf("End balance: %d\n", sum)
			}

			return nil
		},
	}

	cmd.Flags().String("profile", "", "Profile to interact with")

	_ = cmd.MarkFlagRequired("profile")

	return cmd
}
