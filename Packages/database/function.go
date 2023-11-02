package database

import (
	"RaphaelGo/Packages/HandleError"
	"RaphaelGo/Packages/Struct"
	"context"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var db = DB()

func CreatePlayer(s *discordgo.Session, player Struct.Player) {

	p := GetPlayerByName(s, player.Username)
	if p.Username == player.Username {
		fmt.Println("user already exist")
		errMessage := fmt.Errorf("user already exist").Error()
		HandleError.SendDBError(s, errMessage)
		return
	} else {
		var id int
		err := db.QueryRow(context.Background(),
			`insert into public.players (username,level,exp,po,job_id,rank,city_id) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id;`,
			player.Username, player.Level, player.EXP, player.PO, player.Job.ID, player.Rank, player.City.ID).Scan(&id)
		if err != nil {
			errMessage := fmt.Errorf("cannot create Player: %s", err).Error()
			HandleError.SendDBError(s, errMessage)
		}
		_, err = db.Exec(context.Background(),
			`insert into public.attributes (player_id, hp,mp,strength,agility,intelligence,spirit,constitution,speed,luck) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`,
			id, player.Attribute.HP, player.Attribute.MP, player.Attribute.Strength, player.Attribute.Agility, player.Attribute.Intelligence, player.Attribute.Spirit, player.Attribute.Constitution, player.Attribute.Speed, player.Attribute.Luck)
		if err != nil {
			errMessage := fmt.Errorf("cannot create Attributes: %s", err).Error()
			HandleError.SendDBError(s, errMessage)
		}
		_, err = db.Exec(context.Background(),
			`insert into public.drop_rank (player_id,plants,animals,minerals,magic,special) values ($1,$2,$3,$4,$5,$6);`,
			id, player.DropRank.Plants, player.DropRank.Animals, player.DropRank.Minerals, player.DropRank.Magic, player.DropRank.Special)
		if err != nil {
			errMessage := fmt.Errorf("cannot create DropRank: %s", err).Error()
			HandleError.SendDBError(s, errMessage)
		}
	}
}

func GetPlayerById(s *discordgo.Session, id int) Struct.Player {
	var player Struct.Player
	err := db.QueryRow(context.Background(),
		`select * from public.players where id = $1`, id).Scan(&player.ID, &player.Username, &player.Level, &player.EXP, &player.PO, &player.Job.ID, &player.Rank, &player.City.ID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return player
		}
		errMessage := fmt.Errorf("cannot get user by id: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	player.Job = getJobById(s, player.Job.ID)
	player.City = getCityById(s, player.City.ID)
	player.Attribute = getPlayerAttributes(s, player.ID)
	player.DropRank = getPlayerDropRank(s, player.ID)

	return player
}

func GetPlayerByName(s *discordgo.Session, name string) Struct.Player {
	var player Struct.Player
	err := db.QueryRow(context.Background(),
		`select * from public.players where username = $1`, name).Scan(&player.ID, &player.Username, &player.Level, &player.EXP, &player.PO, &player.Job.ID, &player.Rank, &player.City.ID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return player
		}
		errMessage := fmt.Errorf("cannot get user by name: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	player.Job = getJobById(s, player.Job.ID)
	player.City = getCityById(s, player.City.ID)
	player.Attribute = getPlayerAttributes(s, player.ID)
	player.DropRank = getPlayerDropRank(s, player.ID)

	return player
}

func getJobById(s *discordgo.Session, id int) Struct.Job {
	var job Struct.Job
	err := db.QueryRow(context.Background(),
		`select * from public.jobs where id = $1`, id).Scan(&job.ID, &job.Name, &job.Strength, &job.Agility, &job.Intelligence, &job.Spirit, &job.Constitution, &job.Speed, &job.Description)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return job
		}
		errMessage := fmt.Errorf("cannot get job by id: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	return job
}

func GetJobByName(s *discordgo.Session, name string) Struct.Job {
	var job Struct.Job
	err := db.QueryRow(context.Background(),
		`select * from public.jobs where name = $1`, name).Scan(&job.ID, &job.Name, &job.Strength, &job.Agility, &job.Intelligence, &job.Spirit, &job.Constitution, &job.Speed, &job.Description)
	if err != nil {
		errMessage := fmt.Errorf("cannot get job by name: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	return job
}

func getCityById(s *discordgo.Session, id int) Struct.City {
	var city Struct.City
	err := db.QueryRow(context.Background(),
		`select * from public.cities where id = $1`, id).Scan(&city.ID, &city.Name, &city.Rank, &city.Abandoned, &city.Description)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return city
		}
		errMessage := fmt.Errorf("cannot get city by id: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	return city
}

func GetCityByName(s *discordgo.Session, name string) Struct.City {
	var city Struct.City
	err := db.QueryRow(context.Background(),
		`select * from public.cities where name = $1`, name).Scan(&city.ID, &city.Name, &city.Rank, &city.Abandoned, &city.Description)
	if err != nil {
		errMessage := fmt.Errorf("cannot get city by name: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	return city
}

func getPlayerAttributes(s *discordgo.Session, id int) Struct.Attribute {
	var attribute Struct.Attribute
	var playerId int
	var itemId sql.NullInt16
	err := db.QueryRow(context.Background(),
		`select * from public.attributes where player_id = $1`, id).Scan(&playerId, &itemId, &attribute.HP, &attribute.MP, &attribute.Strength, &attribute.Agility, &attribute.Intelligence, &attribute.Spirit, &attribute.Constitution, &attribute.Speed, &attribute.Luck)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return attribute
		}
		errMessage := fmt.Errorf("cannot get attribute by id: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	return attribute
}

func getPlayerDropRank(s *discordgo.Session, id int) Struct.DropRank {
	var dropRank Struct.DropRank
	var player_id int
	err := db.QueryRow(context.Background(),
		`select * from public.drop_rank where player_id = $1`, id).Scan(&player_id, &dropRank.Plants, &dropRank.Animals, &dropRank.Minerals, &dropRank.Magic, &dropRank.Special)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return dropRank
		}
		errMessage := fmt.Errorf("cannot get drop rank by id: %s", err).Error()
		HandleError.SendDBError(s, errMessage)
	}
	return dropRank
}
