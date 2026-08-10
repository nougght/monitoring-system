package cli

import (
	"agent/internal/app"
	"agent/internal/config"
	"agent/internal/model"
	"fmt"

	"github.com/nougght/monitoring-system/shared/go/util"
	"github.com/spf13/cobra"
)

func newRunCmd(setupConfigPath *string) *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "run monitoring agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if setupConfigPath == nil {
				return fmt.Errorf("run command didn't get receive setup config path")
			}
			var cfg *config.SetupConfig
			if err := util.ReadYaml(*setupConfigPath, cfg); err != nil {
				return fmt.Errorf("failed to read yaml setup config: %w", err)
			}

			if cfg.EnrollmentKey != model.EnrollmentKeyUsed {
				return fmt.Errorf(`agent is not enrolled, use 'agent enroll --setupconfig="/path/to/setup/config"`)
			}

			if err := app.RunAgent(cfg); err != nil {
				return fmt.Errorf("run agent error: %w", err)
			}
			return nil
		},
	}
	return runCmd
}
