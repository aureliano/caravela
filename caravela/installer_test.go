package caravela

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	idir := filepath.Join(os.TempDir(), "14-bis", "test-install")
	_ = os.MkdirAll(idir, fs.ModePerm)

	bin := filepath.Join(idir, "qtbis")
	readme := filepath.Join(idir, "README.md")

	file, err := os.Create(bin)
	if err != nil {
		t.Fatal(err)
	}

	_, _ = file.WriteString("binary")
	file.Close()

	file, err = os.Create(readme)
	if err != nil {
		t.Fatal(err)
	}

	_, _ = file.WriteString("read-me")
	file.Close()

	file, err = os.Create(filepath.Join(idir, "14-bis_Linux_x86_64.tar.gz"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	file, err = os.Create(filepath.Join(idir, "14-bis_Linux_x86_64.zip"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	target, err := homeDir()
	if err != nil {
		t.Fatal(err)
	}

	err = install(idir)
	if err != nil {
		t.Fatal(err)
	}

	file, err = os.Open(filepath.Join(target, "qtbis"))
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	actual := string(bytes)
	expected := "binary"
	assert.Equal(t, expected, actual)

	file, err = os.Open(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	bytes, err = io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	actual = string(bytes)
	expected = "read-me"
	assert.Equal(t, expected, actual)

	_, err = os.Stat("14-bis_Linux_x86_64.tar.gz")
	assert.NotNil(t, err, err)

	_, err = os.Stat("14-bis_Linux_x86_64.zip")
	assert.NotNil(t, err, err)

	os.RemoveAll(idir)
}

func TestHomeDirOsExecutableFail(t *testing.T) {
	osExecutable = func() (string, error) { return "", fmt.Errorf("some error") }
	_, err := homeDir()

	actual := err.Error()
	expected := "error getting running process: some error"

	assert.Equal(t, expected, actual)
}

func TestHomeDirUnknownPath(t *testing.T) {
	osExecutable = func() (string, error) { return "/unknown/path", nil }
	_, err := homeDir()

	actual := err.Error()
	expected := "error getting information from process: lstat /unknown/path: no such file or directory"

	assert.Equal(t, expected, actual)
}

func TestHomeDir(t *testing.T) {
	osExecutable = os.Executable
	dir, err := homeDir()
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(dir, string(filepath.Separator))
	if len(parts) < 3 {
		t.Fatal("Expected at least 2 tokens.")
	}

	actual := parts[1]
	expected := filepath.Base(os.TempDir())
	assert.Equal(t, expected, actual)

	actual = parts[2]
	expected = "go-build"
	assert.Contains(t, actual, expected)
}

func TestInstallFileNew(t *testing.T) {
	file, err := os.Create(filepath.Join(os.TempDir(), "file.txt"))
	if err != nil {
		t.Fatal(err)
	}

	_, _ = file.WriteString("12345")
	file.Close()

	dest := filepath.Join(os.TempDir(), "install-new.txt")
	err = installFile(dest, file.Name())
	if err != nil {
		t.Fatal(err)
	}

	file, err = os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	actual := string(bytes)
	expected := "12345"

	assert.Equal(t, expected, actual)

	os.Remove(filepath.Join(os.TempDir(), "file.txt"))
	os.Remove(dest)
}

func TestInstallFileReplace(t *testing.T) {
	file, err := os.Create(filepath.Join(os.TempDir(), "file.txt"))
	if err != nil {
		t.Fatal(err)
	}

	_, _ = file.WriteString("54321")
	file.Close()
	source := file.Name()

	dest := filepath.Join(os.TempDir(), "install-replace.txt")
	file, err = os.Create(dest)
	if err != nil {
		t.Fatal(err)
	}

	_, _ = file.WriteString("12345")
	file.Close()

	err = installFile(dest, source)
	if err != nil {
		t.Fatal(err)
	}

	file, err = os.Open(dest)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	actual := string(bytes)
	expected := "54321"

	assert.Equal(t, expected, actual)

	os.Remove(filepath.Join(os.TempDir(), "file.txt"))
	os.Remove(dest)
}
