# Heal: Comprehensive Mental Health Support Assistant

**Heal** is a web application designed to provide personalized support, guidance, and education for individuals experiencing emotional distress, trauma, or challenges in their mental well-being, with a specific focus on survivors of Sexual and Gender-Based Violence (SGBV).

This application serves as an empathetic, non-judgmental, and empowering AI mental health companion that users can rely on for active listening, education, trauma-sensitive guidance, and crisis navigation.

---

## Key Features

### 1. Active Listening and Validation
- A safe space for users to share their feelings or vent.
- Empathetic, reflective, and supportive responses.
- Non-judgmental acknowledgment of emotions.

### 2. Education and Empowerment
- Information on mental health, coping mechanisms, and healing processes.
- Self-care strategies and mindfulness exercises tailored to user needs.
- Rights education for SGBV survivors.

### 3. Trauma-Sensitive Guidance
- Gradual and safe conversations to process trauma.
- Encouragement of self-compassion and healing.
- Resources for professional help when necessary.

### 4. Crisis Navigation and Reporting
- Emergency guidance with regional resource suggestions (hotlines, shelters).
- Steps to create safety plans for users in danger.

### 5. Privacy and Security
- Secure and anonymized communication channels.
- Local data processing with minimal logging to ensure privacy.

---

## Tech Stack

### **Frontend**
- **HTML, CSS, and JavaScript**: Core technologies for creating a dynamic and user-friendly interface.
- **Fetch API**: For communication between the frontend and backend.(login,signout,get user name, logout, speechify and gemini)

### **Backend**
- **Go (Golang)**: Handles server-side logic, routing, and AI-powered responses.
- **RESTful API**: Built with Go's `net/http` package for client-server communication.
- **Data Storage**: SQLLite In-memory for session data to maintain user privacy.

---

## Features and Functionalities

### **Frontend**
- A simple, responsive UI where users can:
  - Start conversations with the AI assistant.
  - Access educational resources and exercises.
  - Receive tailored guidance and support.

### **Backend**
- Real-time, AI-powered conversational logic implemented in Go.
- RESTful endpoints to handle user inputs and return responses.
- Pre-configured logic for mental health support and trauma-sensitive conversations.

---

## Getting Started

### Prerequisites
- **Go** (v1.20 or higher) installed on your system.
- Basic knowledge of **HTML, CSS, and JavaScript**.

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Anne-Okingo/Heal.git
   cd Heal

2. **Run the program from the root directory*:
   ```bash
   go run ./cmd/

   ```

   or 
   ```bash
   make
   ```
