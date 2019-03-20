package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/macaron.v1"
)

type testData struct {
	method   string
	url      string
	status   int
	handler  macaron.Handler
	path     string
	body     io.Reader
	response string
}

func init() {
	StartXormEngine()
	CreateInitialWorkerProfile()
}

func TestShowAllWorkers(t *testing.T) {

	test := []testData{
		{
			"GET",
			"/appscode/workers",
			200,
			ShowAllWorkers,
			"/appscode/workers",
			nil,
			`[{"username":"masud","firstname":"Masudur","lastname":"Rahman","city":"Madaripur","division":"Dhaka","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1},{"username":"fahim","firstname":"Fahim","lastname":"Abrar","city":"Chittagong","division":"Chittagong","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1},{"username":"tahsin","firstname":"Tahsin","lastname":"Rahman","city":"Chittagong","division":"Chittagong","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1},{"username":"jenny","firstname":"Jannatul","lastname":"Ferdows","city":"Chittagong","division":"Chittagong","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1}]`,
		},
	}

	for _, data := range test {
		runTest(data, t)
	}

}

func TestShowSingleWorker(t *testing.T) {
	test := []testData{
		{
			"GET",
			"/appscode/workers/masud",
			200,
			ShowSingleWorker,
			"/appscode/workers/:username",
			nil,
			`{"username":"masud","firstname":"Masudur","lastname":"Rahman","city":"Madaripur","division":"Dhaka","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1}`,
		},
		{
			"GET",
			"/appscode/workers/fahim",
			200,
			ShowSingleWorker,
			"/appscode/workers/:username",
			nil,
			`{"username":"fahim","firstname":"Fahim","lastname":"Abrar","city":"Chittagong","division":"Chittagong","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1}`,
		},
		{
			"GET",
			"/appscode/workers/tahsin",
			200,
			ShowSingleWorker,
			"/appscode/workers/:username",
			nil,
			`{"username":"tahsin","firstname":"Tahsin","lastname":"Rahman","city":"Chittagong","division":"Chittagong","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1}`,
		},
		{
			"GET",
			"/appscode/workers/jenny",
			200,
			ShowSingleWorker,
			"/appscode/workers/:username",
			nil,
			`{"username":"jenny","firstname":"Jannatul","lastname":"Ferdows","city":"Chittagong","division":"Chittagong","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T18:17:07+06:00","UpdatedAt":"2019-03-20T18:17:07+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1}`,
		},
		{
			"GET",
			"/appscode/workers/abcd",
			404,
			ShowSingleWorker,
			"/appscode/workers/:username",
			nil,
			`404 - Content Not Found`,
		},
	}

	for _, data := range test {
		runTest(data, t)

	}

}

func TestAddNewWorker(t *testing.T) {
	test := []testData{
		{
			"POST",
			"/appscode/workers",
			409,
			AddNewWorker,
			"/appscode/workers",
			strings.NewReader(`{"username":"masud","firstname":"Masudur","lastname":"Rahman","city":"Madaripur","division":"Dhaka","position":"Software Engineer","salary":55}`),
			`409 - username already exists`,
		},
		{
			"POST",
			"/appscode/workers",
			201,
			AddNewWorker,
			"/appscode/workers",
			strings.NewReader(`{"username":"masudur","firstname":"Masudur","lastname":"Rahman","city":"Madaripur","division":"Dhaka","position":"Software Engineer","salary":55}`),
			`{"username":"masudur","firstname":"Masudur","lastname":"Rahman","city":"Madaripur","division":"Dhaka","position":"Software Engineer","salary":55,"CreatedAt":"2019-03-20T19:16:49.9130127+06:00","UpdatedAt":"2019-03-20T19:16:49.913024478+06:00","DeletedAt":"0001-01-01T00:00:00Z","Version":1}`,
		},
	}

	for _, data := range test {
		runTest(data, t)
	}

}

func TestUpdateWorkerProfile(t *testing.T) {
	test := []testData{
		{
			"POST",
			"/appscode/workers/masud",
			405,
			UpdateWorkerProfile,
			"/appscode/workers/:username",
			strings.NewReader(`{"username":"masudd","firstname":"Masudur","lastname":"Rahman","city":"Madaripur","division":"Dhaka","position":"Software Engineer","salary":55}`),
			`405 - Username can't be changed`,
		},
		{
			"POST",
			"/appscode/workers/masudd",
			404,
			UpdateWorkerProfile,
			"/appscode/workers/:username",
			strings.NewReader(`{"username":"masudd","firstname":"Masudur","lastname":"Rahman","city":"Madaripur","division":"Dhaka","position":"Software Engineer","salary":55}`),
			`404 - Content Not Found`,
		},
		{
			"POST",
			"/appscode/workers/masud",
			201,
			UpdateWorkerProfile,
			"/appscode/workers/:username",
			strings.NewReader(`{"username":"masud","firstname":"Masudur","lastname":"Rahman","city":"M","division":"D","position":"Software Engineer","salary":55}`),
			`201 - Updated successfully`,
		},
	}

	for _, data := range test {
		runTest(data, t)
	}
}

func TestDeleteWorker(t *testing.T) {
	test := []testData{
		{
			"DELETE",
			"/appscode/workers/masud",
			200,
			DeleteWorker,
			"/appscode/workers/:username",
			nil,
			`200 - Deleted Successfully`,
		},
		{
			"DELETE",
			"/appscode/workers/fahim",
			200,
			DeleteWorker,
			"/appscode/workers/:username",
			nil,
			`200 - Deleted Successfully`,
		},
		{
			"DELETE",
			"/appscode/workers/hello",
			404,
			DeleteWorker,
			"/appscode/workers/:username",
			nil,
			`404 - Content Not Found`,
		},
	}
	for _, data := range test {
		runTest(data, t)
	}
}

func runTest(test testData, t *testing.T) {
	req, err := http.NewRequest(test.method, test.url, test.body)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()

	m := macaron.Classic()

	m.Handle(test.method, test.path, []macaron.Handler{test.handler})
	m.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != test.status {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, test.status)
	}

}
