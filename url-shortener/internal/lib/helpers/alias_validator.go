package utils

import (
	"regexp"
)

// IsValidAlias проверяет, содержит ли строка только разрешенные символы для alias.
func IsValidAlias(alias string) bool {
	// Регулярное выражение, разрешающее только буквы (в любом регистре), цифры и подчеркивание.
	// Можно дополнить или изменить в соответствии с требованиями к alias.
	allowedChars := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

	return allowedChars.MatchString(alias)
}
