#!/bin/sh
set -e

DOTTS_REPO="arthur404dev/dotts"
DOTTS_VERSION="${DOTTS_VERSION:-latest}"
DOTTS_DIR="${DOTTS_DIR:-$HOME/.local/bin}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info() {
    printf "%b[INFO]%b %s\n" "$CYAN" "$NC" "$1"
}

success() {
    printf "%b[OK]%b %s\n" "$GREEN" "$NC" "$1"
}

error() {
    printf "%b[ERROR]%b %s\n" "$RED" "$NC" "$1" >&2
    exit 1
}

warn() {
    printf "%b[WARN]%b %s\n" "$YELLOW" "$NC" "$1"
}

check_cmd() {
    command -v "$1" >/dev/null 2>&1
}

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *)       error "Unsupported OS: $(uname -s)" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)             error "Unsupported architecture: $(uname -m)" ;;
    esac
}

get_latest_version() {
    curl -fsSL "https://api.github.com/repos/$DOTTS_REPO/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

main() {
    echo ""
    printf "%b" "$CYAN"
    echo "  ╔═══════════════════════════════════════╗"
    echo "  ║           dotts installer             ║"
    echo "  ║     Universal Dotfiles Manager        ║"
    echo "  ╚═══════════════════════════════════════╝"
    printf "%b" "$NC"
    echo ""

    if ! check_cmd curl; then
        error "curl is required but not installed"
    fi

    if ! check_cmd tar; then
        error "tar is required but not installed"
    fi

    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected: $OS/$ARCH"

    VERSION="$DOTTS_VERSION"
    if [ "$VERSION" = "latest" ]; then
        info "Fetching latest version..."
        VERSION=$(get_latest_version)
        if [ -z "$VERSION" ]; then
            error "Could not determine latest version"
        fi
    fi
    info "Version: $VERSION"

    VERSION_NUM="${VERSION#v}"
    FILENAME="dotts_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/$DOTTS_REPO/releases/download/$VERSION/$FILENAME"

    TMPDIR=$(mktemp -d)
    trap 'rm -rf "$TMPDIR"' EXIT

    info "Downloading..."
    if ! curl -fsSL "$URL" -o "$TMPDIR/dotts.tar.gz"; then
        error "Download failed: $URL"
    fi

    info "Extracting..."
    tar -xzf "$TMPDIR/dotts.tar.gz" -C "$TMPDIR"

    info "Installing to $DOTTS_DIR..."
    mkdir -p "$DOTTS_DIR"

    if [ -f "$TMPDIR/dotts" ]; then
        mv "$TMPDIR/dotts" "$DOTTS_DIR/dotts"
        chmod +x "$DOTTS_DIR/dotts"
    else
        error "Binary not found in archive"
    fi

    success "dotts $VERSION installed to $DOTTS_DIR/dotts"

    if [ -z "${DOTTS_NO_MODIFY_PATH:-}" ]; then
        case "$PATH" in
            *"$DOTTS_DIR"*) ;;
            *)
                SHELL_NAME=$(basename "${SHELL:-sh}")
                case "$SHELL_NAME" in
                    bash)
                        if [ -f "$HOME/.bashrc" ]; then
                            SHELL_RC="$HOME/.bashrc"
                        elif [ -f "$HOME/.bash_profile" ]; then
                            SHELL_RC="$HOME/.bash_profile"
                        fi
                        ;;
                    zsh)  SHELL_RC="$HOME/.zshrc" ;;
                    fish) SHELL_RC="$HOME/.config/fish/config.fish" ;;
                esac

                if [ -n "${SHELL_RC:-}" ] && ! grep -q "$DOTTS_DIR" "$SHELL_RC" 2>/dev/null; then
                    info "Adding $DOTTS_DIR to PATH in $SHELL_RC"
                    echo "" >> "$SHELL_RC"
                    if [ "$SHELL_NAME" = "fish" ]; then
                        echo "fish_add_path $DOTTS_DIR" >> "$SHELL_RC"
                    else
                        echo "export PATH=\"$DOTTS_DIR:\$PATH\"" >> "$SHELL_RC"
                    fi
                    warn "Restart your shell or run: source $SHELL_RC"
                fi
                ;;
        esac
    fi

    echo ""
    success "Installation complete!"
    echo ""
    echo "  Get started:"
    printf "    %bdotts init%b     - Bootstrap your dotfiles\n" "$CYAN" "$NC"
    printf "    %bdotts --help%b   - Show available commands\n" "$CYAN" "$NC"
    echo ""
    echo "  Documentation: https://dotts.4o4.sh"
    echo ""
}

main
