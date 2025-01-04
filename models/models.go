package models

type PlayerData struct {
	Username         string                `json:"username"`
	Online           bool                  `json:"online"`
	Server           *string               `json:"server"`
	ActiveCharacter  *string               `json:"activeCharacter"`
	Nickname         *string               `json:"nickname"`
	UUID             string                `json:"uuid"`
	Rank             string                `json:"rank"`
	RankBadge        *string               `json:"rankBadge"`
	LegacyRankColour *RankColour           `json:"legacyRankColour"`
	ShortenedRank    *string               `json:"shortenedRank"`
	SupportRank      *string               `json:"supportRank"`
	Veteran          *bool                 `json:"veteran"`
	FirstJoin        string                `json:"firstJoin"`
	LastJoin         string                `json:"lastJoin"`
	Playtime         float64               `json:"playtime"`
	Guild            *Guild                `json:"guild"`
	GlobalData       GlobalData            `json:"globalData"`
	ForumLink        *string               `json:"forumLink"`
	Ranking          Ranking               `json:"ranking"`
	PreviousRanking  Ranking               `json:"previousRanking"`
	PublicProfile    bool                  `json:"publicProfile"`
	Characters       *map[string]Character `json:"characters"`
}

type RankColour struct {
	Main string `json:"main"`
	Sub  string `json:"sub"`
}

type Guild struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Prefix    string `json:"prefix"`
	Rank      string `json:"rank"`
	RankStars string `json:"rankStars"`
}

type GlobalData struct {
	Wars            int            `json:"wars"`
	TotalLevel      int            `json:"totalLevel"`
	KilledMobs      int            `json:"killedMobs"`
	ChestsFound     int            `json:"chestsFound"`
	Dungeons        DungeonSummary `json:"dungeons"`
	Raids           RaidSummary    `json:"raids"`
	CompletedQuests int            `json:"completedQuests"`
	PvP             PvPStats       `json:"pvp"`
}

type DungeonSummary struct {
	Total int            `json:"total"`
	List  map[string]int `json:"list"`
}

type RaidSummary struct {
	Total int            `json:"total"`
	List  map[string]int `json:"list"`
}

type PvPStats struct {
	Kills  int `json:"kills"`
	Deaths int `json:"deaths"`
}

type Ranking struct {
	OrphionSrPlayers       int `json:"orphionSrPlayers"`
	GrootslangSrPlayers    int `json:"grootslangSrPlayers"`
	NamelessSrPlayers      int `json:"namelessSrPlayers"`
	OrphionCompletion      int `json:"orphionCompletion"`
	GrootslangCompletion   int `json:"grootslangCompletion"`
	FarmingLevel           int `json:"farmingLevel"`
	FishingLevel           int `json:"fishingLevel"`
	ColossusSrPlayers      int `json:"colossusSrPlayers"`
	WoodcuttingLevel       int `json:"woodcuttingLevel"`
	GlobalPlayerContent    int `json:"globalPlayerContent"`
	MiningLevel            int `json:"miningLevel"`
	PlayerContent          int `json:"playerContent"`
	WarsCompletion         int `json:"warsCompletion"`
	ColossusCompletion     int `json:"colossusCompletion"`
	ScribingLevel          int `json:"scribingLevel"`
	NamelessCompletion     int `json:"namelessCompletion"`
	TotalSoloLevel         int `json:"totalSoloLevel"`
	ProfessionsSoloLevel   int `json:"professionsSoloLevel"`
	TotalGlobalLevel       int `json:"totalGlobalLevel"`
	ProfessionsGlobalLevel int `json:"professionsGlobalLevel"`
	JewelingLevel          int `json:"jewelingLevel"`
	CombatGlobalLevel      int `json:"combatGlobalLevel"`
	CombatSoloLevel        int `json:"combatSoloLevel"`
	CraftsmanContent       int `json:"craftsmanContent"`
	CookingLevel           int `json:"cookingLevel"`
	WeaponsmithingLevel    int `json:"weaponsmithingLevel"`
	TailoringLevel         int `json:"tailoringLevel"`
	AlchemismLevel         int `json:"alchemismLevel"`
	WoodworkingLevel       int `json:"woodworkingLevel"`
	ArmouringLevel         int `json:"armouringLevel"`
}

