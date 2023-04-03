package caravela

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aureliano/caravela/file"
	"github.com/aureliano/caravela/http"
	"github.com/aureliano/caravela/provider"
	"github.com/aureliano/caravela/release"
)

var downloadRelease func(http.HttpClientPlugin, *release.Release, string) (string, string, error) = release.DownloadTo
var decompress func(src string) (int, error) = file.Decompress
var checksumRelease func(binPath string, checksumsPath string) error = file.Checksum
var installRelease func(srcDir string) error = file.Install

func CheckForUpdates(client http.HttpClientPlugin, provider provider.UpdaterProvider, currver string) (*release.Release, error) {
	rel, err := provider.RestoreCacheRelease()

	if err != nil {
		rel, err = provider.FetchLastRelease(client)
		if err != nil {
			return nil, err
		}

		_ = provider.CacheRelease(*rel)
	}

	if rel.CompareTo(&release.Release{Name: currver}) == 1 {
		return rel, nil
	} else {
		return nil, nil
	}
}

func Update(client http.HttpClientPlugin, provider provider.UpdaterProvider, pname, currver string) error {
	fmt.Println("Verificando a existência de versão mais recente.")
	rel, err := CheckForUpdates(client, provider, currver)
	if err != nil {
		return err
	}

	if rel == nil {
		return fmt.Errorf("already on the edge")
	}

	fmt.Printf("Encontrada a versão %s\n", rel.Name)
	fmt.Println(rel.Description)

	dir := filepath.Join(os.TempDir(), pname)
	os.MkdirAll(dir, os.ModePerm)
	fmt.Printf("Define %s como diretório de instalação.\n", dir)

	bin, checksums, err := downloadRelease(client, rel, dir)
	if err != nil {
		return err
	}

	fmt.Println("A descomprimir arquivos.")
	num, err := decompress(bin)
	if err != nil {
		return err
	}
	fmt.Printf("%d arquivos descomprimidos de %s.\n", num, filepath.Base(bin))

	err = checksumRelease(bin, checksums)
	if err != nil {
		return err
	}

	err = installRelease(dir)
	if err != nil {
		return err
	}

	fmt.Println("Apagando arquivos de instalação.")
	os.RemoveAll(dir)

	return nil
}
