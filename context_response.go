package gin

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	contractshttp "github.com/goravel/framework/contracts/http"
)

type ContextResponse struct {
	instance *gin.Context
	origin   contractshttp.ResponseOrigin
}

func NewContextResponse(instance *gin.Context, origin contractshttp.ResponseOrigin) *ContextResponse {
	return &ContextResponse{instance, origin}
}

func (r *ContextResponse) Data(code int, contentType string, data []byte) contractshttp.Response {
	return &DataResponse{code, contentType, data, r.instance}
}

func (r *ContextResponse) Download(filepath, filename string) contractshttp.Response {
	return &DownloadResponse{filename, filepath, r.instance}
}

func (r *ContextResponse) File(filepath string) contractshttp.Response {
	return &FileResponse{filepath, r.instance}
}

func (r *ContextResponse) Header(key, value string) contractshttp.ContextResponse {
	r.instance.Header(key, value)

	return r
}

func (r *ContextResponse) Json(code int, obj any) contractshttp.Response {
	return &JsonResponse{code, obj, r.instance}
}

func (r *ContextResponse) Origin() contractshttp.ResponseOrigin {
	return r.origin
}

func (r *ContextResponse) Redirect(code int, location string) contractshttp.Response {
	return &RedirectResponse{code, location, r.instance}
}

func (r *ContextResponse) String(code int, format string, values ...any) contractshttp.Response {
	return &StringResponse{code, format, r.instance, values}
}

func (r *ContextResponse) Success() contractshttp.ResponseSuccess {
	return NewGinSuccess(r.instance)
}

func (r *ContextResponse) Status(code int) contractshttp.ResponseStatus {
	return NewStatus(r.instance, code)
}

func (r *ContextResponse) View() contractshttp.ResponseView {
	return NewView(r.instance)
}

func (r *ContextResponse) Writer() http.ResponseWriter {
	return r.instance.Writer
}

func (r *ContextResponse) Flush() {
	r.instance.Writer.Flush()
}

type Success struct {
	instance *gin.Context
}

func NewGinSuccess(instance *gin.Context) contractshttp.ResponseSuccess {
	return &Success{instance}
}

func (r *Success) Data(contentType string, data []byte) contractshttp.Response {
	return &DataResponse{http.StatusOK, contentType, data, r.instance}
}

func (r *Success) Json(obj any) contractshttp.Response {
	return &JsonResponse{http.StatusOK, obj, r.instance}
}

func (r *Success) String(format string, values ...any) contractshttp.Response {
	return &StringResponse{http.StatusOK, format, r.instance, values}
}

type Status struct {
	instance *gin.Context
	status   int
}

func NewStatus(instance *gin.Context, code int) contractshttp.ResponseSuccess {
	return &Status{instance, code}
}

func (r *Status) Data(contentType string, data []byte) contractshttp.Response {
	return &DataResponse{r.status, contentType, data, r.instance}
}

func (r *Status) Json(obj any) contractshttp.Response {
	return &JsonResponse{r.status, obj, r.instance}
}

func (r *Status) String(format string, values ...any) contractshttp.Response {
	return &StringResponse{r.status, format, r.instance, values}
}

func ResponseMiddleware() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		blw := &BodyWriter{body: bytes.NewBufferString("")}
		switch ctx := ctx.(type) {
		case *Context:
			blw.ResponseWriter = ctx.Instance().Writer
			ctx.Instance().Writer = blw
		}

		ctx.WithValue("responseOrigin", blw)
		ctx.Request().Next()
	}
}

type BodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func (w *BodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)

	return w.ResponseWriter.WriteString(s)
}

func (w *BodyWriter) Body() *bytes.Buffer {
	return w.body
}

func (w *BodyWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
