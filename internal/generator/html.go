package generator

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"os"

	"github.com/nephomaniac/testplan_tools_poc/internal/models"
)

//go:embed templates/testplan.html
var templateFS embed.FS

// Generator handles HTML generation from test plans
type Generator struct {
	template *template.Template
}

// New creates a new HTML generator with embedded template
func New() (*Generator, error) {
	funcMap := template.FuncMap{
		"json": func(v interface{}) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case map[string]models.TestCase:
				return len(val)
			case []string:
				return len(val)
			default:
				return 0
			}
		},
	}

	tmpl, err := template.New("testplan.html").
		Funcs(funcMap).
		ParseFS(templateFS, "templates/testplan.html")
	if err != nil {
		return nil, err
	}

	return &Generator{template: tmpl}, nil
}

// GenerateHTML generates an interactive HTML file from a test plan
func (g *Generator) GenerateHTML(testPlan *models.TestPlan, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return g.template.Execute(f, testPlan)
}

// GenerateToWriter generates HTML to an io.Writer (useful for testing)
func (g *Generator) GenerateToWriter(testPlan *models.TestPlan, w io.Writer) error {
	return g.template.Execute(w, testPlan)
}
