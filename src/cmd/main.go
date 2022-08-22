package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"zeabix.com/blog-service/blog"
	"zeabix.com/blog-service/healthcheck"
)

func main() {

	var (
		httpAddr    = flag.String("http.addr", ":8080", "HTTP listen address")
		mongoUrl    = os.Getenv("MONGO_CONNNECTION_URL") //flag.String("mongo.url", "mongodb://localhost:27017", "Connection URL for mongodb")
		mongoDbname = os.Getenv("MONGO_DATABASE_NAME")   //flag.String("mongo.dbname", "blogs", "Mongo Database name")
		mongoCol    = os.Getenv("MONGO_COLLECTION_NAME") //flag.String("mongo.colname", "blogs", "Mongo Collection name")
	)
	flag.Parse()

	var ()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	//defer cancel()

	ctx := context.TODO()
	client, err := makeMongoClient(ctx, mongoUrl)

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Log("Unable to connect to DB, shutdown")
		panic("Unable to connect to DB")
	}

	col := client.Database(mongoDbname).Collection(mongoCol)

	if err != nil {
		panic(err)
	}

	fieldKeys := []string{"method"}

	var s blog.Service
	{
		s = blog.NewMongoBlogService(*col)
		s = blog.NewLoggingMiddleware(logger, s)
		s = blog.NewInstrucmentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "blog_service",
				Name:      "request_count",
				Help:      "number of requests recieved",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "blog_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			s,
		)
	}

	var h http.Handler
	{
		h = blog.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	mux := http.NewServeMux()

	mux.Handle("/blogs/v1/", h)
	mux.Handle("/health", healthcheck.MakeHealthCheckHandler(*client))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()

	logger.Log("exit", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func makeMongoCollection(ctx context.Context, url string, dbname string, collection string) (*mongo.Collection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	return client.Database(dbname).Collection(collection), nil
}

func makeMongoClient(ctx context.Context, url string) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(url))
}
