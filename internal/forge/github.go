package forge

import (
	"fmt"
	"os"
	"path/filepath"
)


func CreateGithubWorkflowDir() (string, error) {
	workflowDir := filepath.Join(".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create workflows directory: %w", err)
	}
	return workflowDir, nil
}

func CreateGithubPullRequestWorkflow(workflowDir string, protectedBranch string) error {
	path := filepath.Join(workflowDir, "genie-pr-check.yml")
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  Workflow file %s already exists. Skipping.\n", path)
	} else {
		script := fmt.Sprintf(githubPullRequestWorkflow, protectedBranch)
		if err := os.WriteFile(path, []byte(script), 0644); err != nil {
			return fmt.Errorf("failed to write workflow file: %w", err)
		}
		fmt.Printf("✨ PR Checker scaffolded at %s\n", path)
		fmt.Println("   Recipe: Go to repo Settings > Branches > Add protection rule for 'master'. Require status check 'check-pr-title'.")
	}

	return nil
}

func CreateGithubReleaseWorkflow(workflowDir string) error {
	path := filepath.Join(workflowDir, "genie-release.yml")
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  Workflow file %s already exists. Skipping.\n", path)
	} else {
		script := githubReleaseWorkflow
		if err := os.WriteFile(path, []byte(script), 0644); err != nil {
			return fmt.Errorf("failed to write workflow file: %w", err)
		}
		fmt.Printf("✨ Automated Releaser scaffolded at %s\n", path)
		fmt.Println("   Usage: Go to Actions tab > 'Genie Automated Release' > Click 'Run workflow'.")
	}

	return nil
}

const githubPullRequestWorkflow = `name: Genie PR Title Check

on:
  pull_request:
    branches:
      - %s
    types: [opened, edited, synchronize, reopened]

jobs:
  check-pr-title:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install Genie
        run: go install github.com/wilblik/genie/cmd/genie@latest

      - name: Validate PR Title
        run: genie check-msg "${{ github.event.pull_request.title }}"
`

const githubReleaseWorkflow = `name: Genie Automated Release

on:
  workflow_dispatch:

jobs:
  publish-release:
    runs-on: ubuntu-latest
    permissions:
      contents: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Setup Git
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install Genie
        run: go install github.com/wilblik/genie/cmd/genie@latest

      - name: Create Release Tag
        run: |
          # 1. Ask Genie to calculate version, generate notes, create tag, and push it
          genie release --push

          # 2. Get the tag name that was just created
          NEW_TAG=$(git describe --tags --abbrev=0)
          echo "NEW_TAG=$NEW_TAG" >> $GITHUB_ENV

      - name: Publish GitHub Release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # Extract the release notes from the Git tag's message body
          git tag -l --format='%(contents:body)' $NEW_TAG > release_notes.md

          # Create the official GitHub Release UI page
          gh release create $NEW_TAG \
            --title "Release $NEW_TAG" \
            --notes-file release_notes.md
`
