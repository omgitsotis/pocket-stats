package database

import (
	"strconv"
	"time"

	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/sirupsen/logrus"
)

const (
	TagRead  = "_read"
	TagSport = "_sport"
	TagNews  = "_news"
)

// Article is the object that will be saved in the database
type Article struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Tag       string `json:"tag"`
	WordCount int64  `json:"word_count"`
	DateAdded int64  `json:"date_added"`
	DateRead  int64  `json:"date_read"`
}

// ConvertArticles converts the pocket article to the database article
func ConvertArticles(pa pocket.Article) Article {
	var (
		added, read int
		err         error
	)

	// Convert the added time to a number
	added, err = strconv.Atoi(pa.TimeAdded)
	if err != nil {
		logrus.Errorf(
			"pocket article [%d] has a bad added time [%s] the fuck?",
			pa.ItemID,
			pa.TimeAdded,
		)
	}

	// Check we have a read time, and then convert that to a number
	if pa.TimeRead != "" {
		read, err = strconv.Atoi(pa.TimeRead)
		if err != nil {
			logrus.Errorf(
				"pocket article [%d] has a bad read time [%s] the fuck?",
				pa.ItemID,
				pa.TimeRead,
			)
		}
	}

	// Remove the time part of the read and added times
	sAdded := StripTime(added)
	sRead := StripTime(read)

	var tag string
	for key := range pa.Tags {
		if key == TagRead || key == TagSport || key == TagNews {
			continue
		}

		tag = key
		break
	}

	return Article{
		ID:        pa.ItemID,
		URL:       pa.ResolvedURL,
		Title:     pa.ResolvedTitle,
		WordCount: int64(pa.WordCount),
		DateAdded: sAdded,
		DateRead:  sRead,
		Tag:       tag,
	}
}

// StripTime converts the time part of a time stamp into 00:00:00.
func StripTime(fullTime int) int64 {
	if fullTime == 0 {
		return 0
	}

	t := time.Unix(int64(fullTime), 0)
	stripped := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	return stripped.Unix()
}
