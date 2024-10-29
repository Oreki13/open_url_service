package http

//func StartServiceListener(cfg *config.Config) {
//	name := cfg.NoSleepConfig.FlagName
//	topics := strings.Split(cfg.NoSleepConfig.Topics, "|")
//	mController := controller.NewLogController()
//
//	subs := bootstrap.RegistryRabbitMQSubscriber(name, cfg, mController)
//
//	err := subs.Listen(topics)
//
//	if err != nil {
//		logger.Fatal(fmt.Sprintf("Failed to Listen Topic Cause: %v", err))
//	}
//}
