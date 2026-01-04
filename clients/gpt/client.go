package gpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/pkg/errors"
)

type Client interface {
	StreamChatCompletions(body []byte) (<-chan string, error)
	EmbedText(text string) ([]float32, error)
	TranscribeAudio(fileData []byte, fileName string) (*TranscriptionResult, error)
	ClassifyPost(postText string, categories []string) ([]string, error)
}

type client struct {
	BaseUrl    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseUrl, token string) Client {
	return &client{
		BaseUrl: baseUrl,
		Token:   token,
		HTTPClient: &http.Client{
			// Increased timeout for streaming
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *client) StreamChatCompletions(body []byte) (<-chan string, error) {
	resp, err := c.doStreamRequest(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform streaming request")
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	streamChan := make(chan string)
	go c.processStream(resp, streamChan)

	return streamChan, nil
}

func (c *client) EmbedText(text string) ([]float32, error) {
	// 1. Define the request payload.
	// We use a modern and efficient model, but this can be changed.
	embeddingReq := EmbeddingRequest{
		Input: text,
		Model: "text-embedding-3-small",
	}

	// 2. Marshal the payload into a JSON byte slice.
	body, err := json.Marshal(embeddingReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal embedding request")
	}

	// 3. Make the HTTP POST request to the embeddings endpoint.
	endpoint := "/embeddings"
	resp, err := c.doPost(endpoint, body, nil) // No special headers are needed for this request.
	if err != nil {
		return nil, errors.Wrap(err, "failed to make embedding request")
	}
	defer resp.Body.Close()

	// 4. Handle non-successful status codes.
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 5. Decode the JSON response from the response body.
	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode embedding response")
	}

	// 6. Validate the response and extract the embedding vector.
	if len(embeddingResp.Data) == 0 || len(embeddingResp.Data[0].Embedding) == 0 {
		return nil, errors.New("received an empty or invalid embedding from the API")
	}

	return embeddingResp.Data[0].Embedding, nil
}

func (c *client) TranscribeAudio(fileData []byte, fileName string) (*TranscriptionResult, error) {
	// 1. Create a buffer to hold the multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 2. Create the file field
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create form file")
	}

	// 3. Write the file data to the form
	if _, err = part.Write(fileData); err != nil {
		return nil, errors.Wrap(err, "failed to write file data")
	}

	// 4. Add the model field (whisper-1 is the standard model)
	_ = writer.WriteField("model", "whisper-1")

	// 5. Add response_format to get verbose JSON with language info
	_ = writer.WriteField("response_format", "verbose_json")

	_ = writer.WriteField("prompt",
		"The audio may contain speech, music, or singing. "+
			"Please transcribe every spoken or sung word exactly as heard, writing each phrase or lyric line on a new line. "+
			"If there is any part where someone is singing (not just background music), prefix those lines with a 🎵 emoji. "+
			"Do not summarize or skip lyrics.")

	// 6. Close the writer to finalize the multipart form
	if err = writer.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close multipart writer")
	}

	// 7. Make the HTTP POST request using doPost
	endpoint := "/audio/transcriptions"
	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	resp, err := c.doPost(endpoint, buf.Bytes(), headers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transcription request")
	}
	defer resp.Body.Close()

	// 8. Handle non-successful status codes
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("transcription request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 9. Decode the JSON response
	var transcriptionResp TranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&transcriptionResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode transcription response")
	}

	return &TranscriptionResult{
		Text:     transcriptionResp.Text,
		Language: transcriptionResp.Language,
	}, nil
}

func (c *client) ClassifyPost(postText string, categories []string) ([]string, error) {
	systemPrompt := `You are an intelligent content categorization assistant.
You will receive a post text and a list of existing categories.
Each category is in the format: 
"Main Category/Subcategory[/Subcategory2]".

Your task:
- Choose one or more categories that best match the post.
- If needed, propose new **subcategories**, but you must use one of the existing main categories as the root.
- Do NOT invent a new main category.
- Return the list of suggested categories in plain text, one per line.`

	userPrompt := fmt.Sprintf(`Post:
"%s"

Existing categories:
%s

Return format:
MainCategory/Subcategory
MainCategory2/Subcategory/Subcategory`, postText, strings.Join(categories, "\n"))

	body, err := json.Marshal(ChatRequest{
		Model: "gpt-5-mini",
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	resp, err := c.doPost("/v1/chat/completions", body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result ChatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.Wrap(err, "failed to parse response")
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in GPT response")
	}

	// Split by newline and trim
	lines := strings.Split(result.Choices[0].Message.Content, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return cleaned, nil
}

func (c *client) doPost(endpoint string, body []byte, headers map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseUrl, endpoint)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set required headers for all requests
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Only set Content-Type if not already provided in headers
	if headers == nil || headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set any additional headers (this will override Content-Type if provided)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	return resp, nil
}

func (c *client) doStreamRequest(body []byte) (*http.Response, error) {
	endpoint := "/chat/completions"
	headers := map[string]string{
		"Accept":     "text/event-stream",
		"Connection": "keep-alive",
	}
	return c.doPost(endpoint, body, headers)
}

// processStream reads the streaming response body and sends content chunks to a channel.
func (c *client) processStream(resp *http.Response, streamChan chan string) {
	defer resp.Body.Close()
	defer close(streamChan)

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				logger.Info("Error reading stream", err)
			}
			break
		}

		dataPrefix := "data: "
		if !strings.HasPrefix(string(line), dataPrefix) {
			continue
		}

		jsonStr := strings.TrimPrefix(string(line), dataPrefix)
		jsonStr = strings.TrimSpace(jsonStr)

		if jsonStr == "[DONE]" {
			break
		}

		var chunk GPTStreamChunk
		if err := json.Unmarshal([]byte(jsonStr), &chunk); err != nil {
			// It's better to log this than to panic.
			logger.Error("Error unmarshalling stream chunk", err)
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			streamChan <- chunk.Choices[0].Delta.Content
		}
	}
}
