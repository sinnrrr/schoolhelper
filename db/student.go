package db

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	tb "gopkg.in/tucnak/telebot.v2"
)

func CreateStudent(user *tb.User, classID int64) (neo4j.Node, error) {
	var student neo4j.Node

	result, err := Session.Run(
		"MATCH (c:Class { tg_id: $class_id })"+
			"\n"+
			"MERGE (s:Student {"+
			"tg_id: $tg_id,"+
			"first_name: $first_name,"+
			"last_name: $last_name,"+
			"username: $username,"+
			"language_code: $language_code,"+
			"dialogue_state: $dialogue_state"+
			"})-[:STUDYING_IN]->(c)"+
			"\n"+
			"RETURN s",
		map[string]interface{}{
			"tg_id":          user.ID,
			"first_name":     user.FirstName,
			"last_name":      user.LastName,
			"username":       user.Username,
			"language_code":  user.LanguageCode,
			"class_id":       classID,
			"dialogue_state": 0,
		},
	)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		student = result.Record().GetByIndex(0).(neo4j.Node)
	}

	return student, result.Err()
}

func StudentSession(studentID int) (neo4j.Node, error) {
	var session neo4j.Node

	result, err := Session.Run(
		"MATCH (s:Student { tg_id: $tg_id }) RETURN s",
		map[string]interface{}{
			"tg_id": studentID,
		},
	)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		session = result.Record().GetByIndex(0).(neo4j.Node)
	}

	return session, result.Err()
}

func DialogueState(studentID int) (int8, error) {
	session, err := StudentSession(studentID)
	if err != nil {
		panic(err)
	}
	if session == nil {
		return -1, fmt.Errorf("failed to find user with the ID %d", studentID)
	}

	return int8(session.Props()["dialogue_state"].(int64)), err
}

func SetDialogueState(studentID int, state int8) error {
	_, err := Session.Run(
		"MATCH (s:Student { tg_id: $tg_id }) SET s.dialogue_state = $dialogue_state",
		map[string]interface{}{
			"tg_id":          studentID,
			"dialogue_state": state,
		},
	)

	return err
}