# Zsh configuration
# Managed by dotts - https://dotts.sh

# ══════════════════════════════════════════════════════════════════════════════
# Environment
# ══════════════════════════════════════════════════════════════════════════════

export EDITOR=nvim
export VISUAL=$EDITOR
export PAGER="less -R"

# XDG Base Directories
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_CACHE_HOME="$HOME/.cache"
export XDG_STATE_HOME="$HOME/.local/state"

# History
HISTFILE="$XDG_STATE_HOME/zsh/history"
HISTSIZE=50000
SAVEHIST=50000
setopt SHARE_HISTORY
setopt HIST_IGNORE_DUPS
setopt HIST_IGNORE_SPACE

# ══════════════════════════════════════════════════════════════════════════════
# Path
# ══════════════════════════════════════════════════════════════════════════════

typeset -U path
path=(
    $HOME/.local/bin
    $HOME/.cargo/bin
    $HOME/go/bin
    $path
)

# ══════════════════════════════════════════════════════════════════════════════
# Aliases
# ══════════════════════════════════════════════════════════════════════════════

# Navigation
alias ..="cd .."
alias ...="cd ../.."
alias ....="cd ../../.."

# Modern replacements
alias ls="eza --icons --group-directories-first"
alias ll="eza -la --icons --group-directories-first"
alias lt="eza -T --icons --level=2"
alias cat="bat --paging=never"

# Git shortcuts
alias g="git"
alias gs="git status -sb"
alias ga="git add"
alias gc="git commit"
alias gp="git push"
alias gl="git pull"
alias gd="git diff"
alias gco="git checkout"
alias gb="git branch"
alias lg="lazygit"

# Quick edits
alias vim="nvim"
alias v="nvim"
alias e="$EDITOR"

# Safety nets
alias rm="rm -i"
alias mv="mv -i"
alias cp="cp -i"

# ══════════════════════════════════════════════════════════════════════════════
# Integrations
# ══════════════════════════════════════════════════════════════════════════════

# Starship prompt
if command -v starship &> /dev/null; then
    eval "$(starship init zsh)"
fi

# Zoxide (smart cd)
if command -v zoxide &> /dev/null; then
    eval "$(zoxide init zsh)"
fi

# fzf
if command -v fzf &> /dev/null; then
    source <(fzf --zsh)
fi

# asdf version manager
if [[ -f "$HOME/.asdf/asdf.sh" ]]; then
    source "$HOME/.asdf/asdf.sh"
fi

# ══════════════════════════════════════════════════════════════════════════════
# Local overrides (not tracked by git)
# ══════════════════════════════════════════════════════════════════════════════

[[ -f "$XDG_CONFIG_HOME/zsh/local.zsh" ]] && source "$XDG_CONFIG_HOME/zsh/local.zsh"
