package main

const zshCompletion = `#compdef asciizer

# Zsh completion script for asciizer
# Install: cp _asciizer ~/.oh-my-zsh/custom/completions/
#   — or — asciizer -completion zsh > /path/to/completions/_asciizer

_arguments -s \
  '(-h -help)'{-h,-help}'[Show help message]' \
  '-w[Output width in characters (0 = no resize)]:width:' \
  '-o[Output file path]:output file:_files' \
  '-stdout[Print to stdout instead of file]' \
  '-invert[Reverse brightness mapping]' \
  '-color[ANSI 256-color output]' \
  '-full-ramp[Use 70-char gradient instead of 10-char]' \
  '-version[Print version and exit]' \
  '-completion[Print shell completion script]:shell:(zsh bash)' \
  '*:image file:_files -g "*.(jpg|jpeg|png|gif)"'
`

const bashCompletion = `# Bash completion script for asciizer
# Install: asciizer -completion bash > /etc/bash_completion.d/asciizer
#   — or — asciizer -completion bash >> ~/.bashrc

_asciizer() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="-w -o -stdout -invert -color -full-ramp -version -completion -h -help"

    case "$prev" in
        -o)
            COMPREPLY=( $(compgen -f -- "$cur") )
            return 0
            ;;
        -w)
            return 0
            ;;
        -completion)
            COMPREPLY=( $(compgen -W "zsh bash" -- "$cur") )
            return 0
            ;;
    esac

    if [[ "$cur" == -* ]]; then
        COMPREPLY=( $(compgen -W "$opts" -- "$cur") )
        return 0
    fi

    COMPREPLY=( $(compgen -f -X '!*.@(jpg|jpeg|png|gif)' -- "$cur") )
}

complete -F _asciizer asciizer
`
