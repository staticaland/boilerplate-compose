package processor

import (
	"fmt"
	"path/filepath"

	"boilerplate-compose/config"
)

type TemplateProcessor struct {
	config     *config.ComposeConfig
	configPath string
}

func NewTemplateProcessor(cfg *config.ComposeConfig, configPath string) *TemplateProcessor {
	return &TemplateProcessor{
		config:     cfg,
		configPath: configPath,
	}
}

type ProcessingJob struct {
	Name     string
	Template config.Template
	Args     []string
}

func (tp *TemplateProcessor) BuildProcessingJobs() ([]ProcessingJob, error) {
	var jobs []ProcessingJob

	for name, template := range tp.config.Templates {
		args, err := tp.buildBoilerplateArgs(template)
		if err != nil {
			return nil, fmt.Errorf("failed to build args for template '%s': %w", name, err)
		}

		jobs = append(jobs, ProcessingJob{
			Name:     name,
			Template: template,
			Args:     args,
		})
	}

	return jobs, nil
}

func (tp *TemplateProcessor) buildBoilerplateArgs(template config.Template) ([]string, error) {
	var args []string

	// Add template URL
	args = append(args, "--template-url", template.TemplateURL)

	// Add output folder (resolve relative to config file)
	outputPath := tp.resolveOutputPath(template.OutputFolder)
	args = append(args, "--output-folder", outputPath)

	// Add variables
	for key, value := range template.Vars {
		args = append(args, "--var", fmt.Sprintf("%s=%s", key, value))
	}

	// Add var-file(s)
	if template.VarFile != nil {
		switch v := template.VarFile.(type) {
		case string:
			args = append(args, "--var-file", v)
		case []interface{}:
			for _, file := range v {
				if fileStr, ok := file.(string); ok {
					args = append(args, "--var-file", fileStr)
				}
			}
		}
	}

	// Add boolean flags
	if template.NonInteractive {
		args = append(args, "--non-interactive")
	}

	if template.NoHooks {
		args = append(args, "--no-hooks")
	}

	if template.NoShell {
		args = append(args, "--no-shell")
	}

	if template.DisableDependencyPrompt {
		args = append(args, "--disable-dependency-prompt")
	}

	// Add action flags
	if template.MissingKeyAction != "" {
		args = append(args, "--missing-key-action", template.MissingKeyAction)
	}

	if template.MissingConfigAction != "" {
		args = append(args, "--missing-config-action", template.MissingConfigAction)
	}

	return args, nil
}

func (tp *TemplateProcessor) resolveOutputPath(outputFolder string) string {
	// If path is already absolute, return as-is
	if filepath.IsAbs(outputFolder) {
		return outputFolder
	}
	
	// Get directory containing the config file
	configDir := filepath.Dir(tp.configPath)
	
	// Join config directory with relative output path
	return filepath.Join(configDir, outputFolder)
}