#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Dapper variables helpers. These functions need the
# following variables:
#
#    DAPPER_VERSION   -  The dapper version for running, default is v0.4.2.

function ecm::dapper::install() {
  local version=${DAPPER_VERSION:-"v0.4.2"}
  curl -fL "https://github.com/rancher/dapper/releases/download/${version}/dapper-$(uname -s)-$(uname -m)" -o /tmp/dapper
  chmod +x /tmp/dapper && mv /tmp/dapper /usr/local/bin/dapper
}

function ecm::dapper::validate() {
  if [[ -n "$(command -v dapper)" ]]; then
    return 0
  fi

  ecm::log::info "installing dapper"
  if ecm::dapper::install; then
    ecm::log::info "dapper: $(dapper -v)"
    return 0
  fi
  ecm::log::error "no dapper available"
  return 1
}

function ecm::dapper::run() {
  if ! ecm::docker::validate; then
    ecm::log::fatal "docker hasn't been installed"
  fi
  if ! ecm::dapper::validate; then
    ecm::log::fatal "dapper hasn't been installed"
  fi

  dapper "$@"
}
