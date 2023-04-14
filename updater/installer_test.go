package updater

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
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

	exec, err := processFilePath()
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Dir(exec)

	err = install(idir, target)
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

func TestShouldIgoreFile(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected bool
	}
	testCases := []testCase{
		{
			name:     "should ignore package file",
			input:    "file.zip",
			expected: true,
		},
		{
			name:     "should ignore checksums file",
			input:    "checksums.txt",
			expected: true,
		},
		{
			name:     "should not ignore file",
			input:    "file.md",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := shouldIgoreFile(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, actual)
			}
		})
	}
}
