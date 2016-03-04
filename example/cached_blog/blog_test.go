package blog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/cache"
	"github.com/ezbuy/ezorm/codec"
	"github.com/ezbuy/ezorm/db"
	"github.com/golang/groupcache"
	"gopkg.in/mgo.v2"
)

var (
	listenPort = ":8080"
)

func dbInit() {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	db.Setup(conf)

	http.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		k := r.URL.Query().Get("key")

		p, err := BlogMgr.FindByID(k)
		if err != nil {
			if err != mgo.ErrNotFound {
				panic(err.Error())
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		data, err := json.Marshal(p)
		if err != nil {
			panic(err)
		}
		w.Write(data)
	})

	http.HandleFunc("/blog/delete", func(w http.ResponseWriter, r *http.Request) {
		k := r.URL.Query().Get("key")

		err := BlogMgr.RemoveByID(k)
		if err != nil {
			panic(err.Error())
		}

		w.Write([]byte("ok"))
	})

	initCache("http://127.0.0.1:8001", []string{"http://127.0.0.1:8001", "http://127.0.0.1:8002"}, 64<<20)
	fmt.Printf("start listening on port [%s]\n", listenPort)
	go http.ListenAndServe(listenPort, nil)
}

func initCache(selfAddr string, peerAddrs []string, cacheBytes int64) {
	peers := groupcache.NewHTTPPool(selfAddr)
	group := groupcache.NewGroup("BlogCache", cacheBytes, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result, err := BlogMgr.FindByIDFromDB(key)
			if err != nil {
				return err
			}

			data, err := json.Marshal(result)
			dest.SetBytes((data))
			return nil
		}))

	peers.Set(peerAddrs...)

	go http.ListenAndServe(selfAddr, peers)
	codec := codec.NewJSONCodec()
	InitCache(cache.NewGroupCache(group, codec))
}

func TestBlog(t *testing.T) {
	dbInit()
	p := BlogMgr.NewBlog()
	p.Title = "I like ezorm"
	p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())
	_, err := p.Save()
	if err != nil {
		t.Fatal(err)
	}

	id := p.Id()

	b, err := getBlogByID(id)
	if err != nil {
		// handle error
		t.Fatal(err.Error())
	}

	fmt.Printf("get blog ok: %#v", b)

	b, err = getBlogByID(id)
	if err != nil {
		// handle error
		t.Fatal(err.Error())
	}
	fmt.Printf("get blog ok: %#v", b)
}

func getBlogByID(id string) (*Blog, error) {
	u := fmt.Sprintf("http://localhost%s/blog?key=%s", listenPort, id)
	println("request url", u)
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("http response error:%s\n", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}
	var b Blog
	err = json.Unmarshal(body, &b)
	return &b, err
}
