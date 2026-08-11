package cli

import (
	"agent/internal/app"
	"agent/internal/config"
	"agent/internal/service/agentcore"
	"fmt"

	"github.com/nougght/monitoring-system/shared/go/cert_store"
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

			cfg, err := config.LoadSetupConfig(*setupConfigPath)
			if err != nil {
				return fmt.Errorf("failed to load setup config: %w", err)
			}
			if err := agentcore.EnrollAgent(cmd.Context(), cfg,
				cert_store.NewCertStore(cfg.CertPath, cfg.KeyPath, cfg.CaPath)); err != nil {
				return fmt.Errorf("enroll agent error: %w", err)
			}

			// default run after enrollment
			if err := app.RunAgent(cfg); err != nil {
				return fmt.Errorf("run agent error: %w", err)
			}
			return nil
		},
	}
	return runCmd
}
