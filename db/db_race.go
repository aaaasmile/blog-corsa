package db

import (
	"blog-corsa/idl"
	"log"
	"time"
)

func (ld *LiteDB) GetRaces() ([]idl.RaceItem, error) {
	log.Println("[LiteDB - GetRaces] select all races")

	q := `SELECT id, name, title, meter_length, ascending_meter, descending_meter, 
		rank_global, rank_gender, rank_class, class_name, 
		race_start_datetime, sport_type_id, result_time, 
		pace_kmh, pace_minkm, comment, race_subtype_id, km_length, loop_number, loop_length 
		FROM race ORDER BY race_start_datetime DESC;`
	if ld.debugSQL {
		log.Println("Query is", q)
	}
	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []idl.RaceItem{}
	for rows.Next() {
		item := idl.RaceItem{}
		var ts int64
		if err := rows.Scan(&item.Id,
			&item.Name,
			&item.Title,
			&item.MeterLength,
			&item.AscendingMeter,
			&item.DescendingMeter,
			&item.RankGlobal,
			&item.RankGender,
			&item.RankClass,
			&item.ClassName,
			&ts,
			&item.SportTypeId,
			&item.ResultTime,
			&item.PaceKmh,
			&item.PaceMinkm,
			&item.Comment,
			&item.RaceSubtypeId,
			&item.KmLength,
			&item.LoopNumber,
			&item.LoopLength); err != nil {
			return nil, err
		}
		item.RaceStartDateTime = time.Unix(ts, 0)
		res = append(res, item)
	}
	log.Printf("[LiteDB - GetRaces] races read %d", len(res))
	return res, nil
}

func (ld *LiteDB) InsertRace(item *idl.RaceItem) error {
	log.Println("[LiteDB - INSERT] insert new race", item.Name)

	q := `INSERT INTO race(name, title, distance, date, ascending_meter, descending_meter, 
		rank_global, rank_gender, rank_class, class_name, sport_type_id, result_time, 
		pace_kmh, pace_minkm, comment, race_subtype_id, km_length, loop_number, loop_length) 
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	if ld.debugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(
		item.Name,
		item.Title,
		item.MeterLength,
		item.AscendingMeter,
		item.DescendingMeter,
		item.RankGlobal,
		item.RankGender,
		item.RankClass,
		item.ClassName,
		item.RaceStartDateTime.Local().Unix(),
		item.SportTypeId,
		item.ResultTime,
		item.PaceKmh,
		item.PaceMinkm,
		item.Comment,
		item.RaceSubtypeId,
		item.KmLength,
		item.LoopNumber,
		item.LoopLength)
	if err != nil {
		return err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.Id = insertedId
	log.Println("[LiteDB - INSERT] race added OK, id:", item.Id)
	return nil
}

func (ld *LiteDB) UpdateRace(item *idl.RaceItem) error {
	log.Println("[LiteDB - UPDATE] update race id", item.Id)

	q := `UPDATE race SET name=?, title=?, distance=?, date=?, ascending_meter=?, descending_meter=?, 
		rank_global=?, rank_gender=?, rank_class=?, class_name=?, sport_type_id=?, result_time=?, 
		pace_kmh=?, pace_minkm=?, comment=?, race_subtype_id=?, km_length=?, loop_number=?, loop_length=? 
		WHERE id=?;`
	if ld.debugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		item.Name,
		item.Title,
		item.MeterLength,
		item.AscendingMeter,
		item.DescendingMeter,
		item.RankGlobal,
		item.RankGender,
		item.RankClass,
		item.ClassName,
		item.RaceStartDateTime.Local().Unix(),
		item.SportTypeId,
		item.ResultTime,
		item.PaceKmh,
		item.PaceMinkm,
		item.Comment,
		item.RaceSubtypeId,
		item.KmLength,
		item.LoopNumber,
		item.LoopLength,
		item.Id)
	if err != nil {
		return err
	}
	if ld.debugSQL {
		ra, err := res.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("Row affected: ", ra)
	}
	log.Println("[LiteDB - UPDATE] race updated OK, id:", item.Id)
	return nil
}

func (ld *LiteDB) DeleteRace(id int64) error {
	log.Println("[LiteDB - DELETE] delete race id", id)

	q := `DELETE FROM race WHERE id=?;`
	if ld.debugSQL {
		log.Println("SQL is", q)
	}
	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	if ld.debugSQL {
		ra, err := res.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("Row affected: ", ra)
	}
	log.Println("[LiteDB - DELETE] race deleted OK, id:", id)
	return nil
}
