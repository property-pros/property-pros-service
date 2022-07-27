package dependencyInjectionManagement

import (
	"github.com/jinzhu/gorm"
	dig "go.uber.org/dig"

	//according to gorm docs, this import is necessary
	_ "github.com/jinzhu/gorm/dialects/postgres"

	i "github.com/vireocloud/property-pros-service/src/interfaces"

	controllers "github.com/vireocloud/property-pros-service/src/controllers"
	gateways "github.com/vireocloud/property-pros-service/src/data/gateways"
	dr "github.com/vireocloud/property-pros-service/src/data/repositories"
	domainModels "github.com/vireocloud/property-pros-service/src/domain/models"
	domainServices "github.com/vireocloud/property-pros-service/src/domain/services"
)

var Dic IDependencyInjectionContainer
var dic *dig.Container
var db *gorm.DB

type IDependencyInjectionContainer interface {
	GetInstance(function interface{}) error
	Teardown()
}

type DependencyInjectionContainer struct {
	IDependencyInjectionContainer
}

func (c *DependencyInjectionContainer) GetInstance(callback interface{}) error {
	return dic.Invoke(callback)
}

func GetContainer() IDependencyInjectionContainer {
	return Dic
}

func migrate() {
	// db.Exec(fmt.Sprintf("CREATE  SCHEMA IF NOT EXISTS %s;", "localize"))
	// db.Exec(fmt.Sprintf("SET search_path TO %s;", "localize"))

	// Migrate the schema
	db.AutoMigrate(&domainModels.Group{}, &domainModels.Activity{})

	// db.DB().SetMaxIdleConns(0)
}

func (*DependencyInjectionContainer) Teardown() {
	// db.DropTable(&domainModels.Group{}, &domainModels.Activity{})

	defer db.Close()
	// db.Exec(fmt.Sprintf("DROP  SCHEMA IF EXISTS %s;", "localize"))
}

func Initialize(dbConnection *gorm.DB) IDependencyInjectionContainer {

	db = dbConnection

	migrate()

	Dic = &DependencyInjectionContainer{}

	dic = dig.New()

	dic.Provide(func() *gorm.DB {
		return db
	})

	dic.Provide(func(db *gorm.DB) i.IRepository {

		repo := &dr.GormRepository{}

		repo.SetDb(db)

		return repo
	})

	dic.Provide(func(repo i.IRepository, baseGateway i.IEntityGateway) (groupsGateway i.IGroupsGateway) {

		gGateway := &gateways.GroupsGateway{BaseGateway: baseGateway.(*gateways.BaseGateway)}

		gGateway.SetRepo(repo)

		return gGateway
	})

	dic.Provide(func(repo i.IRepository, baseGateway i.IEntityGateway) i.IActivitiesGateway {

		aGateway := &gateways.ActivitiesGateway{BaseGateway: baseGateway.(*gateways.BaseGateway)}

		aGateway.SetRepo(repo)

		return aGateway
	})

	dic.Provide(func(repo i.IRepository, baseGateway i.IEntityGateway) i.IMeetingAgendaItemsGateway {

		aGateway := &gateways.MeetingAgendaItemsGateway{BaseGateway: baseGateway.(*gateways.BaseGateway)}

		aGateway.SetRepo(repo)

		return aGateway
	})

	dic.Provide(func(repo i.IRepository) i.IEntityGateway {

		gateway := &gateways.BaseGateway{}

		gateway.SetRepo(repo)

		return gateway
	})

	dic.Provide(func(g i.IGroupsGateway, a i.IActivitiesGateway) i.IGroupsService {

		service := &domainServices.GroupsService{}

		service.SetGroupsGateway(g)
		service.SetActivitiesGateway(a)

		return service
	})

	dic.Provide(func(g i.IActivitiesGateway, m i.IMeetingAgendaItemsGateway) i.IActivitiesService {

		service := &domainServices.ActivitiesService{}

		service.SetMeetingAgendaItemsGateway(g)

		service.SetActivitiesGateway(m)

		return service
	})

	dic.Provide(func(g i.IMeetingAgendaItemsGateway) i.IMeetingAgendaItemsService {

		service := &domainServices.MeetingAgendaItemsService{}

		service.SetMeetingAgendaItemsGateway(g)

		return service
	})

	dic.Provide(func() i.IFactory {

		return &Factory{dic: Dic}
	})

	dic.Provide(func(service i.IGroupsService, factory i.IFactory) *controllers.GroupsController {

		domainModels.SetFactory(factory)

		return &controllers.GroupsController{
			Service: service,
			Factory: factory,
		}
	})

	dic.Provide(func(service i.IActivitiesService, factory i.IFactory) *controllers.ActivitiesController {

		domainModels.SetFactory(factory)

		return &controllers.ActivitiesController{
			Service: service,
			Factory: factory,
		}
	})

	dic.Provide(func(service i.IMeetingAgendaItemsService, factory i.IFactory) *controllers.MeetingAgendaItemsController {

		domainModels.SetFactory(factory)

		return &controllers.MeetingAgendaItemsController{
			Service: service,
			Factory: factory,
		}
	})
	//
	// dic.Provide(func() (i.IController, i.IController, i.IController) {
	//
	// 	var groupsController *controllers.GroupsController
	// 	var activitiesController *controllers.ActivitiesController
	// 	var meetingAgendasController *controllers.MeetingAgendaItemsController
	//
	// 	dic.Invoke(func(groupsControllerInstance *controllers.GroupsController, activitiesControllerInstance *controllers.ActivitiesController, meetingAgendasControllerInstance *controllers.MeetingAgendaItemsController) {
	// 		groupsController = groupsControllerInstance
	// 		activitiesController = activitiesControllerInstance
	// 		meetingAgendasController = meetingAgendasControllerInstance
	// 	})
	//
	// 	return groupsController, activitiesController, meetingAgendasController
	// }, dig.Group("IControler"))

	return Dic
}

type ControllerCollection struct {
	Controllers []i.IController `group:"IController"`
}
