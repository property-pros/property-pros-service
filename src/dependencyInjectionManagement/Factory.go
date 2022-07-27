package dependencyInjectionManagement

import (
	"fmt"

	gateways "github.com/vireocloud/property-pros-service/src/data/gateways"
	models "github.com/vireocloud/property-pros-service/src/domain/models"
	"github.com/vireocloud/property-pros-service/src/grpc"
	i "github.com/vireocloud/property-pros-service/src/interfaces"
)

type Factory struct {
	i.IFactory
	dic IDependencyInjectionContainer
}

func (f *Factory) CreateGroup(id uint32, name string, title string, activities interface{}) i.IGroup {
	group := &models.Group{
		BaseModel:  &models.BaseModel{ID: id},
		Name:       name,
		Title:      title,
		Activities: activities.([]*models.Activity),
	}

	f.dic.GetInstance(func(groupsGateway i.IGroupsGateway) {
		group.SetGateway(groupsGateway)
	})

	return group
}

func (f *Factory) CreateActivity(
	id uint32,
	Name string,
	Description string,
	Title string,
	GroupID uint32,
	MeetingAgendaItems interface{}) i.IActivity {

	activity := &models.Activity{
		BaseModel:          &models.BaseModel{ID: id},
		Name:               Name,
		Description:        Description,
		Title:              Title,
		GroupID:            GroupID,
		MeetingAgendaItems: MeetingAgendaItems.([]*models.MeetingAgendaItem),
	}

	return activity
}

func (f *Factory) CreateGroupModelFromDeleteGroupRequest(req *grpc.DeleteGroupRequest) i.IGroup {
	return f.CreateGroupModel(req.Group)
}

func (f *Factory) CreateGroupModelsFromSaveGroupsRequest(req *grpc.SaveGroupsRequest) []i.IGroup {

	groups := []i.IGroup{}

	for _, group := range req.Groups {
		groups = append(groups, f.CreateGroupModel(group))
	}

	return groups
}

func (f *Factory) CreateGroupModel(g *grpc.Group) i.IGroup {

	var r i.IGroupsGateway

	err := f.dic.GetInstance(func(groupsGateway i.IGroupsGateway) {
		r = groupsGateway
	})

	if err != nil {
		fmt.Printf("failed to get instance of IGroupsGateway: %+v\n", err)
	}

	group := &models.Group{
		BaseModel:  &models.BaseModel{ID: uint32(g.Id)},
		Title:      g.Title,
		Name:       g.Name,
		Activities: f.createConcreteActivityModels(g.Activities),
	}

	group.SetGateway(r)

	return group
}

func (f *Factory) createConcreteActivityModels(req []*grpc.Activity) []*models.Activity {

	activities := []*models.Activity{}

	for _, activity := range req {
		activities = append(activities, f.CreateActivityModel(activity).(*models.Activity))

	}

	return activities
}

func (f *Factory) createActivityModels(req []*grpc.Activity) []i.IActivity {

	activities := []i.IActivity{}

	for _, activity := range req {
		activities = append(activities, f.CreateActivityModel(activity))

	}

	return activities
}

func (f *Factory) CreateActivityModel(a *grpc.Activity) i.IActivity {

	var r *gateways.ActivitiesGateway

	f.dic.GetInstance(func(activitiesGateway i.IActivitiesGateway) {
		r = activitiesGateway.(*gateways.ActivitiesGateway)
	})

	activity := f.CreateActivity(uint32(a.Id), a.Name, a.Description, a.Title, a.GroupId, a.MeetingAgendaItems).(*models.Activity)

	activity.SetGateway(r)

	return activity
}

func (f *Factory) CreateActivityModelsFromSaveActivitiesRequest(request *grpc.SaveActivitiesRequest) []i.IActivity {
	return f.createActivityModels(request.Activities)
}

// CreateActivityModelFromSaveActivityRequest(req *grpc.SaveActivityRequest) IActivity
// CreateActivityModelFromDeleteActivityRequest(req *grpc.DeleteActivityRequest) IActivity
// CreateActivityModels(*grpc.SaveGroupRequest) []IActivity
//
// CreateMeetingAgendaItemModelsFromSaveMeetingAgendaItemsRequest(req *grpc.SaveMeetingAgendaItemsRequest) []IMeetingAgendaItem
// CreateMeetingAgendaItemModelsFromSaveMeetingAgendaItemRequest(req *grpc.SaveMeetingAgendaItemRequest) []IMeetingAgendaItem
// CreateMeetingAgendaItemModelFromDeleteMeetingAgendaItemRequest(req *grpc.DeleteMeetingAgendaItemRequest) IMeetingAgendaItem
