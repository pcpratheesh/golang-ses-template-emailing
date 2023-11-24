package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

const (
	templateName = "demo-sample-template"
)

func TestSESEmailing(t *testing.T) {
	err := godotenv.Load()
	require.Nil(t, err)

	t.Run("init templator", func(t *testing.T) {
		templator, err := NewTemplate(
			os.Getenv("REGION"),
			os.Getenv("SENDER_EMAIL"),
		)
		require.Nil(t, err)

		// create template
		t.Run("create-template", func(t *testing.T) {
			f, err := os.ReadFile("sample.html")
			require.Nil(t, err)

			err = templator.CreateTemplate(SESTemplate{
				TemplateName: templateName,
				Subject:      "Welcome Letter",
				HTMLBody:     string(f),
			})
			require.Nil(t, err)

		})
		t.Run("get-all-templates", func(t *testing.T) {
			templates, err := templator.GetAllTemplates("1")
			require.Nil(t, err)

			_ = templates
		})

		t.Run("get-template-by-id", func(t *testing.T) {
			template, err := templator.GetTemplateByName("demo-sample-template")
			require.Nil(t, err)

			_ = template
		})

		t.Run("send-email-with-template", func(t *testing.T) {
			out, err := templator.SendEmailWithTemplate(
				os.Getenv("RECIPIENT"),
				"demo-sample-template",
				map[string]interface{}{
					"name":           "John Doe",
					"favoriteanimal": "None",
				},
			)
			require.Nil(t, err)

			_ = out
		})

		t.Run("delete-template-with-id", func(t *testing.T) {
			err := templator.DeleteTemplateByName("demo-sample-template")
			require.Nil(t, err)
		})
	})

}
