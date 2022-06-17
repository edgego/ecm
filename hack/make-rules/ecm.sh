#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CURR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

# The root of the ecm directory
ROOT_DIR="${CURR_DIR}"

source "${ROOT_DIR}/hack/lib/init.sh"
source "${CURR_DIR}/hack/lib/constant.sh"

mkdir -p "${CURR_DIR}/bin"
mkdir -p "${CURR_DIR}/dist"

function mod() {
  [[ "${1:-}" != "only" ]]
  pushd "${ROOT_DIR}" >/dev/null || exist 1
  ecm::log::info "downloading dependencies for ecm..."

  if [[ "$(go env GO111MODULE)" == "off" ]]; then
    ecm::log::warn "go mod has been disabled by GO111MODULE=off"
  else
    ecm::log::info "tidying"
    go mod tidy
  fi

  ecm::log::info "...done"
  popd >/dev/null || return
}


function lint() {
  [[ "${1:-}" != "only" ]] && mod
  ecm::log::info "linting ecm..."

  local targets=(
    "${CURR_DIR}/cmd/..."
    "${CURR_DIR}/pkg/..."
    "${CURR_DIR}/test/..."
  )
  ecm::lint::lint "${targets[@]}"

  ecm::log::info "...done"
}

function cross_build() {
  [[ "${1:-}" != "only" ]] && lint
  ecm::log::info "building ecm(${GIT_VERSION},${GIT_COMMIT},${GIT_TREE_STATE},${BUILD_DATE})..."

  local version_flags="
    -X main.gitVersion=${GIT_VERSION}
    -X main.gitCommit=${GIT_COMMIT}
    -X main.buildDate=${BUILD_DATE}
    -X k8s.io/client-go/pkg/version.gitVersion=${GIT_VERSION}
    -X k8s.io/client-go/pkg/version.gitCommit=${GIT_COMMIT}
    -X k8s.io/client-go/pkg/version.gitTreeState=${GIT_TREE_STATE}
    -X k8s.io/client-go/pkg/version.buildDate=${BUILD_DATE}
    -X k8s.io/component-base/version.gitVersion=${GIT_VERSION}
    -X k8s.io/component-base/version.gitCommit=${GIT_COMMIT}
    -X k8s.io/component-base/version.gitTreeState=${GIT_TREE_STATE}
    -X k8s.io/component-base/version.buildDate=${BUILD_DATE}"
  local flags="
    -w -s"
  local ext_flags="
    -extldflags '-static'"

  ecm::log::info "crossed building"
  local platforms=("${SUPPORTED_PLATFORMS[@]}")

  for platform in "${platforms[@]}"; do
    ecm::log::info "building ${platform}"

    local os_arch
    IFS="/" read -r -a os_arch <<<"${platform}"

    local os=${os_arch[0]}
    local arch=${os_arch[1]}
    if [[ "$os" == "windows" ]]; then
        if [[ "$arch" == "amd64" ]]; then
            GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-posix CXX=x86_64-w64-mingw32-g++-posix go build \
              -ldflags "${version_flags} ${flags} ${ext_flags}" \
              -tags netgo \
              -o "${CURR_DIR}/bin/ecm_${os}_${arch}.exe" \
              "${CURR_DIR}/main.go"
            cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}.exe" "${CURR_DIR}/dist/ecm_${os}_${arch}.exe"
        else
            GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 CC=i686-w64-mingw32-gcc-posix CXX=i686-w64-mingw32-g++-posix go build \
              -ldflags "${version_flags} ${flags} ${ext_flags}" \
              -tags netgo \
              -o "${CURR_DIR}/bin/ecm_${os}_${arch}.exe" \
              "${CURR_DIR}/main.go"
            cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}.exe" "${CURR_DIR}/dist/ecm_${os}_${arch}.exe"
        fi
    elif [[ "$arch" == "arm" ]]; then
        GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 GOARM=7 CC=arm-linux-gnueabihf-gcc-5 CXX=arm-linux-gnueabihf-g++-5 CGO_CFLAGS="-march=armv7-a -fPIC" CGO_CXXFLAGS="-march=armv7-a -fPIC" go build \
          -ldflags "${version_flags} ${flags} ${ext_flags}" \
          -tags netgo \
          -o "${CURR_DIR}/bin/ecm_${os}_${arch}" \
          "${CURR_DIR}/main.go"
        cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}" "${CURR_DIR}/dist/ecm_${os}_${arch}"
    elif [[ "$arch" == "arm64" ]]; then
        GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc-5 CXX=aarch64-linux-gnu-g++-5 go build \
          -ldflags "${version_flags} ${flags} ${ext_flags}" \
          -tags netgo \
          -o "${CURR_DIR}/bin/ecm_${os}_${arch}" \
          "${CURR_DIR}/main.go"
        cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}" "${CURR_DIR}/dist/ecm_${os}_${arch}"
    elif [[ "$os" == "darwin" ]]; then
        MACOSX_DEPLOYMENT_TARGET=11.1.0 GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ HOST=x86_64-apple-darwin19.6.0 go build \
          -ldflags "${version_flags} ${flags}" \
          -tags netgo \
          -o "${CURR_DIR}/bin/ecm_${os}_${arch}" \
          "${CURR_DIR}/main.go"
        cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}" "${CURR_DIR}/dist/ecm_${os}_${arch}"
    else
        GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 go build \
          -ldflags "${version_flags} ${flags} ${ext_flags}" \
          -tags netgo \
          -o "${CURR_DIR}/bin/ecm_${os}_${arch}" \
          "${CURR_DIR}/main.go"
        cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}" "${CURR_DIR}/dist/ecm_${os}_${arch}"
    fi
  done

  ecm::log::info "...done"
}

