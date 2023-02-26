package server

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/ricoschulte/go-myapps/sysclient"
	log "github.com/sirupsen/logrus"
)

//go:embed templates/*
var TemplateFS embed.FS

func GetServerMux(sysclient *sysclient.Sysclient, fs http.FileSystem) *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files on app path
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Serve a file for the admin.htm path
	mux.HandleFunc("/admin.xml", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).Info("HTTP Request")

		w.WriteHeader(http.StatusOK)
		err := RenderTemplate(w, "templates/admin.html", sysclient)
		if err != nil {
			log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).Errorf("error handling request: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/CMD0/mod_cmd.xml", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")

		var body []byte
		r.Body.Read(body)

		for key, value := range r.URL.Query() {

			log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField(key, value).Debug("POST body parameter")
		}

		response_text := "ok\n"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response_text))
	})

	mux.HandleFunc("/!config activate", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!config write", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/LOG0/mod_cmd.xml", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/MEDIA/mod_cmd.xml", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/PBX0/mod_cmd.xml", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	/*

	   WARNING[2023-02-26T06:28:48+01:00] HTTP Request 404                              id=f19033480af4 path=/PBX0/ADMIN/mod_cmd_login.xml query="map[cmd:[show] reg:[*]]"
	   WARNING[2023-02-26T06:28:48+01:00] HTTP Request 404                              id=f19033480af4 path=/PBX0/ADMIN/mod_cmd_login.xml query="map[cd:[pbx.csv] cmd:[download] format:[csv]]"
	   WARNING[2023-02-26T06:28:49+01:00] HTTP Request 404                              id=f19033480af4 path=/PBX0/ADMIN/mod_cmd_login.xml query="map[cd:[pbx.xml] cmd:[download] format:[xml]]"
	*/
	mux.HandleFunc("/PBX0/ADMIN/mod_cmd_login.xml", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/cfg.txt", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/cfg-standard.txt", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	/*
	   WARNING[2023-02-26T06:28:47+01:00] HTTP Request 404                              id=f19033480af4 path=/LOG0/FAULT/mod_cmd.xml query="map[cmd:[xml-alarms]]"
	   WARNING[2023-02-26T06:28:47+01:00] HTTP Request 404                              id=f19033480af4 path=/LOG0/FAULT/mod_cmd.xml query="map[cmd:[xml-faults] xsl:[fault_log.xsl]]"
	*/
	mux.HandleFunc("/LOG0/FAULT/mod_cmd.xml", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/log.txt", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!buf", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mod", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mod cmd CPU mips-usage", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mod CMD FLASHMAN0 info", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mod cmd IP0 tcp-sockets", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mod cmd IP0 udp-sockets", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mem", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mem info tcp_socket", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/!mem info udp_socket", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// catch not existing path with error log
		if r.URL.Path != "/" {
			response_text := fmt.Sprintf("Path %s not found", r.URL.Path)
			log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Warning("HTTP Request 404")

			//content-length: 54138
			w.Header().Add("Content-Length", fmt.Sprint(len(response_text)))
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(response_text))

			return
		}
		log.WithField("id", sysclient.Identity.Id).WithField("path", r.URL.Path).WithField("query", r.URL.Query()).Info("HTTP Request")

		response_text := fmt.Sprintf("index page on %s", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response_text))

	})

	return mux
}

func RenderTemplate(w io.Writer, templateFile string, env any) error {
	log.WithField("templateFile", templateFile).Debug("render template")

	tpllib, err_read_lib := TemplateFS.ReadFile("templates/main.html")
	if err_read_lib != nil {
		log.WithField("templateFile", templateFile).Errorf("error while reading template: %v", err_read_lib)
		return err_read_lib
	}

	tpl, err_read := TemplateFS.ReadFile(templateFile)
	if err_read != nil {
		log.WithField("templateFile", templateFile).Errorf("error while reading template: %v", err_read_lib)
		return err_read
	}

	t, _ := template.New("tpl").Parse(string(tpl))
	t.New("tpllib").Parse(string(tpllib))
	err_write := t.Execute(w, env)
	if err_write != nil {
		log.WithField("templateFile", templateFile).Errorf("error while writing template to response: %v", err_read_lib)
		return err_write
	}
	return nil
}
