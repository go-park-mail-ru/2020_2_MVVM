package pdfGenerator

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"io/ioutil"
	"os"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}

//new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

//parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

//generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {
	// write whole the body
	file, err := ioutil.TempFile("/tmp", "temp-resume.*.html")
	if err != nil {
		return false, err
	}
	defer os.Remove(file.Name())
	fileName := file.Name()

	err1 := ioutil.WriteFile(fileName, []byte(r.body), 0644)
	if err1 != nil {
		return false, err1
	}

	f, err := os.Open(fileName)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return false, err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return false, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return false, err
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		return false, err
	}
	return true, nil
}
