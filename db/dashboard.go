package db

var dashboardKey = "dashboard"

type Dashboard struct {
	Count    int `json:"count"`
	Projects []*Project
}

func (db *RedisDatabase) GetDashboard() (*Dashboard, error) {
	scores := db.Client.ZRangeWithScores(Ctx, dashboardKey, 0, -1)
	if scores == nil {
		return nil, ErrNil
	}
	count := len(scores.Val())
	projects := make([]*Project, count)
	for idx, member := range scores.Val() {
		projects[idx] = &Project{
			Title:       member.Member.(string),
			Storypoints: int(member.Score),
			Rank:        idx,
		}
	}
	dashboard := &Dashboard{
		Count:    count,
		Projects: projects,
	}
	return dashboard, nil
}
