package presence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	defaultUserID = "531484240114876416"
	defaultAPIURL = "https://api.lanyard.rest/v1/users/%s"
)

type Activity struct {
	Type       int    `json:"type"`
	Name       string `json:"name"`
	Details    string `json:"details"`
	State      string `json:"state"`
	Timestamps struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	} `json:"timestamps"`
}

type Spotify struct {
	Song        string `json:"song"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	AlbumArtURL string `json:"album_art_url"`
	Timestamps  struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	} `json:"timestamps"`
}

type Presence struct {
	Status             string
	ListeningToSpotify bool
	Spotify            *Spotify
	Activities         []Activity
	UpdatedAt          time.Time
}

func Fetch(ctx context.Context) (Presence, error) {
	userID := os.Getenv("LANYARD_USER_ID")
	if userID == "" {
		userID = defaultUserID
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(defaultAPIURL, userID), nil)
	if err != nil {
		return Presence{}, err
	}

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return Presence{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Presence{}, fmt.Errorf("lanyard returned %s", res.Status)
	}

	var payload struct {
		Success bool `json:"success"`
		Data    struct {
			DiscordStatus      string     `json:"discord_status"`
			ListeningToSpotify bool       `json:"listening_to_spotify"`
			Spotify            *Spotify   `json:"spotify"`
			Activities         []Activity `json:"activities"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return Presence{}, err
	}
	if !payload.Success {
		return Presence{}, fmt.Errorf("lanyard response was not successful")
	}

	return Presence{
		Status:             payload.Data.DiscordStatus,
		ListeningToSpotify: payload.Data.ListeningToSpotify && hasCompleteSpotify(payload.Data.Spotify),
		Spotify:            completeSpotify(payload.Data.Spotify),
		Activities:         visibleActivities(payload.Data.Activities),
		UpdatedAt:          time.Now(),
	}, nil
}

func (p Presence) Game() *Activity {
	for _, activity := range p.Activities {
		if activity.Type == 0 {
			return &activity
		}
	}
	return nil
}

func (p Presence) VisibleActivities(limit int) []Activity {
	if limit < 1 {
		return nil
	}

	activities := make([]Activity, 0, limit)
	for _, activity := range p.Activities {
		if activity.Name == "" || activity.Name == "Spotify" {
			continue
		}
		if activity.Type == 2 {
			continue
		}
		activities = append(activities, activity)
		if len(activities) == limit {
			return activities
		}
	}
	return activities
}

func (p Presence) FirstVisibleActivity() *Activity {
	activities := p.VisibleActivities(1)
	if len(activities) == 0 {
		return nil
	}
	return &activities[0]
}

func visibleActivities(activities []Activity) []Activity {
	filtered := make([]Activity, 0, len(activities))
	for _, activity := range activities {
		if activity.Type == 4 || activity.Name == "" {
			continue
		}
		filtered = append(filtered, activity)
	}
	return filtered
}

func hasCompleteSpotify(spotify *Spotify) bool {
	return spotify != nil &&
		spotify.Song != "" &&
		spotify.Artist != "" &&
		spotify.Album != "" &&
		spotify.AlbumArtURL != "" &&
		spotify.Timestamps.Start != 0 &&
		spotify.Timestamps.End != 0
}

func completeSpotify(spotify *Spotify) *Spotify {
	if !hasCompleteSpotify(spotify) {
		return nil
	}
	return spotify
}
