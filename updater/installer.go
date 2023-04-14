package updater

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func install(srcDir, destDir string) error {
	err := filepath.Walk(srcDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if srcDir == path {
				return nil
			}

			fname := filepath.Base(path)
			cleanPath := filepath.Clean(path)
			relPath := strings.ReplaceAll(cleanPath, srcDir, "")

			if shouldIgoreFile(fname) {
				return nil
			}

			if info.IsDir() {
				err = os.Mkdir(filepath.Join(destDir, relPath), info.Mode())
				if err != nil {
					return err
				}
			} else {
				src := filepath.Join(srcDir, relPath)
				dest := filepath.Join(destDir, relPath)

				err = installFile(dest, src)
				if err != nil {
					return err
				}
			}

			return nil
		})

	return err
}

func installFile(dest, src string) error {
	destInfo, err := os.Stat(dest)
	fileExist := true
	var fm fs.FileMode

	if os.IsNotExist(err) {
		fileExist = false
		const permFile = 0644
		fm = fs.FileMode(permFile)
	} else if err != nil {
		return err
	}

	if fileExist {
		if err = os.Remove(dest); err != nil {
			return err
		}
		fm = destInfo.Mode()
	}

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fm)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func shouldIgoreFile(fname string) bool {
	types := []string{
		".deb",
		".zip.sbom",
		".tar.gz.sbom",
		".tar.gz",
		".apk",
		".tar.zst",
		".zip",
		".rpm",
	}

	if fname == "checksums.txt" {
		return true
	}

	for _, tp := range types {
		if strings.HasSuffix(fname, tp) {
			return true
		}
	}

	return false
}
