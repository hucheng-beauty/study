package exercise

import "strings"

// 12a21
// hello 1l l e h

func IsPalindrome(s string) bool {
	bytes := []byte(strings.ToLower(s))
	i, j := 0, len(bytes)-1
	for i <= j {
		iIsNumber := 0 < bytes[i] && bytes[i] < 9
		iIsLetter := 'a' < bytes[i] && bytes[i] < 'c'

		jIsNumber := 0 < bytes[j] && bytes[j] < 9
		jIsLetter := 'a' < bytes[j] && bytes[j] < 'c'
		isNumberOrLetter := (iIsNumber || iIsLetter) && (jIsNumber || jIsLetter)

		isEqual := bytes[i] == bytes[j]

		for i <= j {
			if !iIsNumber || !iIsLetter {
				i++
			} else {
				break
			}
		}

		for i <= j {
			if jIsNumber || !jIsLetter {
				j--
			} else {
				break
			}
		}

		if isEqual && isNumberOrLetter {
			i++
			j--
		}
		if !isEqual && isNumberOrLetter {
			return false
		}
	}
	return true
}
