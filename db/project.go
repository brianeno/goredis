package db

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Project struct {
	Title       string `json:"title" binding:"required"`
	Storypoints int    `json:"storypoints" binding:"required"`
	Rank        int    `json:"rank"`
}

func (db *RedisDatabase) SaveProject(project *Project) error {
	member := &redis.Z{
		Score:  float64(project.Storypoints),
		Member: project.Title,
	}
	pipe := db.Client.TxPipeline()
	pipe.ZAdd(Ctx, "dashboard", member)

	rank := pipe.ZRank(Ctx, dashboardKey, project.Title)
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return err
	}

	// example of a set, we will retrieve it GetProject
	err = db.Client.Set(Ctx, project.Title, project.Storypoints, 0).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(rank.Val(), err)
	project.Rank = int(rank.Val())
	return nil
}

func (db *RedisDatabase) GetProject(title string) (*Project, error) {
	pipe := db.Client.TxPipeline()
	sps := pipe.ZScore(Ctx, dashboardKey, title)
	rank := pipe.ZRank(Ctx, dashboardKey, title)
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return nil, err
	}
	if sps == nil {
		return nil, ErrNil
	}

	// example of a get
	val, newErr := db.Client.Get(Ctx, title).Result()
	if newErr != nil {
		panic(newErr)
	}
	fmt.Println("Storypoints for", title, "is", val)

	// return our project
	return &Project{
		Title:       title,
		Storypoints: int(sps.Val()),
		Rank:        int(rank.Val()),
	}, nil
}
