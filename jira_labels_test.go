package main

import (
	"encoding/json"
	"testing"

	"github.com/michael-go/go-jsn/jsn"

	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/assert"
)

// Test GetJiraTicket function
func TestGetJiraTicketFunc(t *testing.T) {
	expectedTestURL := "/v1/org/123/project/123/jira-issues"
	assert := assert.New(t)
	server := HTTPResponseCheckAndStub(expectedTestURL, "existingJiraTickets")

	defer server.Close()

	response := getJiraTickets(server.URL, "123", "123", "123")

	opts := jsondiff.DefaultConsoleOptions()
	marshalledResp, _ := json.Marshal(response)
	comparison, _ := jsondiff.Compare(readFixture("./fixtures/results/ticketRefs.json"), marshalledResp, &opts)
	assert.Equal("FullMatch", comparison.String())

	return
}

// Test openJiraTickets function
func TestOpenJiraTicketWithLabelsFunc(t *testing.T) {
	assert := assert.New(t)
	server := HTTPResponseStubAndMirrorRequest("/v1/org/123/project/12345678-1234-1234-1234-123456789012/issue/SNYK-JS-MINIMIST-559764/jira-issue", "", "")

	defer server.Close()

	projectInfo, _ := jsn.NewJson(readFixture("./fixtures/project.json"))
	vulnsForJira := make(map[string]interface{})
	err := json.Unmarshal(readFixture("./fixtures/vulnsForJira.json"), &vulnsForJira)
	if err != nil {
		panic(err)
	}
	jiraResponse, err := openJiraTickets(server.URL, "123", "123", "123", "Bug", "", "Label1,Label2", projectInfo, vulnsForJira, false)
	if err != nil {
		panic(err)
	}
	var mirroredResponse mirroredResponse
	if err := json.Unmarshal([]byte(jiraResponse), &mirroredResponse); err != nil {
		panic(err)
	}
	assert.Equal(string(readFixture("./fixtures/results/jiraTicketWithLabels.json")), string(mirroredResponse.Body))

	//expectedTestURL := "/v1/org/"+orgID+"/project/"+projectInfo.K("id").String().Value+"/issue/"+vulnID+"/jira-issue"
	//"/v1/org/123/project/123/jira-issues"

	return
}

func TestOpenJiraTicketWithoutLabelsFunc(t *testing.T) {
	assert := assert.New(t)
	server := HTTPResponseStubAndMirrorRequest("/v1/org/123/project/12345678-1234-1234-1234-123456789012/issue/SNYK-JS-MINIMIST-559764/jira-issue", "", "")

	defer server.Close()

	projectInfo, _ := jsn.NewJson(readFixture("./fixtures/project.json"))
	vulnsForJira := make(map[string]interface{})
	err := json.Unmarshal(readFixture("./fixtures/vulnsForJira.json"), &vulnsForJira)
	if err != nil {
		panic(err)
	}
	jiraResponse, err := openJiraTickets(server.URL, "123", "123", "123", "Bug", "", "", projectInfo, vulnsForJira, false)
	if err != nil {
		panic(err)
	}
	var mirroredResponse mirroredResponse
	if err := json.Unmarshal([]byte(jiraResponse), &mirroredResponse); err != nil {
		panic(err)
	}
	assert.Equal(string(readFixture("./fixtures/results/jiraTicketWithoutLabels.json")), string(mirroredResponse.Body))

	//expectedTestURL := "/v1/org/"+orgID+"/project/"+projectInfo.K("id").String().Value+"/issue/"+vulnID+"/jira-issue"
	//"/v1/org/123/project/123/jira-issues"

	return
}
