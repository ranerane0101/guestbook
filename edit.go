package guestbook

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
)

var editTmpl = template.Must(template.ParseFiles("./view/edit.html"))

type EditTemplate struct {
	Title   string
	Name    string
	Message string
	ID      int64
}

func Edit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	//Datastoreのクライアントを作成する
	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	//リクエストからキー情報を取得する
	k := r.FormValue("key")
	keyID, err := strconv.ParseInt(k, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	var msg Message
	//キーを作成し、データを取得する
	key := datastore.IDKey(r.Host, keyID, nil)
	if err = client.Get(ctx, key, &msg); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	edt := &EditTemplate{
		Title:   title,
		Name:    msg.Name,
		Message: msg.Message,
		ID:      keyID,
	}

	if err := editTmpl.Execute(w, edt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
