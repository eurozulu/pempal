package finder

type LocationParser interface {
	Parse(path string, data []byte) (Location, error)
}
