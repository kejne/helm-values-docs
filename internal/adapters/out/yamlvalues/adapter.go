package yamlvalues

// Adapter exposes YAML loading through the values port.
type Adapter struct{}

func (Adapter) Load(paths []string) (map[string]any, error) { return LoadValues(paths) }
