package site

type variableDeclarationError struct {
	op   string
	path string
}

func (e variableDeclarationError) Error() string {
	return e.op + " " + e.path
}

type ambiguousLayoutNameError struct {
	name string
}

func (e ambiguousLayoutNameError) Error() string {
	return "Ambiguous layout name: " + e.name
}
