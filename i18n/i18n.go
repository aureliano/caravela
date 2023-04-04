package i18n

import "fmt"

const (
	EN    = 0
	PT_BR = 1
)

type I18nConf struct {
	Verbose bool
	Locale  int
}

var (
	msg    map[int]string
	config I18nConf
)

func Wmsg(key int, parameters ...interface{}) (int, error) {
	if !config.Verbose {
		return -1, nil
	}

	format := msg[key]
	if format == "" {
		return -1, fmt.Errorf("key '%d' not found", key)
	}

	message := fmt.Sprintf(format, parameters...)

	return fmt.Println(message)
}

func PrepareI18n(conf I18nConf) error {
	err := validateLocale(conf.Locale)
	if err != nil {
		return err
	}

	config = conf

	if conf.Locale == PT_BR {
		msg = ptBrMessages
	} else {
		msg = enMessages
	}

	return nil
}

func validateLocale(locale int) error {
	if PT_BR != locale && EN != locale {
		return fmt.Errorf("invalid locale %d", locale)
	}

	return nil
}
