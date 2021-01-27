package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/steve-nzr/goff-server/internal/domain/aggregates"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces/repositories"
)

type gameCharacterRepository struct {
	client redis.Cmdable
}

func (g *gameCharacterRepository) Get(id customtypes.ID) (interfaces.Character, error) {
	res := g.client.Get(context.TODO(), g.characterKey(id))
	if res.Err() != nil {
		return nil, res.Err()
	}

	char := new(aggregates.Character)
	data, _ := res.Bytes()
	_ = json.Unmarshal(data, char)
	return char, nil
}

func (g *gameCharacterRepository) Save(char interfaces.Character) error {
	key := g.characterKey(char.GetCharacterID())
	data, _ := json.Marshal(char)

	res := g.client.Set(context.TODO(), key, data, 0)
	if res.Err() != nil {
		return res.Err()
	}

	member := fmt.Sprintf("%d", char.GetCharacterID())
	geoAddRes := g.client.GeoAdd(context.TODO(), "POSITIONS", &redis.GeoLocation{
		Name:      member,
		Latitude:  char.GetLocation().X / 10000.0,
		Longitude: char.GetLocation().Z / 10000.0,
	})
	if geoAddRes.Err() != nil {
		return geoAddRes.Err()
	}

	return nil
}

func (g *gameCharacterRepository) GetAround(id customtypes.ID) ([]interfaces.Character, error) {
	geoSearchRes := g.client.GeoRadiusByMember(context.Background(), "POSITIONS", fmt.Sprintf("%d", id), &redis.GeoRadiusQuery{
		Radius: 1.0,
		Unit:   "km",
	})
	if geoSearchRes.Err() != nil {
		return nil, geoSearchRes.Err()
	}

	locations := geoSearchRes.Val()
	aroundCharacters := make([]interfaces.Character, 0, len(locations))
	for _, location := range locations {
		id, _ := strconv.ParseInt(location.Name, 10, 32)
		if char, err := g.Get((customtypes.ID)(id)); err == nil {
			logrus.Infof("Found characer around ! %d", id)
			aroundCharacters = append(aroundCharacters, char)
		}
	}
	return aroundCharacters, nil
}

func (g *gameCharacterRepository) GetSameMap(id customtypes.ID) ([]interfaces.Character, error) {
	panic("implement me")
}

func (g *gameCharacterRepository) Delete(id customtypes.ID) error {
	panic("implement me")
}

// private

func (g *gameCharacterRepository) characterKey(id customtypes.ID) string {
	return fmt.Sprintf("CHAR:%d", id)
}

func NewGameCharacterRepository(c redis.Cmdable) repositories.GameCharacter {
	return &gameCharacterRepository{
		client: c,
	}
}
