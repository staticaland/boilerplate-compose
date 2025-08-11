package config

type ComposeConfig struct {
	Templates map[string]Template `yaml:"templates"`
	Include   []IncludeConfig     `yaml:"include,omitempty"`
	Extends   *ExtendsConfig      `yaml:"extends,omitempty"`
}

type Template struct {
	TemplateURL             string            `yaml:"template-url"`
	OutputFolder            string            `yaml:"output-folder"`
	Vars                    map[string]string `yaml:"vars,omitempty"`
	VarFile                 interface{}       `yaml:"var-file,omitempty"` // string or []string
	NonInteractive          bool              `yaml:"non-interactive,omitempty"`
	MissingKeyAction        string            `yaml:"missing-key-action,omitempty"`
	MissingConfigAction     string            `yaml:"missing-config-action,omitempty"`
	NoHooks                 bool              `yaml:"no-hooks,omitempty"`
	NoShell                 bool              `yaml:"no-shell,omitempty"`
	DisableDependencyPrompt bool              `yaml:"disable-dependency-prompt,omitempty"`
}

type IncludeConfig struct {
	Path string `yaml:"path"`
}

type ExtendsConfig struct {
	File     string `yaml:"file"`
	Template string `yaml:"template"`
}