package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(conventionalCmd)
}

var conventionalCmd = &cobra.Command{
	Use:   "conventional",
	Short: "Cheat sheet for Conventional Commits standard",
	Long:  `Displays a quick reference guide for the Conventional Commits specification, including formatting, types, and examples.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cheatSheet)
	},
}

const cheatSheet = `Conventional Commits Cheat Sheet

The Conventional Commits specification is a lightweight convention on top of commit messages.
It provides an easy set of rules for creating an explicit commit history.

-------------------------------------------------------------------------------
1. COMMIT MESSAGE FORMAT
-------------------------------------------------------------------------------
<type>([optional scope]): <description>

[optional body]

[optional footer(s)]

-------------------------------------------------------------------------------
2. TYPES
-------------------------------------------------------------------------------
  feat     : A new feature for the user, not a new feature for build script.
             (Correlates with MINOR in Semantic Versioning).
  fix      : A bug fix for the user, not a fix to a build script.
             (Correlates with PATCH in Semantic Versioning).

  Other common types:
  docs     : Documentation only changes (e.g., README, comments).
  style    : Formatting, missing semi colons, etc; no production code change.
  refactor : Refactoring production code, eg. renaming a variable.
  perf     : A code change that improves performance.
  test     : Adding missing tests, refactoring tests; no production code change.
  chore    : Updating grunt tasks etc; no production code change.
  revert   : Reverting a previous commit.

-------------------------------------------------------------------------------
3. BREAKING CHANGES (MAJOR version bump)
-------------------------------------------------------------------------------
A commit that has a footer 'BREAKING CHANGE:', or appends a '!' after the 
type/scope, introduces a breaking API change.

-------------------------------------------------------------------------------
4. EXAMPLES
-------------------------------------------------------------------------------
Commit message with description and breaking change footer:
  feat: allow provided config object to extend other configs
  
  BREAKING CHANGE: 'extends' key in config file is now used for extending other config files

Commit message with ! to draw attention to breaking change:
  feat!: send an email to the customer when a product is shipped

Commit message with scope:
  feat(lang): add Polish language

Commit message for a fix:
  fix: prevent racing of requests
  
  Introduce a request id and a reference to latest request. Dismiss
  incoming responses other than from latest request.`
