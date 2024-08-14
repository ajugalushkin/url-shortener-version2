package v1

import (
	"context"
	"errors"

	url_shortenerv1 "github.com/ajugalushkin/url-shortener-version2/gen/url_shortener/v1"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	userErr "github.com/ajugalushkin/url-shortener-version2/internal/errors"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
)

// URLSServer Структура
type URLSServer struct {
	url_shortenerv1.UnimplementedURLShortenerV1ServiceServer
	ctx     context.Context
	servAPI *service.Service
}

// NewHandler конструктор
func NewHandler(ctx context.Context, servAPI *service.Service) *URLSServer {
	return &URLSServer{
		ctx:     ctx,
		servAPI: servAPI}
}

// ShortenV1 метод реализует метод Shorten
func (s *URLSServer) ShortenV1(ctx context.Context, in *url_shortenerv1.ShortenV1Request) (*url_shortenerv1.ShortenV1Response, error) {
	var response url_shortenerv1.ShortenV1Response

	shorten, err := s.servAPI.Shorten(s.ctx, dto.Shortening{
		OriginalURL: in.Input.Url,
		UserID:      ""})

	if err != nil {
		if errors.Is(err, userErr.ErrorDuplicateURL) {
			return nil, err
		}
		return nil, err
	}
	response.Output.ShortUrl = shorten.ShortURL
	return &response, nil
}

// ShortenBatchV1 метод реализует метод ShortenBatch
func (s *URLSServer) ShortenBatchV1(ctx context.Context, in *url_shortenerv1.ShortenBatchV1Request) (*url_shortenerv1.ShortenBatchV1Response, error) {
	var response url_shortenerv1.ShortenBatchV1Response

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
		response.Output = append(response.Output, &url_shortenerv1.ShortenBatchV1Response_ShortenBatchOutput{
			CorrelationId: output.CorrelationID,
			ShortUrl:      output.ShortURL})
	}
	return &response, nil
}

// GetV1 метод реализует метод Get GRPC
func (s *URLSServer) GetV1(ctx context.Context, in *url_shortenerv1.GetV1Request) (*url_shortenerv1.GetV1Response, error) {
	var response url_shortenerv1.GetV1Response

	redirect, err := s.servAPI.Redirect(s.ctx, in.ShortUrl)
	if err != nil {
		return nil, err
	}

	response.OriginalUrl = redirect.OriginalURL

	return &response, nil
}

// PingV1 метод реализует метод Ping GRPC
func (s *URLSServer) PingV1(ctx context.Context, in *url_shortenerv1.PingV1Request) (*url_shortenerv1.PingV1Response, error) {
	var response url_shortenerv1.PingV1Response

	return &response, nil
}

// UserUrlsV1 метод реализует метод UserUrls GRPC
func (s *URLSServer) UserUrlsV1(ctx context.Context, in *url_shortenerv1.UserUrlsV1Request) (*url_shortenerv1.UserUrlsV1Response, error) {
	var response url_shortenerv1.UserUrlsV1Response

	shortList, err := s.servAPI.GetUserURLS(s.ctx, 0)
	if err != nil || len(*shortList) == 0 {
		return nil, err
	}

	for _, shortURL := range *shortList {
		response.Output = append(response.Output, &url_shortenerv1.UserUrlsV1Response_UserUrls{
			OriginalUrl: shortURL.OriginalURL,
			ShortUrl:    shortURL.ShortURL,
		})
	}

	return &response, nil
}

// UserUrlsDeleteV1 метод реализует метод UserUrlsDelete GRPC
func (s *URLSServer) UserUrlsDeleteV1(ctx context.Context, in *url_shortenerv1.UserUrlsDeleteV1Request) (*url_shortenerv1.UserUrlsDeleteV1Response, error) {
	var response url_shortenerv1.UserUrlsDeleteV1Response

	s.servAPI.DeleteUserURL(s.ctx, in.Urls, 0)

	response.Result = "URLS Delete OK"

	return &response, nil
}

// StatsV1 метод реализует метод Stats GRPC
func (s *URLSServer) StatsV1(ctx context.Context, in *url_shortenerv1.StatsV1Request) (*url_shortenerv1.StatsV1Response, error) {
	stats := s.servAPI.GetStats(s.ctx)

	response := url_shortenerv1.StatsV1Response{
		Urls:  int64(stats.URLS),
		Users: int64(stats.Users),
	}
	return &response, nil
}
