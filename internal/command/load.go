package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/finance-planner/internal/application/service"
)

func Load(profileService service.Profile) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finance",
		Short: "Finance iss a quick and easy financial planner.",
		Long:  `A quick and easy financial planner for the month.`,
	}

	cmd.PersistentFlags().String("profile", "", "Profile to interact with")

	_ = cmd.MarkPersistentFlagRequired("profile")

	cmd.AddCommand(ListTransactions(profileService))
	cmd.AddCommand(AddTransaction(profileService))

	return cmd
}
