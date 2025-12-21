package tftools

import (
	"context"
	"fmt"

	oscarcfgpbv1 "github.com/opensourcecorp/oscar/internal/generated/opensourcecorp/oscar/config/v1"
	"github.com/opensourcecorp/oscar/internal/oscarcfg"
	iprint "github.com/opensourcecorp/oscar/internal/print"
	"github.com/opensourcecorp/oscar/internal/system"
	taskutil "github.com/opensourcecorp/oscar/internal/tasks/util"
)

type (
	tfFormat   struct{ taskutil.Tool }
	tfLint     struct{ taskutil.Tool }
	tfValidate struct{ taskutil.Tool }
)

// NewTasksForCI returns the list of CI tasks.
func NewTasksForCI(repo taskutil.Repo) []taskutil.Tasker {
	if repo.HasTerraform {
		return []taskutil.Tasker{
			tfFormat{
				Tool: taskutil.Tool{
					RunArgs: []string{"terraform", "fmt", "-recursive", "."},
				},
			},
			tfLint{
				Tool: taskutil.Tool{
					RunArgs: []string{"tflint", "--recursive", "."},
				},
			},
			// Logic is complex so is in its Exec() method below
			tfValidate{},
		}
	}

	return nil
}

// InfoText implements [taskutil.Tasker.InfoText].
func (t tfFormat) InfoText() string { return "Format" }

// Exec implements [taskutil.Tasker.Exec].
func (t tfFormat) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t tfFormat) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t tfLint) InfoText() string { return "Lint" }

// Exec implements [taskutil.Tasker.Exec].
func (t tfLint) Exec(ctx context.Context) error {
	if _, err := system.RunCommand(ctx, t.RunArgs); err != nil {
		return err
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t tfLint) Post(_ context.Context) error { return nil }

// InfoText implements [taskutil.Tasker.InfoText].
func (t tfValidate) InfoText() string { return "Validate" }

// Exec implements [taskutil.Tasker.Exec].
func (t tfValidate) Exec(ctx context.Context) error {
	cfg, err := oscarcfg.Get()
	if err != nil {
		return err
	}

	envs := getTFMapping(cfg)

	for envName, env := range envs {
		iprint.Debugf("Terraform env %s has %d module roots", envName, len(envs[envName].GetModuleRootDirs()))
		for _, root := range env.GetModuleRootDirs() {
			iprint.Debugf("Processing Terraform module root %s", root)
			args := []string{"bash", "-c", fmt.Sprintf(`
				terraform -chdir=%s init -backend=false
				terraform -chdir=%s validate`,
				root,
				root,
			)}

			if _, err := system.RunCommand(ctx, args); err != nil {
				return err
			}
		}
	}

	return nil
}

// Post implements [taskutil.Tasker.Post].
func (t tfValidate) Post(_ context.Context) error { return nil }

// TODO:
func getTFMapping(cfg *oscarcfgpbv1.Config) map[string]*oscarcfgpbv1.TerraformEnv {
	tf := cfg.GetDeployables().GetTerraform()
	return map[string]*oscarcfgpbv1.TerraformEnv{
		"Dev":     tf.GetDev(),
		"Nonprod": tf.GetNonprod(),
		"Prod":    tf.GetProd(),
	}
}
