package server

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/omgitsotis/pocket-stats/pkg/database"
)

func createTotalStats(start, end int64, articles []*database.Article) (*StatTotals, error) {
	st := StatTotals{}

	for _, a := range articles {
		log.Debugf("checking article [%d]", a.ID)
		// Check to see if the article was added in the date range
		if isInRange(start, end, a.DateAdded) {
			log.Debugf(
				"Article [%d]: Start date [%d] < Article add date [%d] < End date [%d]",
				a.ID, start, a.DateAdded, end,
			)

			timeReading := convertWordsToTime(a.WordCount)

			// Update total values
			st.ArticlesAdded++
			st.WordsAdded += a.WordCount
			st.TimeAdded += timeReading
		}

		// Check to see if the article is read
		if a.DateRead != 0 && isInRange(start, end, a.DateRead) {
			log.Debugf("Article [%d] read [%d]", a.ID, a.DateRead)
			timeReading := convertWordsToTime(a.WordCount)

			// Update total values
			st.ArticlesRead++
			st.WordsRead += a.WordCount
			st.TimeRead += timeReading
		}
	}

	return &st, nil
}

func createTagStats(start, end int64, articles []*database.Article) (*TagStats, error) {
	tags := make(TagStats)

	for _, a := range articles {
		log.Debugf("checking article [%d]", a.ID)

		if isInRange(start, end, a.DateAdded) {
			timeReading := convertWordsToTime(a.WordCount)
			if _, ok := tags[a.Tag]; !ok {
				tags[a.Tag] = &StatTotals{
					ArticlesAdded: 1,
					WordsRead:     a.WordCount,
					TimeAdded:     timeReading,
				}
			} else {
				tags[a.Tag].WordsAdded += a.WordCount
				tags[a.Tag].TimeAdded += timeReading
				tags[a.Tag].ArticlesAdded++
			}
		}

		// Check to see if the article is read
		if a.DateRead != 0 && isInRange(start, end, a.DateRead) {
			// Update the tag values
			timeReading := convertWordsToTime(a.WordCount)

			if _, ok := tags[a.Tag]; !ok {
				tags[a.Tag] = &StatTotals{
					ArticlesRead: 1,
					WordsRead:    a.WordCount,
					TimeRead:     timeReading,
				}
			} else {
				tags[a.Tag].ArticlesRead++
				tags[a.Tag].WordsRead += a.WordCount
				tags[a.Tag].TimeRead += timeReading
			}
		}
	}

	return &tags, nil
}

func createItemisedStats(start, end int64, articles []*database.Article) (*ItemisedStats, error) {
	itemised := make(ItemisedStats)

	// Populate the itemised map. We want all the dates in the range, including
	// the days with no updates
	t := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	for {
		itemised[t.Unix()] = &StatTotals{}
		t = t.AddDate(0, 0, 1)
		log.Debugf("%d greater than %d", t.Unix(), endTime.Unix())
		if t.Unix() > endTime.Unix() {
			break
		}
	}

	for _, a := range articles {
		log.Debugf("checking article [%d]", a.ID)
		// Check to see if the article was added in the date range
		if isInRange(start, end, a.DateAdded) {
			log.Debugf(
				"Article [%d]: Start date [%d] less than Article add date [%d] less than End date [%d]",
				a.ID, start, a.DateAdded, end,
			)

			dayAddedTotal, ok := itemised[a.DateAdded]
			if !ok {
				return nil, fmt.Errorf("what the fuck, date [%d] not created", a.DateAdded)
			}

			timeReading := convertWordsToTime(a.WordCount)

			// Update itemised values
			dayAddedTotal.ArticlesAdded++
			dayAddedTotal.WordsAdded += a.WordCount
			dayAddedTotal.TimeAdded += timeReading
		}

		// Check to see if the article is read
		if a.DateRead != 0 && isInRange(start, end, a.DateRead) {
			log.Debugf("Article [%d] read [%d]", a.ID, a.DateRead)
			dayReadTotal, ok := itemised[a.DateRead]
			if !ok {
				return nil, fmt.Errorf("what the fuck, date [%d] not created", a.DateRead)
			}

			timeReading := convertWordsToTime(a.WordCount)

			// Update itemised values
			dayReadTotal.ArticlesRead++
			dayReadTotal.WordsRead += a.WordCount
			dayReadTotal.TimeRead += timeReading
		}
	}

	return &itemised, nil
}

// isInRange checks to see if the article was added within the date range
func isInRange(start, end, added int64) bool {
	return start <= added && added <= end
}

func convertWordsToTime(words int64) int64 {
	timeReading := float64(words / WordsPerMinute)
	rounded := math.Round(timeReading)
	return int64(rounded)
}

// createDBTime takes two date strings, converts them to Time objects and normalises them todo
// the start of the day
func createDBTime(start, end string) (int64, int64, error) {
	if start == "" {
		return 0, 0, fmt.Errorf("no start date provided")
	}

	if end == "" {
		return 0, 0, fmt.Errorf("no end date provided")
	}

	startInt, err := strconv.Atoi(start)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not convert start date [%s]: %w", start, err)
	}

	endInt, err := strconv.Atoi(end)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not convert end date [%s]: %w", end, err)
	}

	return stripTime(startInt), stripTime(endInt), nil
}

func getPreviousDate(s, e int64) (int64, int64) {
	st := time.Unix(s, 0)
	et := time.Unix(e, 0)

	dur := st.Sub(et)
	log.Debugf("calculated duration is [%s]", dur.String())

	start := st.Add(dur).Unix()
	end := et.Add(dur).Unix()

	return stripTime(int(start)), stripTime(int(end))
}
