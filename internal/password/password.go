package password

type Hasher interface {
	Hash(password string) (string, error)
}

type Comparer interface {
	Compare(hash, password string) bool
}

type HasherComparer interface {
	Hasher
	Comparer
}
