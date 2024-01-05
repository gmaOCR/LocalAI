package assets

import (
<<<<<<< HEAD
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/rs/zerolog/log"
)

func ListFiles(content *rice.Box) (files []string) {
	err := content.Walk("", func(path string, info os.FileInfo, err error) error {
=======
	"embed"
	"io/fs"
)

func ListFiles(content embed.FS) (files []string) {
	fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
>>>>>>> fda6bf56 (feat: embedded model configurations, add popular model examples, refactoring (#1532))
		if err != nil {
			return err
		}

<<<<<<< HEAD
		if info.IsDir() {
=======
		if d.IsDir() {
>>>>>>> fda6bf56 (feat: embedded model configurations, add popular model examples, refactoring (#1532))
			return nil
		}

		files = append(files, path)
		return nil
	})
<<<<<<< HEAD
	if err != nil {
		log.Error().Err(err).Msg("error walking the rice box")
	}
=======
>>>>>>> fda6bf56 (feat: embedded model configurations, add popular model examples, refactoring (#1532))
	return
}
