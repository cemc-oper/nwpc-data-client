#!/bin/sh
# version.sh — the single implementation of version computation for `make` builds.
#
# Rules (VERSION file holds the NEXT release version, e.g. "0.11.0"):
#   checkout exactly on tag v$(cat VERSION)  ->  0.11.0            (release; -dirty if tree is dirty)
#   git available but not on that tag        ->  0.11.0-dev.N+g<sha>[-dirty]
#   no git (source tarball)                  ->  0.11.0-dev
#
# GoReleaser builds do NOT use this script: release versions come from the git
# tag (validated against VERSION by the release workflow), snapshot versions
# use .goreleaser.yaml's snapshot.version_template.
set -u

version_file=$(tr -d '[:space:]' < VERSION 2>/dev/null || echo "0.0.0")

if git rev-parse --git-dir >/dev/null 2>&1; then
    exact=$(git describe --tags --exact-match 2>/dev/null || true)
    sha=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    dirty=""
    if ! git diff --quiet 2>/dev/null || ! git diff --cached --quiet 2>/dev/null; then
        dirty="-dirty"
    fi
    if [ "$exact" = "v$version_file" ]; then
        echo "${version_file}${dirty}"
    else
        last_tag=$(git describe --tags --abbrev=0 2>/dev/null || true)
        if [ -n "$last_tag" ]; then
            n=$(git rev-list --count "${last_tag}..HEAD")
        else
            n=$(git rev-list --count HEAD)
        fi
        echo "${version_file}-dev.${n}+g${sha}${dirty}"
    fi
else
    echo "${version_file}-dev"
fi
