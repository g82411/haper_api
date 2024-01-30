package bussinessLogic

//func FindOrCreateTags(ctx context.Context, tagNames []string) ([]models.Tag, error) {
//	//db := ctx.Value("db").(*gorm.DB)
//	//err := db.AutoMigrate(&models.Tag{})
//	//if err != nil {
//	//	panic(fmt.Errorf("Migrate Tag Table: %v", err))
//	//}
//	var existTags []models.Tag
//db.Where("name IN ?", tagNames).Find(&existTags)
//	//existingTagMapping := make(map[string]bool)
//	//for _, tag := range existTags {
//	//	existingTagMapping[tag.Name] = true
//	//}
//	//var newTags []models.Tag
//	//for _, tagName := range tagNames {
//	//	if _, ok := existingTagMapping[tagName]; !ok {
//	//		id, _ := uuid.NewUUID()
//	//		newTags = append(newTags, models.Tag{
//	//			ID:   id.String(),
//	//			Name: tagName,
//	//		})
//	//	}
//	//}
//	//if len(newTags) > 0 {
//	//	err = db.Create(&newTags).Error
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//}
//	//var allTags []models.Tag
//	//db.Debug().Where("name IN ?", tagNames).Find(&allTags)
//	return existTags, nil
//}