function build() {
  [[ "${1:-}" != "only" ]] && lint
  #ecm::log::info "building ecm(${GIT_VERSION},${GIT_COMMIT},${GIT_TREE_STATE},${BUILD_DATE})..."

  local version_flags="
    -X main.gitVersion=${GIT_VERSION}
    -X main.gitCommit=${GIT_COMMIT}
    -X main.buildDate=${BUILD_DATE}
    -X k8s.io/client-go/pkg/version.gitVersion=${GIT_VERSION}
    -X k8s.io/client-go/pkg/version.gitCommit=${GIT_COMMIT}
    -X k8s.io/client-go/pkg/version.buildDate=${BUILD_DATE}
    -X k8s.io/component-base/version.gitVersion=${GIT_VERSION}
    -X k8s.io/component-base/version.gitCommit=${GIT_COMMIT}
    -X k8s.io/component-base/version.buildDate=${BUILD_DATE}"
  local flags="
    -w -s"
  local ext_flags="
    -extldflags '-static'"

  local os="${OS:-$(go env GOOS)}"
  local arch="${ARCH:-$(go env GOARCH)}"

  ecm::log::info "building ${os}/${arch}"
  if [[ "$os" == "windows" ]]; then
      GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 go build \
        -ldflags "${version_flags} ${flags} ${ext_flags}" \
        -tags netgo \
        -o "${CURR_DIR}/bin/ecm_${os}_${arch}.exe" \
        "${CURR_DIR}/main.go"
      cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}.exe" "${CURR_DIR}/dist/ecm_${os}_${arch}.exe"
  elif [[ "$os" == "darwin" ]]; then
      GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 go build \
        -ldflags "${version_flags} ${flags}" \
        -tags netgo \
        -o "${CURR_DIR}/bin/ecm_${os}_${arch}" \
        "${CURR_DIR}/main.go"
      cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}" "${CURR_DIR}/dist/ecm_${os}_${arch}"
  else
      GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 go build \
        -ldflags "${version_flags} ${flags} ${ext_flags}" \
        -tags netgo \
        -o "${CURR_DIR}/bin/ecm_${os}_${arch}" \
        "${CURR_DIR}/main.go"
      cp -f "${CURR_DIR}/bin/ecm_${os}_${arch}" "${CURR_DIR}/dist/ecm_${os}_${arch}"
  fi

  ecm::log::info "...done"
}

function package() {
  [[ "${1:-}" != "only" ]] && build
  ecm::log::info "packaging ecm..."

  REPO=${REPO:-cnrancher}
  TAG=${TAG:-${GIT_VERSION}}

  local os="$(go env GOOS)"
  if [[ "$os" == "windows" || "$os" == "darwin" ]]; then
    ecm::log::warn "package into Darwin/Windows OS image is unavailable, use OS=linux ARCH=amd64 env to containerize linux/amd64 image"
    return
  fi

  ARCH=${ARCH:-$(go env GOARCH)}
  SUFFIX="-linux-${ARCH}"
  IMAGE_NAME=${REPO}/ecm:${TAG}${SUFFIX}

  docker build --build-arg ARCH=${ARCH} -t ${IMAGE_NAME} .

  ecm::log::info "...done"
}

