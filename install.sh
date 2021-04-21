#!/bin/bash
# I stole a lot of this from people smarter than me: https://github.com/rancher/k3d/blob/main/install.sh

APP_NAME="env-replace"
REPO_URL="https://github.com/mzrimsek/env-replace"

: ${USE_SUDO:="true"}
: ${ENV_REPLACE_INSTALL_DIR:="/usr/local/bin"}

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    aarch64) ARCH="arm64";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(uname|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='Windows';;
    darwin*) OS='Darwin';;
    linux*) OS='Linux';;
  esac
}

# runs the given command as root (detects if we are root already)
runAsRoot() {
  local CMD="$*"

  if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
    CMD="sudo $CMD"
  fi

  $CMD
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  local supported="Darwin_arm64\nDarwin_x86_64\nLinux_arm64\nLinux_i386\nLinux_x86_64\nWindows_i386\nWindows_x86_64"
  if ! echo "${supported}" | grep -q "${OS}_${ARCH}"; then
    echo "No prebuilt binary for ${OS}_${ARCH}."
    echo "To build from source, go to $REPO_URL"
    exit 1
  fi

  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# checkLatestVersion grabs the latest version string from the releases
checkLatestVersion() {
  local latest_release_url="$REPO_URL/releases/latest"
  if type "curl" > /dev/null; then
    TAG=$(curl -Ls -o /dev/null -w %{url_effective} $latest_release_url | grep -oE "[^/]+$" )
  elif type "wget" > /dev/null; then
    TAG=$(wget $latest_release_url --server-response -O /dev/null 2>&1 | awk '/^\s*Location: /{DEST=$2} END{ print DEST}' | grep -oE "[^/]+$")
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  VERSION=$(echo ${TAG} | sed 's/v//')
  ENV_REPLACE_DIST="${APP_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
  DOWNLOAD_URL="$REPO_URL/releases/download/$TAG/$ENV_REPLACE_DIST"
  ENV_REPLACE_TMP_ROOT="$(mktemp -dt env-replace-binary-XXXXXX)"
  ENV_REPLACE_TMP_TAR_FILE="$ENV_REPLACE_TMP_ROOT/$ENV_REPLACE_DIST"
  if type "curl" > /dev/null; then
    curl -SsL "$DOWNLOAD_URL" -o "$ENV_REPLACE_TMP_TAR_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$ENV_REPLACE_TMP_TAR_FILE" "$DOWNLOAD_URL"
  fi
}

unpackTar() {
  tar xzf $ENV_REPLACE_TMP_TAR_FILE -C $ENV_REPLACE_TMP_ROOT $APP_NAME
  ENV_REPLACE_TMP_FILE="$ENV_REPLACE_TMP_ROOT/$APP_NAME"
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  echo "Preparing to install $APP_NAME into ${ENV_REPLACE_INSTALL_DIR}"
  runAsRoot chmod +x "$ENV_REPLACE_TMP_FILE"
  runAsRoot cp "$ENV_REPLACE_TMP_FILE" "$ENV_REPLACE_INSTALL_DIR/$APP_NAME"
  echo "$APP_NAME installed into $ENV_REPLACE_INSTALL_DIR/$APP_NAME"
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    if [[ -n "$INPUT_ARGUMENTS" ]]; then
      echo "Failed to install $APP_NAME with the arguments provided: $INPUT_ARGUMENTS"
      help
    else
      echo "Failed to install $APP_NAME"
    fi
    echo -e "\tFor support, go to $REPO_URL."
  fi
  cleanup
  exit $result
}

# cleanup temporary files
cleanup() {
  if [[ -d "${ENV_REPLACE_TMP_ROOT:-}" ]]; then
    rm -rf "$ENV_REPLACE_TMP_ROOT"
  fi
}

#Stop execution on any error
trap "fail_trap" EXIT
set -e

# Parsing input arguments (if any)
export INPUT_ARGUMENTS="${@}"
set -u
while [[ $# -gt 0 ]]; do
  case $1 in
    '--no-sudo')
       USE_SUDO="false"
       ;;
    *) exit 1
       ;;
  esac
  shift
done
set +u

initArch
initOS
verifySupported
checkLatestVersion
downloadFile
unpackTar
installFile
cleanup