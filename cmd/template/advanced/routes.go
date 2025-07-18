package advanced

import (
	_ "embed"
)

//go:embed files/htmx/hello.templ.tmpl
var helloTemplTemplate []byte

//go:embed files/htmx/base.templ.tmpl
var baseTemplTemplate []byte

//go:embed files/react/tailwind/index.css.tmpl
var inputCssTemplateReact []byte

//go:embed files/react/tailwind/vite.config.ts.tmpl
var viteTailwindConfigFile []byte

//go:embed files/react/tailwind/app.tsx.tmpl
var reactTailwindAppFile []byte

//go:embed files/react/app.tsx.tmpl
var reactAppFile []byte

//go:embed files/tailwind/input.css.tmpl
var inputCssTemplate []byte

//go:embed files/tailwind/output.css.tmpl
var outputCssTemplate []byte

//go:embed files/htmx/htmx.min.js.tmpl
var htmxMinJsTemplate []byte

//go:embed files/htmx/efs.go.tmpl
var efsTemplate []byte

//go:embed files/htmx/hello.go.tmpl
var helloGoTemplate []byte

//go:embed files/htmx/hello_fiber.go.tmpl
var helloFiberGoTemplate []byte

//go:embed files/htmx/routes/http_router.tmpl
var httpRouterHtmxTemplRoutes []byte

//go:embed files/htmx/routes/standard_library.tmpl
var stdLibHtmxTemplRoutes []byte

//go:embed files/htmx/imports/standard_library.tmpl
var stdLibHtmxTemplImports []byte

//go:embed files/websocket/imports/standard_library.tmpl
var stdLibWebsocketImports []byte

//go:embed files/htmx/routes/chi.tmpl
var chiHtmxTemplRoutes []byte

//go:embed files/htmx/routes/gin.tmpl
var ginHtmxTemplRoutes []byte

//go:embed files/htmx/imports/gin.tmpl
var ginHtmxTemplImports []byte

//go:embed files/htmx/routes/gorilla.tmpl
var gorillaHtmxTemplRoutes []byte

//go:embed files/htmx/routes/echo.tmpl
var echoHtmxTemplRoutes []byte

//go:embed files/htmx/routes/fiber.tmpl
var fiberHtmxTemplRoutes []byte

//go:embed files/htmx/imports/fiber.tmpl
var fiberHtmxTemplImports []byte

//go:embed files/websocket/imports/fiber.tmpl
var fiberWebsocketTemplImports []byte

//go:embed files/kafka/consumer.go.tmpl
var kafkaConsumerTemplate []byte

//go:embed files/kafka/consumer_test.go.tmpl
var kafkaConsumerTestTemplate []byte

//go:embed files/kafka/env.tmpl
var kafkaEnvTemplate []byte

//go:embed files/kafka/cmd/consumer/main.go.tmpl
var kafkaConsumerMainTemplate []byte

// Helper functions for React files
//
//go:embed files/worker/cmd/worker/main.go.tmpl
var workerMainTemplate []byte

//go:embed files/worker/cmd/worker/tasks/hello_world_task.go.tmpl
var workerHelloWorldTaskTemplate []byte

//go:embed files/worker/env.tmpl
var workerEnvTemplate []byte

// Swagger template files
//
//go:embed files/swagger/imports/standard_library.tmpl
var stdLibSwaggerImports []byte

//go:embed files/swagger/routes/standard_library.tmpl
var stdLibSwaggerRoutes []byte

//go:embed files/swagger/imports/gin.tmpl
var ginSwaggerImports []byte

//go:embed files/swagger/routes/gin.tmpl
var ginSwaggerRoutes []byte

//go:embed files/swagger/imports/echo.tmpl
var echoSwaggerImports []byte

//go:embed files/swagger/routes/echo.tmpl
var echoSwaggerRoutes []byte

//go:embed files/swagger/imports/fiber.tmpl
var fiberSwaggerImports []byte

//go:embed files/swagger/routes/fiber.tmpl
var fiberSwaggerRoutes []byte

//go:embed files/swagger/imports/chi.tmpl
var chiSwaggerImports []byte

//go:embed files/swagger/routes/chi.tmpl
var chiSwaggerRoutes []byte

//go:embed files/swagger/imports/gorilla.tmpl
var gorillaSwaggerImports []byte

//go:embed files/swagger/routes/gorilla.tmpl
var gorillaSwaggerRoutes []byte

