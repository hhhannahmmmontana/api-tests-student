package token

func DistortToken(token string) string {
	lastChar := token[len(token)-1]
	replacement := "a"
	if lastChar == 'a' {
		replacement = "b"
	}
	return token[:len(token)-1] + replacement
}
