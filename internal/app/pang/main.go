package pang

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/pangbox/pangfiles/internal/cli/commands/pak"
	"github.com/pangbox/pangfiles/internal/cli/commands/updatelist"
	"github.com/pangbox/pangfiles/internal/cli/commands/updatelistsrv"
	"github.com/pangbox/pangfiles/internal/cli/commands/version"
)

func Main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(version.Command(), "")
	subcommands.Register(pak.ExtractCommand(), "paks")
	subcommands.Register(pak.MountCommand(), "paks")
	subcommands.Register(updatelist.EncryptCommand(), "updatelists")
	subcommands.Register(updatelist.DecryptCommand(), "updatelists")
	subcommands.Register(updatelistsrv.ServeCommand(), "updatelists")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
