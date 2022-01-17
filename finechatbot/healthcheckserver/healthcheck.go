package healthcheckserver

import (
	"net/http"

	"github.com/disturb16/finechatbot/broker"
)

func New(b broker.MessageBroker) *http.Server {
	return &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/healthcheck" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if b.IsClosed() {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
		}),
	}
}
