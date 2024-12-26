package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (r Response) Json(w http.ResponseWriter) {
	res, _ := json.Marshal(r)

	if _, err := fmt.Fprintf(w, string(res)); err != nil {
		panic(err)
	}

	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN"))
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)

		res := Response{
			Error: "Wrong Method, only accepting 'POST'",
		}

		res.Json(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)

		res := Response{
			Error: "Wrong Content-Type, only accepting 'application/json'",
		}

		res.Json(w)
		return
	}

	var d struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		panic(err)
	}

	if d.Email == "" || d.Message == "" {
		w.WriteHeader(http.StatusBadRequest)

		res := Response{
			Error: "Empty data, 'email' or 'message'",
		}

		res.Json(w)
		return
	}

	mg := mailgun.NewMailgun(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
	)

	mg.SetAPIBase(os.Getenv("MAILGUN_API_BASE"))

	sender := d.Email

	if d.Name != "" {
		sender = fmt.Sprintf("%s <%s>", d.Name, d.Email)
	}
	subject := os.Getenv("MAIL_SUBJECT")
	recipient := os.Getenv("MAIL_RECIPIENT")
	name := fmt.Sprintf("\n\n%s", d.Name)
	phone := fmt.Sprintf("\n%s", d.Phone)
	text := fmt.Sprintf("%s%s%s", d.Message, name, phone)

	message := mg.NewMessage(sender, subject, text, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, _, err := mg.Send(ctx, message)

	if err != nil {
		log.Printf("Could not send email: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		res := Response{
			Error: err.Error(),
		}

		res.Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)

	res := Response{
		Message: "All good!",
	}

	res.Json(w)
	return
}

func main() {
	log.Print("starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
