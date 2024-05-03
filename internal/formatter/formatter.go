package formatter

import (
	"fmt"
	"strings"
)

func MaxLengthAfterSplit(str string) int {
	substrings := strings.Split(str, "\n")
	maxLength := 0

	if strings.Contains(str, "\n") {
		for _, sub := range substrings {
			if len(sub) > maxLength {
				maxLength = len(sub)
			}
		}
	} else {
		maxLength = len(str)
	}

	return maxLength
}

// LineSplit adds the "\n" to the headers' rows
func LineSplit(data [][]string) {
	every := 7
	for n, colName := range data[0] {
		var result string
		cnt := 0
		for i, char := range colName {
			result += string(char)
			cnt++
			if ((i+1) > every || (i+1) < 2*every) && char == ' ' && cnt > 8 {
				result += "\n"
				cnt = 0
			}
		}
		data[0][n] = strings.ReplaceAll(result, "|", "\n")
	}
}

func VersionFormatter(localVer []byte, env string) string {
	localVerStr := string(localVer)
	year := localVerStr[:4]
	month := localVerStr[4:6]
	day := localVerStr[6:8]
	hour := localVerStr[8:10]
	minute := localVerStr[10:12]
	second := localVerStr[12:14]
	return fmt.Sprintf("%s.%s.%s %s:%s:%s\n%s mode", year, month, day, hour, minute, second, env)
}
