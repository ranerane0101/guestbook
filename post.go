package guestbook

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
)

func Post(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "NO NAME"
	}

	//フォームから送信されたメッセージを取得
	message := r.FormValue("message")
	if message == "" {
		message = "NO MESSAGE"
	}

	//新しいメッセージのエンティティを作成
	msg := &Message{
		Name:      name,
		Message:   message,
		CreatedAt: time.Now(),
	}

	var key *datastore.Key
	k := r.FormValue("key")

	//キーが指定されていない場合、新しいキーを作成(新規作成か編集の条件分岐)
	if k == "" {
		key = datastore.IncompleteKey(r.Host, nil)
	} else {
		keyID, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		key = datastore.IDKey(r.Host, keyID, nil)
	}

	//データストアにメッセージを保存
	if _, err := client.Put(ctx, key, msg); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
