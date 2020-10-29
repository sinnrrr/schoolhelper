package db

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func CreateClass(classID int64, name string) (neo4j.Node, error) {
	var class neo4j.Node

	result, err := Session.Run(
		"MERGE (c:Class { tg_id: $tg_id, name: $name }) RETURN c",
		map[string]interface{}{
			"tg_id": classID,
			"name":  name,
		},
	)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		class = result.Record().GetByIndex(0).(neo4j.Node)
	}

	return class, result.Err()
}