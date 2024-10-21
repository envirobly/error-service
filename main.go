package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"
    "time"
)

// Define the HTML template
var tmpl = template.Must(template.New("status").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Code}} {{.Message}}</title>
</head>
<body>
    <h1>⚠️ {{.Code}}</h1>
    <p>{{.Message}}</p>
</body>
</html>
`))

// StatusData holds data to be passed to the template
type StatusData struct {
    Code    int
    Message string
}

// Custom message based on status code
func getStatusMessage(code int) string {
    switch code {
    case 200:
        return "OK - Everything is good!"
    case 404:
        return "Not Found - Sorry, the page you are looking for doesn't exist."
    case 500:
        return "Internal Server Error - Something went wrong on our side."
    case 502:
        return "Bad Gateway - Can't get a response from the container handling this route. Check the service logs for startup errors."
    case 503:
        return "Service Unavailable - Currently there is no service configured to respond to this request."
    default:
        return http.StatusText(code) // Default message from net/http
    }
}

// Handler function for serving custom status codes
func statusHandler(w http.ResponseWriter, r *http.Request) {
    // Extract the status code from the URL path, e.g., "/404.html"
    path := r.URL.Path
    if len(path) < 2 || path[0] != '/' {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    codeStr := path[1 : len(path)-5] // Remove the ".html" suffix
    code, err := strconv.Atoi(codeStr)
    if err != nil || len(path) < 6 || path[len(path)-5:] != ".html" {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    // Set the no-cache headers
    w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
    w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
    w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

    // Set the HTTP status code
    w.WriteHeader(code)

    // Get the custom message
    message := getStatusMessage(code)

    // Render the HTML template
    data := StatusData{
        Code:    code,
        Message: message,
    }
    if err := tmpl.Execute(w, data); err != nil {
        log.Printf("Template execution error: %v", err)
    }
}

func main() {
    http.HandleFunc("/", statusHandler)
    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
