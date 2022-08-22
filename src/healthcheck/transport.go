package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

func MakeHealthCheckHandler(c mongo.Client) http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").Path("/health").HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		defer cancel()

		err := c.Ping(ctx, readpref.Primary())

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(healthCheckResponse{Status: "ok"})
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(json)
	})

	return r
}
