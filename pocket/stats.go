package pocket

import (
	"github.com/omgitsotis/pocket-stats/pocket/model"
)

var readingSpeed int64 = 247

// createStats generates stats based on the list of articles provided
func createStats(sp model.StatsParams, articles []model.Article) *model.Stats {
	days := make(map[int64]*model.Stat)
	tags := make(map[string]*model.Stat)

	var wAdded, wRead, aAdded, aRead int64
	for _, a := range articles {
		if a.Status == model.Archived {
			// Update article read counts
			aRead++
			wRead += a.WordCount

			s, ok := days[a.DateRead]
			if ok {
				s.ArticleRead++
				s.WordsRead += a.WordCount
			} else {
				newStat := model.Stat{
					ArticleRead: 1,
					WordsRead:   a.WordCount,
				}

				days[a.DateRead] = &newStat
			}

			// Update tag read counts. We only update reads for tags
			tag, ok := tags[a.Tag]
			if ok {
				tag.ArticleRead++
				tag.WordsRead += a.WordCount
			} else {
				newStat := model.Stat{
					ArticleRead: 1,
					WordsRead:   a.WordCount,
				}
				tags[a.Tag] = &newStat
			}

			// If the article was added in the current time frame, update added
			// count
			if a.DateAdded >= sp.Start && a.DateAdded <= sp.End {
				aAdded++
				wAdded += a.WordCount

				s, ok = days[a.DateAdded]
				if ok {
					s.ArticleAdded++
					s.WordsAdded += a.WordCount
				} else {
					newStat := model.Stat{
						ArticleAdded: 1,
						WordsAdded:   a.WordCount,
					}

					days[a.DateAdded] = &newStat
				}
			}

		} else {
			// Update added count
			aAdded++
			wAdded += a.WordCount

			s, ok := days[a.DateAdded]
			if ok {
				s.ArticleAdded++
				s.WordsAdded += a.WordCount
			} else {
				newStat := model.Stat{
					ArticleAdded: 1,
					WordsAdded:   a.WordCount,
				}

				days[a.DateAdded] = &newStat
			}
		}
	}

	totals := model.TotalStats{
		ArticlesAdded: aAdded,
		ArticlesRead:  aRead,
		WordsAdded:    wAdded,
		WordsRead:     wRead,
	}

	stats := &model.Stats{
		Start:      sp.Start,
		End:        sp.End,
		DateValues: days,
		TagValues:  tags,
		Totals:     totals,
	}

	getTimeReading(stats)

	return stats
}

// getTimeReading calulates the total time reading as well as time reading for
// each day
func getTimeReading(s *model.Stats) {
	s.Totals.TimeReading = s.Totals.WordsRead / readingSpeed
	logger.Debugf("Total time spend reading [%d]", s.Totals.TimeReading)

	for _, stat := range s.DateValues {
		stat.TimeReading = stat.WordsRead / readingSpeed
	}

	for _, stat := range s.TagValues {
		stat.TimeReading = stat.WordsRead / readingSpeed
	}
}
