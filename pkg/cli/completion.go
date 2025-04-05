package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

type completionCommand struct {
	logE   *logrus.Entry
	stdout io.Writer
}

func (cc *completionCommand) command() *cli.Command {
	// https://cli.urfave.org/v2/#bash-completion
	return &cli.Command{
		Name:  "completion",
		Usage: "Output shell completion script for bash, zsh, or fish",
		Description: `Output shell completion script for bash, zsh, or fish.
Source the output to enable completion.

e.g.

.bash_profile

source <(ghatm completion bash)

.zprofile

source <(ghatm completion zsh)

fish

ghatm completion fish > ~/.config/fish/completions/ghatm.fish
`,
		Commands: []*cli.Command{
			{
				Name:   "bash",
				Usage:  "Output shell completion script for bash",
				Action: cc.bashCompletionAction,
			},
			{
				Name:   "zsh",
				Usage:  "Output shell completion script for zsh",
				Action: cc.zshCompletionAction,
			},
			{
				Name:   "fish",
				Usage:  "Output shell completion script for fish",
				Action: cc.fishCompletionAction,
			},
		},
	}
}

func (cc *completionCommand) bashCompletionAction(context.Context, *cli.Command) error {
	// https://github.com/urfave/cli/blob/main/autocomplete/bash_autocomplete
	// https://github.com/urfave/cli/blob/c3f51bed6fffdf84227c5b59bd3f2e90683314df/autocomplete/bash_autocomplete#L5-L20
	fmt.Fprintln(cc.stdout, `
_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-shell-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-shell-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete ghatm`)
	return nil
}

func (cc *completionCommand) zshCompletionAction(context.Context, *cli.Command) error {
	// https://github.com/urfave/cli/blob/main/autocomplete/zsh_autocomplete
	// https://github.com/urfave/cli/blob/947f9894eef4725a1c15ed75459907b52dde7616/autocomplete/zsh_autocomplete
	fmt.Fprintln(cc.stdout, `#compdef ghatm

_ghatm() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-shell-completion)}")
  else
    opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-shell-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}

if [ "$funcstack[1]" = "_ghatm" ]; then
  _ghatm "$@"
else
  compdef _ghatm ghatm
fi`)
	return nil
}

func (cc *completionCommand) fishCompletionAction(_ context.Context, c *cli.Command) error {
	s, err := c.ToFishCompletion()
	if err != nil {
		return fmt.Errorf("generate fish completion: %w", err)
	}
	fmt.Fprintln(cc.stdout, s)
	return nil
}
