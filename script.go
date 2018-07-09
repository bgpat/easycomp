package easycomp

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/mattn/go-isatty"
)

const (
	scriptTemplate = `# for bash
if type complete &>/dev/null; then
  _{{ .Name }} () {
    local words cword
    if type _get_comp_words_by_ref &>/dev/null; then
      _get_comp_words_by_ref -n = -n @ -n : -w words -i cword
    else
      cword="$COMP_CWORD"
      words=("${COMP_WORDS[@]}")
    fi

    local si="$IFS"
    IFS=$'\n' COMPREPLY=($(COMP_CWORD="$cword" \
                           COMP_LINE="$COMP_LINE" \
                           COMP_POINT="$COMP_POINT" \
                           {{ .Command }} -- "${words[@]}" \
                           2>/dev/null)) || return $?
    IFS="$si"
    if type __ltrim_colon_completions &>/dev/null; then
      __ltrim_colon_completions "${words[cword]}"
    fi
  }
  complete -o default -F _{{ .Name }} {{ .Name }}

#for zsh
elif type compdef &>/dev/null; then
  _{{ .Name }} () {
    local si=$IFS
    compadd -- $(COMP_CWORD=$((CURRENT-1)) \
                 COMP_LINE=$BUFFER \
                 COMP_POINT=0 \
                 {{ .Command }} -- "${words[@]}" \
                 2>/dev/null)
    IFS=$si
  }
  compdef _{{ .Name }} {{ .Name }}
elif type compctl &>/dev/null; then
  _{{ .Name }} () {
    local cword line point words si
    read -Ac words
    read -cn cword
    let cword-=1
    read -l line
    read -ln point
    si="$IFS"
    IFS=$'\n' reply=($(COMP_CWORD="$cword" \
                       COMP_LINE="$line" \
                       COMP_POINT="$point" \
                       OFS=$'\n' \
					   {{ .Command }} -- "${words[@]}" \
                       2>/dev/null)) || return $?
    IFS="$si"
  }
  compctl -K _{{ .Name }} {{ .Name }}
fi
`
)

func dumpScript(name, cmd string) (string, error) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return "", fmt.Errorf("usage:\n\tsource <(%s)", strings.Join(os.Args, " "))
	}
	t := template.Must(template.New("script").Parse(scriptTemplate))
	obj := map[string]string{
		"Name":    name,
		"Command": cmd,
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, obj); err != nil {
		return "", err
	}
	return buf.String(), nil
}
