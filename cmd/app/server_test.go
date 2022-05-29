package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muhammadkhon-abdulloev/load-balancer-go/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestServer_Lb(t *testing.T) {
	targetOne := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, one")
	}))
	defer targetOne.Close()

	targetTwo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, two")
	}))
	defer targetTwo.Close()

	lb := new(LoadBalancer)
	lb.Port = ":8000"
	lb.Services = []service.Service{
		{URL: targetOne.URL},
		{URL: targetTwo.URL},
	}

	lbs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lb.Lb(w, r)
	}))
	defer lbs.Close()

	res, err := http.Get(lbs.URL)
	if err != nil {
		t.Fatal(err)
	}

	result, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	expected := "Hello, one"

	require.Equal(t, expected, string(result), "They must be equal")
}
