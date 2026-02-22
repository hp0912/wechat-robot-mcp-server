package tools

import (
	"context"
	"testing"
)

func TestRequestSong(t *testing.T) {
	_, _, err := RequestSong(context.Background(), nil, &RequestSongInput{
		SongTitle: "夜曲",
		Singer:    "周杰伦",
	})
	if err != nil {
		t.Errorf("RequestSong failed: %v", err)
	}
}
