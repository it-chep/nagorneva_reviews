package password

import "golang.org/x/crypto/bcrypt"

func Hash(v string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
	return string(b), err
}
func Check(hash, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
}
