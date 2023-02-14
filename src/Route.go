package src

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func SetRoute(r *gin.Engine, FS embed.FS) {
	tmpl := template.Must(template.New("").ParseFS(FS, "static/index.html"))
	r.SetHTMLTemplate(tmpl)

	fStatic, _ := fs.Sub(FS, "static")

	r.StaticFS("/static", http.FS(fStatic))
	r.GET("/", index)

	r.GET("/chat", index)
	r.GET("/ws", HandleWS)
	if err := r.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
