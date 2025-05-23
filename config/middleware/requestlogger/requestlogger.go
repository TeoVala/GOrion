package requestlogger

import (
	"GOrion/internal/helpers/terminal"
	"fmt"
	"log"
	"net/http"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// TODO THIS CAN BE SIMPLIFIED SO WE DONT HAVE BOTH LOGMSG and then print again
		// PUT IT IN A SLICE AND GET IT FROM THERE

		// Print request method and URL to the console
		var logMsg string

		// TODO Check if tls and set accordingly
		var tls string = "http"

		logMsg += "\n--------------- Request Details ---------------\n"
		logMsg += fmt.Sprintf("Method: %s\n", r.Method)
		logMsg += fmt.Sprintf("URL: %s\n", r.URL.RequestURI())
		logMsg += fmt.Sprintf("Host: %s://%s\n", tls, r.Host)
		logMsg += fmt.Sprintf("Protocol: %s\n", r.Proto)        // HTTP/1.1, HTTP/2.0, etc.
		logMsg += fmt.Sprintf("RemoteAddr: %s\n", r.RemoteAddr) // Client IP and portlogMsg +=
		logMsg += "Headers:"
		for name, values := range r.Header {
			for _, value := range values {
				logMsg += fmt.Sprintf("  %s: %s\n", name, value)
			}
		}

		logMsg += "----------------------------------------------"

		log.Printf(logMsg)

		terminal.CW(true, terminal.NWhite, "\n--------------- Request Details ---------------\n")

		terminal.CW(true, terminal.NBlue, "Method: %s\n", r.Method)
		terminal.CW(true, terminal.BMagenta, "URL: %s\n", r.URL.RequestURI())
		terminal.CW(true, terminal.BCyan, "Host: %s://%s\n", tls, r.Host)
		terminal.CW(true, terminal.NYellow, "Protocol: %s\n", r.Proto)
		terminal.CW(true, terminal.BRed, "RemoteAddr: %s\n", r.RemoteAddr)

		terminal.CW(true, terminal.NGreen, "Headers:\n")
		for name, values := range r.Header {
			for _, value := range values {
				terminal.CW(true, terminal.BGreen, "  %s: %s\n", name, value)
			}
		}

		terminal.CW(false, terminal.NWhite, "----------------------------------------------\n")

		next.ServeHTTP(w, r)
	})
}
