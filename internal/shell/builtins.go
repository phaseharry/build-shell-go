package shell

func (s *Shell) builtRegistry() map[string]builtin {
	return map[string]builtin{
		"exit": s.exit,
	}
}

func (s *Shell) exit(args []string) Result {
	// exiting the shell with status code of 0 because there was no error
	return Result{
		Exit:   true,
		Status: 0,
	}
}
