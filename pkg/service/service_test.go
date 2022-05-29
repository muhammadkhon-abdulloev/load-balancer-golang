package service

import (
	// "sync"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_SetDead(t *testing.T) {
	testTable := []struct {
		testService Service
		expected    bool
	}{
		{
			testService: Service{
				URL: "alive",
			},
			expected: false,
		},
		{
			testService: Service{
				URL: "dead",
			},
			expected: true,
		},
	}

	for _, testCase := range testTable {
		if testCase.testService.URL == "alive" {
			testCase.testService.SetDead(false)
		} else {
			testCase.testService.SetDead(true)
		}
		if ok := assert.Equal(t, testCase.expected, testCase.testService.IsDead); !ok {
			t.Errorf("Incorrect result. expected %v, got %v", testCase.expected, testCase.testService.IsDead)
		}
	}

}

func TestService_GetIsDead(t *testing.T) {
	testTable := []struct {
		testService Service
		expected    bool
	}{
		{
			testService: Service{
				URL: "alive",
			},
			expected: false,
		},
		{
			testService: Service{
				URL: "dead",
				IsDead: true,
			},
			expected: true,
		},
	}

	for _, testCase := range testTable {
		var result bool
		if testCase.testService.URL == "alive" {
			result = testCase.testService.GetIsDead()
		} else {
			result = testCase.testService.GetIsDead()
		}
		
		if ok := assert.Equal(t, testCase.expected, result); !ok {
			t.Errorf("Incorrect result. expected %v, got %v", testCase.expected, testCase.testService.IsDead)
		}
	}
}

func TestService_IsAlive(t *testing.T) {
	testService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, one")
	}))
	defer testService.Close()

	testTable := []struct {
		url      string
		expected bool
	}{
		{
			url:      testService.URL,
			expected: true,
		},
		{
			url:      "http://localhost:5000",
			expected: false,
		},
	}

	for _, testCase := range testTable {
		url, _ := url.Parse(testCase.url)
		result := IsAlive(url)

		if ok := assert.Equal(t, testCase.expected, result); !ok {
			t.Errorf("Incorrect result. expected %v, got %v", testCase.expected, result)
		}
	}
}
