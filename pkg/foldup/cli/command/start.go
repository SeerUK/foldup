package command

import (
	"context"
	"os"
	"path"

	gstorage "cloud.google.com/go/storage"
	"github.com/SeerUK/foldup/pkg/archive"
	"github.com/SeerUK/foldup/pkg/storage"
	"github.com/SeerUK/foldup/pkg/storage/gcs"
	"github.com/SeerUK/foldup/pkg/xioutil"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StartCommand creates a command to trigger periodic backups.
func StartCommand() *console.Command {
	var dirname string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&dirname),
			Spec:  "DIRNAME",
			Desc:  "The directory to archive folders from.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		dirs, err := xioutil.ReadDirsInDir(dirname, false)
		if err != nil {
			return err
		}

		relativePaths := []string{}

		for _, d := range dirs {
			relativePaths = append(relativePaths, path.Join(dirname, d.Name()))
		}

		archives, err := archive.Dirsf(relativePaths, "backup-%s-%d", archive.TarGz)
		if err != nil {
			return err
		}

		storageClient, err := gstorage.NewClient(context.Background())
		if err != nil {
			return err
		}

		client := gcs.NewGoogleClient(storageClient)
		gateway := storage.NewGCSGateway(client, "backups-sierra")

		for _, a := range archives {
			in, err := os.Open(a)
			if err != nil {
				return err
			}

			output.Printf("Going to store '%s'\n", a)

			err = gateway.Store(context.Background(), a, in)
			if err != nil {
				return err
			}

			err = os.Remove(a)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return &console.Command{
		Name:        "start",
		Description: "Begin periodically backing up.",
		Configure:   configure,
		Execute:     execute,
	}
}
