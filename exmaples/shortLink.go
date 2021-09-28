package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var tplContent = `
<html>
<head>
<title></title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/clipboard@2.0.8/dist/clipboard.min.js"></script>

<style>
.card {
	margin: 30px;
  }
</style>
</head>
<body>
<div class="card">
  <div class="card-body">
  <form action="/s" method="get" class="form-inline">
  <div class="form-group mx-sm-3">
    <input type="url" name="s" class="form-control" id="url" placeholder="url" style="width: 600px;">
  </div>
  <button type="submit" class="btn btn-primary">生成短链</button>
  </form>
  </div>
</div>
</body>
</html>
`

var textContent = `
<html>
<head>
<title></title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/clipboard@2.0.8/dist/clipboard.min.js"></script>
<style>
.card {
	margin: 30px;
  }
 textarea {
	overflow: scroll;
    min-height: 700px;
}
</style>
</head>
<body>
<div class="card">
  <div class="card-body">
  <form action="/t" method="get">
  <div class="form-group mx-sm-3">
    <textarea type="text" name="t" class="form-control" id="text" placeholder="text" style="width: 1000px;"></textarea>
  </div>
  <button type="submit" class="btn btn-primary" style="margin-left: 14px">生成短链</button>
  </form>
  </div>
</div>
</body>
</html>
`

var shortContent = `
<html>
  <head>
    <title></title>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css"></link>
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.bundle.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/clipboard@2.0.8/dist/clipboard.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/clipboard.js/2.0.0/clipboard.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
	<script type="text/javascript">
	  $(document).ready(function() {
	    new ClipboardJS('.btn');
	  });
	</script>
  </head>
  <body>
    <div class="page-content page-container" id="page-content">
      <div class="padding">
        <div class="row container d-flex justify-content-center">
          <div class="col-12 grid-margin">
            <div class="card">
              <div class="row">
                <div class="col-md-6">
                  <div class="card-body">
                    <p class="card-description">点击 Copy 复制短链内容到粘贴板</p>
                    <input type="text" id="clipboardExample1" class="form-control" value="{value}">
                    <div class="mt-3"> <button type="button" class="btn btn-info btn-clipboard" data-clipboard-action="copy" data-clipboard-target="#clipboardExample1">Copy</button></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </body>
</html>
`

var ctx = context.Background()

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, tplContent)
	} else {
		fmt.Println("url:", r.Form)
	}
}

func tindex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, textContent)
	} else {
		fmt.Println("url:", r.Form)
	}
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	http.HandleFunc("/index", index)
	http.HandleFunc("/e", tindex)
	http.HandleFunc("/t", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["t"]
		if !ok || len(keys[0]) < 1 {
			fmt.Fprintf(w, "Url Param 'text' is missing")
			return
		} else {
			text := keys[0]
			fmt.Println(text)
			encodeText := GetMD5Hash(text)
			text = base64.StdEncoding.EncodeToString([]byte(text))
			val, err := rdb.Get(ctx, encodeText).Result()
			if val != "" && err == nil {
				fmt.Fprintln(w, strings.Replace(shortContent, "{value}", "http://"+val, 1))
				return
			}
			host := r.Host
			shortName := RandStringRunes(6)
			dbKey := host + "/" + shortName
			for val, err := rdb.Get(ctx, dbKey).Result(); val != "" && err == nil; {
				shortName := RandStringRunes(6)
				dbKey = host + "/" + shortName
			}
			err = rdb.Set(ctx, dbKey, text, time.Hour*24*7).Err()
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, "Failed to generate short chain, please contact SRE.")
				return
			}
			err = rdb.Set(ctx, encodeText, dbKey, time.Hour*24*7).Err()
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, "Failed to generate short chain, please contact SRE.")
				return
			}
			fmt.Fprintln(w, strings.Replace(shortContent, "{value}", "http://"+r.Host+"/"+shortName, 1))
		}
	})
	http.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["s"]
		if !ok || len(keys[0]) < 1 {
			fmt.Fprintf(w, "Url Param 's' is missing")
			return
		} else {
			sourceURL := keys[0]
			encodeURL := base64.StdEncoding.EncodeToString([]byte(sourceURL))
			val, err := rdb.Get(ctx, encodeURL).Result()
			log.Println(val, err)
			if val != "" && err == nil {
				fmt.Fprintln(w, strings.Replace(shortContent, "{value}", "http://"+val, 1))
				return
			}
			host := r.Host
			shortName := RandStringRunes(6)
			dbKey := host + "/" + shortName
			for val, err := rdb.Get(ctx, dbKey).Result(); val != "" && err == nil; {
				shortName := RandStringRunes(6)
				dbKey = host + "/" + shortName
			}
			err = rdb.Set(ctx, dbKey, encodeURL, 0).Err()
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, "Failed to generate short chain, please contact SRE.")
				return
			}
			err = rdb.Set(ctx, encodeURL, dbKey, 0).Err()
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, "Failed to generate short chain, please contact SRE.")
				return
			}
			fmt.Fprintln(w, strings.Replace(shortContent, "{value}", "http://"+r.Host+"/"+shortName, 1))
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dbKey := r.Host + "/" + r.URL.Path[1:]
		val, err := rdb.Get(ctx, dbKey).Result()
		if val == "nil" || err != nil {
			fmt.Fprintf(w, "No matching url found")
			log.Println(err)
			return
		}
		content, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			fmt.Fprint(w, "content is corrupted")
			return
		} else {
			contentStr := string(content)
			if strings.HasPrefix(contentStr, "http") {
				http.Redirect(w, r, contentStr, 301)
			} else {
				fmt.Fprint(w, contentStr)
			}
		}

	})
	err := http.ListenAndServe("127.0.0.1:9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
