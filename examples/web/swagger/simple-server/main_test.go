package main

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	go main()
}

func TestController(t *testing.T) {
	time.Sleep(time.Second)
	testApp := web.NewTestApp(t).Run(t)

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/123").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/999/name").
			Expect().Status(http.StatusOK)
	})

	t.Run("should delete employee ", func(t *testing.T) {
		testApp.Delete("/employee/333").
			Expect().Status(http.StatusOK)
	})

	t.Run("should report 404 when employee does not exist", func(t *testing.T) {
		testApp.Get("/employee/100").
			Expect().Status(http.StatusNotFound)
	})

	t.Run("should list employee", func(t *testing.T) {
		testApp.Get("/employee").
			Expect().Status(http.StatusOK)
	})

	t.Run("should update employee", func(t *testing.T) {
		testApp.Put("/employee").
			WithJSON(&UpdateEmployeeRequest{
			Employee: Employee{
				Id:        12345,
				FirstName: "foo",
				LastName:  "bar",
				Manger: Manager{
					ID:   23345,
					Name: "baz",
				},
				Assets: []Asset{
					{
						ID:   1234,
						Name: "abc",
					},
					{
						ID:   5678,
						Name: "def",
					},
				},
			},
		}).Expect().Status(http.StatusOK)
	})

	t.Run("should create employee by post /employee", func(t *testing.T) {
		testApp.Post("/employee").
			WithJSON(&CreateEmployeeRequest{
			Employee: Employee{
				Id:        12345,
				FirstName: "foo",
				LastName:  "bar",
				Manger: Manager{
					ID:   23345,
					Name: "baz",
				},
				Assets: []Asset{
					{
						ID:   1234,
						Name: "abc",
					},
					{
						ID:   5678,
						Name: "def",
					},
				},
			},
		}).Expect().Status(http.StatusOK)
	})

	t.Run("should add asset", func(t *testing.T) {
		testApp.Post("/employee/add-assets").
			WithJSON([]*Asset{
				{
					ID: 1,
					Name: "foo",
				},
				{
					ID: 2,
					Name: "bar",
				},
			}).Expect().Status(http.StatusOK)
	})

	t.Run("should report 500 error if create employee without request body", func(t *testing.T) {
		testApp.Post("/employee").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should get employees", func(t *testing.T) {
		testApp.Get("/employee").
			Expect().Status(http.StatusOK)
	})

}
