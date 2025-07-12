# Swagger Documentation

This project includes Swagger documentation for the API endpoints.

## Usage

1. **Add Swagger annotations to your handlers**

Example for a Hello World handler:

```go
// HelloWorldHandler returns a hello world message
// @Summary      Show hello world
// @Description  get hello world message
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       / [get]
func (s *Server) HelloWorldHandler(c *gin.Context) {
    resp := make(map[string]string)
    resp["message"] = "Hello World"
    c.JSON(http.StatusOK, resp)
}
```

2. **Generate Swagger documentation**

Run the following command to generate documentation:

```bash
swag init
```

This will generate the `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml` files.

3. **Access Swagger UI**

Once your server is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Configuration

### Base URL

You can configure the base URL for your API documentation by setting the `API_BASE_URL` environment variable:

```bash
export API_BASE_URL="api.example.com"
```

Or in your `.env` file:

```
API_BASE_URL=api.example.com
```

When set, this will override the default `localhost:8080` host in your Swagger documentation, making it useful for production deployments.

## Available Swagger Comments

- `@Summary`: A short summary of the API
- `@Description`: A longer description of the API
- `@Tags`: Tags for grouping operations
- `@Accept`: The MIME types that the API can consume
- `@Produce`: The MIME types that the API can produce
- `@Param`: Parameters for the API
- `@Success`: A successful response
- `@Failure`: An error response
- `@Router`: The route definition

For more information about Swagger annotations, visit: https://github.com/swaggo/swag
