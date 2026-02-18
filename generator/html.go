package generator

import (
	"embed"
	"encoding/json"
	"html/template"
	"os"

	"github.com/rhobs/testplan-viewer/models"
)

//go:embed templates/testplan.html
var templateFS embed.FS

// GenerateHTML generates an interactive HTML file from a test plan
func GenerateHTML(testPlan *models.TestPlan, outputPath string) error {
	// Create template with custom functions
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
		return err
	}

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute template
	return tmpl.Execute(f, testPlan)
}
