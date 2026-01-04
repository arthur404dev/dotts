#!/usr/bin/env bash
# dotts installer script
# Usage: curl -fsSL https://dotts.4o4.sh/install.sh | bash
#
# Environment variables:
#   DOTTS_VERSION   - Specific version to install (default: latest)
#   DOTTS_DIR       - Installation directory (default: ~/.local/bin)
#   DOTTS_NO_MODIFY_PATH - If set, don't modify shell config to add to PATH

set -euo pipefail

# ══════════════════════════════════════════════════════════════════════════════
# Configuration
# ══════════════════════════════════════════════════════════════════════════════

DOTTS_REPO="arthur404dev/dotts"
DOTTS_VERSION="${DOTTS_VERSION:-latest}"
DOTTS_DIR="${DOTTS_DIR:-$HOME/.local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# ══════════════════════════════════════════════════════════════════════════════
# Helper functions
# ══════════════════════════════════════════════════════════════════════════════

info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}[OK]${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
}

fatal() {
    error "$1"
    exit 1
}

# ══════════════════════════════════════════════════════════════════════════════
# System detection
# ══════════════════════════════════════════════════════════════════════════════

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)       fatal "Unsupported operating system: $(uname -s)" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        armv7l)        echo "arm" ;;
        *)             fatal "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Check for required commands
check_requirements() {
    local missing=()
    
    for cmd in curl tar; do
        if ! command -v "$cmd" &> /dev/null; then
            missing+=("$cmd")
        fi
    done
    
    if [[ ${#missing[@]} -gt 0 ]]; then
        fatal "Missing required commands: ${missing[*]}"
    fi
}

# ══════════════════════════════════════════════════════════════════════════════
# Installation
# ══════════════════════════════════════════════════════════════════════════════

get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${DOTTS_REPO}/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

download_and_install() {
    local os="$1"
    local arch="$2"
    local version="$3"
    
    # Construct download URL
    local filename="dotts_${version#v}_${os}_${arch}.tar.gz"
    local url="https://github.com/${DOTTS_REPO}/releases/download/${version}/${filename}"
    
    info "Downloading dotts ${version} for ${os}/${arch}..."
    
    # Create temp directory
    local tmpdir
    tmpdir=$(mktemp -d)
    trap "rm -rf ${tmpdir}" EXIT
    
    # Download
    if ! curl -fsSL "$url" -o "${tmpdir}/dotts.tar.gz"; then
        fatal "Failed to download from ${url}"
    fi
    
    # Extract
    info "Extracting..."
    tar -xzf "${tmpdir}/dotts.tar.gz" -C "${tmpdir}"
    
    # Install
    info "Installing to ${DOTTS_DIR}..."
    mkdir -p "${DOTTS_DIR}"
    
    if [[ -f "${tmpdir}/dotts" ]]; then
        mv "${tmpdir}/dotts" "${DOTTS_DIR}/dotts"
        chmod +x "${DOTTS_DIR}/dotts"
    else
        fatal "Binary not found in archive"
    fi
    
    success "dotts ${version} installed to ${DOTTS_DIR}/dotts"
}

# ══════════════════════════════════════════════════════════════════════════════
# PATH configuration
# ══════════════════════════════════════════════════════════════════════════════

add_to_path() {
    if [[ -n "${DOTTS_NO_MODIFY_PATH:-}" ]]; then
        return
    fi
    
    # Check if already in PATH
    if [[ ":$PATH:" == *":${DOTTS_DIR}:"* ]]; then
        return
    fi
    
    local shell_config=""
    local export_line="export PATH=\"${DOTTS_DIR}:\$PATH\""
    
    # Detect shell and config file
    case "${SHELL:-}" in
        */bash)
            if [[ -f "$HOME/.bashrc" ]]; then
                shell_config="$HOME/.bashrc"
            elif [[ -f "$HOME/.bash_profile" ]]; then
                shell_config="$HOME/.bash_profile"
            fi
            ;;
        */zsh)
            shell_config="$HOME/.zshrc"
            ;;
        */fish)
            shell_config="$HOME/.config/fish/config.fish"
            export_line="fish_add_path ${DOTTS_DIR}"
            ;;
    esac
    
    if [[ -n "$shell_config" ]]; then
        if ! grep -q "${DOTTS_DIR}" "$shell_config" 2>/dev/null; then
            info "Adding ${DOTTS_DIR} to PATH in ${shell_config}"
            echo "" >> "$shell_config"
            echo "# Added by dotts installer" >> "$shell_config"
            echo "$export_line" >> "$shell_config"
            warn "Restart your shell or run: source ${shell_config}"
        fi
    else
        warn "Could not detect shell config. Add ${DOTTS_DIR} to your PATH manually."
    fi
}

# ══════════════════════════════════════════════════════════════════════════════
# Main
# ══════════════════════════════════════════════════════════════════════════════

main() {
    echo ""
    printf "${PURPLE}${BOLD}"
    echo "  ╔═══════════════════════════════════════╗"
    echo "  ║           dotts installer             ║"
    echo "  ║   Universal Dotfiles Manager          ║"
    echo "  ╚═══════════════════════════════════════╝"
    printf "${NC}"
    echo ""
    
    # Check requirements
    check_requirements
    
    # Detect system
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)
    
    info "Detected: ${os}/${arch}"
    
    # Get version
    local version="$DOTTS_VERSION"
    if [[ "$version" == "latest" ]]; then
        info "Fetching latest version..."
        version=$(get_latest_version)
        if [[ -z "$version" ]]; then
            fatal "Could not determine latest version"
        fi
    fi
    
    info "Version: ${version}"
    
    # Download and install
    download_and_install "$os" "$arch" "$version"
    
    # Add to PATH
    add_to_path
    
    echo ""
    success "Installation complete!"
    echo ""
    echo "  Get started:"
    printf "    ${CYAN}dotts init${NC}     # Bootstrap your dotfiles\n"
    printf "    ${CYAN}dotts --help${NC}   # Show available commands\n"
    echo ""
    echo "  Documentation: https://dotts.4o4.sh"
    echo ""
}

# Run if not sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
