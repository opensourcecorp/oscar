package nodejsci

import (
	ciutil "github.com/opensourcecorp/oscar/internal/ci/util"
)

var tasks = []ciutil.Tasker{}

// Tasks returns the list of CI tasks.
func Tasks(repo ciutil.Repo) []ciutil.Tasker {
	if repo.HasNodejs {
		return tasks
	}

	return nil
}
