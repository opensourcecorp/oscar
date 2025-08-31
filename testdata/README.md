# Test data & files

In this subdirectory tree, you will find various files used during more comprehensive testing. There
are a few things to keep in mind when adding files to this tree:

* Since `oscar` runs its own CI against itself, and since its file discovery is usually based on
  file extensions, you must make sure that any files added here are suffixed with `.test`. For
  example, if you want to add a new file for testing CI against Python files, a new file should be
  named something like `main.py.test`.

* Any new files added here should ensure that they are referenced in `/scripts/test-bootstrap.sh`,
  for both setup and teardown operations. It is in this script that you should handle renames to
  remove the `.test` file extensions, so that `oscar` can pick them up. Be sure to add any needed
  code to also remove those files during teardown.
