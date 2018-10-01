package webhook

import (
	"fmt"
	"time"
	"vector-sms/JiraAPI/Jiraauthentication"
	"vector-sms/JiraAPI/Jiraconfig"
	"vector-sms/dao"

	jira "github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2/bson"
)

//JiraWebhook - to represent the JIRA webhook structure
type JiraWebhook struct {
	Id        bson.ObjectId `bson:"_id"`
	TimeStamp string        `json:"timestamp" bson:"timestamp"`
	Issue     struct {
		Id     string `json:"id" bson:"id"`
		Self   string `json:"self" bson:"self"`
		Key    string `json:"key" bson:"key"`
		Fields struct {
			Summary     string    `json:"summary" bson:"summary"`
			Created     time.Time `json:"created" bson:"created"`
			Description string    `json:"description" bson:"description"`
			Labels      []string  `json:"labels" bson:"labels"`
			Priority    string    `json:"priority" bson:"priority"`
		} `json:"feilds" bson:"feilds"`
	} `json:"issue" bson:"issue" `
	User struct {
		Self         string `json:"self" bson:"self"`
		Name         string `json:"name" bson:"name"`
		Key          string `json:"key" bson:"key"`
		EmailAddress string `json:"emailAddress" bson:"emailAddress"`
	} `json:"user" bson:"user"`
	ChangeLog struct {
		Items []struct {
			ToString   string `json:"toString" bson:"toString"`
			To         string `json:"to" bson:"to"`
			FromString string `json:"fromString" bson:"fromString"`
			From       string `json:"from" bson:"from"`
			FeildType  string `json:"fieldtype" bson:"fieldtype"`
			Feild      string `json:"field" bson:"field"`
		} `json:"items" bson:"items"`
		Id string `json:"id" bson:"id"`
	} `json:"changelog" bson:"changelog"`
	Comment struct {
		Body         string `json:"body"`
		UpdateAuthor struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls   struct {
				One6X16  string `json:"16x16"`
				Four8X48 string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			Self        string `json:"self"`
		} `json:"updateAuthor"`
		Created string `json:"created"`
		Updated string `json:"updated"`
		Self    string `json:"self"`
		ID      string `json:"id"`
		Author  struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls   struct {
				One6X16  string `json:"16x16"`
				Four8X48 string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			Self        string `json:"self"`
		} `json:"author"`
	} `json:"comment"`
	WebhookEvent   string `json:"webhookEvent" bson:"webhookEvent"`
	IssueEventType string `json:"issue_event_type_name" bson:"issue_event_type_name"`
}

//HandleWebhook -
func HandleWebhook(jsonString string, m *dao.DevicesDAO) error {

	///////////////////
	ConfStruct := Jiraconfig.LoadConfiguration("./config/jira_config.json")

	jiraClient := Jiraauthentication.GetJiraClientBasicAuthentication(ConfStruct)
	/////////////////////

	webHookEvent := gjson.Get(jsonString, "webhookEvent").String()
	switch webHookEvent {
	case "jira:issue_created":
		err := handleIssue("Create", m, jsonString, jiraClient)
		return err
	case "jira:issue_updated":
		err := handleIssue("Update", m, jsonString, jiraClient)
		return err
	case "jira:issue_deleted":
		err := handleIssue("Delete", m, jsonString, jiraClient)
		return err
	case "comment_created":
		err := handleComment("Create", m, jsonString, jiraClient)
		return err
	case "comment_updated":
		err := handleComment("Update", m, jsonString, jiraClient)
		return err
	case "comment_deleted":
		err := handleComment("Delete", m, jsonString, jiraClient)
		return err

	}
	return nil
}

//handleIssue -
func handleIssue(event string, m *dao.DevicesDAO, jsonString string, jiraClient *jira.Client) error {

	issueKey := gjson.Get(jsonString, "issue.key").String()
	issueURL := gjson.Get(jsonString, "issue.self").String()
	fmt.Print("HandleIssue Function: " + issueKey)
	contextLogger := log.WithFields(log.Fields{
		"URL": issueURL,
	})
	switch event {
	case "Create":
		contextLogger.Info("A new ticket was Created in Jira")
	case "Update":
		dao.UpdateTicketStatus(issueKey, m, jiraClient)
		contextLogger.Info("A ticket was Updated.")
	case "Delete":
		ticket, _ := m.FindTicketByID(issueKey)
		ticket.TicketStatus = "DELETED"
		m.UpdateTicket(ticket)
		contextLogger.Info("A ticket was Deleted.")
	}

	return nil
}

func handleComment(event string, m *dao.DevicesDAO, jsonString string, jiraClient *jira.Client) error {
	issueKey := gjson.Get(jsonString, "issue.key").String()
	issueURL := gjson.Get(jsonString, "issue.self").String()

	fmt.Print("HandleIssue Function: " + issueKey)
	contextLogger := log.WithFields(log.Fields{
		"URL": issueURL,
	})

	switch event {
	case "Create":
		contextLogger.Info("A new comment was Created in Jira")

	case "Update":
		contextLogger.Info("A Comment was Updated.")

	case "Delete":
		contextLogger.Info("A Comment was Deleted.")
	}
	return nil

}
