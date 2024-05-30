package main

import (
	"errors"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Middleware to handle CORS and check the access token
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Middleware: Received request")

		// Handle CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Party")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			log.Warn("Middleware: CORS error")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Check for access token
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Warn("Middleware: Missing Authorization token")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		log.Debug("Middleware: Authorization token found")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func getAltinnPlatformToken(idPortenToken string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://platform.tt02.altinn.no/authentication/api/v1/exchange/id-porten", nil)
	if err != nil {
		log.Warn("AltinnPlatformToken: Error creating new request: ", err)
		return "", errors.New("error creating Altinn exchange API request")
	}

	req.Header.Set("Authorization", idPortenToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Warn("AltinnPlatformToken: Error calling Altinn exchange API: ", err)
		return "", errors.New("error calling Altinn exchange API")
	}

	if resp.StatusCode != http.StatusOK {
		log.Warn("AltinnPlatformToken: Altinn exchange API responded with " + resp.Status)
		return "", errors.New("altinn exchange API responded with " + resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("AltinnPlatformToken: Error parsing response body: ", err)
		return "", errors.New("failed to parse response body from Altinn exchange API")
	}

	newToken := "Bearer " + string(body)
	log.Debug("AltinnPlatformToken: Altinn Plattform Token: ", newToken)
	return newToken, nil
}

// Handler function to make the GET request to the external API
func handleCompanyRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug("CompanyHandler: Received request to handle external API call")

	// Extract the access token from the header
	idPortenToken := r.Header.Get("Authorization")
	log.Debug("CompanyHandler: Extracted token: ", idPortenToken)

	altinnPlattformToken, err := getAltinnPlatformToken(idPortenToken)
	if err != nil {
		log.Warn("CompanyHandler: Error authentication to Altinn API: ", err)
		http.Error(w, "Error authentication to Altinn API", http.StatusInternalServerError)
		return
	}

	// Make the GET request to the external API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://brg.apps.tt02.altinn.no/brg/lpid-wallet-2024/api/v1/parties?allowedtoinstantiatefilter=true", nil)
	if err != nil {
		log.Warn("CompanyHandler: Error creating new request: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", altinnPlattformToken)
	req.Header.Set("Content-Type", "application/json")

	log.Debug("CompanyHandler: Making request to external API")

	resp, err := client.Do(req)
	if err != nil {
		log.Warn("CompanyHandler: Error making request to external API: ", err)
		http.Error(w, "Failed to make request to external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	log.Debug("CompanyHandler: response body from Altinn allowed parties", string(body))
	if err != nil {
		log.Warn("CompanyHandler: Error reading response from external API: ", err)
		http.Error(w, "Failed to read response from external API", http.StatusInternalServerError)
		return
	}

	log.Info("CompanyHandler: Received response from external API with status: ", resp.Status)

	// Write the response back to the user
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Debug("CompanyHandler: Response sent back to the client")
}

// Handler function to make the GET request to the external API
func handleLpidRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug("LpidHandler: Received request to handle external API call")

	// Extract the access token from the header
	idPortenToken := r.Header.Get("Authorization")
	partyId := r.Header.Get("Party")
	log.Debug("LpidHandler: Extracted Auth token: ", idPortenToken)
	log.Debug("LpidHandler: Extracted PartyID: ", partyId)

	altinnPlattformToken, err := getAltinnPlatformToken(idPortenToken)
	if err != nil {
		log.Warn("LpidHandler: Error authentication to Altinn API: ", err)
		http.Error(w, "Error authentication to Altinn API", http.StatusInternalServerError)
		return
	}

	// Make the GET request to the external API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://brg.apps.tt02.altinn.no/brg/lpid-wallet-2024/v1/data?dataType=model&includeRowId=true&language=nb", nil)
	if err != nil {
		log.Warn("LpidHandler: Error creating new request: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", altinnPlattformToken)
	req.Header.Set("Party", partyId)
	req.Header.Set("Content-Type", "application/json")

	log.Debug("LpidHandler: Making request to external API")

	resp, err := client.Do(req)
	if err != nil {
		log.Warn("LpidHandler: Error making request to external API: ", err)
		http.Error(w, "Failed to make request to external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	log.Debug("LpidHandler: response body from Altinn allowed parties", string(body))
	if err != nil {
		log.Warn("LpidHandler: Error reading response from external API: ", err)
		http.Error(w, "Failed to read response from external API", http.StatusInternalServerError)
		return
	}

	log.Info("LpidHandler: Received response from external API with status: ", resp.Status)

	// Write the response back to the user
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	log.Debug("LpidHandler: Response sent back to the client")
}

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("Starting server")

	// Create a new HTTP server
	mux := http.NewServeMux()
	mux.Handle("/api/v1/company", middleware(http.HandlerFunc(handleCompanyRequest)))
	mux.Handle("/api/v1/lpid", middleware(http.HandlerFunc(handleLpidRequest)))

	// Start the server on port 8080
	log.Info("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Error("Could not start server: ", err)
	}
}
