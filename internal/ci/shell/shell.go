package shellci

// import "errors"

// func initShell() error {
// 	return errors.New("initShell not implemented")
// }

// func getShellConfig(repo Repo) Config {
// 	if repo.HasShell {
// 		return Config{
// 			LanguageName: "Shell",
// 			Tasks: []Task{
// 				{
// 					InfoText: "Init",
// 					InitFunc: initShell,
// 				},
// 				{
// 					InfoText:    "Shellcheck",
// 					CommandArgs: []string{"bash", "-c", "shopt -s globstar && shellcheck **/*.sh"},
// 				},
// 			},
// 		}
// 	}

// 	return Config{}
// }
