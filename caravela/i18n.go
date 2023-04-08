package caravela

import "fmt"

const (
	En   = 0
	PtBr = 1
)

type I18nConf struct {
	Verbose bool
	Locale  int
}

var (
	msg    map[int]string
	config I18nConf
)

func wmsg(key int, parameters ...interface{}) int {
	if !config.Verbose {
		return -1
	}

	format := msg[key]
	if format == "" {
		return -1
	}

	message := fmt.Sprintf(format, parameters...)

	n, _ := fmt.Println(message)

	return n
}

func prepareI18n(conf I18nConf) error {
	err := validateLocale(conf.Locale)
	if err != nil {
		return err
	}

	config = conf

	if conf.Locale == PtBr {
		msg = ptBrMessages
	} else {
		msg = enMessages
	}

	return nil
}

func validateLocale(locale int) error {
	if PtBr != locale && En != locale {
		return fmt.Errorf("invalid locale %d", locale)
	}

	return nil
}
