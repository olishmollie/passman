_passman_entries() {
    prefix="${PASSMAN_DIR:-$HOME/.passman/}"
    # prefix="${prefix%/}/"
    autoexpand=${1:-0}

    local IFS=$'\n'
    local items=( $(compgen -f $prefix$cur) )

    local firstitem=""
    local i=0

    for item in ${items[@]}; do
        [[ $item =~ /\.[^/]*$ ]] && continue

        if [[ ${#items[@]} -eq 1 && autoexpand -eq 1 ]]; then
            while [[ -d $item ]]; do
                local subitems=($(compgen -f "$item/"))
                local filtereditems=( )
                for item2 in "${subitems[@]}"; do
                    [[ $item2 =~ /\.[^/]*$ ]] && continue
                    filtereditems+=( "$item2" )
                done
                if [[ ${#filtereditems[@]} -eq 1 ]]; then
                    item="${filtereditems[0]}"
                else
                    break
                fi
            done
        fi

        [[ -d $item ]] && item="$item/"
        COMPREPLY+=("${item#$prefix}")
        if [[ $i -eq 0 ]]; then
            firstitem=$item
        fi
        let i+=1
    done

    if [[ $i -gt 1 || $i -eq 1 && -d $firstitem ]]; then
        compopt -o nospace
    fi
}

_passman_folders() {
    prefix="${PASSMAN_DIR:-$HOME/.passman/}"

    local IFS=$'\n'
    local items=( $(compgen -d $prefix$cur) )
    for item in ${items[@]}; do
        [[ $item == $prefix.* ]] && continue
        COMPREPLY+=("${item#$prefix}/")
        compopt -o nospace
    done
}

_passman() {
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD - 1]}"

    local root options commands 
    COMPREPLY=()
    root=$HOME/.passman
    commands="touch rm edit"
    options="-copy"

    if [[ $COMP_CWORD -gt 1 ]]; then
        case "${COMP_WORDS[1]}" in
            touch)
                _passman_folders
                return 0;;
            rm|-copy)
                _passman_entries 1
                return 0;;
        esac
    else
        COMPREPLY+=( $(compgen -W "${commands} ${options}" -- ${cur}) )
        _passman_entries 1
    fi

    return 0
}

complete -o filenames -F _passman passman
# _passman_entries 1

