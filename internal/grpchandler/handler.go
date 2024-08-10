package grpchandler

import (
	"context"
	"errors"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	userErr "github.com/ajugalushkin/url-shortener-version2/internal/errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/proto"
)

type URLSServer struct {
	proto.UnimplementedURLShortenerServiceServer
	ctx     context.Context
	cache   map[string]*dto.User
	servAPI *service.Service
}

// NewHandler конструктор
func NewHandler(ctx context.Context, servAPI *service.Service) *URLSServer {
	return &URLSServer{
		ctx:     ctx,
		cache:   make(map[string]*dto.User),
		servAPI: servAPI}
}

func (s *URLSServer) Save(ctx context.Context, in *proto.SaveRequest) (*proto.SaveResponse, error) {
	var response proto.SaveResponse

	shorten, err := s.servAPI.Shorten(s.ctx, dto.Shortening{
		OriginalURL: in.Url,
		UserID:      ""})

	if err != nil {
		if errors.Is(err, userErr.ErrorDuplicateURL) {
			return nil, err
		}
		return nil, err
	}
	response.ShortUrl = shorten.ShortURL
	return &response, nil
}

func (s *URLSServer) Shorten(ctx context.Context, in *proto.ShortenRequest) (*proto.ShortenResponse, error) {
	var response proto.ShortenResponse

	shortenURL, err := s.servAPI.Shorten(s.ctx, dto.Shortening{OriginalURL: in.Input.GetUrl()})
	if err != nil {
		if errors.Is(err, userErr.ErrorDuplicateURL) {
			return nil, err
		}
		return nil, err
	}

	response.Output.ShortUrl = shortenURL.ShortURL
	return &response, nil
}

func (s *URLSServer) ShortenBatch(ctx context.Context, in *proto.ShortenBatchRequest) (*proto.ShortenBatchResponse, error) {
	var response proto.ShortenBatchResponse

	inputList := in.GetInput()
	var inputListParse dto.ShortenListInput

	for _, input := range inputList {
		inputListParse = append(inputListParse, dto.ShortenListItemInput{
			CorrelationID: input.CorrelationId,
			OriginalURL:   input.OriginalUrl})
	}

	listOutput, err := s.servAPI.ShortenList(s.ctx, inputListParse)
	if err != nil {
		return nil, err
	}

	for _, output := range *listOutput {
		response.Output = append(response.Output, &proto.ShortenBatchResponse_ShortenBatchOutput{
			CorrelationId: output.CorrelationID,
			ShortUrl:      output.ShortURL})
	}
	return &response, nil
}

func (s *URLSServer) Redirect(ctx context.Context, in *proto.RedirectRequest) (*proto.RedirectResponse, error) {
	var response proto.RedirectResponse

	redirect, err := s.servAPI.Redirect(s.ctx, in.ShortUrl)
	if err != nil {
		return nil, err
	}

	response.OriginalUrl = redirect.OriginalURL

	return &response, nil
}

func (s *URLSServer) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	var response proto.PingResponse

	return &response, nil
}

func (s *URLSServer) UserUrls(ctx context.Context, in *proto.UserUrlsRequest) (*proto.UserUrlsResponse, error) {
	var response proto.UserUrlsResponse

	shortList, err := s.servAPI.GetUserURLS(s.ctx, 0)
	if err != nil || len(*shortList) == 0 {
		return nil, err
	}

	for _, shortURL := range *shortList {
		response.Output = append(response.Output, &proto.UserUrlsResponse_UserUrls{
			OriginalUrl: shortURL.OriginalURL,
			ShortUrl:    shortURL.ShortURL,
		})
	}

	return &response, nil
}

func (s *URLSServer) UserUrlsDelete(ctx context.Context, in *proto.UserUrlsDeleteRequest) (*proto.UserUrlsDeleteResponse, error) {
	var response proto.UserUrlsDeleteResponse

	s.servAPI.DeleteUserURL(s.ctx, in.Urls, 0)

	response.Result = "URLS Delete OK"

	return &response, nil
}

func (s *URLSServer) Stats(ctx context.Context, in *proto.StatsRequest) (*proto.StatsResponse, error) {
	stats := s.servAPI.GetStats(s.ctx)

	response := proto.StatsResponse{
		Urls:  int64(stats.URLS),
		Users: int64(stats.Users),
	}
	return &response, nil
}
