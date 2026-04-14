package common

import "github.com/gagliardetto/solana-go"

var (
	FarmProgramIDV3 = mustPublicKey("EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q")
	FarmProgramIDV4 = mustPublicKey("CBuCnLe26faBpcBP2fktp4rp8abpcAnTWft6ZrP5Q4T")
	FarmProgramIDV5 = mustPublicKey("9KEPoZmtHUrBbhWN1v1KWLMkkvwY6WLtAVUCPRtRjP4z")
	FarmProgramIDV6 = mustPublicKey("FarmqiPv5eAj3j1GMdMCMUGXqPUvmquZtMy86QH6rzhG")

	UTIL1216 = mustPublicKey("CLaimxFqjHzgTJtAGHU47NPhg6qrc5sCnpC4tBLyABQS")

	OpenBookProgram = mustPublicKey("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")
	SerumProgramIDV3 = mustPublicKey("9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin")

	AMMV4  = mustPublicKey("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	AMMStable = mustPublicKey("5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h")
	LiquidityPoolProgramIDV5Model = mustPublicKey("CDSr3ssLcRB6XYPJwAfFt18MZvEZp4LjHcvzBVZ45duo")
	CLMMProgramID     = mustPublicKey("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")
	CLMMLockProgramID = mustPublicKey("LockrWmn6K5twhz3y9w1dQERbmgSaRkfnTeTKbpofwE")
	CLMMLockAuthID    = mustPublicKey("kN1kEznaF5Xbd8LYuqtEFcxzWSBk5Fv6ygX6SqEGJVy")

	ModelDataPubKey = mustPublicKey("CDSr3ssLcRB6XYPJwAfFt18MZvEZp4LjHcvzBVZ45duo")
	RouterProgram   = mustPublicKey("routeUGWgWzqBWFcrCfv8tritsqukccJPu3q5GPP3xS")
	FeeDestinationID = mustPublicKey("7YttLkHDoNj9wyDur5pM1ejNaAvT9X4eqaYcHQqtj2G5")

	IDOProgramIDV1 = mustPublicKey("6FJon3QE27qgPVggARueB22hLvoh22VzJpXv4rBEoSLF")
	IDOProgramIDV2 = mustPublicKey("CC12se5To1CdEuw7fDS27B7Geo5jJyL7t5UK2B44NgiH")
	IDOProgramIDV3 = mustPublicKey("9HzJyW1qZsEiSfMUf6L2jo3CcTKAyBmSyKdwQeYisHrC")
	IDOProgramIDV4 = mustPublicKey("DropEU8AvevN3UrXWXTMuz3rqnMczQVNjq3kcSdW2SQi")

	CreateCPMMPoolProgram = mustPublicKey("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")
	CreateCPMMPoolAuth    = mustPublicKey("GpMZbSM2GgvTKHJirzeGfMFoaZ8UR2X7F4v8vHTvxFbL")
	CreateCPMMPoolFeeAcc  = mustPublicKey("DNXgeM9EiiaAbaWvwjHj9fQQLAX5ZsfHyvmYUNRAdNC8")

	LockCPMMProgram = mustPublicKey("LockrWmn6K5twhz3y9w1dQERbmgSaRkfnTeTKbpofwE")
	LockCPMMAuth    = mustPublicKey("3f7GcQFG397GAaEnv51zR6tsTVihYRydnydDD1cXekxH")

	LaunchpadProgram = mustPublicKey("LanMV9sAd7wArD4vJFi2qDdfnVhFxYSUg6eADduJ3uj")
	LaunchpadAuth    = mustPublicKey("WLHv2UAZm6z4KyaaELi5pjdbJh6RESMva1Rnn8pJVVh")
	DevLaunchpadProgram = mustPublicKey("DRay6fNdQ5J82H7xV6uq2aV3mNrUZ1J4PgSKsWgptcm6")
	DevLaunchpadAuth    = mustPublicKey("5xqNaZXX5eUi4p5HU4oz9i5QnwRNT2y6oN7yyn4qENeq")
	LaunchpadPlatform = mustPublicKey("4Bu96XjU84XjPDSpveTVf6LYGCkfW5FK7SNkREWcEfV4")
	LaunchpadConfig   = mustPublicKey("6s1xP3hpbAfFoNtUNF8mfHsjr2Bd97JxFJRWLbL6aHuX")
)

type ProgramIDConfig map[string]solana.PublicKey

var IDOAllProgram = ProgramIDConfig{
	"IDO_PROGRAM_ID_V1": IDOProgramIDV1,
	"IDO_PROGRAM_ID_V2": IDOProgramIDV2,
	"IDO_PROGRAM_ID_V3": IDOProgramIDV3,
	"IDO_PROGRAM_ID_V4": IDOProgramIDV4,
}

var AllProgramID = ProgramIDConfig{
	"AMM_V4":                   AMMV4,
	"AMM_STABLE":               AMMStable,
	"CLMM_PROGRAM_ID":          CLMMProgramID,
	"CLMM_LOCK_PROGRAM_ID":     CLMMLockProgramID,
	"CLMM_LOCK_AUTH_ID":        CLMMLockAuthID,
	"FARM_PROGRAM_ID_V3":       FarmProgramIDV3,
	"FARM_PROGRAM_ID_V4":       FarmProgramIDV4,
	"FARM_PROGRAM_ID_V5":       FarmProgramIDV5,
	"FARM_PROGRAM_ID_V6":       FarmProgramIDV6,
	"OPEN_BOOK_PROGRAM":        OpenBookProgram,
	"SERUM_PROGRAM_ID_V3":      SerumProgramIDV3,
	"UTIL1216":                 UTIL1216,
	"Router":                   RouterProgram,
	"CREATE_CPMM_POOL_PROGRAM": CreateCPMMPoolProgram,
	"CREATE_CPMM_POOL_AUTH":    CreateCPMMPoolAuth,
	"CREATE_CPMM_POOL_FEE_ACC": CreateCPMMPoolFeeAcc,
	"LOCK_CPMM_PROGRAM":        LockCPMMProgram,
	"LOCK_CPMM_AUTH":           LockCPMMAuth,
	"LAUNCHPAD_PROGRAM":        LaunchpadProgram,
	"LAUNCHPAD_AUTH":           LaunchpadAuth,
	"LAUNCHPAD_PLATFORM":       LaunchpadPlatform,
	"LAUNCHPAD_CONFIG":         LaunchpadConfig,
	"FEE_DESTINATION_ID":       FeeDestinationID,
	"MODEL_DATA_PUBKEY":        ModelDataPubKey,
}

var DevnetProgramID = ProgramIDConfig{
	"OPEN_BOOK_PROGRAM":        mustPublicKey("EoTcMgcDRTJVZDMZWBoU6rhYHZfkNTVEAfz3uUJRcYGj"),
	"SERUM_PROGRAM_ID_V3":      mustPublicKey("Ray1111111111111111111111111111111111111111"),
	"AMM_V4":                   mustPublicKey("DRaya7Kj3aMWQSy19kSjvmuwq9docCHofyP9kanQGaav"),
	"AMM_STABLE":               mustPublicKey("DRayDdXc1NZQ9C3hRWmoSf8zK4iapgMnjdNZWrfwsP8m"),
	"CLMM_PROGRAM_ID":          mustPublicKey("DRayAUgENGQBKVaX8owNhgzkEDyoHTGVEGHVJT1E9pfH"),
	"CLMM_LOCK_PROGRAM_ID":     mustPublicKey("DRay25Usp3YJAi7beckgpGUC7mGJ2cR1AVPxhYfwVCUX"),
	"CLMM_LOCK_AUTH_ID":        mustPublicKey("6Aoh8h2Lw2m5UGxYR8AdAL87jTWYeKoxM52mJRzfYwN"),
	"CREATE_CPMM_POOL_PROGRAM": mustPublicKey("DRaycpLY18LhpbydsBWbVJtxpNv9oXPgjRSfpF2bWpYb"),
	"CREATE_CPMM_POOL_AUTH":    mustPublicKey("CXniRufdq5xL8t8jZAPxsPZDpuudwuJSPWnbcD5Y5Nxq"),
	"CREATE_CPMM_POOL_FEE_ACC": mustPublicKey("3oE58BKVt8KuYkGxx8zBojugnymWmBiyafWgMrnb6eYy"),
	"LOCK_CPMM_PROGRAM":        mustPublicKey("DRay25Usp3YJAi7beckgpGUC7mGJ2cR1AVPxhYfwVCUX"),
	"LOCK_CPMM_AUTH":           mustPublicKey("7qWVV8UY2bRJfDLP4s37YzBPKUkVB46DStYJBpYbQzu3"),
	"UTIL1216":                 solana.PublicKey{},
	"Router":                   mustPublicKey("DRaybByLpbUL57LJARs3j8BitTxVfzBg351EaMr5UTCd"),
	"FARM_PROGRAM_ID_V3":       mustPublicKey("DRayWyrLmEW5KEeqs8kdTMMaBabapqagaBC7KWpGtJeZ"),
	"FARM_PROGRAM_ID_V4":       mustPublicKey("Ray1111111111111111111111111111111111111111"),
	"FARM_PROGRAM_ID_V5":       mustPublicKey("DRayiCGSZgku1GTK6rXD6mVDdingXy6APAH1R6R5L2LC"),
	"FARM_PROGRAM_ID_V6":       mustPublicKey("DRayzbYakXs45ELHkzH6vC3fuhQqTAnv5A68gdFuvZyZ"),
	"LAUNCHPAD_PROGRAM":        mustPublicKey("DRay6fNdQ5J82H7xV6uq2aV3mNrUZ1J4PgSKsWgptcm6"),
	"LAUNCHPAD_AUTH":           mustPublicKey("5xqNaZXX5eUi4p5HU4oz9i5QnwRNT2y6oN7yyn4qENeq"),
	"LAUNCHPAD_PLATFORM":       mustPublicKey("2Jx4KTDrVSdWNazuGpcA8n3ZLTRGGBDxAWhuKe2Xcj2a"),
	"LAUNCHPAD_CONFIG":         mustPublicKey("7ZR4zD7PYfY2XxoG1Gxcy2EgEeGYrpxrwzPuwdUBssEt"),
	"FEE_DESTINATION_ID":       mustPublicKey("9y8ENuuZ3b19quffx9hQvRVygG5ky6snHfRvGpuSfeJy"),
	"MODEL_DATA_PUBKEY":        mustPublicKey("Ray1111111111111111111111111111111111111111"),
}
