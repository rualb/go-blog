package service

func mustCreateRepository(appService AppService) {

	db := appService.Repository() // not full inited

	for _, x := range []any{
		&BlogPost{},
	} {

		if err := db.AutoMigrate(x); err != nil {
			panic(err)
		}

	}

	mustInitRepositoryMasterData(appService) // not full inited

}

func mustInitRepositoryMasterData(_ AppService) {

}
