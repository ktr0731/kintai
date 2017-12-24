package main

import (
	"html/template"
	"net/http"
)

const tmpl = `
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <title>kintai</title>
</head>
<body>
  <h1>kintai</h1>
  <form method="post" action="/report" onsubmit="window.open('about:blank','_self').close()">
    <p>Description</p><input type="text" name="description" autofocus>
    <!--<p>Project</p><input type="text" name="project">-->
    <input type="submit" value="submit">
  </form>
</body>
</html>
`

type Report struct {
	Description, Project string
}

type handler struct {
	closeCh chan<- *Report
	t       *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		description, project := r.FormValue("description"), r.FormValue("project")
		h.closeCh <- &Report{Description: description, Project: project}
		w.WriteHeader(http.StatusCreated)
		return
	}

	h.t.Execute(w, nil)
}

func newHandler(closeCh chan<- *Report) (*handler, error) {
	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	return &handler{closeCh: closeCh, t: t}, nil
}
