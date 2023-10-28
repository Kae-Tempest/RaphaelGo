package rpg

type Player struct {
	ID          int
	Username    string
	Level       int
	EXP         int
	PO          int
	Job         string
	Rank        string
	Attribut    Attribut
	Inventory   []Items
	Equipements Equipements
	DropRank    DropRank
	City        City
}

type Monster struct {
	ID       int
	Name     string
	Rank     string
	DropedPo int
	Attribut Attribut
	Drop     []Items
}

type Attribut struct {
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
	Name     string
	Attribut Attribut
	Quantity int
	Rank     string
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
	Abandonned  bool
	Description string
}

type Dongeon struct {
	ID          int
	Name        string
	Rank        string
	Floor       int
	InCity      bool
	Description string
}
