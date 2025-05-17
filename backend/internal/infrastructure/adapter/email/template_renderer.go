package email

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/huydq/test/internal/pkg/logger"
)

// TemplateRenderer handles email template loading and rendering
type TemplateRenderer struct {
	templateDir string
	logger      logger.Logger
}

// NewTemplateRenderer creates a new instance of TemplateRenderer
func NewTemplateRenderer(templateDir string, logger logger.Logger) *TemplateRenderer {
	return &TemplateRenderer{
		templateDir: templateDir,
		logger:      logger,
	}
}

// RenderTemplate renders an email template with the provided data
func (r *TemplateRenderer) RenderTemplate(templateFolder, templateFile string, data map[string]any) (string, error) {
	mainTemplatePath := filepath.Join(r.templateDir, templateFolder, templateFile)
	footerTemplatePath := filepath.Join(r.templateDir, "common", "footer.tmpl")

	// Read main template file
	mainContent, err := os.ReadFile(mainTemplatePath)
	if err != nil {
		r.logger.Error("Failed to read main template", map[string]any{
			"path":  mainTemplatePath,
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to read main template: %w", err)
	}

	// Read footer template (optional)
	footerContent, err := os.ReadFile(footerTemplatePath)
	if err != nil {
		r.logger.Error("Footer template not found or couldn't be read", map[string]any{
			"path":  footerTemplatePath,
			"error": err.Error(),
		})
		footerContent = []byte("") // Use empty footer if not found
	}

	// Combine templates
	combinedContent := fmt.Sprintf(`{{define "main"}}%s{{end}}
{{define "footer"}}%s{{end}}

{{template "main" .}}
{{template "footer" .}}`, string(mainContent), string(footerContent))

	// Parse and execute template
	tmpl, err := template.New("email").Parse(combinedContent)
	if err != nil {
		r.logger.Error("Failed to parse template", map[string]any{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		r.logger.Error("Failed to execute template", map[string]any{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
