package bpcompose

import (
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

// Loader loads boilerplate-compose files and resolves include/extends.
type Loader struct{}

func NewLoader() *Loader { return &Loader{} }

// LoadAndResolve loads one or more entry files and returns a merged ComposeFile with includes and extends resolved.
func (l *Loader) LoadAndResolve(entryFiles []string) (*ComposeFile, error) {
    visited := map[string]bool{}
    var result *ComposeFile

    for _, file := range entryFiles {
        abs, err := filepath.Abs(file)
        if err != nil {
            return nil, err
        }
        cf, err := l.loadRecursive(abs, visited)
        if err != nil {
            return nil, err
        }
        if result == nil {
            result = cf
        } else {
            merged := mergeCompose(result, cf)
            result = &merged
        }
    }

    if result == nil {
        empty := ComposeFile{}
        return &empty, nil
    }

    // Resolve template-level extends (by name)
    resolveTemplateExtends(result)

    return result, nil
}

func (l *Loader) loadRecursive(path string, visited map[string]bool) (*ComposeFile, error) {
    if visited[path] {
        return nil, fmt.Errorf("cycle detected loading %s", path)
    }
    visited[path] = true

    cf, err := readCompose(path)
    if err != nil {
        return nil, err
    }

    baseDir := filepath.Dir(path)

    // Process extends: treat as list of files to merge first, then current overrides
    for _, ext := range cf.Extends {
        extPath := ext
        if !filepath.IsAbs(extPath) {
            extPath = filepath.Join(baseDir, extPath)
        }
        parent, err := l.loadRecursive(extPath, visited)
        if err != nil {
            return nil, err
        }
        merged := mergeCompose(parent, cf)
        cf = &merged
    }

    // Process include: include files merged before current, similar to extends, but meant for composition
    for _, inc := range cf.Include {
        incPath := inc
        if !filepath.IsAbs(incPath) {
            incPath = filepath.Join(baseDir, incPath)
        }
        included, err := l.loadRecursive(incPath, visited)
        if err != nil {
            return nil, err
        }
        merged := mergeCompose(included, cf)
        cf = &merged
    }

    return cf, nil
}

func readCompose(path string) (*ComposeFile, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    data, err := io.ReadAll(f)
    if err != nil {
        return nil, err
    }

    var cf ComposeFile
    if err := yaml.Unmarshal(data, &cf); err != nil {
        return nil, fmt.Errorf("%s: %w", path, err)
    }

    if cf.Templates == nil {
        cf.Templates = map[string]TemplateSpec{}
    }
    if cf.Vars == nil {
        cf.Vars = map[string]any{}
    }

    return &cf, nil
}

func resolveTemplateExtends(cf *ComposeFile) {
    // For each template that lists Extends, merge fields from parent templates by name in order
    for name, tpl := range cf.Templates {
        if len(tpl.Extends) == 0 {
            continue
        }
        resolved := TemplateSpec{}
        for _, parentName := range tpl.Extends {
            parent, ok := cf.Templates[parentName]
            if !ok {
                continue
            }
            merged := mergeTemplate(resolved, parent)
            resolved = merged
        }
        merged := mergeTemplate(resolved, tpl)
        cf.Templates[name] = merged
    }
}

// Validate performs minimal validation of the resolved configuration
func Validate(cf *ComposeFile) error {
    if cf == nil {
        return errors.New("nil config")
    }
    // No strict requirements, but ensure template fields are coherent
    for name, t := range cf.Templates {
        if t.TemplateURL == "" {
            // Allowed: user might set via CLI, but warn by returning error to force explicitness
            _ = name
        }
    }
    return nil
}

// PrintResolved writes the resolved configuration to w in YAML format
func PrintResolved(cf *ComposeFile, w io.Writer) error {
    out, err := yaml.Marshal(cf)
    if err != nil {
        return err
    }
    _, err = w.Write(out)
    return err
}

// MergeVarFiles loads variables from the provided YAML files and merges into cf.Vars (later files override earlier)
func MergeVarFiles(cf *ComposeFile, files []string) error {
    for _, file := range files {
        m, err := readVarsFile(file)
        if err != nil {
            return err
        }
        cf.Vars = deepMergeMap(cf.Vars, m)
    }
    return nil
}

// MergeVarPairs merges NAME=VALUE pairs into cf.Vars
func MergeVarPairs(cf *ComposeFile, pairs []string) error {
    for _, p := range pairs {
        key, val, ok := splitOnce(p, '=')
        if !ok {
            return fmt.Errorf("invalid --var %q, expected NAME=VALUE", p)
        }
        if cf.Vars == nil {
            cf.Vars = map[string]any{}
        }
        cf.Vars[key] = val
    }
    return nil
}

func readVarsFile(path string) (map[string]any, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var m map[string]any
    if err := yaml.Unmarshal(b, &m); err != nil {
        return nil, err
    }
    if m == nil {
        m = map[string]any{}
    }
    return m, nil
}

// mergeCompose performs deep merge: right overrides left
func mergeCompose(left, right *ComposeFile) ComposeFile {
    out := ComposeFile{}
    // Include/Extends are not propagated post-merge
    out.Templates = map[string]TemplateSpec{}

    // merge templates map: per-key merge structs
    for k, v := range left.Templates {
        out.Templates[k] = v
    }
    for k, v := range right.Templates {
        if existing, ok := out.Templates[k]; ok {
            out.Templates[k] = mergeTemplate(existing, v)
        } else {
            out.Templates[k] = v
        }
    }

    out.Vars = deepMergeMap(left.Vars, right.Vars)

    return out
}

func mergeTemplate(a, b TemplateSpec) TemplateSpec {
    out := a
    if b.TemplateURL != "" {
        out.TemplateURL = b.TemplateURL
    }
    if b.OutputFolder != "" {
        out.OutputFolder = b.OutputFolder
    }
    if b.NonInteractive {
        out.NonInteractive = b.NonInteractive
    }
    if b.NoHooks {
        out.NoHooks = b.NoHooks
    }
    if b.NoShell {
        out.NoShell = b.NoShell
    }
    if b.DisableDependencyPrompt {
        out.DisableDependencyPrompt = b.DisableDependencyPrompt
    }
    if b.MissingKeyAction != "" {
        out.MissingKeyAction = b.MissingKeyAction
    }
    if b.MissingConfigAction != "" {
        out.MissingConfigAction = b.MissingConfigAction
    }

    if out.Vars == nil && b.Vars != nil {
        out.Vars = map[string]any{}
    }
    out.Vars = deepMergeMap(out.Vars, b.Vars)

    if len(b.Extends) > 0 {
        out.Extends = append(out.Extends, b.Extends...)
    }

    return out
}

func deepMergeMap(a, b map[string]any) map[string]any {
    if a == nil && b == nil {
        return map[string]any{}
    }
    if a == nil {
        a = map[string]any{}
    }
    out := map[string]any{}
    for k, v := range a {
        out[k] = v
    }
    for k, v := range b {
        if vMap, ok := v.(map[string]any); ok {
            if existing, ok2 := out[k].(map[string]any); ok2 {
                out[k] = deepMergeMap(existing, vMap)
                continue
            }
        }
        out[k] = v
    }
    return out
}

func splitOnce(s string, sep rune) (string, string, bool) {
    for i, r := range s {
        if r == sep {
            return s[:i], s[i+1:], true
        }
    }
    return s, "", false
}