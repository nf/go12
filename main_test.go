/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type StatusHandler int

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(*h))
}

func TestIsTagged(t *testing.T) {
	var status StatusHandler = http.StatusNotFound
	s := httptest.NewServer(&status)
	defer s.Close()
	if isTagged(s.URL) {
		t.Fatal("isTagged == true, want false")
	}
	status = http.StatusOK
	if !isTagged(s.URL) {
		t.Fatal("isTagged == false, want true")
	}
}

func TestIntegration(t *testing.T) {
	var status StatusHandler = http.StatusNotFound
	ts := httptest.NewServer(&status)
	defer ts.Close()

	s := NewServer(ts.URL, 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if b := w.Body.String(); !strings.Contains(b, "No.") {
		t.Fatalf("body = %q, want no", b)
	}

	status = http.StatusOK
	time.Sleep(20 * time.Millisecond)

	w = httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if b := w.Body.String(); !strings.Contains(b, "YES!") {
		t.Fatalf("body = %q, want yes", b)
	}
}
