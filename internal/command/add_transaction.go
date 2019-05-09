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

			profile, err := profileService.LoadOrCreateProfileByName(profileName)
			if err != nil {
				return err
			}

			// create the transaction.
			t := domain.NewTransaction()
			t.Label = label
			t.Amount = amount
			t.ProfileID = profile.ID
			t.Tags = tags

			// save the profile.
			if err := profileService.CreateTransaction(t); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("profile", "", "Profile to interact with")
	cmd.Flags().String("label", "", "Transaction label")
	cmd.Flags().Int64("amount", 0, "Transaction amount")
	cmd.Flags().StringArray("tags", []string{}, "Tags to group the transaction")

	_ = cmd.MarkFlagRequired("profile")
	_ = cmd.MarkFlagRequired("label")
	_ = cmd.MarkFlagRequired("amount")

	return cmd
}