function deploy() {
  [[ "${1:-}" != "only" ]] && package
  ecm::log::info "deploying ecm..."

  local repo=${REPO:-cnrancher}
  local image_name=${IMAGE_NAME:-ecm}
  local tag=${TAG:-${GIT_VERSION}}

  local platforms
  if [[ "${CROSS:-false}" == "true" ]]; then
    ecm::log::info "crossed deploying"
    platforms=("${SUPPORTED_PLATFORMS[@]}")
  else
    local os="${OS:-$(go env GOOS)}"
    local arch="${ARCH:-$(go env GOARCH)}"
    platforms=("${os}/${arch}")
  fi
  local images=()
  for platform in "${platforms[@]}"; do
    if [[ "${platform}" =~ darwin/* || "${platform}" =~ windows/* || "${platform}" == "linux/arm" ]]; then
      ecm::log::warn "package into Darwin/Windows OS image is unavailable, please use CROSS=true env to containerize multiple arch images or use OS=linux ARCH=amd64 env to containerize linux/amd64 image"
    else
      images+=("${repo}/${image_name}:${tag}-${platform////-}")
    fi
  done

  local without_manifest=${WITHOUT_MANIFEST:-false}
  local ignore_missing=${IGNORE_MISSING:-false}

  # docker manifest
  if [[ "${without_manifest}" == "false" ]]; then
    if [[ "${ignore_missing}" == "false" ]]; then
      ecm::docker::manifest "${repo}/${image_name}:${tag}" "${images[@]}"
    else
      ecm::manifest_tool::push from-args \
        --ignore-missing \
        --target="${repo}/${image_name}:${tag}" \
        --template="${repo}/${image_name}:${tag}-OS-ARCH" \
        --platforms="$(ecm::util::join_array "," "${platforms[@]}")"
    fi
  else
    ecm::log::warn "deploying manifest images has been stopped by WITHOUT_MANIFEST"
  fi

  ecm::log::info "...done"
}

function unit() {
  [[ "${1:-}" != "only" ]] && build
  ecm::log::info "running unit tests for ecm..."

  local unit_test_targets=(
    "${CURR_DIR}/cmd/..."
    "${CURR_DIR}/pkg/..."
    "${CURR_DIR}/test/..."
  )

  if [[ "${CROSS:-false}" == "true" ]]; then
    ecm::log::warn "crossed test is not supported"
  fi

  local os="${OS:-$(go env GOOS)}"
  local arch="${ARCH:-$(go env GOARCH)}"
  if [[ "${arch}" == "arm" ]]; then
    # NB(thxCode): race detector doesn't support `arm` arch, ref to:
    # - https://golang.org/doc/articles/race_detector.html#Supported_Systems
    GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 go test \
      -tags=test \
      -cover -coverprofile "${CURR_DIR}/dist/coverage_${os}_${arch}.out" \
      "${unit_test_targets[@]}"
  else
    GOOS=${os} GOARCH=${arch} CGO_ENABLED=1 go test \
      -tags=test \
      -race \
      -cover -coverprofile "${CURR_DIR}/dist/coverage_${os}_${arch}.out" \
      "${unit_test_targets[@]}"
  fi

  ecm::log::info "...done"
}

function verify() {
  [[ "${1:-}" != "only" ]] && unit
  ecm::log::info "running integration tests for ecm..."

  ecm::ginkgo::test "${CURR_DIR}/test/integration"

  ecm::log::info "...done"
}

function e2e() {
  [[ "${1:-}" != "only" ]] && verify
  ecm::log::info "running E2E tests for ecm..."

  # execute the E2E testing as ordered.
  #ecm::ginkgo::test "${CURR_DIR}/test/e2e/installation"
  #ecm::ginkgo::test "${CURR_DIR}/test/e2e/usability"

  ecm::log::info "...done"
}

function entry() {
  local stages="${1:-build}"
  shift $(($# > 0 ? 1 : 0))

  IFS="," read -r -a stages <<<"${stages}"
  local commands=$*
  if [[ ${#stages[@]} -ne 1 ]]; then
    commands="only"
  fi

  for stage in "${stages[@]}"; do
    ecm::log::info "# make ecm ${stage} ${commands}"
    case ${stage} in
    m | mod) mod "${commands}" ;;
    l | lint) lint "${commands}" ;;
    b | build) build "${commands}" ;;
    p | pkg | package) package "${commands}" ;;
    d | dep | deploy) deploy "${commands}" ;;
    u | unit) unit "${commands}" ;;
    v | ver | verify) verify "${commands}" ;;
    e | e2e) e2e "${commands}" ;;
    cb | cross_build) cross_build "${commands}" ;;
    *) ecm::log::fatal "unknown action '${stage}', select from mod,lint,build,unit,verify,package,deploy,e2e" ;;
    esac
  done
}

if [[ ${BY:-} == "dapper" ]]; then
  ecm::dapper::run -C "${CURR_DIR}" -f ${DAPPER_FILE:-Dockerfile.dapper} "$@"
else
  entry "$@"
fi
