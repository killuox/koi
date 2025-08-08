package config

type Config struct {
	API       API                 `yaml:"api"`
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

type API struct {
	BaseURL string `yaml:"baseUrl"`
	Version string `yaml:"version"`
	Auth    Auth   `yaml:"auth"`
}

type Auth struct {
	Type string `yaml:"type"`
}

type Endpoint struct {
	Type       string               `yaml:"type"`
	Method     string               `yaml:"method"`
	Mode       string               `yaml:"mode"`
	Path       string               `yaml:"path"`
	Parameters map[string]Parameter `yaml:"parameters"`
	Defaults   map[string]any       `yaml:"defaults"`
}

type Parameter struct {
	Mode        string     `yaml:"mode"`
	In          string     `yaml:"in"`
	Description string     `yaml:"description"`
	Required    bool       `yaml:"required"`
	Validation  Validation `yaml:"validation"`
}

type Validation struct {
	MinLength int `yaml:"minLength"`
	MaxLength int `yaml:"maxLength"`
}

func (p Parameter) GetValue(cfg Config) string {
	// Check if we have a flag first
	// else use the default
	// else do nothing or throw error?
	return ""
}

// type ValueGetter struct {
// 	Key  string
// 	Type string
// }

// func (v ValueGetter) GetValue() interface{} {
// 	cmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
// 	var value interface{}
// 	switch v.Type {
// 	case "string":
// 		value = v.getStringValue(cmd)
// 	case "bool":
// 		value = v.getBoolValue(cmd)
// 	case "int":
// 		value = v.getIntValue(cmd)
// 	default:
// 		value = nil
// 	}

// 	cmd.Parse(os.Args[2:])

// 	if value == nil {
// 		fmt.Println("Error: --slug is required for get-pokemon command.")
// 		cmd.Usage()
// 		os.Exit(1)
// 	}

// 	return value
// }

// func (v ValueGetter) getStringValue(cmd *flag.FlagSet) string {
// 	return *cmd.String(v.Key, "", "")
// }

// func (v ValueGetter) getBoolValue(cmd *flag.FlagSet) bool {
// 	return *cmd.Bool(v.Key, false, "")
// }

// func (v ValueGetter) getIntValue(cmd *flag.FlagSet) int {
// 	return *cmd.Int(v.Key, 0, "")
// }
