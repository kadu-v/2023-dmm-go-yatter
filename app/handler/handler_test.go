package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/config"
	"yatter-backend-go/app/dao"

	"github.com/stretchr/testify/assert"
)

/*---------- Test for the handler of account ----------------------------------*/
func TestAccountRegistration(t *testing.T) {
	c := setup(t)
	defer c.Close()

	func() {
		resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()

	func() {
		resp, err := c.Get("/v1/accounts/john")
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()
}

func TestAccountFetch(t *testing.T) {
	c := setup(t)
	defer c.Close()

	username := "john"

	// ユーザーjohnを追加
	resp, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	// 登録したユーザーを取得
	resp, err = c.Get(fmt.Sprintf("/v1/accounts/%s", username))
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var j map[string]interface{}
	if assert.NoError(t, json.Unmarshal(body, &j)) {
		assert.Equal(t, "john", j["username"])
		assert.Equal(t, "john", j["display_name"])
	}
}

/*---------- Test for the handler of the status -------------------------------*/
func TestStatusPost(t *testing.T) {
	c := setup(t)
	defer c.Close()

	username := "john"
	content := `{
			"status": "ピタ ゴラ スイッチ♪",
			"medias": [
			  {
				"media_id": 123,
				"description": "hoge hoge"
			  }
			]
		  }`

	// ユーザーjohnを追加
	resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	// 認証つきでStatusをPost
	resp, err = c.PostJsonWithAuthentication("/v1/statuses", username, content)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var j map[string]interface{}
	if assert.NoError(t, json.Unmarshal(body, &j)) {
		assert.Equal(t, "ピタ ゴラ スイッチ♪", j["content"])
	}
}

func TestStatusFetch(t *testing.T) {
	c := setup(t)
	defer c.Close()

	func() {
		username := "john"
		content := `{
			"status": "ピタ ゴラ スイッチ♪",
			"medias": [
			  {
				"media_id": 123,
				"description": "hoge hoge"
			  }
			]
		  }`

		// ユーザーjohnを追加
		resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
			return
		}

		// 認証つきでStatusをPost
		resp, err = c.PostJsonWithAuthentication("/v1/statuses", username, content)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
			return
		}

		// PostしたStatusをGet
		resp, err = c.Get("/v1/statuses/1")
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "ピタ ゴラ スイッチ♪", j["content"])
		}
	}()
}

/*---------- Test for the handler of the public timeline ----------------------*/
func TestTimelineFetch(t *testing.T) {
	c := setup(t)
	defer c.Close()

	username := "john"
	content := `{
			"status": "ピタ ゴラ スイッチ♪",
			"medias": [
			  {
				"media_id": 123,
				"description": "hoge hoge"
			  }
			]
		  }`

	// ユーザーjohnを追加
	resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	// 認証つきでStatusを40件Post
	for i := 0; i < 40; i++ {
		resp, err = c.PostJsonWithAuthentication("/v1/statuses", username, content)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
			return
		}
	}

	// PostしたStatusをGet
	queryPrams := map[string]string{"only_media": "false", "max_id": "100", "since_id": "0", "limit": "40"}
	resp, err = c.GetWithQueryParams("/v1/timelines/public", queryPrams)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var j map[string]([](map[string]interface{}))
	if assert.NoError(t, json.Unmarshal(body, &j)) {
		k := "Statuses"
		assert.Equal(t, 40, len(j[k]))
		assert.Equal(t, "ピタ ゴラ スイッチ♪", j[k][0]["content"])
	}

}

/*---------- Auxiliary functions for tests ---------------------------------*/
func setup(t *testing.T) *C {
	db, err := dao.NewDB(config.MySQLConfig())
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
			log.Fatalln(err)
		}
	}()

	for _, table := range []string{"account", "status"} {
		if _, err := db.Exec("TRUNCATE TABLE " + table); err != nil {
			log.Fatalln(err)
		}
	}
	server := httptest.NewServer(NewRouter(
		dao.NewAccount(db), dao.NewStatus(db),
	))

	return &C{
		Server: server,
	}
}

type C struct {
	Server *httptest.Server
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

func (c *C) PostJsonWithAuthentication(apiPath string, auth string, payload string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.asURL(apiPath), bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// 認証ヘッダを追加
	req.Header.Set("Authentication", "username "+auth)
	req.Header.Set("Content-Type", "application/json")

	return c.Server.Client().Do(req)
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) GetWithQueryParams(apiPath string, queries map[string]string) (*http.Response, error) {
	params := url.Values{}
	for k, v := range queries {
		params.Add(k, v)
	}

	fullURL := c.asURL(apiPath) + "?" + params.Encode()

	return c.Server.Client().Get(fullURL)
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
