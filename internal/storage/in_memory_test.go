package storage

import (
	"sync"
	"testing"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

func setupStorage(m *sync.Map) {
	m.Store("Xnrr2Mt", dto.Shortening{Key: "Xnrr2Mt", URL: "https://practicum.yandex.ru"})
}

func TestInMemory_Put(t *testing.T) {
	tests := []struct {
		name    string
		input   dto.Shortening
		want    *dto.Shortening
		wantErr bool
	}{
		{name: "Test Exists",
			input:   dto.Shortening{Key: "Xnrr2Mt", URL: "https://practicum.yandex.ru"},
			want:    &dto.Shortening{Key: "", URL: ""},
			wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &InMemory{}
			setupStorage(&s.m)

			_, err := s.Put(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, test.wantErr)
				return
			}
		})
	}
}
