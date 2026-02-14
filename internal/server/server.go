package server

import (
	"encoding/json"
	"io/fs"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
	"visitor/internal/dashboard"
	"visitor/internal/geoip"
	"visitor/internal/hash"
	"visitor/internal/model"
	"visitor/internal/storage"
	"visitor/web"

	"github.com/mssola/useragent"
)

var screenSizeRe = regexp.MustCompile(`^\d+x\d+$`)

type Server struct {
	addr			string
	db 				*storage.DB
	hasher 			*hash.Manager
	mux 			*http.ServeMux
	geoip			*geoip.Resolver
	password 		string
	allowedDomains 	map[string]bool
	limiter			*rateLimiter
}

func New(addr string, db *storage.DB, hasher *hash.Manager, geoip *geoip.Resolver,password string, allowedDomains string) *Server {
	mux := http.NewServeMux()

	domains := make(map[string]bool)
	for d := range strings.SplitSeq(allowedDomains, ",") {
		d = strings.TrimSpace(d)
		if d != "" {
			domains[d] = true
		}
	}

	s := &Server{
		addr: 			addr,
		db: 			db,
		hasher: 		hasher,
		mux: 			mux,
		geoip: 			geoip,
		password: 		password,
		allowedDomains: domains,
		limiter: 		newRateLimiter(5, 10),
	}

	s.mux.Handle("POST /api/event", s.limiter.middleware(http.HandlerFunc(s.handleEvent)))
	s.mux.HandleFunc("GET /tracker.js", s.handleTracker)

	dash := dashboard.NewHandler(dashboard.NewQueries(db.Pool()))
	s.mux.Handle("GET /api/stats/summary", s.auth(http.HandlerFunc(dash.HandleSummary)))
	s.mux.Handle("GET /api/stats/pages", s.auth(http.HandlerFunc(dash.HandlePages)))
	s.mux.Handle("GET /api/stats/referrers", s.auth(http.HandlerFunc(dash.HandleReferrers)))
	s.mux.Handle("GET /api/stats/locations", s.auth(http.HandlerFunc(dash.HandleLocations)))
	s.mux.Handle("GET /api/stats/sizes", s.auth(http.HandlerFunc(dash.HandleSizes)))
	s.mux.Handle("GET /api/stats/browsers", s.auth(http.HandlerFunc(dash.HandleBrowsers)))
	s.mux.Handle("GET /api/stats/systems", s.auth(http.HandlerFunc(dash.HandleSystems)))



	staticFS, _ := fs.Sub(web.StaticFS, "static")
	s.mux.Handle("GET /static/", s.auth(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	s.mux.Handle("GET /dashboard", s.auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		data, _ := web.StaticFS.ReadFile("static/dashboard.html")
		w.Write(data)
	})))

	return s
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr: 			s.addr,
		Handler: 		s.cors(s.mux),
		ReadTimeout: 	5 * time.Second,
		WriteTimeout: 	10 * time.Second,
		IdleTimeout: 	120 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<10) // 10KB

	var event model.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if !validateInputData(event.Domain, event.Path, event.Referrer, event.ScreenSize) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if !s.isAllowedDomain(event.Domain) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	origin := r.Header.Get("Origin")
	if origin != "http://"+event.Domain && origin != "https://"+event.Domain {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	userAgent := r.Header.Get("User-Agent")

	visitorHash, err := s.hasher.GetHash(r.Context(), event.Domain, ip, userAgent)
	if err != nil {
		log.Printf("failed to get visitor hash: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	countryCode := s.geoip.Country(ip)
	ua := useragent.New(userAgent)
	browser, _ := ua.Browser()
	os := ua.OS()

	pv := &model.PageView{
		Domain: 		event.Domain,
		Path: 			event.Path,
		Referrer: 		event.Referrer,
		ScreenSize: 	event.ScreenSize,
		CountryCode: 	countryCode,	
		Browser:        browser,
		OS: 			os,	
		VisitorHash: 	visitorHash,
	}

	if err := s.db.InsertPageView(r.Context(), pv); err != nil {
		log.Printf("Failed to insert page view: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleTracker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(web.TrackerJS)
}

func (s *Server) isAllowedDomain(domain string) bool {
    if len(s.allowedDomains) == 0 {
        return true
    }
    return s.allowedDomains[domain]
}

func validateInputData(domain string, path string, referrer string, screenSize string) bool {
	if domain == "" || len(domain) > 253 {
		return false
	}
	if !strings.HasPrefix(path, "/") || len(path) > 2048 {
		return false
	}
	if len(referrer) > 2048 {
		return false
	}
	if screenSize != "" && !screenSizeRe.MatchString(screenSize) {
		return false
	}
	return true
}