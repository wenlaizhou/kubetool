package kubetool

type Dev struct {
	Env       []string
	Port      []string
	Image     string
	Namespace string
	Volumes   map[string]string
	Command   []string
	Args      []string
}
