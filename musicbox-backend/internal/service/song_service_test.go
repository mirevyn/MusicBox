package service

import "testing"

func TestSongOrderClauseAddsDeterministicIDTieBreaker(t *testing.T) {
	tests := []struct {
		name   string
		sortBy string
		order  string
		want   string
	}{
		{name: "upload time", sortBy: "upload_at", order: "desc", want: "upload_at desc, id desc"},
		{name: "title", sortBy: "title", order: "asc", want: "title asc, id asc"},
		{name: "id", sortBy: "id", order: "desc", want: "id desc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := songOrderClause(tt.sortBy, tt.order); got != tt.want {
				t.Fatalf("songOrderClause(%q, %q) = %q, want %q", tt.sortBy, tt.order, got, tt.want)
			}
		})
	}
}
