package hash

type Hash interface {
	Hash(password string) (string, error)
	Verify(hash, password string) (bool, error)
}
