package pocket

import "github.com/omgitsotis/pocket-stats/pocket/model"

// createStats generates stats based on the list of articles provided
func createStats(sp model.StatsParams, arts []model.Article) *model.Stats {
	stats := make(map[int64]*model.Stat)
	var wAdded, wRead, aAdded, aRead int64
	for _, a := range arts {
		if a.Status == model.Archived {
			// Update article read counts
			aRead++
			wRead += a.WordCount

			s, ok := stats[a.DateRead]
			if ok {
				s.ArticleRead++
				s.WordRead += a.WordCount
			} else {
				newStat := model.Stat{
					ArticleRead: 1,
					WordRead:    a.WordCount,
				}

				stats[a.DateRead] = &newStat
			}

			// If the article was added in the current time frame, update added
			// count
			if a.DateAdded >= sp.Start && a.DateAdded <= sp.End {
				aAdded++
				wAdded += a.WordCount

				s, ok = stats[a.DateAdded]
				if ok {
					s.ArticleAdded++
					s.WordAdded += a.WordCount
				} else {
					newStat := model.Stat{
						ArticleAdded: 1,
						WordAdded:    a.WordCount,
					}

					stats[a.DateAdded] = &newStat
				}
			}

		} else {
			// Update added count
			aAdded++
			wAdded += a.WordCount

			s, ok := stats[a.DateAdded]
			if ok {
				s.ArticleAdded++
				s.WordAdded += a.WordCount
			} else {
				newStat := model.Stat{
					ArticleAdded: 1,
					WordAdded:    a.WordCount,
				}

				stats[a.DateAdded] = &newStat
			}
		}
	}

	totals := model.TotalStats{
		ArticlesAdded: aAdded,
		ArticlesRead:  aRead,
		WordsAdded:    wAdded,
		WordsRead:     wRead,
	}

	return &model.Stats{
		Start:  sp.Start,
		End:    sp.End,
		Value:  stats,
		Totals: totals,
	}
}
