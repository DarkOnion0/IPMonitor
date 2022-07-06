#!/usr/bin/env just --justfile

# Just SETTINGS (vars...)
set dotenv-load

VERSION := "latest"
CONTAINER_BUILDER := "docker"
CONTAINER_NAME := "ipmonitor-dev"
export GH_TOKEN := ""
export GH_REPO := env_var_or_default("GH_REPO", "DarkOnion0/IPMonitor")

#Change the default just behaviour
default:
  @just --list

# Build app for all plateform
build: install
    @echo -e "\nBuild the app binary for Linux, Mac and Windows"
    ./build.sh {{VERSION}}

# Build app's container image
build_container: format lint
    @echo -e "\nBuild the app container for a single arch"
    {{CONTAINER_BUILDER}} build . -t {{CONTAINER_NAME}}

# Clean the remote GHCR container registry
cleanc:
    ./delete_remote_images.sh

# Clean the binary folder
cleanb:
    rm -rf ./bin

# Lint the project files
lint: install
    @echo -e "\nLint all go files"
    golangci-lint run --verbose --fix --timeout 5m .

# Format all the project files
format:
    @echo -e "\nFormat go code"
    gofmt -w .

    @echo -e "\nFormat other code with prettier (yaml, md...)"
    prettier -w .

# Check the go.mod and the go.sum files
check: install format lint
    @echo -e "\nVerify dependencies have expected content"
    go mod verify
    
    @echo -e "\nCheck if go.mod and go.sum are up to date"
    go mod tidy

# Build & release app, it needs GH_TOKEN to be overwritten and UNSTABLE set to unstable to publish a pre-release
release_full $UNSTABLE="stable": build generate_changelog
    #!/usr/bin/env bash
    if [ "${UNSTABLE}" = "unstable" ]; then
        gh release create --generate-notes --prerelease {{VERSION}} ./bin/*.zip
    else; then
        gh release create --generate-notes {{VERSION}} ./bin/*.zip
    fi

# Upload the release binary to an existing release, it needs GH_TOKEN to be overwritten
release_ci: build
    gh release upload {{VERSION}} ./bin/*.zip

# Generate the changelog
generate_changelog:
    @echo -e "\nGenerate the changelog with git-cliff"
    git-cliff --date-order --output CHANGELOG.md

# App dev command, binary mode
dev: check
    @echo -e "\nRun IPMonitor in dev mode (binary)"
    go run main.go -debug true -cron "*/1 * * * *"

# App dev command, container mode
dev_container: build_container
    @echo -e "\nRun IPMonitor in dev mode (container)"
    {{CONTAINER_BUILDER}} run -e DEBUG="true" -e CRON="*/1 * * * *" -p 8080:8080 {{CONTAINER_NAME}}:latest

# Run the prerequisites to install all the missing deps that nix can't cover
install:
    @echo -e "\nInstall the go prerequisites"
    go mod download
