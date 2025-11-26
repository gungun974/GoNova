package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeResourceCmd)
}

var makeResourceCmd = &cobra.Command{
	Use:   "make:resource (ResourceName)",
	Short: "Create a new Entity with controller, usecase, model, repository and presenter",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeResource,
}

func MakeResource(cmd *cobra.Command, args []string) {
	entityName := ""
	if len(args) == 0 {
		entityName = form.AskInputWithPlaceholder("The Resource name :", "Post")
	} else {
		entityName = args[0]
	}

	entityName = helpers.CapitalizeFirstLetter(entityName)

	entities := analyzer.AnalyzeProjectEntities()

	var entity *analyzer.AnalyzedEntity

	for _, e := range entities {
		if e.Name == entityName {
			entity = &e
			break
		}
	}

	if entity == nil {
		err := actions.MakeEntity(entityName)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Make Entity : %v", err)
		}
	}

	entities = analyzer.AnalyzeProjectEntities()

	for _, e := range entities {
		if e.Name == entityName {
			entity = &e
			break
		}
	}

	if entity == nil {
		logger.MainLogger.Fatalf("Failed to find Entity")
	}

	models := analyzer.AnalyzeProjectModels(entities)

	var model *analyzer.AnalyzedModel

	for _, m := range models {
		if m.Entity.Name == entity.Name {
			model = &m
			break
		}
	}

	if model == nil {
		err := actions.MakeModel(*entity)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Make Model : %v", err)
		}
	}

	models = analyzer.AnalyzeProjectModels(entities)

	for _, m := range models {
		if m.Entity.Name == entity.Name {
			model = &m
			break
		}
	}

	if model == nil {
		logger.MainLogger.Fatalf("Failed to find Model")
	}

	shouldLinkRepository := false

	repositories := analyzer.AnalyzeProjectRepositories()

	var repository *analyzer.AnalyzedRepository

	for _, r := range repositories {
		if r.Name == entity.Name+"Repository" {
			repository = &r
			break
		}
	}

	if repository == nil {
		_, err := actions.MakeRepository(entity.Name, model)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Make Repository : %v", err)
		}
		shouldLinkRepository = true
	}

	repositories = analyzer.AnalyzeProjectRepositories()

	for _, r := range repositories {
		if r.Name == entity.Name+"Repository" {
			repository = &r
			break
		}
	}

	if repository == nil {
		logger.MainLogger.Fatalf("Failed to find Repository")
	}

	shouldLinkPresenter := false

	presenters := analyzer.AnalyzeProjectPresenters()

	var presenter *analyzer.AnalyzedPresenter

	for _, p := range presenters {
		if p.Name == entity.Name+"Presenter" {
			presenter = &p
			break
		}
	}

	if presenter == nil {
		_, err := actions.MakePresenter(entity.Name)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Make Presenter : %v", err)
		}
		shouldLinkPresenter = true
	}

	presenters = analyzer.AnalyzeProjectPresenters()

	for _, p := range presenters {
		if p.Name == entity.Name+"Presenter" {
			presenter = &p
			break
		}
	}

	if presenter == nil {
		logger.MainLogger.Fatalf("Failed to find Presenter")
	}

	shouldLinkUsecase := false

	storages := analyzer.AnalyzeProjectStorages()
	usecases := analyzer.AnalyzeProjectUsecases(repositories, storages, presenters)

	var usecase *analyzer.AnalyzedUsecase

	for _, u := range usecases {
		if u.Name == entity.Name+"Usecase" {
			usecase = &u
			break
		}
	}

	if usecase == nil {
		_, err := actions.MakeUsecase(entity.Name)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Make Usecase : %v", err)
		}
		shouldLinkUsecase = true
	}

	usecases = analyzer.AnalyzeProjectUsecases(repositories, storages, presenters)

	for _, u := range usecases {
		if u.Name == entity.Name+"Usecase" {
			usecase = &u
			break
		}
	}

	if usecase == nil {
		logger.MainLogger.Fatalf("Failed to find Usecase")
	}

	shouldLinkController := false

	controllers := analyzer.AnalyzeProjectControllers(usecases)

	var controller *analyzer.AnalyzedController

	for _, c := range controllers {
		if c.Name == entity.Name+"Controller" {
			controller = &c
			break
		}
	}

	if controller == nil {
		_, err := actions.MakeController(entity.Name)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Make Controller : %v", err)
		}
		shouldLinkController = true
	}

	controllers = analyzer.AnalyzeProjectControllers(usecases)

	for _, c := range controllers {
		if c.Name == entity.Name+"Controller" {
			controller = &c
			break
		}
	}

	if controller == nil {
		logger.MainLogger.Fatalf("Failed to find Controller")
	}

	if shouldLinkRepository {
		err := actions.LinkRepository(*repository, *usecase)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Link Repository : %v", err)
		}

		repositories = analyzer.AnalyzeProjectRepositories()
		presenters = analyzer.AnalyzeProjectPresenters()
		usecases = analyzer.AnalyzeProjectUsecases(repositories, storages, presenters)
		analyzer.DeepAnalyzeProjectController(controller, usecases)
	}

	if shouldLinkPresenter {
		err := actions.LinkPresenter(*presenter, *usecase)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Link Presenter : %v", err)
		}

		repositories = analyzer.AnalyzeProjectRepositories()
		presenters = analyzer.AnalyzeProjectPresenters()
		usecases = analyzer.AnalyzeProjectUsecases(repositories, storages, presenters)
		analyzer.DeepAnalyzeProjectController(controller, usecases)
	}

	if shouldLinkUsecase {
		err := actions.LinkUsecase(*usecase, *controller)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Link Usecase : %v", err)
		}

		repositories = analyzer.AnalyzeProjectRepositories()
		presenters = analyzer.AnalyzeProjectPresenters()
		usecases = analyzer.AnalyzeProjectUsecases(repositories, storages, presenters)
		analyzer.DeepAnalyzeProjectController(controller, usecases)
	}

	if shouldLinkController {
		err := actions.LinkController(*controller)
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Link Controller : %v", err)
		}
	}
}
