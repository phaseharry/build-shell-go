package parser

type Command struct {
	Name string
	Args []string
	Redirects []Redirect
}

type Redirect struct {
	FileDescriptor int    // 1 for stdout, 2 for stderr
	Target         string // filename
	Append         bool   // >> is append to file while > is overwrite
}
