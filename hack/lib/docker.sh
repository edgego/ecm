#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Docker variables helpers. These functions need the
# following variables:
#
#    DOCKER_VERSION  -  The docker version for running, default is 19.03.

function ecm::docker::install() {
  local version=${DOCKER_VERSION:-"19.03"}
  curl -SfL "https://get.docker.com" | sh -s VERSION="${version}"
}

function ecm::docker::validate() {
  if [[ -n "$(command -v docker)" ]]; then
    return 0
  fi

  ecm::log::info "installing docker"
  if ecm::docker::install; then
    ecm::log::info "docker: $(docker version --format '{{.Server.Version}}' 2>&1)"
    return 0
  fi
  ecm::log::error "no docker available"
  return 1
}

function ecm::docker::login() {
  if [[ -n ${DOCKER_USERNAME} ]] && [[ -n ${DOCKER_PASSWORD} ]]; then
    if ! docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}" >/dev/null 2>&1; then
      return 1
    fi
  fi
  return 0
}

function ecm::docker::prebuild() {
  docker run --rm --privileged multiarch/qemu-user-static --reset -p yes i
  DOCKER_CLI_EXPERIMENTAL=enabled docker buildx create --name multi-builder
  DOCKER_CLI_EXPERIMENTAL=enabled docker buildx inspect multi-builder --bootstrap
  DOCKER_CLI_EXPERIMENTAL=enabled docker buildx use multi-builder
}

function ecm::docker::build() {
  if ! ecm::docker::validate; then
    ecm::log::fatal "docker hasn't been installed"
  fi
  # NB(thxCode): use Docker buildkit to cross build images, ref to:
  # - https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope#buildkit
  DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1 docker buildx build "$@"
}

function ecm::docker::manifest() {
  if ! ecm::docker::validate; then
    ecm::log::fatal "docker hasn't been installed"
  fi
  if ! ecm::docker::login; then
    ecm::log::fatal "failed to login docker"
  fi

  # NB(thxCode): use Docker manifest needs to enable client experimental feature, ref to:
  # - https://docs.docker.com/engine/reference/commandline/manifest_create/
  # - https://docs.docker.com/engine/reference/commandline/cli/#experimental-features#environment-variables
  ecm::log::info "docker manifest create --amend $*"
  DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create --amend "$@"

  # NB(thxCode): use Docker manifest needs to enable client experimental feature, ref to:
  # - https://docs.docker.com/engine/reference/commandline/manifest_push/
  # - https://docs.docker.com/engine/reference/commandline/cli/#experimental-features#environment-variables
  ecm::log::info "docker manifest push --purge ${1}"
  DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push --purge "${1}"
}

function ecm::docker::push() {
  if ! ecm::docker::validate; then
    ecm::log::fatal "docker hasn't been installed"
  fi
  if ! ecm::docker::login; then
    ecm::log::fatal "failed to login docker"
  fi

  for image in "$@"; do
    ecm::log::info "docker push ${image}"
    docker push "${image}"
  done
}
