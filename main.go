package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/spf13/cast"
)

type Template struct {
	sess   *session.Session
	svc    *ses.SES
	Sender string
}

type SESTemplate struct {
	TemplateName string `json:"templateName"`
	Subject      string `json:"subject"`
	HTMLBody     string `json:"htmlBody"`
	TextBody     string `json:"textBody"`
}

func NewSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
}

func NewTemplate(region string, sender string) (*Template, error) {
	sess, err := NewSession(region)
	if err != nil {
		return nil, err
	}

	svc := ses.New(sess)

	return &Template{
		sess:   sess,
		svc:    svc,
		Sender: sender,
	}, nil

}

// Get All Templates by page number
func (t *Template) GetAllTemplates(page string) ([]*ses.TemplateMetadata, error) {
	var itemsFrom int64
	itemsFrom = (cast.ToInt64(page) * 10) - 9
	listTemplatesInput := ses.ListTemplatesInput{
		MaxItems: &itemsFrom,
	}

	listTemplatesOutput, err := t.svc.ListTemplates(&listTemplatesInput)
	if err != nil {
		return nil, err
	}
	return listTemplatesOutput.TemplatesMetadata, nil
}

// Get Template by name
func (t *Template) GetTemplateByName(name string) (interface{}, error) {
	templateOutput, err := t.svc.GetTemplate(&ses.GetTemplateInput{TemplateName: aws.String(name)})
	if err != nil {
		return nil, err
	}

	return SESTemplate{
		TemplateName: cast.ToString(templateOutput.Template.TemplateName),
		Subject:      cast.ToString(templateOutput.Template.SubjectPart),
		HTMLBody:     cast.ToString(templateOutput.Template.HtmlPart),
		TextBody:     cast.ToString(templateOutput.Template.TextPart),
	}, nil
}

// Create new html template
func (t *Template) CreateTemplate(template SESTemplate) error {
	createTemplateInput := &ses.CreateTemplateInput{
		Template: t.FormatTemplate(template),
	}
	_, err := t.svc.CreateTemplate(createTemplateInput)
	if err != nil {
		return err
	}

	return nil
}

// FormatTemplate
func (t *Template) FormatTemplate(template SESTemplate) *ses.Template {
	return &ses.Template{
		TemplateName: aws.String(template.TemplateName),
		SubjectPart:  aws.String(template.Subject),
		HtmlPart:     aws.String(template.HTMLBody),
		TextPart:     aws.String(template.TextBody),
	}
}

// UpdateTemplate
func (t *Template) UpdateTemplate(template SESTemplate) error {

	updateTemplateInput := &ses.UpdateTemplateInput{
		Template: t.FormatTemplate(template),
	}
	_, err := t.svc.UpdateTemplate(updateTemplateInput)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTemplateByName
func (t *Template) DeleteTemplateByName(name string) error {

	_, err := t.svc.DeleteTemplate(&ses.DeleteTemplateInput{
		TemplateName: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// SendEmailWithTemplate
func (t *Template) SendEmailWithTemplate(Recipient string, template string, data interface{}) (*ses.SendTemplatedEmailOutput, error) {

	// json marshal the data
	encode, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	input := &ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Source:       aws.String(t.Sender),
		Template:     aws.String(template),
		TemplateData: aws.String(string(encode)),
	}
	result, err := t.svc.SendTemplatedEmail(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
