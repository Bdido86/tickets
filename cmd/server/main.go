package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	apiPkg "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	postgre "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgre"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	tokenHeader = "Token"
	authPathRPC = "UserAuth"

	swaggerDir = "./third_party/swagger-ui"
)

var depsRepo apiPkg.Deps

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := config.GetConfig()
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DbHost(), c.DbPort(), c.DbUser(), c.DbPassword(), c.DbName())
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	defer pool.Close()

	// настраиваем
	//config := pool.Config()
	//config.MaxConnIdleTime = MaxConnIdleTime
	//config.MaxConnLifetime = MaxConnLifetime
	//config.MinConns = MinConns
	//config.MaxConns = MaxConns

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	serverAddress := ":" + c.ServerPort()
	restAddress := ":" + c.RestPort()

	go runREST(ctx, serverAddress, restAddress, c.RequestTimeOutInMilliSecond())
	runGRPCServer(ctx, pool, serverAddress)
}

func runGRPCServer(_ context.Context, pool *pgxpool.Pool, serverAddress string) {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error server connect tcp: %v", err)
	}
	defer listener.Close()

	depsRepo = apiPkg.Deps{CinemaRepository: postgre.NewRepository(pool)}

	option := grpc.UnaryInterceptor(AuthInterceptor)
	grpcServer := grpc.NewServer(option)
	pb.RegisterCinemaServer(grpcServer, apiPkg.NewServer(depsRepo))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST(ctx context.Context, serverAddress, restAddress string, requestTimeoutInMilliSecond time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	runtime.DefaultContextTimeout = requestTimeoutInMilliSecond
	rmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
	)

	// Swagger server
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	fs := http.FileServer(http.Dir(swaggerDir))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterCinemaHandlerFromEndpoint(ctx, rmux, serverAddress, opts); err != nil {
		panic(err)
	}

	log.Println("Serving gRPC-Gateway on " + restAddress)
	log.Fatalln(http.ListenAndServe(restAddress, mux))
}

func headerMatcherREST(key string) (string, bool) {
	switch key {
	case tokenHeader:
		return key, true
	default:
		return key, false
	}
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if isAuthPath(info.FullMethod) {
		return handler(ctx, req)
	}

	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Header [token] is required")
	}

	tokens := metaData.Get(tokenHeader)
	if len(tokens) == 0 {
		return nil, status.Error(codes.PermissionDenied, "Header [token] is required")
	}

	for _, token := range tokens {
		userId, err := depsRepo.CinemaRepository.GetUserIdByToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "Header [Token] is invalid. See 'auth' method")
		}

		ctx = context.WithValue(ctx, "userId", userId)
		return handler(ctx, req)
	}

	return nil, status.Error(codes.PermissionDenied, "Header [Token] is invalid. See 'auth' method")
}

func isAuthPath(fullMethod string) bool {
	paths := strings.Split(fullMethod, "/")
	for _, path := range paths {
		if path == authPathRPC {
			return true
		}
	}
	return false
}
