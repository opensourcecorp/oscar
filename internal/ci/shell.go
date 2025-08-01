package ci

import "errors"

func initShell() error {
	return errors.New("initShell not implemented")
}

func getShellConfigs(repo Repo) []Config {
	if repo.HasShell {
		return []Config{
			{
				LanguageName: "Shell",
				Tasks: []Task{
					{
						InfoText: "Init",
						InitFunc: initShell,
					},
					{
						InfoText:  "Shellcheck",
						RunScript: []string{"bash", "-c", "shopt -s globstar && shellcheck **/*.sh"},
					},
				},
			},
		}
	}

	return nil
}
