package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// Login starts a local server, opens the browser to the OAuth URL,
// waits for the backend to redirect back with the token cookie,
// then fetches the token from /auth/token.
// Returns the token string.
func Login(apiURL string) (string, error) {
	tokenCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Start local server on a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("could not start local server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	server := &http.Server{}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			errCh <- fmt.Errorf("no token in callback")
			fmt.Fprintf(w, "<html><body><h2>Authentication failed. You can close this tab.</h2></body></html>")
			return
		}

		tokenCh <- token
		fmt.Fprintf(w, "<html><body><h2>Authenticated! You can close this tab.</h2></body></html>")
	})

	go func() {
		server.Serve(listener)
	}()

	// The backend needs to redirect to our local callback
	// For CLI we use a slightly different OAuth flow:
	// Open browser → user logs in → backend redirects to localhost:PORT/callback
	loginURL := fmt.Sprintf("%s/auth/google?redirect_uri=http://127.0.0.1:%d/callback", apiURL, port)

	fmt.Printf("Opening browser for authentication...\n")
	if err := openBrowser(loginURL); err != nil {
		fmt.Printf("Could not open browser automatically.\nPlease open this URL manually:\n\n%s\n\n", loginURL)
	}

	// Wait for token or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	defer server.Close()

	select {
	case token := <-tokenCh:
		return token, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", fmt.Errorf("login timed out after 2 minutes")
	}
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default: // linux
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

func parseJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
