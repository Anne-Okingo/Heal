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

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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

func GetUsernameHandler(w http.ResponseWriter, r *http.Request) {
	// Get session ID from the cookies
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized: Missing session ID", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	// Open the database connection
	db, err := sql.Open("sqlite3", "./Heal.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query to fetch the username based on the session ID
	var username string
	query := `
		SELECT users.username 
		FROM users 
		INNER JOIN sessions ON users.id = sessions.user_id 
		WHERE sessions.session_id = ?
	`
	err = db.QueryRow(query, sessionID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized: Invalid session ID", http.StatusUnauthorized)
			return
		}
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with the username in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"username": username})
}

// login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./Heal.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if user.Name == "" || user.Password == "" {
		http.Error(w, "Name and password are required", http.StatusBadRequest)
		return
	}

	// Query database for user
	var storedPassword string
	var userID int
	err = db.QueryRow("SELECT id,password FROM users WHERE username = ?", user.Name).Scan(&userID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("Database error:", err)
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a session ID
	sessionID := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour) // Session expires in 24 hours

	// Store session in database
	_, err = db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("Database error:", err)
		return
	}

	// Set the session ID as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true, // Prevent JavaScript access
	})

	// Successful login
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}

func Getdb(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create the users table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	createsessionquerry := `
	CREATE TABLE IF NOT EXISTS sessions (
		session_id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := db.Exec(createsessionquerry); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return db, nil
}

// singup handler
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./Heal.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	Getdb("./Heal.db")
	defer db.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var creds struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if creds.Name == "" || creds.Password == "" {
		http.Error(w, `{"message": "Name and password are required"}`, http.StatusBadRequest)
		return
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err = db.Exec(query, creds.Name, hashedPassword)
	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, `{"message": "Failed to create account. Username may already exist."}`, http.StatusConflict)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Signup successful",
	})
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error in SessionHandler: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		}
	}()

	// Check if a session_id cookie exists
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		// Generate a guest session ID
		guestSessionID := generateGuestSessionID()

		// Store the guest session (for example, tracking purposes)
		guestSessions[guestSessionID] = "guest"

		// Set the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    guestSessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Set to true if using HTTPS
			Expires:  time.Now().Add(24 * time.Hour),
		})
		log.Printf("Assigned guest session ID: %s", guestSessionID)
	}

	renders.RenderTemplate(w, "Sessions.page.html", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error in HomeHandler: %v", err)
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

	renders.RenderTemplate(w, "Index.page.html", nil)
}

func DataPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error in PrivacyHandler: %v", err)
			http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		}
	}()

	renders.RenderTemplate(w, "Privacy.page.html", nil)
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var guestSessions = make(map[string]string)

// NotFoundHandler handles unknown routes; 404 status
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renders.RenderTemplate(w, "notfound.page.html", nil)
}

// BadRequestHandler handles bad requests routes
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	renders.RenderTemplate(w, "badrequest.page.html", nil)
}

// ServerErrorHandler handles server failures that result in status 500
func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	renders.RenderTemplate(w, "serverError.page.html", nil)
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error in HomeHandler: %v", err)
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
			fmt.Println("hello")
			// Valid session: redirect to /session
			http.Redirect(w, r, "/session", http.StatusFound)
			return
		}
	}

	renders.RenderTemplate(w, "chat.page.html", nil)
}