package settings

type ReadOptions struct {
	BasePath string
}

func Options() ReadOptions {
	return ReadOptions{}
}

func (ro ReadOptions) SetBasePath(path string) ReadOptions {
	ro.BasePath = path
	return ro
}
