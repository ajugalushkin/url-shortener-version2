// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto(in *jlexer.Lexer, out *ShorteningList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ShorteningList, 0, 1)
			} else {
				*out = ShorteningList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Shortening
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto(out *jwriter.Writer, in ShorteningList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ShorteningList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShorteningList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShorteningList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShorteningList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto1(in *jlexer.Lexer, out *Shortening) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "correlation_id":
			out.CorrelationId = string(in.String())
		case "short_url":
			out.ShortURL = string(in.String())
		case "original_url":
			out.OriginalURL = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto1(out *jwriter.Writer, in Shortening) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"correlation_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.CorrelationId))
	}
	{
		const prefix string = ",\"short_url\":"
		out.RawString(prefix)
		out.String(string(in.ShortURL))
	}
	{
		const prefix string = ",\"original_url\":"
		out.RawString(prefix)
		out.String(string(in.OriginalURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Shortening) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Shortening) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Shortening) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Shortening) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto1(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto2(in *jlexer.Lexer, out *ShortenOutput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "result":
			out.Result = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto2(out *jwriter.Writer, in ShortenOutput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"result\":"
		out.RawString(prefix[1:])
		out.String(string(in.Result))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ShortenOutput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenOutput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenOutput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenOutput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto2(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto3(in *jlexer.Lexer, out *ShortenInput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "url":
			out.URL = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto3(out *jwriter.Writer, in ShortenInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ShortenInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto3(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto4(in *jlexer.Lexer, out *ShortenListOutput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ShortenListOutput, 0, 2)
			} else {
				*out = ShortenListOutput{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 struct {
				CorrelationId string `json:"correlation_id"`
				ShortUrl      string `json:"short_url"`
			}
			easyjsonC80ae7adDecode(in, &v4)
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto4(out *jwriter.Writer, in ShortenListOutput) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			easyjsonC80ae7adEncode(out, v6)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ShortenListOutput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenListOutput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenListOutput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenListOutput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto4(l, v)
}
func easyjsonC80ae7adDecode(in *jlexer.Lexer, out *struct {
	CorrelationId string `json:"correlation_id"`
	ShortUrl      string `json:"short_url"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "correlation_id":
			out.CorrelationId = string(in.String())
		case "short_url":
			out.ShortUrl = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncode(out *jwriter.Writer, in struct {
	CorrelationId string `json:"correlation_id"`
	ShortUrl      string `json:"short_url"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"correlation_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.CorrelationId))
	}
	{
		const prefix string = ",\"short_url\":"
		out.RawString(prefix)
		out.String(string(in.ShortUrl))
	}
	out.RawByte('}')
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto5(in *jlexer.Lexer, out *ShortenListInput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ShortenListInput, 0, 2)
			} else {
				*out = ShortenListInput{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 struct {
				CorrelationId string `json:"correlation_id"`
				OriginalURL   string `json:"original_url"`
			}
			easyjsonC80ae7adDecode1(in, &v7)
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto5(out *jwriter.Writer, in ShortenListInput) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			easyjsonC80ae7adEncode1(out, v9)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ShortenListInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortenListInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinUrlShortenerVersion2InternalDto5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortenListInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortenListInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinUrlShortenerVersion2InternalDto5(l, v)
}
func easyjsonC80ae7adDecode1(in *jlexer.Lexer, out *struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "correlation_id":
			out.CorrelationId = string(in.String())
		case "original_url":
			out.OriginalURL = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncode1(out *jwriter.Writer, in struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"correlation_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.CorrelationId))
	}
	{
		const prefix string = ",\"original_url\":"
		out.RawString(prefix)
		out.String(string(in.OriginalURL))
	}
	out.RawByte('}')
}
