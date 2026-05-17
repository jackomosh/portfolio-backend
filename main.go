package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/submit", contactHandler)

	log.Printf("Production server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for your GitHub Pages domain
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid Form Data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	message := r.FormValue("comment")

	if name == "" || email == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	err := sendEmail(name, email, message)
	if err != nil {
		log.Println("Email Error:", err)
		http.Error(w, "Email failed to send", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message received successfully!"))
}

func sendEmail(name, email, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", os.Getenv("RECEIVER_EMAIL"))
	m.SetHeader("Subject", "New Contact from " + name)

	htmlBody := `
	<html>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
		<div style="max-width: 600px; margin: auto; background: white; padding: 20px; border-radius: 10px; border: 1px solid #ddd;">
			<h2 style="color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px;">New Inquiry</h2>
			<p><strong>Name:</strong> ` + name + `</p>
			<p><strong>Email:</strong> ` + email + `</p>
			<div style="margin-top: 20px; padding: 15px; background: #f9f9f9; border-left: 5px solid #007bff;">
				<p><strong>Message:</strong></p>
				<p>` + message + `</p>
			</div>
			<p style="font-size: 12px; color: #888; margin-top: 20px;">Sent via Portfolio Backend • ` + time.Now().Format("Jan 02, 2006") + `</p>
		</div>
	</body>
	</html>`

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASS"))
	return d.DialAndSend(m)
}