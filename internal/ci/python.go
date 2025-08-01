package ci

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/opensourcecorp/oscar/internal/consts"
)

func initPython() error {
	// Install uv
	var uvArch, uvOS, uvKernel string

	switch runtime.GOARCH {
	case "amd64":
		uvArch = "x86_64"
	case "arm64":
		uvArch = "aarch64"
	default:
		return fmt.Errorf("unsupported CPU architecture '%s'", runtime.GOARCH)
	}

	switch runtime.GOOS {
	case "darwin":
		uvOS = "apple"
		uvKernel = "darwin"
	case "linux":
		uvOS = "unknown"
		uvKernel = "linux-gnu"
	default:
		return fmt.Errorf("unsupported operating system '%s'", runtime.GOOS)
	}

	// This will also be the name of the directory once extracted from the archive
	uvPackagedName := fmt.Sprintf("uv-%s-%s-%s", uvArch, uvOS, uvKernel)
	uvReleaseURL := fmt.Sprintf(
		"https://github.com/astral-sh/uv/releases/download/%s/%s.tar.gz",
		consts.PythonCIVersions.UV, uvPackagedName,
	)
	uvDownloadedFile := filepath.Join(os.TempDir(), "uv.tar.gz")

	// NOTE: yes, I know, but this is WAY easier than doing a whole Go song & dance with downloading
	// & unpacking a targz archive. System deps are called out in the README, don't @ me.
	installUVCmd := []string{"bash", "-c",
		fmt.Sprintf(`
			curl -fsSL -o %s %s
			tar -C %s -xzf %s
			mv %s/%s/{uv,uvx} %s/
		`,
			uvDownloadedFile, uvReleaseURL,
			os.TempDir(), uvDownloadedFile,
			os.TempDir(), uvPackagedName, consts.OscarHomeBin,
		),
	}

	if err := runInitCommand(installUVCmd); err != nil {
		return err
	}

	return nil
}

func getPythonConfigs(repo Repo) []Config {
	if repo.HasPython {
		return []Config{
			{
				LanguageName: "Python",
				Tasks: []Task{
					{
						InfoText: "Init",
						InitFunc: initPython,
					},
					// TODO: uncomment once a test Python tree is in place
					// {
					// 	InfoText: "Build",
					// 	RunScript: []string{"uv", "build"},
					// 	},
					// },
					{
						InfoText:  "Linter (ruff)",
						RunScript: []string{"uvx", "ruff", "check", "--fix", "./src"},
					},
					{
						InfoText:  "Linter (pydoclint)",
						RunScript: []string{"uvx", "pydoclint", "./src"},
					},
					{
						InfoText:  "Linter (mypy)",
						RunScript: []string{"uvx", "mypy", "./src"},
					},
					{
						InfoText:  "Formatter (ruff)",
						RunScript: []string{"uvx", "ruff", "format", "./src"},
					},
				},
			},
		}
	}

	return nil
}
