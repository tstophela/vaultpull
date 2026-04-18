package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// TemplateApplyOptions configures how a template is applied to secrets.
type TemplateApplyOptions struct {
	TemplatePath string
	Strict       bool
	Output       io.Writer
}

// ApplyTemplate renders a .env template file using the provided secrets
// and writes the result in KEY=VALUE format to opts.Output.
func ApplyTemplate(secrets map[string]string, opts TemplateApplyOptions) error {
	if opts.TemplatePath == "" {
		return fmt.Errorf("template: path is required")
	}
	if opts.Output == nil {
		return fmt.Errorf("template: output writer is required")
	}

	renderer := NewTemplateRenderer(opts.Strict)
	resolved, err := renderer.RenderFile(opts.TemplatePath, secrets)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(resolved))
	for k := range resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, resolved[k]))
	}

	_, err = fmt.Fprint(opts.Output, sb.String())
	return err
}
