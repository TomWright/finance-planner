package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/finance-planner/internal/application/service"
	"github.com/tomwright/finance-planner/internal/errs"
)

func UpdateTransaction(profileService service.Profile) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-transaction",
		Short: "Update a transaction in the profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName, _ := cmd.Flags().GetString("profile")
			id, _ := cmd.Flags().GetString("id")
			label, _ := cmd.Flags().GetString("label")
			amount, _ := cmd.Flags().GetInt64("amount")
			tags, _ := cmd.Flags().GetStringArray("tags")

			profile, err := profileService.LoadProfileByName(profileName)
			if err != nil {
				return err
			}

			t, err := profileService.LoadTransactionByID(id)
			if err != nil {
				return err
			}

			if t.ProfileID != profile.ID {
				return errs.New().
					WithCode(errs.ErrUnknownTransaction).
					WithMessage("unknown transaction")
			}

			if label != "" {
				t.Label = label
			}
			if amount != 0 {
				t.Amount = amount
			}
			if len(tags) > 0 {
				t.Tags = tags
			}

			// save the transaction.
			if err := profileService.UpdateTransaction(t); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("profile", "", "Profile to interact with")
	cmd.Flags().String("id", "", "Transaction ID")
	cmd.Flags().String("label", "", "Transaction label")
	cmd.Flags().Int64("amount", 0, "Transaction amount")
	cmd.Flags().StringArray("tags", nil, "Tags to group the transaction")

	_ = cmd.MarkFlagRequired("profile")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
