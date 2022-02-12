package exercise

import "strings"

func ReverseWords(s string) string {
	strSli := strings.Fields(s)

	for i := 0; i <= (len(strSli)-1)/2; i++ {
		strSli[i], strSli[len(strSli)-1-i] = strSli[len(strSli)-1-i], strSli[i]
	}
	return strings.Join(strSli, " ")
}
