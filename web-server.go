package main

import (
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

func webServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the multipart form data
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get the uploaded image file
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Decode the image
		img, _, err := image.Decode(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Ensure Message is passed
		message := r.FormValue("message")

		if message == "" {
			http.Error(w, "Message is required", http.StatusBadRequest)
			return
		}

		// Encode the message in the image

		encoded := encodeMessage(img, message)

		// Create a new image file
		output, err := os.Create("output.jpg")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer output.Close()

		// Write the new image file
		err = jpeg.Encode(output, encoded, &jpeg.Options{Quality: 100})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the new image file to the client
		http.ServeFile(w, r, "output.jpg")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
