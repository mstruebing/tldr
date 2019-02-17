#!/bin/bash

function _tldr_autocomplete {
  pages=$(tldr --list-all)
  COMPREPLY=()
  if [ "$COMP_CWORD" = 1 ]; then
    COMPREPLY=($(compgen -W "$pages" -- $2))
  fi
}

complete -F _tldr_autocomplete tldr
