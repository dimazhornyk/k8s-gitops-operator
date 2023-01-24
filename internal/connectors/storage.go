package connectors

// TODO: use real storage
type storage struct {
	m map[string]string
}

func NewStorage() Storage {
	return &storage{
		m: make(map[string]string),
	}
}

func (s storage) SaveConfigHash(repo string, hash string) error {
	s.m[repo] = hash

	return nil
}

func (s storage) GetConfigHash(repo string) string {
	return s.m[repo]
}
