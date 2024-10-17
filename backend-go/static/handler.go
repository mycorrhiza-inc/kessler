package static

import (
	"text/template"

	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func RenderStaticSitemap(dbtx_val dbstore.DBTX) {
	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	q := dbstore.New(dbtx_val)
}
