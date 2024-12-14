package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"Heal/internals/renders"

	"github.com/joho/godotenv"
)

type GeminiRequest struct {
	Contents []struct {
		Role  string `json:"role,omitempty"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerationConfig struct {
		MaxOutputTokens int     `json:"maxOutputTokens"`
		Temperature     float64 `json:"temperature"`
		TopP            float64 `json:"topP"`
		TopK            int     `json:"topK"`
	} `json:"generationConfig"`
	SystemInstruction struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"systemInstruction"`
}

func ProxyGeminiRequest(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Retrieve API key securely from environment
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Println("Gemini API key is not set")
		http.Error(w, "API configuration error", http.StatusInternalServerError)
		return
	}

	// Read the incoming request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error processing request", http.StatusBadRequest)
		return
	}

	// Validate and parse the request body
	var geminiRequest GeminiRequest
	if err := json.Unmarshal(body, &geminiRequest); err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Construct request to Gemini API
	geminiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", geminiAPIKey)

	// Forward the request to Gemini
	proxyReq, err := http.NewRequest(http.MethodPost, geminiURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set necessary headers
	proxyReq.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error sending request to Gemini: %v", err)
		http.Error(w, "Error communicating with API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Gemini response: %v", err)
		http.Error(w, "Error reading API response", http.StatusInternalServerError)
		return
	}

	os.Stdout.Write([]byte{responseBody[1]})

	// Copy headers and status code
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	// Write the response back to the client
	w.Write(responseBody)
}

type SpeechifyRequest struct {
	Input       string `json:"input"`
	VoiceID     string `json:"voice_id"`
	AudioFormat string `json:"audio_format"`
}

func ProxySpeechifyRequest(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Retrieve API key securely from environment
	speechifyAPIKey := os.Getenv("SPEECHIFY_API_KEY")
	if speechifyAPIKey == "" {
		log.Println("Speechify API key is not set")
		http.Error(w, "API configuration error", http.StatusInternalServerError)
		return
	}

	// Read the incoming request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error processing request", http.StatusBadRequest)
		return
	}

	// Validate and parse the request body
	var speechifyRequest SpeechifyRequest
	if err := json.Unmarshal(body, &speechifyRequest); err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Construct request to Speechify API
	speechifyURL := "https://api.sws.speechify.com/v1/audio/speech"

	// Forward the request to Speechify
	proxyReq, err := http.NewRequest(http.MethodPost, speechifyURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set necessary headers
	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("Authorization", "Bearer "+speechifyAPIKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error sending request to Speechify: %v", err)
		http.Error(w, "Error communicating with API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Speechify response: %v", err)
		http.Error(w, "Error reading API response", http.StatusInternalServerError)
		return
	}

	// Copy headers and status code
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	// Write the response back to the client
	w.Write(responseBody)
}

// Helper function to generate a guest session ID
func generateGuestSessionID() string {
	return fmt.Sprintf("guest-%d", time.Now().UnixNano())
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error in WelcomeHandler: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		}
	}()

	db, err := sql.Open("sqlite3", "./Heal.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get session ID from cookies
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// If the cookie exists, validate the session
		var expiresAt time.Time
		var userID int
		err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = ?", cookie.Value).Scan(&userID, &expiresAt)
		if err == nil && expiresAt.After(time.Now()) {
			// Valid session: redirect to /session
			http.Redirect(w, r, "/session", http.StatusFound)
			return
		}
	}

	renders.RenderTemplate(w, "Welcome.page.html", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Optionally, delete the session from the database (if necessary)
	db, err := sql.Open("sqlite3", "./Heal.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// You can access the session ID from the deleted cookie if needed
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized: Missing session ID", http.StatusUnauthorized)
		log.Printf("error getting session ID: %v", err)
		return
	}

	sessionID := cookie.Value
	if sessionID == "" {
		http.Error(w, "Unauthorized: Empty session ID", http.StatusUnauthorized)
		log.Printf("empty session ID")
		return
	}

	_, err = db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete the cookie
		HttpOnly: true,
	})

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
}