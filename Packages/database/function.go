package database

import (
	"RaphaelGo/Struct"
	"RaphaelGo/events"
	"context"
	"fmt"
)

var db = DB()

func CreatePlayer(player Struct.Player) {
	err := db.QueryRow(context.Background(),
		`insert into public.players (username,level,exp,po,job_id,rank,city_id) values ($1,$2,$3,$4,$5,$6,$7)`, player.Username, player.Level, player.EXP, player.PO, player.Job.ID, player.Rank, player.City.ID)
	if err != nil {
		errMessage := fmt.Errorf("cannot create Player: %s", err).Error()
		events.SendDBError(errMessage)
	}
	err = db.QueryRow(context.Background(),
		`insert into public.attributes (hp,mp,strength,agility,intelligence,spirit,constitution,speed,luck) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, player.Attribute.HP, player.Attribute.MP, player.Attribute.Strength, player.Attribute.Agility, player.Attribute.Intelligence, player.Attribute.Spirit, player.Attribute.Constitution, player.Attribute.Speed, player.Attribute.Luck)
	if err != nil {
		errMessage := fmt.Errorf("cannot create Attributes: %s", err).Error()
		events.SendDBError(errMessage)
	}
	err = db.QueryRow(context.Background(),
		`insert into public.drop_rank (plants,animals,minerals,magic,special) values ($1,$2,$3,$4,$5)`, player.DropRank.Plants, player.DropRank.Animals, player.DropRank.Minerals, player.DropRank.Magic, player.DropRank.Special)
	if err != nil {
		errMessage := fmt.Errorf("cannot create DropRank: %s", err).Error()
		events.SendDBError(errMessage)
	}
}

func GetPlayerById(id int) Struct.Player {
	var player Struct.Player
	err := db.QueryRow(context.Background(),
		`select * from public.players where id = $1`, id).Scan(&player.ID, &player.Username, &player.Level, &player.EXP, &player.PO, &player.Rank)
	if err != nil {
		errMessage := fmt.Errorf("cannot get user by id: %s", err).Error()
		events.SendDBError(errMessage)
	}
	player.Job = getJobById(player.Job.ID)
	player.City = getCityById(player.City.ID)
	player.Attribute = getPlayerAttributes(player.ID)
	player.DropRank = getPlayerDropRank(player.ID)

	return player
}

func getJobById(id int) Struct.Job {
	var job Struct.Job
	err := db.QueryRow(context.Background(),
		`select * from public.jobs where id = $1`, id).Scan(&job.ID, &job.Name, &job.Strength, &job.Agility, &job.Intelligence, &job.Spirit, &job.Constitution, &job.Speed)
	if err != nil {
		errMessage := fmt.Errorf("cannot get job by id: %s", err).Error()
		events.SendDBError(errMessage)
	}
	return job
}

func getCityById(id int) Struct.City {
	var city Struct.City
	err := db.QueryRow(context.Background(),
		`select * from public.cities where id = $1`, id).Scan(&city.ID, &city.Name, &city.Abandoned, &city.Rank, &city.Description)
	if err != nil {
		errMessage := fmt.Errorf("cannot get city by id: %s", err).Error()
		events.SendDBError(errMessage)
	}
	return city
}

func getPlayerAttributes(id int) Struct.Attribute {
	var attribute Struct.Attribute
	err := db.QueryRow(context.Background(),
		`select * from public.attributes where player_id = $1`, id).Scan(&attribute.HP, &attribute.MP, &attribute.Strength, &attribute.Agility, &attribute.Intelligence, &attribute.Spirit, &attribute.Constitution, &attribute.Speed, &attribute.Luck)
	if err != nil {
		errMessage := fmt.Errorf("cannot get attribute by id: %s", err).Error()
		events.SendDBError(errMessage)
	}
	return attribute
}

func getPlayerDropRank(id int) Struct.DropRank {
	var dropRank Struct.DropRank
	err := db.QueryRow(context.Background(),
		`select * from public.drop_rank where player_id = $1`, id).Scan(&dropRank.Plants, &dropRank.Animals, &dropRank.Minerals, &dropRank.Magic, &dropRank.Special)
	if err != nil {
		errMessage := fmt.Errorf("cannot get drop rank by id: %s", err).Error()
		events.SendDBError(errMessage)
	}
	return dropRank
}
