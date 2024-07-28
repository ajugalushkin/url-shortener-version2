package parse

// GetJSONDataFromBatch получение данных из контекста.
//func GetJSONDataFromBatch(ctx context.Context, echoCtx echo.Context) (dto.ShortenListInput, error) {
//	var shortList dto.ShortenListInput
//
//	body, err := io.ReadAll(echoCtx.Request().Body)
//	if err != nil {
//		return shortList, echoCtx.String(http.StatusBadRequest, validate.URLParseError)
//	}
//
//	err = shortList.UnmarshalJSON(body)
//	if err != nil {
//		return shortList, echoCtx.String(http.StatusBadRequest, validate.JSONParseError)
//	}
//
//	return shortList, nil
//}

// SetJSONDataToBody внесение данных в контекст.
//func SetJSONDataToBody(ctx context.Context, echoCtx echo.Context, list *dto.ShorteningList) ([]byte, error) {
//	var shortenListOut dto.ShortenListOutput
//	flag := config.GetConfig()
//	for _, item := range *list {
//		shortWithHost, _ := url.JoinPath(flag.BaseURL, item.ShortURL)
//		shortenListOut = append(
//			shortenListOut,
//			dto.ShortenListOutputLine{
//				CorrelationID: item.CorrelationID,
//				ShortURL:      shortWithHost,
//			},
//		)
//	}
//
//	newBody, err := shortenListOut.MarshalJSON()
//	if err != nil {
//		return newBody, echoCtx.String(http.StatusBadRequest, validate.JSONNotCreate)
//	}
//
//	return newBody, nil
//}
