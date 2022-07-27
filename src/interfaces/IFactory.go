package interfaces

import (
	"github.com/vireocloud/property-pros-service/src/grpc"
)

type IFactory interface {
	CreateGroup(Id uint32, name string, title string, activities interface{}) IGroup
	CreateActivity(
		Id uint32,
		Name string,
		Title string,
		Description string,
		GroupID uint32,
		MeetingAgendaItems interface{}) IActivity

	CreateMeetingAgendaItem(
		Id uint32,
		Name string,
		Description string,
		ProSummary string,
		ConSummary string,
		TalkingPoints string,
		MessageTemplate string,
		RelevantContacts string) IMeetingAgendaItem
	CreateGroupModelsFromSaveGroupsRequest(req *grpc.SaveGroupsRequest) []IGroup
	CreateGroupModelsFromSaveGroupRequest(req *grpc.SaveGroupRequest) []IGroup
	CreateGroupModelFromDeleteGroupRequest(req *grpc.DeleteGroupRequest) IGroup

	CreateActivityModelsFromSaveActivitiesRequest(req *grpc.SaveActivitiesRequest) []IActivity
	CreateActivityModelFromSaveActivityRequest(req *grpc.SaveActivityRequest) IActivity
	CreateActivityModelFromDeleteActivityRequest(req *grpc.DeleteActivityRequest) IActivity
	// CreateActivityModels(*grpc.SaveGroupRequest) []IActivity

	CreateMeetingAgendaItemModelsFromSaveMeetingAgendaItemsRequest(req *grpc.SaveMeetingAgendaItemsRequest) []IMeetingAgendaItem
	CreateMeetingAgendaItemModelsFromSaveMeetingAgendaItemRequest(req *grpc.SaveMeetingAgendaItemRequest) []IMeetingAgendaItem
	CreateMeetingAgendaItemModelFromDeleteMeetingAgendaItemRequest(req *grpc.DeleteMeetingAgendaItemRequest) IMeetingAgendaItem
}
