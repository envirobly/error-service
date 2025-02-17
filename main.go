package main

import (
    "flag"
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
  <title>{{.Title}}</title>

  <style type="text/css">
    html {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
      background: hwb(53 96% 10%);
      color: hwb(53 16% 99%);
    }

    .message {
      margin: 4rem auto 0;
      text-align: center;
      font-size: 1.3rem;
      max-width: 30rem;
    }

    .response-code-and-title {
      position: relative;
    }

    .response-code {
      font-size: 11rem;
      font-weight: 900;
      color: hwb(53 96% 27%);
      margin: 0;
    }

    .response-title {
      width: 100%;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      margin: 0;
      font-size: 1.9rem;
      text-shadow: 0px 0px 0.4rem hwb(53 96% 10%);
    }

    .explanation {
      line-height: 1.5;
      text-wrap: balance;
    }
  </style>
</head>
<body>
  <div class="message">
    <div class="response-code-and-title">
      <h1 class="response-code">{{.Code}}</h1>
      <h2 class="response-title">{{.Title}}</h2>
    </div>
    <p class="explanation">{{.Message}}</p>
  </div>
</body>
</html>
`))

// StatusData holds data to be passed to the template
type StatusData struct {
    Code    int
    Title   string
    Message string
}

// Custom message based on status code
func getStatusMessage(code int) string {
    switch code {
    case 502:
        return "Can't get a response from the container handling this route. Check the service logs for errors."
    case 503:
        return "Currently there is no service configured to respond to this request."
    default:
        return ""
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
    title := http.StatusText(code) // Default message from net/http

    // Render the HTML template
    data := StatusData{
        Code:    code,
        Title:   title,
        Message: message,
    }
    if err := tmpl.Execute(w, data); err != nil {
        log.Printf("Template execution error: %v", err)
    }
}

func main() {
    listenAddress := flag.String("listen-address", ":63108", "Address and port to listen on (e.g., :63108 or 0.0.0.0:63108)")
    flag.Parse()

    http.HandleFunc("/", statusHandler)
    fmt.Printf("Listening on %s\n", *listenAddress)
    log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
