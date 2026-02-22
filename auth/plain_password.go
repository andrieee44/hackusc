package auth

type PlainPassword struct{}

func (PlainPassword) Hash(plain string) ([]byte, error) {
	return []byte(plain), nil
}

func (PlainPassword) ComparePlainToHash(plain string, hash []byte) (bool, error) {
	return plain == string(hash), nil
}
