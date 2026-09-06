package main

import "net/http"
import "net/http/httptest"
import "testing"
import "io"
import "sync"

func TestSetHandler(t *testing.T) {
	db := make(map[string]string)
	lock := new(sync.RWMutex)
	var req *http.Request
	var rr *httptest.ResponseRecorder
	
	t.Run("set foo=bar", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/set?foo=bar", nil)
		rr = httptest.NewRecorder()
		SetHandler(db, lock, rr, req)
		if db["foo"] != "bar" {
			t.Errorf("key foo not found (should be bar)")
		}
	})
	
	t.Run("overwrite foo=quux", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/set?foo=quux", nil)
		rr = httptest.NewRecorder()
		SetHandler(db, lock, rr, req)
		if db["foo"] != "quux" {
			t.Errorf("key foo not quux")
		}
	})


	t.Run("no keys on set", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/set", nil)
		rr = httptest.NewRecorder()
		SetHandler(db, lock, rr, req)
		if rr.Code != 400 {
			t.Errorf("no keys on set should return 400")
		}
	})

	t.Run("no val for set", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/set?foo", nil)
		rr = httptest.NewRecorder()
		SetHandler(db, lock, rr, req)
		if rr.Code != 400 {
			t.Errorf("no value on set should return 400")
		}
	})

}

func TestGetHandler(t *testing.T) {
	db := map[string]string{"foo": "bar"}
	lock := new(sync.RWMutex)
	var req *http.Request
	var rr *httptest.ResponseRecorder


	t.Run("get foo", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/get?foo", nil)
		rr = httptest.NewRecorder()
		GetHandler(db, lock, rr, req)
		resp, err := io.ReadAll(rr.Result().Body)
		if err != nil {
			t.Errorf("error reading response: %s", err)
		}
		
		if string(resp) != "bar" {
			t.Errorf("key foo not found (should be bar)")
		}
	})

	t.Run("get nonexistent 404s", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/get?baz", nil)
		rr = httptest.NewRecorder()
		GetHandler(db, lock, rr, req)
		if rr.Result().StatusCode != 404 {
			t.Errorf("get nonexistent key should 404")
		}
		
	})
	
	t.Run("get can't have value", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/get?baz=bar", nil)
		rr = httptest.NewRecorder()
		GetHandler(db, lock, rr, req)
		if rr.Result().StatusCode != 400 {
			t.Errorf("get with value should 400")
		}
		
	})

	t.Run("get can't have multiple keys", func(t *testing.T) {
		req = httptest.NewRequest("GET", "/get?baz&quux", nil)
		rr = httptest.NewRecorder()
		GetHandler(db, lock, rr, req)
		if rr.Result().StatusCode != 400 {
			t.Errorf("get with multiple keys should 400")
		}
		
	})



}
	
