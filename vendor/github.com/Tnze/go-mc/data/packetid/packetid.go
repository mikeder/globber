// This file is automatically generated by gen_packetIDs.go. DO NOT EDIT.

package packetid

// Login state
const (
	// Clientbound
	Disconnect                 = 0x0
	EncryptionBeginClientbound = 0x1
	Success                    = 0x2
	Compress                   = 0x3
	LoginPluginRequest         = 0x4

	// Serverbound
	LoginStart                 = 0x0
	EncryptionBeginServerbound = 0x1
	LoginPluginResponse        = 0x2
)

// Ping state
const (
	// Clientbound
	ServerInfo      = 0x0
	PingClientbound = 0x1

	// Serverbound
	PingStart       = 0x0
	PingServerbound = 0x1
)

// Play state
const (
	// Clientbound
	SpawnEntity                = 0x0
	SpawnEntityExperienceOrb   = 0x1
	SpawnEntityLiving          = 0x2
	SpawnEntityPainting        = 0x3
	NamedEntitySpawn           = 0x4
	SculkVibrationSignal       = 0x5
	Animation                  = 0x6
	Statistics                 = 0x7
	AcknowledgePlayerDigging   = 0x8
	BlockBreakAnimation        = 0x9
	TileEntityData             = 0xa
	BlockAction                = 0xb
	BlockChange                = 0xc
	BossBar                    = 0xd
	Difficulty                 = 0xe
	ChatClientbound            = 0xf
	ClearTitles                = 0x10
	TabCompleteClientbound     = 0x11
	DeclareCommands            = 0x12
	CloseWindowClientbound     = 0x13
	WindowItems                = 0x14
	CraftProgressBar           = 0x15
	SetSlot                    = 0x16
	SetCooldown                = 0x17
	CustomPayloadClientbound   = 0x18
	NamedSoundEffect           = 0x19
	KickDisconnect             = 0x1a
	EntityStatus               = 0x1b
	Explosion                  = 0x1c
	UnloadChunk                = 0x1d
	GameStateChange            = 0x1e
	OpenHorseWindow            = 0x1f
	InitializeWorldBorder      = 0x20
	KeepAliveClientbound       = 0x21
	MapChunk                   = 0x22
	WorldEvent                 = 0x23
	WorldParticles             = 0x24
	UpdateLight                = 0x25
	Login                      = 0x26
	Map                        = 0x27
	TradeList                  = 0x28
	RelEntityMove              = 0x29
	EntityMoveLook             = 0x2a
	EntityLook                 = 0x2b
	VehicleMoveClientbound     = 0x2c
	OpenBook                   = 0x2d
	OpenWindow                 = 0x2e
	OpenSignEntity             = 0x2f
	Ping                       = 0x30
	CraftRecipeResponse        = 0x31
	AbilitiesClientbound       = 0x32
	EndCombatEvent             = 0x33
	EnterCombatEvent           = 0x34
	DeathCombatEvent           = 0x35
	PlayerInfo                 = 0x36
	FacePlayer                 = 0x37
	PositionClientbound        = 0x38
	UnlockRecipes              = 0x39
	DestroyEntity              = 0x3a
	RemoveEntityEffect         = 0x3b
	ResourcePackSend           = 0x3c
	Respawn                    = 0x3d
	EntityHeadRotation         = 0x3e
	MultiBlockChange           = 0x3f
	SelectAdvancementTab       = 0x40
	ActionBar                  = 0x41
	WorldBorderCenter          = 0x42
	WorldBorderLerpSize        = 0x43
	WorldBorderSize            = 0x44
	WorldBorderWarningDelay    = 0x45
	WorldBorderWarningReach    = 0x46
	Camera                     = 0x47
	HeldItemSlotClientbound    = 0x48
	UpdateViewPosition         = 0x49
	UpdateViewDistance         = 0x4a
	SpawnPosition              = 0x4b
	ScoreboardDisplayObjective = 0x4c
	EntityMetadata             = 0x4d
	AttachEntity               = 0x4e
	EntityVelocity             = 0x4f
	EntityEquipment            = 0x50
	Experience                 = 0x51
	UpdateHealth               = 0x52
	ScoreboardObjective        = 0x53
	SetPassengers              = 0x54
	Teams                      = 0x55
	ScoreboardScore            = 0x56
	SetTitleSubtitle           = 0x57
	UpdateTime                 = 0x58
	SetTitleText               = 0x59
	SetTitleTime               = 0x5a
	EntitySoundEffect          = 0x5b
	SoundEffect                = 0x5c
	StopSound                  = 0x5d
	PlayerlistHeader           = 0x5e
	NbtQueryResponse           = 0x5f
	Collect                    = 0x60
	EntityTeleport             = 0x61
	Advancements               = 0x62
	EntityUpdateAttributes     = 0x63
	EntityEffect               = 0x64
	DeclareRecipes             = 0x65
	Tags                       = 0x66

	// Serverbound
	TeleportConfirm            = 0x0
	QueryBlockNbt              = 0x1
	SetDifficulty              = 0x2
	ChatServerbound            = 0x3
	ClientCommand              = 0x4
	Settings                   = 0x5
	TabCompleteServerbound     = 0x6
	EnchantItem                = 0x7
	WindowClick                = 0x8
	CloseWindowServerbound     = 0x9
	CustomPayloadServerbound   = 0xa
	EditBook                   = 0xb
	QueryEntityNbt             = 0xc
	UseEntity                  = 0xd
	GenerateStructure          = 0xe
	KeepAliveServerbound       = 0xf
	LockDifficulty             = 0x10
	PositionServerbound        = 0x11
	PositionLook               = 0x12
	Look                       = 0x13
	Flying                     = 0x14
	VehicleMoveServerbound     = 0x15
	SteerBoat                  = 0x16
	PickItem                   = 0x17
	CraftRecipeRequest         = 0x18
	AbilitiesServerbound       = 0x19
	BlockDig                   = 0x1a
	EntityAction               = 0x1b
	SteerVehicle               = 0x1c
	Pong                       = 0x1d
	DisplayedRecipe            = 0x1e
	RecipeBook                 = 0x1f
	NameItem                   = 0x20
	ResourcePackReceive        = 0x21
	AdvancementTab             = 0x22
	SelectTrade                = 0x23
	SetBeaconEffect            = 0x24
	HeldItemSlotServerbound    = 0x25
	UpdateCommandBlock         = 0x26
	UpdateCommandBlockMinecart = 0x27
	SetCreativeSlot            = 0x28
	UpdateJigsawBlock          = 0x29
	UpdateStructureBlock       = 0x2a
	UpdateSign                 = 0x2b
	ArmAnimation               = 0x2c
	Spectate                   = 0x2d
	BlockPlace                 = 0x2e
	UseItem                    = 0x2f
)
