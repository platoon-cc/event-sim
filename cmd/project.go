package cmd

import (
	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/platoon-cc/platoon-cli/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	projectCmd := &cobra.Command{
		Use: "project",
	}

	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(&cobra.Command{
		Use: "choose",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := "project"
			activeProject := settings.GetActive(key)

			platoon := client.New(settings.GetAuthToken())
			projects, err := platoon.GetProjectList()
			if err != nil {
				return err
			}

			list := ui.NewList("Choose Active Project")
			for _, t := range projects {
				list.AddItem(t.Id, t.Name, t.Id == activeProject)
			}

			return list.Run(
				func(i ui.ListItem) {
					settings.SetActive(key, i.Key)
				})
		},
	})
}
