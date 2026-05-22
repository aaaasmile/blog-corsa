package admin

import (
	"blog-corsa/idl"
	"encoding/json"
	"fmt"
	"log"
)

type RaceParam struct {
	Cmd string `json:"cmd"`
	Id  int64  `json:"id"`
}

type RaceDataParam struct {
	Cmd  string       `json:"cmd"`
	Race idl.RaceItem `json:"race"`
}

type RaceReq struct {
	Params RaceParam
}

type RaceDataReq struct {
	Params RaceDataParam
}

func (ah *AdminHandler) doRace() error {
	raceReq := RaceReq{}
	if err := json.Unmarshal(ah.rawbody, &raceReq); err != nil {
		return err
	}
	var err error
	racePara := raceReq.Params
	switch racePara.Cmd {
	case "list":
		err = ah.doRaceList()
	case "insert":
		err = ah.doRaceInsert()
	case "update":
		err = ah.doRaceUpdate()
	case "delete":
		err = ah.doRaceDelete()
	default:
		return fmt.Errorf("[doRace] race command not supported: %v", racePara)
	}
	if err != nil {
		return err
	}
	return nil
}

func (ah *AdminHandler) doRaceList() error {
	races, err := ah.liteMainDB.GetRaces()
	if err != nil {
		return err
	}
	resp := struct {
		Races []idl.RaceItem
	}{
		Races: races,
	}
	return writeResponse(ah._w, resp)
}

func (ah *AdminHandler) doRaceInsert() error {
	raceDataReq := RaceDataReq{}
	if err := json.Unmarshal(ah.rawbody, &raceDataReq); err != nil {
		return err
	}
	race := raceDataReq.Params.Race
	log.Println("[doRaceInsert] insert race", race.Name)
	if err := ah.liteMainDB.InsertRace(&race); err != nil {
		return err
	}
	return ah.doRaceList()
}

func (ah *AdminHandler) doRaceUpdate() error {
	raceDataReq := RaceDataReq{}
	if err := json.Unmarshal(ah.rawbody, &raceDataReq); err != nil {
		return err
	}
	race := raceDataReq.Params.Race
	log.Println("[doRaceUpdate] update race id", race.Id)
	if err := ah.liteMainDB.UpdateRace(&race); err != nil {
		return err
	}
	return ah.doRaceList()
}

func (ah *AdminHandler) doRaceDelete() error {
	raceReq := RaceReq{}
	if err := json.Unmarshal(ah.rawbody, &raceReq); err != nil {
		return err
	}
	id := raceReq.Params.Id
	log.Println("[doRaceDelete] delete race id", id)
	if err := ah.liteMainDB.DeleteRace(id); err != nil {
		return err
	}
	return ah.doRaceList()
}
