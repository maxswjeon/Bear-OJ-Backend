package server

import (
	"encoding/json"

	"github.com/maxswjeon/contest-backend/schemas"
	"gorm.io/gorm"
)

type HandshakePacket struct {
	Problems  [][]interface{}            `json:"problems"`
	Executors map[string][][]interface{} `json:"executors"`
}

func on_handshake(judge *Judge, db *gorm.DB, data string) error {
	var packet HandshakePacket

	json.Unmarshal([]byte(data), &packet)

	for name := range packet.Executors {
		judge.languages = append(judge.languages, name)
	}

	for _, problemInfo := range packet.Problems {
		var problemCount int64
		db.Model(&schemas.InternalProblem{}).Where("name = ? AND last_update = ?", problemInfo[0].(string), problemInfo[1].(float64)).Count(&problemCount)

		if problemCount == 0 {
			db.Create(&schemas.InternalProblem{Name: problemInfo[0].(string), LastUpdate: problemInfo[1].(float64)})
		} else {
			db.Model(&schemas.InternalProblem{}).Where("name = ? AND last_update = ?", problemInfo[0].(string), problemInfo[1].(float64)).Update("valid", true)
		}

		var problem schemas.InternalProblem
		db.Model(&schemas.InternalProblem{}).Where("name = ? AND last_update = ?", problemInfo[0].(string), problemInfo[1].(float64)).First(&problem)

		judge.problems = append(judge.problems, problem.ID)
	}

	response, err := compressPacket("{\"name\": \"handshake-success\"}")
	if err != nil {
		return err
	}

	size, err := encodeSize(response)
	if err != nil {
		return err
	}

	judge.conn.Write(size)
	judge.conn.Write(response)

	return nil
}
