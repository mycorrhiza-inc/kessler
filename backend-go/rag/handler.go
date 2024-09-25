package rag

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/sashabaranov/go-openai"
)

// Define the structure of our request JSON
type ChatHistory struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type RequestBody struct {
    Model       string        `json:"model"`
    ChatHistory []ChatHistory `json:"chat_history"`
}

func chatHandler(c *openai.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var reqBody RequestBody
        if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Create message slice for OpenAI request
        var messages []openai.ChatCompletionMessage
        for _, history := range reqBody.ChatHistory {
            messages = append(messages, openai.ChatCompletionMessage{
                Role:    history.Role,
                Content: history.Content,
            })
        }

        openaiRequest := openai.ChatCompletionRequest{
            Model:     openai.GPT3Dot5Turbo, // Use default model, will ignore reqBody.Model for simplicity.
            MaxTokens: 20, // For demonstration, we will limit the output to 20 tokens.
            Messages:  messages,
            Stream:    true,
        }

        ctx := context.Background()
        stream, err := c.CreateChatCompletionStream(ctx, openaiRequest)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to create chat completion stream: %v", err), http.StatusInternalServerError)
            return
        }
        defer stream.Close()

        var chatResponse string
        for {
            response, err := stream.Recv()
            if err != nil {
                if err != io.EOF {
                    http.Error(w, fmt.Sprintf("Stream error: %v", err), http.StatusInternalServerError)
                }
                break
            }
            chatResponse += response.Choices[0].Delta.Content
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"response": chatResponse})
    }
}

func main() {
    c := openai.NewClient("your token") // Replace with your actual token

    r := mux.NewRouter()
    r.HandleFunc("/chat", chatHandler(c)).Methods("POST")
    
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    log.Println("Server listening on port 8080")
    log.Fatal(srv.ListenAndServe())
}
```

Now let's zoom in to the details of what's happening:

* **HTTP Server Setup:** We set up a mux router and registered a handler for the `/chat` endpoint for POST requests.
* **Chat Request Handler (chatHandler):** 
    * We define the structures `ChatHistory` and `RequestBody` to mirror the expected JSON format of the request.
    * In `chatHandler`, we decode the incoming JSON into an instance of `RequestBody`.
    * We create a slice of `openai.ChatCompletionMessage` using the information from the `ChatHistory`.
    * We build a `ChatCompletionRequest` to be sent to OpenAI API.
    * We call OpenAI API using the provided client and stream handling functions.
    * We gather the response, combine it into a single string, and send it back as a JSON response.

Don't forget to replace `"your token"` with your actual OpenAI API token. Let me know if you have any clarifications or further questions!



package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SearchRequest struct {
	Query          string   `json:"query"`
	SeaerchFilters Metadata `json:"filters"`
}

func HandleSearchRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Hi there!")
		return
	case http.MethodPost:
		log.Println("Received a search request")

		// Create an instance of RequestData
		var RequestData SearchRequest

		// Decode the JSON body into the struct
		err := json.NewDecoder(r.Body).Decode(&RequestData)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close() // Close the body when done

		data, err := searchQuickwit(RequestData)
		if err != nil {
			log.Printf("Error searching quickwit: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		respString, err := json.Marshal(data)

		if err != nil {
			log.Println("Error marshalling response data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(respString))
	case http.MethodPut:
		fmt.Fprintf(w, "PUT request")
	case http.MethodDelete:
		fmt.Fprintf(w, "DELETE request")
	default:
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
	}
}
