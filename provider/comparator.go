package provider

import (
	"regexp"
	"strconv"
	"strings"
)

var releaseRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)-?([\w.]+)?$`)

func compareVersions(ver1, ver2 string) int {
	match1 := releaseRegex.FindAllStringSubmatch(ver1, -1)
	match2 := releaseRegex.FindAllStringSubmatch(ver2, -1)

	major := compareVersionParts(match1[0][1], match2[0][1])
	if major != 0 {
		return major
	}

	minor := compareVersionParts(match1[0][2], match2[0][2])
	if minor != 0 {
		return minor
	}

	patch := compareVersionParts(match1[0][3], match2[0][3])
	if patch != 0 {
		return patch
	}

	pr1 := match1[0][4]
	pr2 := match2[0][4]

	switch {
	case pr1 != "" && pr2 == "":
		return -1
	case pr1 == "" && pr2 != "":
		return 1
	default:
		return strings.Compare(match1[0][4], match2[0][4])
	}
}

func compareVersionParts(v1, v2 string) int {
	iv1, _ := strconv.Atoi(v1)
	iv2, _ := strconv.Atoi(v2)

	switch {
	case iv1 > iv2:
		return 1
	case iv1 == iv2:
		return 0
	default:
		return -1
	}
}
