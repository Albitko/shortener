package grpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo/postgres"
	"github.com/Albitko/shortener/internal/workers"
	pb "github.com/Albitko/shortener/proto"
)

type urlConverter interface {
	URLToID(context.Context, entity.OriginalURL, string) (entity.URLID, error)
	IDToURL(context.Context, entity.URLID) (entity.OriginalURL, error)
	UserIDToURLs(c context.Context, userID string) (map[string]string, bool)
	PingDB() error
	GetStats(context.Context) (entity.URLStats, error)
}

type grpcHandlers struct {
	pb.UnimplementedShortenerServer
	uc             urlConverter
	baseURL        string
	trustedNetwork string
	q              workers.Queue
}

func getUser(ctx context.Context) string {
	var user string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("user")
		if len(values) > 0 {
			user = values[0]
		}
	}
	return user
}

// ShortenURL short URL
func (g *grpcHandlers) ShortenURL(ctx context.Context, in *pb.PostURLRequest) (*pb.PostURLResponse, error) {
	var resp *pb.PostURLResponse
	user := getUser(ctx)
	shortURL, urlError := g.uc.URLToID(ctx, entity.OriginalURL(in.BaseUrl), user)

	if errors.Is(urlError, postgres.ErrURLAlreadyExists) {
		return resp, urlError
	}
	return &pb.PostURLResponse{
		ShortUrl: g.baseURL + string(shortURL),
	}, nil
}

// ShortenURLBatch short multiple URLs
func (g *grpcHandlers) ShortenURLBatch(
	ctx context.Context, in *pb.PostURLBatchRequest,
) (*pb.PostURLBatchResponse, error) {
	var resp *pb.PostURLBatchResponse
	var shortenURL entity.ModelURLBatchResponse

	response := make([]entity.ModelURLBatchResponse, 0, len(in.RequestUrls))
	user := getUser(ctx)
	for _, requestBatchURL := range in.RequestUrls {
		shortURL, urlError := g.uc.URLToID(ctx, entity.OriginalURL(requestBatchURL.Url), user)
		if errors.Is(urlError, postgres.ErrURLAlreadyExists) {
			return resp, urlError
		}
		shortenURL.ShortURL = g.baseURL + string(shortURL)
		shortenURL.CorrelationID = requestBatchURL.CorrelationId
		response = append(response, shortenURL)
	}

	for i, val := range response {
		resp.ResponseUrls[i].Url = val.ShortURL
		resp.ResponseUrls[i].CorrelationId = val.CorrelationID
	}

	return resp, nil
}

// DeleteURLBatch delete multiple shorten URLs
func (g *grpcHandlers) DeleteURLBatch(
	ctx context.Context, in *pb.DeleteURLBatchRequest,
) (*pb.DeleteURLBatchResponse, error) {
	var resp *pb.DeleteURLBatchResponse
	urlsForDelete := make([]string, 0, len(in.RequestUrls.Urls))

	user := getUser(ctx)

	urlsForDelete = append(urlsForDelete, in.RequestUrls.Urls...)
	g.q.Push(&workers.Task{UserID: user, IDsForDelete: urlsForDelete})

	return resp, nil
}

// GetURL return original URL by the short one
func (g *grpcHandlers) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	var resp *pb.GetURLResponse

	originalURL, err := g.uc.IDToURL(ctx, entity.URLID(in.ShortUrlId))
	if err != nil {
		return resp, err
	}
	fmt.Println(originalURL)
	return &pb.GetURLResponse{
		RedirectTo: string(originalURL),
	}, nil
}

// Ping check db connection
func (g *grpcHandlers) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	var resp *pb.PingResponse
	err := g.uc.PingDB()
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetStats return count of urls and users
func (g *grpcHandlers) GetStats(ctx context.Context, in *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	serviceStats, err := g.uc.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetStatsResponse{
		Users: int64(serviceStats.UsersCount),
		Urls:  int64(serviceStats.URLsCount),
	}, nil
}

// GetURLsByUserID return all urls for user
func (g *grpcHandlers) GetURLsByUserID(
	ctx context.Context, in *pb.GetURLsByUserIDRequest,
) (*pb.GetURLsByUserIDResponse, error) {
	var resp *pb.GetURLsByUserIDResponse

	user := getUser(ctx)
	if user == "" {
		return resp, fmt.Errorf("no user in metadata")
	}

	userURLs, ok := g.uc.UserIDToURLs(ctx, user)
	if ok {
		for shortURL, originalURL := range userURLs {
			responseURL := pb.ResponseURLs{
				BaseUrl:  originalURL,
				ShortUrl: g.baseURL + shortURL,
			}
			resp.ResponseUrls = append(resp.ResponseUrls, &responseURL)
		}
	} else {
		return resp, fmt.Errorf("couldn't find urls for this user")
	}

	return resp, nil
}

// New create instance of `grpcHandlers` struct
func New(u urlConverter, conf entity.Config, queue *workers.Queue) *grpcHandlers {
	baseURL := "http://localhost:8080/"
	if conf.BaseURL != "" {
		baseURL = conf.BaseURL + "/"
	}

	return &grpcHandlers{
		uc:             u,
		baseURL:        baseURL,
		trustedNetwork: conf.TrustedSubnet,
		q:              *queue,
	}
}
