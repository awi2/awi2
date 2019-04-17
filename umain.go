package main

import (
	"code.google.com/p/gorilla/mux"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"log"
	"strconv"
	"strings"
	"time"
	"flag"
	"github.com/boltdb/bolt"
)


var sw = make(map[string]string, 4)

var buck = "tracks"



var posts = ""
var cnt = 20
var last, old string

var localRun = "0"
var user = "0"
var db *bolt.DB

func init() {
// Initialize db.

	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	db.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists([]byte("flats"))
		_, err1 = tx.CreateBucketIfNotExists([]byte("cats"))
		if err1 != nil {
			return fmt.Errorf("creating bucket <?>: %s", err1)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	flag.Parse()
	if flag.Arg(0) == "0" {
		localRun = "1"
		allowedIP = "127.0.0.1"
	}
}


func Log(v ...interface{}) {
	fmt.Println(v)
}

func main() {

	mainLoop()
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {

	titles := ""
	fmt.Fprintf(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p class=\"b\"><strong>List Of Titled Topics</strong></p><br>"+titles+
		"</form></html>")
	return
}

func handlerAllPages(w http.ResponseWriter, r *http.Request) {

	resu := ReadAllFlats()
	fmt.Fprintf(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p class=\"b\"><strong>List Of All Pages</strong></p><br>"+resu+
		"</form></html>")
	return
}



func handlerFront(w http.ResponseWriter, r *http.Request) {

	titles := ""
	fmt.Fprintf(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p class=\"b\"><strong>List Of All Pages</strong></p><br>"+titles+
		"</form></html>")
	return
}

func handlerBrokenTrack(w http.ResponseWriter, r *http.Request) {

	thread := mux.Vars(r)["thr"]
	http.Redirect(w, r, "/page/"+thread, http.StatusFound)
	return
}

func handlerMegaBrokenTrack(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/page/0", http.StatusFound)
	return
}



func handlerShowIndex(w http.ResponseWriter, r *http.Request) {

	user := ""
	topic := getFlat("index")


//	fmt.Println(topic)
//	fmt.Println("user == ", user)

	if user == "Anon" || user == "" {
		user = "<em>~ Guest ~</em>"
	}
	status := "<a href=\"/+\"> ~ "+user+" ~</a>"


	fmt.Fprint(w, "<html><head>"+//<title>"+ttl+"</title>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body>"+
		"<br><div class=\"r\"><p>"+Markdown(topic)+
		"</p></div><br>"+
		"<div class=\"b\"><a href=\"/page/0\">Home</a>&nbsp;&nbsp;[ "+status+
		"&nbsp;&nbsp;] <a href=\"/index\">Index</a></p></body></html>")

	return
}



func handlerFlatShow(w http.ResponseWriter, r *http.Request) {

	by := make([]byte, 22)
	rand.Read(by)
	challenge := hex.EncodeToString(by)
	ref := r.Referer()
	fid := mux.Vars(r)["flt"]
	fmt.Println("fid= ", fid)



	topic := getFlat(fid)
	mark := "."


	ttl := strings.ToUpper(fid)
	user = "Anon"
	raddr := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Println("raddr: ", raddr)
	if raddr == allowedIP { user = "0" }

	if topic != "" { mark = string(topic[0]) }

	if mark == "@" && user != "0" {
		http.Redirect(w, r, "/page/oops", http.StatusFound)
///		return
	}
	fmt.Println("mark=", mark)
	fmt.Println("user=", user)



	if len(topic) == 0 {
		topic = " "
	}
	fmt.Println("referer=", ref)
	replybox := "<textarea name=\"fxt\" id=\"rb\" cols=\"80\" rows=\"30\">"+topic+"</textarea>" +
	"<br><br><input type=\"submit\" value=\"Save\"><br>"
	fmt.Fprint(w, "<html><head><title>"+ttl+"</title>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"function toggle(id) {\n"+
		"var e = document.getElementById(id);\n"+
		"if(e.style.display == 'none')\n"+
		"e.style.display = 'block';\n"+
		"else\n"+
		"e.style.display = 'none';\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body>"+
		"<p id=\"d\">&nbsp;<a href=\"/all\">All_Pages </a>&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"/del/"+fid+"\">&nbsp;&nbsp;&nbsp;Delete Page </a></p>"+
		"<p id=\"z\"><a href=\"#zz\">"+
		"<input type=\"submit\" value=\"Toggle\" onclick=\"toggle('tx');\"></a></p>"+
		"<div class=\"r\" id=\"tx\"><p>"+Markdown(topic)+
		"</p></div><p id=\"z\">"+
		"<a href=\"#z\"><input type=\"submit\" value=\"Edit\" onclick=\"toggle('rbox');\"></a></p>"+
		"<br><a href=\"/print/"+fid+"\">"+
		"Print Flat Page</a>"+
		"<p></p><div id=\"rbox\"><form action=\"/fsubmit\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		replybox+
		"<input type=\"hidden\" name=\"ch\" value=\""+challenge+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<input type=\"hidden\" name=\"user\" value=\""+user+"\">"+
		"<input type=\"hidden\" name=\"flt\" value=\""+fid+"\"></form>"+
		"<br></body></html>")

	return
}



func handlerPrintShow(w http.ResponseWriter, r *http.Request) {

	fid := mux.Vars(r)["flt"]
	fmt.Println("fid= ", fid)

	topic := getFlat(fid)

	if len(topic) == 0 {
		topic = " "
	}
	t := "<html><head><title>"+fid+"</title>"+
	"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
	"</head><body>"+
	"<div class=\"r\" id=\"tx\"><p>"+Markdown(topic)+
	"</p></div>"
	tt := "<p id=\"u\"><a href=\"/page/"+fid+"\">Go Back</a></p>"+
	"<br></body></html>"
	savePrint(fid, t+"</body></html>")
	fmt.Fprint(w, t + tt)

	return
}


func handlerFlatSubmit(w http.ResponseWriter, r *http.Request) {
	txt := r.FormValue("fxt")
////	fmt.Println("flat txt from form=", txt)
	fid := r.FormValue("flt")
	fmt.Println("fid from form=", fid)
	mlen, _ := strconv.Atoi(r.FormValue("mlen"))
	owner := r.FormValue("own")
	if len(txt) < 1 {
		http.Redirect(w, r, "/page/"+fid, http.StatusFound)
		return
	}
	
	if txt[0] == '$' {
		bytetxt := []byte(txt)
		bytetxt = bytetxt[1:]
		txt = "<em>©</em><br><br>" + string(bytetxt)
	}
	user := ""
	fmt.Println("owner from form=", owner)

	if user == "" {
		user = "Anon"
	}
	if localRun == "1" { user = "0" }
	if user == "DrEvil" { goto Adm2 }
	if user == "9" { goto Adm2 }
Adm2:
	meta := ""

	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if mlen != 0 {
		meta = "<a id=\"u\" href=\"/page/"+fid+"/+\">Edit</a></div><hr>"
	}
	if mlen == 0 {
		meta = ""
	}
	txt = meta + txt

	

	err = db.Update(func(tx *bolt.Tx) error {
		bucket1 := tx.Bucket([]byte("flats"))

		err = bucket1.Put([]byte(fid), []byte(txt))

		if err != nil {
			return err
		}
		return nil
	})
	db.Close()
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/page/"+fid, http.StatusFound)
	return
}

func handlerFlatEdit(w http.ResponseWriter, r *http.Request) {
	by := make([]byte, 22)
	rand.Read(by)
	challenge := hex.EncodeToString(by)

	user := ""
	fid := mux.Vars(r)["fid"]
//	postnum := mux.Vars(r)["postnum"]


	metapost := ReadFlatPage(fid)
	if len(metapost) < 26 { return }
	meta := []byte(metapost)[25:]
	
	owner := ""
	for i := 0; i < len(meta); i++ {
		if meta[i] > 62 {
			owner += string(meta[i])
		}
		if meta[i] < 60 {
			owner += string(meta[i])
		}
		if meta[i] == 60 {
			break
		}
	}
	fmt.Println("owner", owner)
	fmt.Println("user in TrackEdut()", user)
	var metalen, postbody string
	for i := 0; i < len(meta); i++ {
		if meta[i] == 114 && meta[i+1] == 62 {
			metalen = strconv.Itoa(len(meta[:i+2]))
			postbody = string(meta[i+2:])
			break
		}
	}
	title := "Editing page : " + fid
	fmt.Fprint(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />\n"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body><p class=\"b\"><em><a>"+title+"</a></em></p>"+
		"<br><form action=\"/fsubmit\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		"<textarea name=\"fxt\" cols=\"80\" rows=\"8\" value=\""+postbody+"\"></textarea><br>"+
		"<p><input type=\"submit\" value=\"Save\"></p>"+
		"<input type=\"hidden\" name=\"ch\" value=\""+challenge+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<input type=\"hidden\" name=\"own\" value=\""+owner+"\">"+
		"<input type=\"hidden\" name=\"flt\" value=\""+fid+"\">"+
		"<input type=\"hidden\" name=\"mlen\" value=\""+metalen+"\">"+
		"</form></body></html>")

	return
}

func handlerFlatDelete(w http.ResponseWriter, r *http.Request) {
	fid := mux.Vars(r)["flt"]
	fmt.Println("fid from form=", fid)


	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})

	err = db.Update(func(tx *bolt.Tx) error {
		bucket1 := tx.Bucket([]byte("flats"))

		err = bucket1.Delete([]byte(fid))
		fmt.Println("Db error=", err)
		if err != nil {
			return err
		}
		return nil
	})
	db.Close()
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/all", http.StatusFound)
	return
}


func mainLoop() {
	if localRun == "1" { ipPort = "127.0.0.1:80" }
	srv := &http.Server{
		Addr: ipPort,
//		Addr: forumIp + ":443",
		//		Addr: ":",
//		ReadTimeout: time.Duration(2) * time.Second,
		//		TLSConfig *tls.Config,
	}



	r := mux.NewRouter()

	r.HandleFunc("/", handlerShowIndex)
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("s"))))

	r.HandleFunc("/del/{flt}", handlerFlatDelete)
	r.HandleFunc("/p/{fid}/+", handlerFlatEdit)
	r.HandleFunc("/print/{flt}", handlerPrintShow)
	r.HandleFunc("/page/{flt}", handlerFlatShow)
	r.HandleFunc("/page/{flt}/", handlerFlatShow)
	r.HandleFunc("/fsubmit", handlerFlatSubmit)
	r.HandleFunc("/all", handlerAllPages)
	http.Handle("/", r)
	Log("Awi_2 ("+revision+") is running on " + srv.Addr)
//	srv.ListenAndServeTLS("site.crt", "site.key")
	srv.ListenAndServe()
}
