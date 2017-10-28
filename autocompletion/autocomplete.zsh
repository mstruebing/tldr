#compdef tldr

local -a pages oses
pages=$(tldr --list-all)
oses='( linux osx sunos )'

_arguments \
  '(- *)'{-h,--help}'[show help]' \
  '(- *)'{-v,--version}'[show version number]' \
  '(- *)'{-a,--list-all}'[list all commands]' \
  '(-r --render)'{-r,--render}'[render a specific markdown file]:markdown file:_files -/' \
  '(-p --platform)'{-p,--platform}"[override operating system]:os:${oses}" \
  '(- *)'{-u,--update}'[update local cache]' \
  "*:page:(${pages})" && return 0
