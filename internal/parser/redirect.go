package parser

type Redirect struct {
	FileDescriptor int    // 1 for stdout, 2 for stderr
	Target         string // filename
	Append         bool   // >> is append to file while > is overwrite
}

func getFileDescriptor(redirectOperator string) int {
	if redirectOperator == ">" || redirectOperator == "1>" {
		return 1
	}
	return 2
}
