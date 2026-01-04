package gpt

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// mock response body similar to OpenAI API
func mockGPTResponse(content string) string {
	resp := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]interface{}{
					"content": content,
				},
			},
		},
	}
	b, _ := json.Marshal(resp)
	return string(b)
}

func TestClassifyPost(t *testing.T) {
	// --- mock GPT server ---
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "Existing categories") {
			t.Errorf("expected prompt to contain category list, got: %s", body)
		}

		// return mock response
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, mockGPTResponse(`Technology/AI
Technology/Automation`))
	}))
	defer server.Close()

	// --- create client ---
	token := os.Getenv("GPT_API_TOKEN")
	client := NewClient(server.URL, token)

	categories := []string{
		"Technology/AI",
		"Technology/Programming",
		"Health/Nutrition",
	}

	post := "Apple just announced new AI features in iOS 19."

	// --- run test ---
	suggestions, err := client.ClassifyPost(post, categories)
	if err != nil {
		t.Fatalf("ClassifyPost returned error: %v", err)
	}

	// --- validate ---
	expected := []string{"Technology/AI", "Technology/Automation"}
	if len(suggestions) != len(expected) {
		t.Fatalf("expected %d suggestions, got %d", len(expected), len(suggestions))
	}

	for i, cat := range expected {
		if suggestions[i] != cat {
			t.Errorf("expected %s, got %s", cat, suggestions[i])
		}
	}
}
