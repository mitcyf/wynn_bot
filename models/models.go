package models

type PlayerData struct {
	Username         string               `json:"username"`
	Online           bool                 `json:"online"`
	Server           string               `json:"server"`
	ActiveCharacter  string               `json:"activeCharacter"`
	Nickname         *string              `json:"nickname"`
	UUID             string               `json:"uuid"`
	Rank             string               `json:"rank"`
	RankBadge        string               `json:"rankBadge"`
	LegacyRankColour RankColour           `json:"legacyRankColour"`
	ShortenedRank    *string              `json:"shortenedRank"`
	SupportRank      string               `json:"supportRank"`
	Veteran          *bool                `json:"veteran"`
	FirstJoin        string               `json:"firstJoin"`
	LastJoin         string               `json:"lastJoin"`
	Playtime         float64              `json:"playtime"`
	Guild            GuildSimple          `json:"guild"`
	GlobalData       GlobalData           `json:"globalData"`
	ForumLink        *string              `json:"forumLink"`
	Ranking          Ranking              `json:"ranking"`
	PreviousRanking  Ranking              `json:"previousRanking"`
	PublicProfile    bool                 `json:"publicProfile"`
	Characters       map[string]Character `json:"characters"`
}

type RankColour struct {
	Main string `json:"main"`
	Sub  string `json:"sub"`
}

type GuildSimple struct {
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
	PreEconomy      *string               `json:"preEconomy"`
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

// Guild represents the complete guild data from the Wynncraft API.
type Guild struct {
	Name        string        `json:"name"`        // Guild name
	Prefix      string        `json:"prefix"`      // Guild prefix
	Level       float64       `json:"level"`       // Guild level
	XP          float64       `json:"xp"`          // Total XP
	Members     []GuildMember `json:"members"`     // List of guild members
	Banner      GuildBanner   `json:"banner"`      // Guild banner details
	Created     string        `json:"created"`     // Guild creation date
	Territories int           `json:"territories"` // Number of territories controlled
	Ranking     int           `json:"ranking"`     // Global ranking
	UUID        string        `json:"uuid"`        // Guild unique identifier
	Stats       GuildStats    `json:"stats"`       // Additional stats about the guild
}

// GuildMember represents a single member of the guild.
type GuildMember struct {
	Name        string  `json:"name"`        // Player name
	Rank        string  `json:"rank"`        // Guild rank (e.g., Recruiter, Captain)
	Contributed float64 `json:"contributed"` // XP contributed by the member
	Joined      string  `json:"joined"`      // Date the player joined the guild
}

// GuildBanner represents the banner used by the guild.
type GuildBanner struct {
	Base     string        `json:"base"`     // Base color of the banner
	Patterns []BannerLayer `json:"patterns"` // List of patterns on the banner
}

// BannerLayer represents a single layer in the guild banner.
type BannerLayer struct {
	Pattern string `json:"pattern"` // Type of pattern
	Color   string `json:"color"`   // Color of the pattern
}

// GuildStats represents additional stats about the guild.
type GuildStats struct {
	MobsKilled    int `json:"mobs_killed"`    // Total mobs killed by the guild
	PlayersKilled int `json:"players_killed"` // Total players killed by the guild
	ChestsOpened  int `json:"chests_opened"`  // Total chests opened by guild members
	Quests        int `json:"quests"`         // Total quests completed by guild members
	Dungeons      int `json:"dungeons"`       // Total dungeons completed by guild members
	Raids         int `json:"raids"`          // Total raids completed by guild members
}
