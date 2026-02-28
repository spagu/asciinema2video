#!/bin/sh
# asciinema2video install script
# Usage: curl -sSL https://raw.githubusercontent.com/spagu/asciinema2video/main/install.sh | sh

set -e

REPO="spagu/asciinema2video"
BINARY_NAME="asciinema2video"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
    exit 1
}

# Detect OS
detect_os() {
    OS="$(uname -s)"
    case "${OS}" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        FreeBSD*)   OS=freebsd;;
        *)          error "Unsupported OS: ${OS}";;
    esac
    echo "${OS}"
}

# Detect architecture
detect_arch() {
    ARCH="$(uname -m)"
    case "${ARCH}" in
        x86_64|amd64)   ARCH=amd64;;
        aarch64|arm64)  ARCH=arm64;;
        *)              error "Unsupported architecture: ${ARCH}";;
    esac
    echo "${ARCH}"
}

# Get latest version from GitHub
get_latest_version() {
    if command -v curl > /dev/null 2>&1; then
        curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | \
            grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget > /dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | \
            grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "curl or wget required"
    fi
}

# Download and install
install() {
    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION=$(get_latest_version)

    if [ -z "${VERSION}" ]; then
        error "Failed to get latest version"
    fi

    info "Installing ${BINARY_NAME} ${VERSION} for ${OS}/${ARCH}"

    # Construct download URL
    FILENAME="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"

    URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

    info "Downloading from ${URL}"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf ${TMP_DIR}" EXIT

    # Download
    if command -v curl > /dev/null 2>&1; then
        curl -sSL "${URL}" -o "${TMP_DIR}/${FILENAME}" || error "Download failed"
    elif command -v wget > /dev/null 2>&1; then
        wget -q "${URL}" -O "${TMP_DIR}/${FILENAME}" || error "Download failed"
    fi

    # Extract
    cd "${TMP_DIR}"
    tar xzf "${FILENAME}" || error "Extraction failed"

    # Install binary
    if [ -w "${INSTALL_DIR}" ]; then
        cp "${BINARY_NAME}" "${INSTALL_DIR}/" || error "Installation failed"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        info "Installing to ${INSTALL_DIR} requires root privileges"
        sudo cp "${BINARY_NAME}" "${INSTALL_DIR}/" || error "Installation failed"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Install man page if available
    if [ -f "${BINARY_NAME}.1" ]; then
        MAN_DIR="/usr/local/share/man/man1"
        if [ -w "${MAN_DIR}" ] 2>/dev/null || [ -w "$(dirname ${MAN_DIR})" ] 2>/dev/null; then
            mkdir -p "${MAN_DIR}" 2>/dev/null || sudo mkdir -p "${MAN_DIR}"
            cp "${BINARY_NAME}.1" "${MAN_DIR}/" 2>/dev/null || sudo cp "${BINARY_NAME}.1" "${MAN_DIR}/"
            info "Man page installed"
        fi
    fi

    info "Successfully installed ${BINARY_NAME} ${VERSION}"
    info "Run '${BINARY_NAME} --help' to get started"

    # Check for ffmpeg
    if ! command -v ffmpeg > /dev/null 2>&1; then
        warn "ffmpeg is not installed. Install it to use ${BINARY_NAME}:"
        case "${OS}" in
            linux)
                warn "  Ubuntu/Debian: sudo apt install ffmpeg"
                warn "  Fedora: sudo dnf install ffmpeg"
                warn "  Arch: sudo pacman -S ffmpeg"
                ;;
            darwin)
                warn "  Homebrew: brew install ffmpeg"
                ;;
            freebsd)
                warn "  pkg install ffmpeg"
                ;;
        esac
    fi
}

# Run installation
install
