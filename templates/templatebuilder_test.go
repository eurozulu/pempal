package templates

import (
	"testing"
)

func TestTemplateBuilder_AddTemplate(t *testing.T) {
	tb := NewTemplateBuilder("../testing/templates").(*templateBuilder)

	tb.AddTemplate("certificate")
	if len(tb.buildStack) != 1 {
		t.Errorf("Expected 1 buildStack template, got %d", len(tb.buildStack))
	}
	if tb.buildStack[0].Name() != "certificate" {
		t.Errorf("Expected buildStack[0].Name to be 'certificate', got %s", tb.buildStack[0].Name())
	}
	if _, ok := tb.buildStack[0].(*CertificateTemplate); !ok {
		t.Errorf("Expected buildStack[0] as %T, got %T", &CertificateTemplate{}, tb.buildStack[0])
	}
	tb.ClearTemplates()

	tb.AddTemplate("servercertificate")
	if len(tb.buildStack) != 2 {
		t.Errorf("Expected 2 buildStack template, got %d", len(tb.buildStack))
	}
	if tb.buildStack[0].Name() != "certificate" {
		t.Errorf("Expected buildStack[0].Name to be 'certificate', got %s", tb.buildStack[0].Name())
	}
	if tb.buildStack[1].Name() != "certificate/servercertificate" {
		t.Errorf("Expected buildStack[0].Name to be 'certificate/servercertificate', got %s", tb.buildStack[0].Name())
	}
	if _, ok := tb.buildStack[1].(*template); !ok {
		t.Errorf("Expected buildStack[1] to be raw template, got %T", tb.buildStack[1])
	}
	if _, ok := tb.buildStack[0].(*CertificateTemplate); !ok {
		t.Errorf("Expected buildStack[0] as %T, got %T", &CertificateTemplate{}, tb.buildStack[0])
	}

	tb.ClearTemplates()

	tb.AddTemplate("testservercertificate")
	if len(tb.buildStack) != 3 {
		t.Errorf("Expected 3 buildStack template, got %d", len(tb.buildStack))
	}
	if tb.buildStack[0].Name() != "certificate" {
		t.Errorf("Expected buildStack[0].Name to be 'certificate', got %s", tb.buildStack[0].Name())
	}
	if tb.buildStack[1].Name() != "certificate/servercertificate" {
		t.Errorf("Expected buildStack[1].Name to be 'certificate/servercertificate', got %s", tb.buildStack[1].Name())
	}
	if tb.buildStack[2].Name() != "certificate/servercertificate/testservercertificate" {
		t.Errorf("Expected buildStack[2].Name to be 'certificate/servercertificate/testservercertificate', got %s", tb.buildStack[2].Name())
	}
	if _, ok := tb.buildStack[2].(*template); !ok {
		t.Errorf("Expected buildStack[2] to be raw template, got %T", tb.buildStack[2])
	}
	if _, ok := tb.buildStack[1].(*template); !ok {
		t.Errorf("Expected buildStack[1] to be raw template, got %T", tb.buildStack[1])
	}
	if _, ok := tb.buildStack[0].(*CertificateTemplate); !ok {
		t.Errorf("Expected buildStack[0] as %T, got %T", &CertificateTemplate{}, tb.buildStack[0])
	}

}

func TestTemplateBuilder_BaseTemplate(t *testing.T) {
	//tb := NewTemplateBuilder("../testing/templates").(*templateBuilder)
	tb := NewTemplateBuilder("").(*templateBuilder)
	tb.AddTemplate("certificate")
	baseTemplate := tb.BaseTemplate()
	if baseTemplate.Name() != "certificate" {
		t.Errorf("Expected baseName to be 'certificate', got %s", baseTemplate.Name())
	}
	tb.ClearTemplates()
	if err := tb.AddTemplate("serverkey"); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	baseTemplate = tb.BaseTemplate()
	if baseTemplate.Name() != "privatekey" {
		t.Errorf("Expected baseName to be 'key', got %s", baseTemplate.Name())
	}
}
