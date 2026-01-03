# Fish shell configuration
# Managed by dotts - https://dotts.sh

# Disable greeting
set -g fish_greeting

# ══════════════════════════════════════════════════════════════════════════════
# Environment
# ══════════════════════════════════════════════════════════════════════════════

set -gx EDITOR nvim
set -gx VISUAL $EDITOR
set -gx PAGER "less -R"

# XDG Base Directories
set -gx XDG_CONFIG_HOME $HOME/.config
set -gx XDG_DATA_HOME $HOME/.local/share
set -gx XDG_CACHE_HOME $HOME/.cache
set -gx XDG_STATE_HOME $HOME/.local/state

# ══════════════════════════════════════════════════════════════════════════════
# Path
# ══════════════════════════════════════════════════════════════════════════════

fish_add_path $HOME/.local/bin
fish_add_path $HOME/.cargo/bin
fish_add_path $HOME/go/bin

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
if command -q starship
    starship init fish | source
end

# Zoxide (smart cd)
if command -q zoxide
    zoxide init fish | source
end

# fzf key bindings
if command -q fzf
    fzf --fish | source
end

# asdf version manager
if test -f $HOME/.asdf/asdf.fish
    source $HOME/.asdf/asdf.fish
end

# ══════════════════════════════════════════════════════════════════════════════
# Local overrides (not tracked by git)
# ══════════════════════════════════════════════════════════════════════════════

if test -f $XDG_CONFIG_HOME/fish/local.fish
    source $XDG_CONFIG_HOME/fish/local.fish
end
