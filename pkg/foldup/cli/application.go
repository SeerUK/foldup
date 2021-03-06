package cli

import (
	"io"

	"github.com/SeerUK/foldup/pkg/foldup"
	"github.com/SeerUK/foldup/pkg/foldup/cli/command"
	"github.com/eidolon/console"
)

// CreateApplication builds the console application instance. Providing it with some basic
// information like the name and version.
func CreateApplication(writer io.Writer) *console.Application {
	application := console.NewApplication("foldup", "0.1.0")
	application.Writer = writer
	application.Logo = `
███████╗ ██████╗ ██╗     ██████╗ ██╗   ██╗██████╗
██╔════╝██╔═══██╗██║     ██╔══██╗██║   ██║██╔══██╗
█████╗  ██║   ██║██║     ██║  ██║██║   ██║██████╔╝
██╔══╝  ██║   ██║██║     ██║  ██║██║   ██║██╔═══╝
██║     ╚██████╔╝███████╗██████╔╝╚██████╔╝██║
╚═╝      ╚═════╝ ╚══════╝╚═════╝  ╚═════╝ ╚═╝
`

	application.AddCommands(buildCommands(foldup.NewCLIFactory()))

	return application
}

// buildCommands instantiates all of the commands registered in the application.
func buildCommands(factory foldup.Factory) []*console.Command {
	return []*console.Command{
		command.BackupCommand(factory),
	}
}
