#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Lint variables helpers. These functions need the
# following variables:
#
#    GOLANGCI_LINT_VERSION  -  The golangci-lint version, default is v1.31.0.
#    DIRTY_CHECK            -  Specify to check the git tree is dirty or not.

function ecm::lint::install() {
  local version=${GOLANGCI_LINT_VERSION:-"v1.38.0"}
  curl -fL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" "${version}"
}

function ecm::lint::validate() {
  if [[ -n "$(command -v golangci-lint)" ]]; then
    return 0
  fi

  ecm::log::info "installing golangci-lint"
  if ecm::lint::install; then
    ecm::log::info "$(golangci-lint --version)"
    return 0
  fi
  return 1
}

function ecm::lint::lint() {
  if [[ "${DIRTY_CHECK:-}" == "true" ]]; then
    if [[ "${GIT_TREE_STATE}" == "dirty" ]]; then
      ecm::log::fatal "the git tree is dirty:\n$(git status --porcelain)"
    fi
  fi

  if ecm::lint::validate; then
    for path in "$@"; do
      golangci-lint run "${path}"
    done
  else
    ecm::log::warn "no golangci-lint available, using go fmt/vet instead"
    for path in "$@"; do
      go fmt "${path}"
      go vet -tags=test "${path}"
    done
  fi
}
