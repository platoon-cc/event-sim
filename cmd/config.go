package cmd

import (
	"fmt"

	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/platoon-cc/platoon-cli/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Choose which team and project you wish to connect to",
		RunE: func(cmd *cobra.Command, args []string) error {
			activeTeam, _ := settings.GetActive("team")

			platoon, err := client.New()
			if err != nil {
				return err
			}

			// force the cache to be clear
			settings.ClearCache("team")
			settings.ClearCache("project")

			teams, err := platoon.GetTeamList()
			if err != nil {
				return err
			}

			list := ui.NewList("Select Active Team")
			for _, t := range teams {
				list.AddItem(t.Id, t.Name, t.Id == activeTeam)
			}

			err = list.Run(
				func(i ui.ListItem) {
					if i.Key != activeTeam {
						settings.ClearCache("project")
						settings.ClearActive("project")
						settings.SetActive("team", i.Key)
					}
				})
			if err != nil {
				return err
			}

			activeProject, _ := settings.GetActive("project")
			projects, err := platoon.GetProjectList()
			if err != nil {
				return err
			}
			list = ui.NewList("Select Active Project")
			for _, t := range projects {
				fmt.Printf("active: %s - this: %s\n", activeProject, t.Id)
				list.AddItem(t.Id, t.Name, t.Id == activeProject)
			}

			return list.Run(
				func(i ui.ListItem) {
					if i.Key != activeProject {
						settings.SetActive("project", i.Key)
					}
				})
		},
	}

	rootCmd.AddCommand(configCmd)
}
