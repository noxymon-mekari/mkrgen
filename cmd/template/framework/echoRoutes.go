package framework

import (
	_ "embed"

	"github.com/noxymon-mekari/mkrgen/cmd/template/advanced"
)

//go:embed files/routes/echo.go.tmpl
var echoRoutesTemplate []byte

//go:embed files/tests/echo-test.go.tmpl
var echoTestHandlerTemplate []byte

// EchoTemplates contains the methods used for building
// an app that uses [github.com/labstack/echo]
type EchoTemplates struct{}

func (e EchoTemplates) Main() []byte {
	return mainTemplate
}
func (e EchoTemplates) Server() []byte {
	return standardServerTemplate
}

func (e EchoTemplates) Routes() []byte {
	return echoRoutesTemplate
}

func (e EchoTemplates) TestHandler() []byte {
	return echoTestHandlerTemplate
}

func (e EchoTemplates) HtmxTemplImports() []byte {
	return advanced.StdLibHtmxTemplImportsTemplate()
}

func (e EchoTemplates) HtmxTemplRoutes() []byte {
	return advanced.EchoHtmxTemplRoutesTemplate()
}

func (e EchoTemplates) WebsocketImports() []byte {
	return advanced.StdLibWebsocketTemplImportsTemplate()
}

func (e EchoTemplates) SwaggerImports() []byte {
	return advanced.EchoSwaggerImportsTemplate()
}

func (e EchoTemplates) SwaggerRoutes() []byte {
	return advanced.EchoSwaggerRoutesTemplate()
}
