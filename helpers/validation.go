package helpers

import "regexp"



func IsValidEmail(email string) bool {

    // Regular expression for validating an Email
    var re = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return re.MatchString(email)
}

// IsStrongPassword checks if the password is strong enough
func IsStrongPassword(password string) bool {
    var passwordRegex = regexp.MustCompile(`^[a-zA-Z\d]{8,}$`)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`\d`).MatchString(password)
    return passwordRegex.MatchString(password) && hasLower && hasUpper && hasDigit
}