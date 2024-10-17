package static

import "text/template"

func RenderStaticSitemap() {
	tmpl := template.Must(template.ParseFiles("templates/post.html"))
}