//go:embed files/swagger/routes/http_router.tmpl
var httpRouterSwaggerRoutes []byte

//go:embed files/swagger/README_SWAGGER.md
var swaggerReadme []byte

// Helper functions for React files
func EchoHtmxTemplRoutesTemplate() []byte {
	return echoHtmxTemplRoutes
}

func GorillaHtmxTemplRoutesTemplate() []byte {
	return gorillaHtmxTemplRoutes
}

func ChiHtmxTemplRoutesTemplate() []byte {
	return chiHtmxTemplRoutes
}

func GinHtmxTemplRoutesTemplate() []byte {
	return ginHtmxTemplRoutes
}

func HttpRouterHtmxTemplRoutesTemplate() []byte {
	return httpRouterHtmxTemplRoutes
}

func StdLibHtmxTemplRoutesTemplate() []byte {
	return stdLibHtmxTemplRoutes
}

func StdLibHtmxTemplImportsTemplate() []byte {
	return stdLibHtmxTemplImports
}

func StdLibWebsocketTemplImportsTemplate() []byte {
	return stdLibWebsocketImports
}

func HelloTemplTemplate() []byte {
	return helloTemplTemplate
}

func BaseTemplTemplate() []byte {
	return baseTemplTemplate
}

func ReactTailwindAppfile() []byte {
	return reactTailwindAppFile
}

func ReactAppfile() []byte {
	return reactAppFile
}

func InputCssTemplateReact() []byte {
	return inputCssTemplateReact
}

func ViteTailwindConfigFile() []byte {
	return viteTailwindConfigFile
}

func InputCssTemplate() []byte {
	return inputCssTemplate
}

func OutputCssTemplate() []byte {
	return outputCssTemplate
}

func HtmxJSTemplate() []byte {
	return htmxMinJsTemplate
}

func EfsTemplate() []byte {
	return efsTemplate
}

func HelloGoTemplate() []byte {
	return helloGoTemplate
}

func HelloFiberGoTemplate() []byte {
	return helloFiberGoTemplate
}

func FiberHtmxTemplRoutesTemplate() []byte {
	return fiberHtmxTemplRoutes
}

func FiberHtmxTemplImportsTemplate() []byte {
	return fiberHtmxTemplImports
}

func FiberWebsocketTemplImportsTemplate() []byte {
	return fiberWebsocketTemplImports
}

func GinHtmxTemplImportsTemplate() []byte {
	return ginHtmxTemplImports
}

func KafkaConsumerTemplate() []byte {
	return kafkaConsumerTemplate
}

func KafkaConsumerTestTemplate() []byte {
	return kafkaConsumerTestTemplate
}

func KafkaEnvTemplate() []byte {
	return kafkaEnvTemplate
}

func KafkaConsumerMainTemplate() []byte {
	return kafkaConsumerMainTemplate
}

func WorkerMainTemplate() []byte {
	return workerMainTemplate
}

func WorkerHelloWorldTaskTemplate() []byte {
	return workerHelloWorldTaskTemplate
}

func WorkerEnvTemplate() []byte {
	return workerEnvTemplate
}

// Swagger template functions
func StdLibSwaggerImportsTemplate() []byte {
	return stdLibSwaggerImports
}

func StdLibSwaggerRoutesTemplate() []byte {
	return stdLibSwaggerRoutes
}

func GinSwaggerImportsTemplate() []byte {
	return ginSwaggerImports
}

func GinSwaggerRoutesTemplate() []byte {
	return ginSwaggerRoutes
}

func EchoSwaggerImportsTemplate() []byte {
	return echoSwaggerImports
}

func EchoSwaggerRoutesTemplate() []byte {
	return echoSwaggerRoutes
}

func FiberSwaggerImportsTemplate() []byte {
	return fiberSwaggerImports
}

func FiberSwaggerRoutesTemplate() []byte {
	return fiberSwaggerRoutes
}

func ChiSwaggerImportsTemplate() []byte {
	return chiSwaggerImports
}

func ChiSwaggerRoutesTemplate() []byte {
	return chiSwaggerRoutes
}

func GorillaSwaggerImportsTemplate() []byte {
	return gorillaSwaggerImports
}

func GorillaSwaggerRoutesTemplate() []byte {
	return gorillaSwaggerRoutes
}

func HttpRouterSwaggerRoutesTemplate() []byte {
	return httpRouterSwaggerRoutes
}

func SwaggerReadmeTemplate() []byte {
	return swaggerReadme
}
