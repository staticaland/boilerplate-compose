package bpcompose

// ComposeFile represents the top-level boilerplate-compose document.
type ComposeFile struct {
    Include   []string                     `yaml:"include,omitempty"`
    Extends   []string                     `yaml:"extends,omitempty"`
    Templates map[string]TemplateSpec      `yaml:"templates,omitempty"`
    Vars      map[string]any               `yaml:"vars,omitempty"`
}

// TemplateSpec defines a single template invocation.
type TemplateSpec struct {
    TemplateURL              string         `yaml:"template_url,omitempty"`
    OutputFolder             string         `yaml:"output_folder,omitempty"`
    NonInteractive           bool           `yaml:"non_interactive,omitempty"`
    NoHooks                  bool           `yaml:"no_hooks,omitempty"`
    NoShell                  bool           `yaml:"no_shell,omitempty"`
    DisableDependencyPrompt  bool           `yaml:"disable_dependency_prompt,omitempty"`
    MissingKeyAction         string         `yaml:"missing_key_action,omitempty"`
    MissingConfigAction      string         `yaml:"missing_config_action,omitempty"`
    Vars                     map[string]any `yaml:"vars,omitempty"`
    Extends                  []string       `yaml:"extends,omitempty"` // extend other templates by name
}