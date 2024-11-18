package server

import (
	"fmt"
	"html"
	"net/http"

	"github.com/vehbiu/dnsly/internal/config"
	"github.com/vehbiu/dnsly/internal/models"
)

type Server struct {
	cfg             *config.Config
	blockedRequests []models.BlockedRequest
}

func NewServer(cfg *config.Config, blockedRequests []models.BlockedRequest) *Server {
	return &Server{
		cfg:             cfg,
		blockedRequests: blockedRequests,
	}
}

func (s *Server) handleWebUI(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	if password != s.cfg.HTTPPassword {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`
			<html>
			<head>
				<title>DNSly - DNS Blocker - Login</title>
				<style>
					body { font-family: Arial, sans-serif; margin: 40px; }
					.login-form { max-width: 400px; margin: 0 auto; }
					.error { color: red; margin-bottom: 20px; }
					input[type="password"] { 
						width: 100%; 
						padding: 8px; 
						margin: 10px 0; 
						border: 1px solid #ddd; 
					}
					button { 
						background-color: #4CAF50; 
						color: white; 
						padding: 10px 15px; 
						border: none; 
						cursor: pointer; 
					}
					button:hover { background-color: #45a049; }
				</style>
			</head>
			<body>
				<div class="login-form">
					<h1><a href="https://github.com/vehbiu/dnsly">DNSly - DNS Blocker</a> Login</h1>
					` + func() string {
			if password != "" {
				return `<div class="error">Invalid password</div>`
			}
			return ""
		}() + `
					<form>
						<input type="password" name="password" placeholder="Enter password" required>
						<button type="submit">Login</button>
					</form>
				</div>
			</body>
			</html>
		`))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<html>
		<head>
			<title>DNSly - DNS Blocker Status</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 20px; }
				table { border-collapse: collapse; width: 100%; margin-top: 20px; }
				th, td { padding: 12px; text-align: left; border: 1px solid #ddd; }
				th { background-color: #f2f2f2; }
				tr:nth-child(even) { background-color: #f9f9f9; }
				.header { 
					display: flex; 
					justify-content: space-between; 
					align-items: center;
					margin-bottom: 20px;
				}
				.logout {
					background-color: #f44336;
					color: white;
					padding: 10px 15px;
					text-decoration: none;
					border-radius: 4px;
				}
				.logout:hover {
					background-color: #da190b;
				}
			</style>
		</head>
		<body>
			<div class="header">
				<h1>Blocked DNS Requests</h1>
				<a href="/" class="logout">Logout</a>
			</div>
			<table>
				<tr>
					<th>Domain</th>
					<th>Blocked IP</th>
					<th>Source</th>
					<th>Time</th>
				</tr>
	`))

	for _, req := range s.blockedRequests {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td></tr>",
			html.EscapeString(req.Domain),
			html.EscapeString(req.IP),
			html.EscapeString(req.Source))
	}

	w.Write([]byte("</table></body></html>"))
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`status: ok`))

}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.handleWebUI)
	http.HandleFunc("/status", s.handleStatus)
	return http.ListenAndServe(addr, nil)
}
