package shorten

// Shortener represents the interface for generic shortener service
type Shortener interface {
	shorten(u string) (string, error)
}

// Shorten executes the shortening service
func Shorten(s Shortener, u string) (string, error) {
	return s.shorten(u)
}
