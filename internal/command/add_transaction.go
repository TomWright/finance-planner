package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/service"
)

func AddTransaction(profileService service.Profile) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-transaction",
		Short: "Add a transaction to the profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName, _ := cmd.Flags().GetString("profile")
			label, _ := cmd.Flags().GetString("label")
			amount, _ := cmd.Flags().GetInt64("amount")
			tags, _ := cmd.Flags().GetStringArray("tags")

			profile, err := profileService.LoadOrCreateProfile(profileName)
			if err != nil {
				return err
			}

			// create the transaction.
			t := domain.NewTransaction().
				WithLabel(label).
				WithAmount(amount)
			if tags != nil {
				t = t.WithTags(tags...)
			}

			// validate the transaction.
			if err := t.Validate(); err != nil {
				return err
			}

			// add the transaction to the profile.
			profile.Transactions.Add(t)

			// save the profile.
			if err := profileService.SaveProfile(profile); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("label", "", "Transaction label")
	cmd.Flags().Int64("amount", 0, "Transaction amount")
	cmd.Flags().StringArray("tags", []string{}, "Tags to group the transaction")

	_ = cmd.MarkFlagRequired("label")
	_ = cmd.MarkFlagRequired("amount")

	return cmd
}
