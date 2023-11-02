package Struct

type Player struct {
	ID        int
	Username  string
	Level     int
	EXP       int
	PO        int
	Job       Job
	Rank      string
	Attribute Attribute
	Inventory []Items
	Equipment Equipements
	DropRank  DropRank
	City      City
}

type Attribute struct {
	HP           int
	MP           int
	Strength     int
	Agility      int
	Intelligence int
	Spirit       int
	Constitution int
	Speed        int
	Luck         int
}

type Items struct {
	Name      string
	Attribute Attribute
	Quantity  int
	Rank      string
}

type Equipements struct {
	Mh         string
	Oh         string
	Chestplate string
	Boots      string
	Broach     string
	Rings      string
	Earrings   string
	Belt       string
}

type DropRank struct {
	Plants   string
	Animals  string
	Minerals string
	Magic    string
	Special  string
}

type Job struct {
	ID           int
	Name         string
	Strength     int
	Agility      int
	Intelligence int
	Spirit       int
	Constitution int
	Speed        int
	Description  string
}

type City struct {
	ID          int
	Name        string
	Rank        string
	Abandoned   bool
	Description string
}

type Dungeon struct {
	ID          int
	Name        string
	Rank        string
	Floor       int
	InCity      bool
	Description string
}

type Monster struct {
	ID        int
	Name      string
	Rank      string
	DroppedPo int
	Attribute Attribute
	Drop      []Items
}
