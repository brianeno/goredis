package db

var dashboardKey = "dashboard"

type Dashboard struct {
	Count    int `json:"count"`
	Projects []*Project
}

func (db *RedisDatabase) GetDashboard() (*Dashboard, error) {
	sps := db.Client.ZRangeWithScores(Ctx, dashboardKey, 0, -1)
	if sps == nil {
		return nil, ErrNil
	}
	count := len(sps.Val())
	projects := make([]*Project, count)
	for idx, project := range sps.Val() {
		projects[idx] = &Project{
			Title:       project.Member.(string),
			Storypoints: int(project.Score),
			Rank:        idx,
		}
	}
	dashboard := &Dashboard{
		Count:    count,
		Projects: projects,
	}
	return dashboard, nil
}
