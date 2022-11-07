package finder

type LocationFilter interface {
	FilepathFilter
	FilterLocation(rl Location) Location
}
