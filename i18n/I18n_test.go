package i18n

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWmsgNotVerbose(t *testing.T) {
	err := PrepareI18n(I18nConf{false, 0})
	assert.Nil(t, err)
	n := Wmsg(100)

	assert.Equal(t, -1, n)
}

func TestWmsgKeyNotFound(t *testing.T) {
	err := PrepareI18n(I18nConf{true, 0})
	assert.Nil(t, err)
	n := Wmsg(0)

	assert.Equal(t, -1, n)
}

func TestWmsg(t *testing.T) {
	err := PrepareI18n(I18nConf{true, 0})
	assert.Nil(t, err)
	n := Wmsg(100)

	assert.Equal(t, 29, n)
}

func TestPrepareI18nInvalidLocale(t *testing.T) {
	err := PrepareI18n(I18nConf{false, -1})
	assert.Equal(t, "invalid locale -1", err.Error())
}

func TestPrepareI18n(t *testing.T) {
	conf := I18nConf{false, PT_BR}
	err := PrepareI18n(conf)

	assert.Nil(t, err)
	assert.Equal(t, conf, config)
	assert.Equal(t, "Baixando pacote de atualização.\n", msg[100])

	conf = I18nConf{false, EN}
	err = PrepareI18n(conf)

	assert.Nil(t, err)
	assert.Equal(t, conf, config)
	assert.Equal(t, "Downloading update package.\n", msg[100])

}

func TestValidateLocale(t *testing.T) {
	type testCase struct {
		name     string
		input    int
		expected error
	}
	testCases := []testCase{
		{
			name:     "-1 is an invalid locale",
			input:    -1,
			expected: fmt.Errorf("invalid locale -1"),
		},
		{
			name:     "2 is an invalid locale",
			input:    2,
			expected: fmt.Errorf("invalid locale 2"),
		},
		{
			name:     "0 is a valid locale",
			input:    0,
			expected: nil,
		},
		{
			name:     "1 is a valid locale",
			input:    1,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := validateLocale(tc.input)
			if (tc.expected != nil && actual == nil) || (tc.expected == nil && actual != nil) {
				t.Errorf("expected %s, got %s", tc.expected, actual)
			} else {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
