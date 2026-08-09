package cli

import (
	"agent/internal/config"
	"agent/internal/model"
	"fmt"

	"github.com/nougght/monitoring-system/shared/go/util"
	"github.com/spf13/cobra"
)

func newEnrollCmd(setupConfigPath *string) *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "enroll",
		Short: "enroll monitoring agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if setupConfigPath == nil {
				return fmt.Errorf("enroll command didn't get receive setup config path")
			}
			var cfg config.SetupConfig
			if err := util.ReadYaml(*setupConfigPath, &cfg); err != nil {
				return fmt.Errorf("failed to read yaml setup config: %w", err)
			}
			if cfg.EnrollmentKey == model.EnrollmentKeyUsed {
				return fmt.Errorf(`agent is already enrolled, use 'agent run --setupconfig="/path/to/setup/config"`)
			}
			if err := enrollAgent(cfg); err != nil {
				return fmt.Errorf("enroll agent error: %w", err)
			}
			return nil
		},
	}
	return runCmd
}

func enrollAgent(setupConfig config.SetupConfig) error {

}