type Character struct {
	Type            string                `json:"type"`
	Reskin          *string               `json:"reskin"`
	Nickname        *string               `json:"nickname"`
	Level           int                   `json:"level"`
	XP              int                   `json:"xp"`
	XPPercent       int                   `json:"xpPercent"`
	TotalLevel      int                   `json:"totalLevel"`
	Wars            int                   `json:"wars"`
	Playtime        float64               `json:"playtime"`
	MobsKilled      int                   `json:"mobsKilled"`
	ChestsFound     int                   `json:"chestsFound"`
	ItemsIdentified *int                  `json:"itemsIdentified"`
	BlocksWalked    int                   `json:"blocksWalked"`
	Logins          int                   `json:"logins"`
	Deaths          int                   `json:"deaths"`
	Discoveries     int                   `json:"discoveries"`
	PvP             PvPStats              `json:"pvp"`
	GameMode        []string              `json:"gamemode"`
	SkillPoints     map[string]int        `json:"skillPoints"`
	Professions     map[string]Profession `json:"professions"`
	Dungeons        DungeonSummary        `json:"dungeons"`
	Raids           RaidSummary           `json:"raids"`
	Quests          []string              `json:"quests"`
}

type Profession struct {
	Level     int `json:"level"`
	XPPercent int `json:"xpPercent"`
}

// GuildData idk what to say
type GuildData struct {
	UUID        string                `json:"uuid"`
	Name        string                `json:"name"`
	Prefix      string                `json:"prefix"`
	Level       int                   `json:"level"`
	XPPercent   int                   `json:"xpPercent"`
	Territories int                   `json:"territories"`
	Wars        int                   `json:"wars"`
	Created     string                `json:"created"`
	Members     Members               `json:"members"`
	Online      int                   `json:"online"`
	Banner      Banner                `json:"banner"`
	SeasonRanks map[string]SeasonRank `json:"seasonRanks"`
}

// Members is a collection of guild member groups by rank plus total count.
type Members struct {
	Total      int                   `json:"total"`
	Owner      map[string]MemberInfo `json:"owner"`
	Chief      map[string]MemberInfo `json:"chief"`
	Strategist map[string]MemberInfo `json:"strategist"`
	Captain    map[string]MemberInfo `json:"captain"`
	Recruiter  map[string]MemberInfo `json:"recruiter"`
	Recruit    map[string]MemberInfo `json:"recruit"`
}

// MemberInfo represents data for an individual member.
// Note that "guildRank" and "contributionRank" are sometimes
// present/absent depending on the rank of the member.
type MemberInfo struct {
	// This field corresponds to "<username/uuid>" in the original JSON.
	// If you want to preserve it, you can rename it or parse it in manually.
	UsernameOrUUID string  `json:"<username/uuid>"`
	Online         bool    `json:"online"`
	Server         *string `json:"server"` // nil if "server": null
	Contributed    int     `json:"contributed"`

	// For non-recruits:
	GuildRank *int `json:"guildRank,omitempty"`
	// For recruits:
	ContributionRank *int `json:"contributionRank,omitempty"`

	Joined string `json:"joined"`
}

// Banner defines the structure of the banner object.
type Banner struct {
	Base      string        `json:"base"`
	Tier      int           `json:"tier"`
	Structure string        `json:"structure"`
	Layers    []BannerLayer `json:"layers"`
}

// BannerLayer describes each layer of the banner.
type BannerLayer struct {
	Colour  string `json:"colour"`
	Pattern string `json:"pattern"`
}

// SeasonRank represents each "season" entry inside "seasonRanks".
type SeasonRank struct {
	Rating           int `json:"rating"`
	FinalTerritories int `json:"finalTerritories"`
}
