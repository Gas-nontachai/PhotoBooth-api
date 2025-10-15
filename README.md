# Photobooth-api

This repository contains the Photobooth API written in Go.

CI and CD
----------

This repository includes two GitHub Actions workflows:

- `.github/workflows/ci.yml` — runs on push and pull requests to `main`. It runs:
	- go build
	- go test ./...
	- go vet ./...
	- staticcheck ./...

- `.github/workflows/cd.yml` — runs on push to `main`. It builds a Docker image and pushes it to GitHub Container Registry (GHCR) with tags `latest` and the commit SHA.

Authentication & secrets
-----------------------

The CD workflow uses the automatically-provided `GITHUB_TOKEN` to authenticate with GHCR. If you prefer to use a Personal Access Token, create a repository secret called `GHCR_TOKEN` with permissions `write:packages` and `read:packages` and update the workflow to use it.

Docker
------

A multi-stage `Dockerfile` is included at the repository root for building a small static container.

Build locally:

```bash
docker build -t ghcr.io/<owner>/<repo>:local .
```

Run locally:

```bash
docker run --rm -p 8080:8080 ghcr.io/<owner>/<repo>:local
```

Notes
-----

- The CI workflow installs `staticcheck`. It assumes `go` modules are configured (this repo contains `go.mod`).
- The Docker build uses `./cmd` as the build target. If your main package is in a different path, update the Dockerfile accordingly.
# Photobooth-api
