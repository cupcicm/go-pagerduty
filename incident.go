package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Acknowledgement is the data structure of an acknowledgement of an incident.
type Acknowledgement struct {
	At           string    `json:"at,omitempty"`
	Acknowledger APIObject `json:"acknowledger,omitempty"`
}

// PendingAction is the data structure for any pending actions on an incident.
type PendingAction struct {
	Type string `json:"type,omitempty"`
	At   string `json:"at,omitempty"`
}

// Assignment is the data structure for an assignment of an incident
type Assignment struct {
	At       string    `json:"at,omitempty"`
	Assignee APIObject `json:"assignee,omitempty"`
}

// AlertCounts is the data structure holding a summary of the number of alerts by status of an incident.
type AlertCounts struct {
	Triggered uint `json:"triggered,omitempty"`
	Resolved  uint `json:"resolved,omitempty"`
	All       uint `json:"all,omitempty"`
}

// Priority is the data structure describing the priority of an incident.
type Priority struct {
	APIObject
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ResolveReason is the data structure describing the reason an incident was resolved
type ResolveReason struct {
	Type     string    `json:"type,omitempty"`
	Incident APIObject `json:"incident"`
}

// IncidentBody is the datastructure containing data describing the incident.
type IncidentBody struct {
	Type    string `json:"type,omitempty"`
	Details string `json:"details,omitempty"`
}

// Assignee is an individual assigned to an incident.
type Assignee struct {
	Assignee APIObject `json:"assignee"`
}

// FirstTriggerLogEntry is the first LogEntry
type FirstTriggerLogEntry struct {
	CommonLogEntryField
	Incident APIObject `json:"incident,omitempty"`
}

// Incident is a normalized, de-duplicated event generated by a PagerDuty integration.
type Incident struct {
	APIObject
	IncidentNumber       uint                 `json:"incident_number,omitempty"`
	Title                string               `json:"title,omitempty"`
	Description          string               `json:"description,omitempty"`
	CreatedAt            string               `json:"created_at,omitempty"`
	PendingActions       []PendingAction      `json:"pending_actions,omitempty"`
	IncidentKey          string               `json:"incident_key,omitempty"`
	Service              APIObject            `json:"service,omitempty"`
	Assignments          []Assignment         `json:"assignments,omitempty"`
	Acknowledgements     []Acknowledgement    `json:"acknowledgements,omitempty"`
	LastStatusChangeAt   string               `json:"last_status_change_at,omitempty"`
	LastStatusChangeBy   APIObject            `json:"last_status_change_by,omitempty"`
	FirstTriggerLogEntry FirstTriggerLogEntry `json:"first_trigger_log_entry,omitempty"`
	EscalationPolicy     APIObject            `json:"escalation_policy,omitempty"`
	Teams                []APIObject          `json:"teams,omitempty"`
	Priority             *Priority            `json:"priority,omitempty"`
	Urgency              string               `json:"urgency,omitempty"`
	Status               string               `json:"status,omitempty"`
	Id                   string               `json:"id,omitempty"`
	ResolveReason        ResolveReason        `json:"resolve_reason,omitempty"`
	AlertCounts          AlertCounts          `json:"alert_counts,omitempty"`
	Body                 IncidentBody         `json:"body,omitempty"`
	IsMergeable          bool                 `json:"is_mergeable,omitempty"`
	ConferenceBridge     *ConferenceBridge    `json:"conference_bridge,omitempty"`
}

// ListIncidentsResponse is the response structure when calling the ListIncident API endpoint.
type ListIncidentsResponse struct {
	APIListObject
	Incidents []Incident `json:"incidents,omitempty"`
}

// ListIncidentsOptions is the structure used when passing parameters to the ListIncident API endpoint.
type ListIncidentsOptions struct {
	APIListObject
	Since       string   `url:"since,omitempty"`
	Until       string   `url:"until,omitempty"`
	DateRange   string   `url:"date_range,omitempty"`
	Statuses    []string `url:"statuses,omitempty,brackets"`
	IncidentKey string   `url:"incident_key,omitempty"`
	ServiceIDs  []string `url:"service_ids,omitempty,brackets"`
	TeamIDs     []string `url:"team_ids,omitempty,brackets"`
	UserIDs     []string `url:"user_ids,omitempty,brackets"`
	Urgencies   []string `url:"urgencies,omitempty,brackets"`
	TimeZone    string   `url:"time_zone,omitempty"`
	SortBy      string   `url:"sort_by,omitempty"`
	Includes    []string `url:"include,omitempty,brackets"`
}

// ConferenceBridge is a struct for the conference_bridge object on an incident
type ConferenceBridge struct {
	ConferenceNumber string `json:"conference_number,omitempty"`
	ConferenceURL    string `json:"conference_url,omitempty"`
}

// ListIncidents lists existing incidents. It's recommended to use
// ListIncidentsWithContext instead.
func (c *Client) ListIncidents(o ListIncidentsOptions) (*ListIncidentsResponse, error) {
	return c.ListIncidentsWithContext(context.Background(), o)
}

// ListIncidentsWithContext lists existing incidents.
func (c *Client) ListIncidentsWithContext(ctx context.Context, o ListIncidentsOptions) (*ListIncidentsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListIncidentsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// createIncidentResponse is returned from the API when creating a response.
type createIncidentResponse struct {
	Incident Incident `json:"incident"`
}

// CreateIncidentOptions is the structure used when POSTing to the CreateIncident API endpoint.
type CreateIncidentOptions struct {
	Type             string        `json:"type"`
	Title            string        `json:"title"`
	Service          *APIReference `json:"service"`
	Priority         *APIReference `json:"priority"`
	Urgency          string        `json:"urgency,omitempty"`
	IncidentKey      string        `json:"incident_key,omitempty"`
	Body             *APIDetails   `json:"body,omitempty"`
	EscalationPolicy *APIReference `json:"escalation_policy,omitempty"`
	Assignments      []Assignee    `json:"assignments,omitempty"`
}

// ManageIncidentsOptions is the structure used when PUTing updates to incidents to the ManageIncidents func
type ManageIncidentsOptions struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Status      string        `json:"status,omitempty"`
	Priority    *APIReference `json:"priority,omitempty"`
	Assignments []Assignee    `json:"assignments,omitempty"`
	Resolution  string        `json:"resolution,omitempty"`
}

// MergeIncidentsOptions is the structure used when merging incidents with MergeIncidents func
type MergeIncidentsOptions struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// CreateIncident creates an incident synchronously without a corresponding
// event from a monitoring service. It's recommended to use
// CreateIncidentWithContext instead.
func (c *Client) CreateIncident(from string, o *CreateIncidentOptions) (*Incident, error) {
	return c.CreateIncidentWithContext(context.Background(), from, o)
}

// CreateIncidentWithContext creates an incident synchronously without a
// corresponding event from a monitoring service.
func (c *Client) CreateIncidentWithContext(ctx context.Context, from string, o *CreateIncidentOptions) (*Incident, error) {
	h := map[string]string{
		"From": from,
	}

	d := map[string]*CreateIncidentOptions{
		"incident": o,
	}

	resp, err := c.post(ctx, "/incidents", d, h)
	if err != nil {
		return nil, err
	}

	var ii createIncidentResponse
	if err = c.decodeJSON(resp, &ii); err != nil {
		return nil, err
	}

	return &ii.Incident, nil
}

// ManageIncidents acknowledges, resolves, escalates, or reassigns one or more
// incidents. It's recommended to use ManageIncidentsWithContext instead.
func (c *Client) ManageIncidents(from string, incidents []ManageIncidentsOptions) (*ListIncidentsResponse, error) {
	return c.ManageIncidentsWithContext(context.Background(), from, incidents)
}

// ManageIncidentsWithContext acknowledges, resolves, escalates, or reassigns
// one or more incidents.
func (c *Client) ManageIncidentsWithContext(ctx context.Context, from string, incidents []ManageIncidentsOptions) (*ListIncidentsResponse, error) {
	d := map[string][]ManageIncidentsOptions{
		"incidents": incidents,
	}

	h := map[string]string{
		"From": from,
	}

	resp, err := c.put(ctx, "/incidents", d, h)
	if err != nil {
		return nil, err
	}

	var result ListIncidentsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// MergeIncidents merges a list of source incidents into a specified incident.
// It's recommended to use MergeIncidentsWithContext instead.
func (c *Client) MergeIncidents(from string, id string, sourceIncidents []MergeIncidentsOptions) (*Incident, error) {
	return c.MergeIncidentsWithContext(context.Background(), from, id, sourceIncidents)
}

// MergeIncidentsWithContext merges a list of source incidents into a specified incident.
func (c *Client) MergeIncidentsWithContext(ctx context.Context, from, id string, sourceIncidents []MergeIncidentsOptions) (*Incident, error) {
	d := map[string][]MergeIncidentsOptions{
		"source_incidents": sourceIncidents,
	}

	h := map[string]string{
		"From": from,
	}

	resp, err := c.put(ctx, "/incidents/"+id+"/merge", d, h)
	if err != nil {
		return nil, err
	}

	var result createIncidentResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Incident, nil
}

// GetIncident shows detailed information about an incident. It's recommended to
// use GetIncidentWithContext instead.
func (c *Client) GetIncident(id string) (*Incident, error) {
	return c.GetIncidentWithContext(context.Background(), id)
}

// GetIncidentWithContext shows detailed information about an incident.
func (c *Client) GetIncidentWithContext(ctx context.Context, id string) (*Incident, error) {
	resp, err := c.get(ctx, "/incidents/"+id)
	if err != nil {
		return nil, err
	}

	var result map[string]Incident
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	i, ok := result["incident"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have incident field")
	}

	return &i, nil
}

// IncidentNote is a note for the specified incident.
type IncidentNote struct {
	ID        string    `json:"id,omitempty"`
	User      APIObject `json:"user,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt string    `json:"created_at,omitempty"`
}

// CreateIncidentNoteResponse is returned from the API as a response to creating an incident note.
type CreateIncidentNoteResponse struct {
	IncidentNote IncidentNote `json:"note"`
}

// ListIncidentNotes lists existing notes for the specified incident. It's
// recommended to use ListIncidentNotesWithContext instead.
func (c *Client) ListIncidentNotes(id string) ([]IncidentNote, error) {
	return c.ListIncidentNotesWithContext(context.Background(), id)
}

// ListIncidentNotesWithContext lists existing notes for the specified incident.
func (c *Client) ListIncidentNotesWithContext(ctx context.Context, id string) ([]IncidentNote, error) {
	resp, err := c.get(ctx, "/incidents/"+id+"/notes")
	if err != nil {
		return nil, err
	}

	var result map[string][]IncidentNote
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	notes, ok := result["notes"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have notes field")
	}

	return notes, nil
}

// IncidentAlert is a alert for the specified incident.
type IncidentAlert struct {
	APIObject
	CreatedAt   string                 `json:"created_at,omitempty"`
	Status      string                 `json:"status,omitempty"`
	AlertKey    string                 `json:"alert_key,omitempty"`
	Service     APIObject              `json:"service,omitempty"`
	Body        map[string]interface{} `json:"body,omitempty"`
	Incident    APIReference           `json:"incident,omitempty"`
	Suppressed  bool                   `json:"suppressed,omitempty"`
	Severity    string                 `json:"severity,omitempty"`
	Integration APIObject              `json:"integration,omitempty"`
}

// IncidentAlertResponse is the response of a sincle incident alert
type IncidentAlertResponse struct {
	IncidentAlert *IncidentAlert `json:"alert,omitempty"`
}

// IncidentAlertList is the generic structure of a list of alerts
type IncidentAlertList struct {
	Alerts []IncidentAlert `json:"alerts,omitempty"`
}

// ListAlertsResponse is the response structure when calling the ListAlert API endpoint.
type ListAlertsResponse struct {
	APIListObject
	Alerts []IncidentAlert `json:"alerts,omitempty"`
}

// ListIncidentAlertsOptions is the structure used when passing parameters to the ListIncidentAlerts API endpoint.
type ListIncidentAlertsOptions struct {
	APIListObject
	Statuses []string `url:"statuses,omitempty,brackets"`
	SortBy   string   `url:"sort_by,omitempty"`
	Includes []string `url:"include,omitempty,brackets"`
}

// ListIncidentAlerts lists existing alerts for the specified incident. It's
// recommended to use ListIncidentAlertsWithContext instead.
func (c *Client) ListIncidentAlerts(id string) (*ListAlertsResponse, error) {
	return c.ListIncidentAlertsWithContext(context.Background(), id, ListIncidentAlertsOptions{})
}

// ListIncidentAlertsWithOpts lists existing alerts for the specified incident.
// It's recommended to use ListIncidentAlertsWithContext instead.
func (c *Client) ListIncidentAlertsWithOpts(id string, o ListIncidentAlertsOptions) (*ListAlertsResponse, error) {
	return c.ListIncidentAlertsWithContext(context.Background(), id, o)
}

// ListIncidentAlertsWithContext lists existing alerts for the specified
// incident. If you don't want to filter any of the results, pass in an empty
// ListIncidentAlertOptions.
func (c *Client) ListIncidentAlertsWithContext(ctx context.Context, id string, o ListIncidentAlertsOptions) (*ListAlertsResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents/"+id+"/alerts?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListAlertsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, err
}

// CreateIncidentNoteWithResponse creates a new note for the specified incident.
// It's recommended to use CreateIncidentNoteWithContext instead.
func (c *Client) CreateIncidentNoteWithResponse(id string, note IncidentNote) (*IncidentNote, error) {
	return c.CreateIncidentNoteWithContext(context.Background(), id, note)
}

// CreateIncidentNoteWithContext creates a new note for the specified incident.
func (c *Client) CreateIncidentNoteWithContext(ctx context.Context, id string, note IncidentNote) (*IncidentNote, error) {
	d := map[string]IncidentNote{
		"note": note,
	}

	h := map[string]string{
		"From": note.User.Summary,
	}

	resp, err := c.post(ctx, "/incidents/"+id+"/notes", d, h)
	if err != nil {
		return nil, err
	}

	var result CreateIncidentNoteResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.IncidentNote, nil
}

// CreateIncidentNote creates a new note for the specified incident.
//
// Deprecated: please use CreateIncidentNoteWithContext going forward
func (c *Client) CreateIncidentNote(id string, note IncidentNote) error {
	data := make(map[string]IncidentNote)
	headers := make(map[string]string)
	headers["From"] = note.User.Summary
	data["note"] = note
	_, err := c.post(context.Background(), "/incidents/"+id+"/notes", data, headers)
	return err
}

// SnoozeIncidentWithResponse sets an incident to not alert for a specified
// period of time. It's recommended to use SnoozeIncidentWithContext instead.
func (c *Client) SnoozeIncidentWithResponse(id string, duration uint) (*Incident, error) {
	return c.SnoozeIncidentWithContext(context.Background(), id, duration)
}

// SnoozeIncidentWithContext sets an incident to not alert for a specified period of time.
func (c *Client) SnoozeIncidentWithContext(ctx context.Context, id string, duration uint) (*Incident, error) {
	d := map[string]uint{
		"duration": duration,
	}

	resp, err := c.post(ctx, "/incidents/"+id+"/snooze", d, nil)
	if err != nil {
		return nil, err
	}

	var result createIncidentResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Incident, nil
}

// SnoozeIncident sets an incident to not alert for a specified period of time.
//
// Deprecated: please use SnoozeIncidentWithContext going forward
func (c *Client) SnoozeIncident(id string, duration uint) error {
	data := make(map[string]uint)
	data["duration"] = duration
	_, err := c.post(context.Background(), "/incidents/"+id+"/snooze", data, nil)
	return err
}

// ListIncidentLogEntriesResponse is the response structure when calling the ListIncidentLogEntries API endpoint.
type ListIncidentLogEntriesResponse struct {
	APIListObject
	LogEntries []LogEntry `json:"log_entries,omitempty"`
}

// ListIncidentLogEntriesOptions is the structure used when passing parameters to the ListIncidentLogEntries API endpoint.
type ListIncidentLogEntriesOptions struct {
	APIListObject
	Includes   []string `url:"include,omitempty,brackets"`
	IsOverview bool     `url:"is_overview,omitempty"`
	TimeZone   string   `url:"time_zone,omitempty"`
	Since      string   `url:"since,omitempty"`
	Until      string   `url:"until,omitempty"`
}

// ListIncidentLogEntries lists existing log entries for the specified incident.
// It's recommended to use ListIncidentLogEntriesWithContext instead.
func (c *Client) ListIncidentLogEntries(id string, o ListIncidentLogEntriesOptions) (*ListIncidentLogEntriesResponse, error) {
	return c.ListIncidentLogEntriesWithContext(context.Background(), id, o)
}

// ListIncidentLogEntriesWithContext lists existing log entries for the
// specified incident.
func (c *Client) ListIncidentLogEntriesWithContext(ctx context.Context, id string, o ListIncidentLogEntriesOptions) (*ListIncidentLogEntriesResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/incidents/"+id+"/log_entries?"+v.Encode())
	if err != nil {
		return nil, err
	}

	var result ListIncidentLogEntriesResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// IncidentResponders contains details about responders to an incident.
type IncidentResponders struct {
	State       string    `json:"state"`
	User        APIObject `json:"user"`
	Incident    APIObject `json:"incident"`
	UpdatedAt   string    `json:"updated_at"`
	Message     string    `json:"message"`
	Requester   APIObject `json:"requester"`
	RequestedAt string    `json:"requested_at"`
}

// ResponderRequestResponse is the response from the API when requesting someone
// respond to an incident.
type ResponderRequestResponse struct {
	ResponderRequest ResponderRequest `json:"responder_request"`
}

// ResponderRequestTarget specifies an individual target for the responder request.
type ResponderRequestTarget struct {
	APIObject
	Responders IncidentResponders `json:"incident_responders"`
}

// ResponderRequestTargets is a wrapper for a ResponderRequestTarget.
type ResponderRequestTargets struct {
	Target ResponderRequestTarget `json:"responder_request_target"`
}

// ResponderRequestOptions defines the input options for the Create Responder function.
type ResponderRequestOptions struct {
	From        string                   `json:"-"`
	Message     string                   `json:"message"`
	RequesterID string                   `json:"requester_id"`
	Targets     []ResponderRequestTarget `json:"responder_request_targets"`
}

// ResponderRequest contains the API structure for an incident responder request.
type ResponderRequest struct {
	Incident    Incident                `json:"incident"`
	Requester   User                    `json:"requester,omitempty"`
	RequestedAt string                  `json:"request_at,omitempty"`
	Message     string                  `json:"message,omitempty"`
	Targets     ResponderRequestTargets `json:"responder_request_targets"`
}

// ResponderRequest will submit a request to have a responder join an incident.
// It's recommended to use ResponderRequestWithContext instead.
func (c *Client) ResponderRequest(id string, o ResponderRequestOptions) (*ResponderRequestResponse, error) {
	return c.ResponderRequestWithContext(context.Background(), id, o)
}

// ResponderRequestWithContext will submit a request to have a responder join an incident.
func (c *Client) ResponderRequestWithContext(ctx context.Context, id string, o ResponderRequestOptions) (*ResponderRequestResponse, error) {
	h := map[string]string{
		"From": o.From,
	}

	resp, err := c.post(ctx, "/incidents/"+id+"/responder_requests", o, h)
	if err != nil {
		return nil, err
	}

	var result ResponderRequestResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIncidentAlert gets the alert that triggered the incident. It's recommended
// to use GetIncidentAlertWithContext instead.
func (c *Client) GetIncidentAlert(incidentID, alertID string) (*IncidentAlertResponse, *http.Response, error) {
	return c.getIncidentAlertWithContext(context.Background(), incidentID, alertID)
}

// GetIncidentAlertWithContext gets the alert that triggered the incident.
func (c *Client) GetIncidentAlertWithContext(ctx context.Context, incidentID, alertID string) (*IncidentAlertResponse, error) {
	iar, _, err := c.getIncidentAlertWithContext(context.Background(), incidentID, alertID)
	return iar, err
}

func (c *Client) getIncidentAlertWithContext(ctx context.Context, incidentID, alertID string) (*IncidentAlertResponse, *http.Response, error) {
	resp, err := c.get(ctx, "/incidents/"+incidentID+"/alerts/"+alertID)
	if err != nil {
		return nil, nil, err
	}

	var result IncidentAlertResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, nil, err
	}

	return &result, resp, nil
}

// ManageIncidentAlerts allows you to manage the alerts of an incident. It's
// recommended to use ManageIncidentAlertsWithContext instead.
func (c *Client) ManageIncidentAlerts(incidentID string, alerts *IncidentAlertList) (*ListAlertsResponse, *http.Response, error) {
	return c.manageIncidentAlertsWithContext(context.Background(), incidentID, alerts)
}

// ManageIncidentAlertsWithContext allows you to manage the alerts of an incident.
func (c *Client) ManageIncidentAlertsWithContext(ctx context.Context, incidentID string, alerts *IncidentAlertList) (*ListAlertsResponse, error) {
	lar, _, err := c.manageIncidentAlertsWithContext(context.Background(), incidentID, alerts)
	return lar, err
}

func (c *Client) manageIncidentAlertsWithContext(ctx context.Context, incidentID string, alerts *IncidentAlertList) (*ListAlertsResponse, *http.Response, error) {
	resp, err := c.put(ctx, "/incidents/"+incidentID+"/alerts/", alerts, nil)
	if err != nil {
		return nil, nil, err
	}

	var result ListAlertsResponse
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, nil, err
	}

	return &result, resp, nil
}

/* TODO: Create Status Updates */
