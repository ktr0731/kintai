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
  <title></title>
</head>
<body>
  <h1>kintai</h1>
  <form method="post" action="/report">
    <p>Description</p>
    <p>Project</p>
  </form>
</body>
</html>
`

type handler struct {
	closeCh chan<- struct{}
	t       *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.closeCh <- struct{}{}
		w.WriteHeader(http.StatusCreated)
		return
	}

	h.t.Execute(w, nil)
}

func newHandler(closeCh chan<- struct{}) (*handler, error) {
	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	return &handler{closeCh: closeCh, t: t}, nil
}
