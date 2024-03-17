package cmd

import (
	"os"
	"path/filepath"
	"time"

	// tea "github.com/charmbracelet/bubbletea"
	// "github.com/liamawhite/nt/pkg/task/models/router"
	// "github.com/liamawhite/nt/pkg/task/client"

	"github.com/liamawhite/nt/pkg/workspace"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use: "task",
    Aliases: []string{"tsk", "t"},
	RunE: func(cmd *cobra.Command, args []string) error {
        defer logFile.Close()

        dir, err := workspaceDir()
        if err != nil {
            return err
        }
        _, err = workspace.New(dir)
        if err != nil {
            return err
        }

        time.Sleep(10 * time.Second)
		// dir, err := persistenceDir()
		// if err != nil {
			// return err
		// }
		// tasksDir, err := ensureDir(filepath.Join(dir, "tasks"))
		// if err != nil {
			// return err
		// }
//
		// taskClient, err := task.NewClient(tasksDir)
		// if err != nil {
			// return err
		// }
//
		// model := router.NewModel(taskClient)
		// p := tea.NewProgram(model, tea.WithAltScreen())
		// if _, err := p.Run(); err != nil {
			// return err
		// }
		return nil
	},
}


func workspaceDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return ensureDir(filepath.Join(home, "notes"))
}

func ensureDir(dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}

func init() {
    rootCmd.AddCommand(taskCmd)
}

