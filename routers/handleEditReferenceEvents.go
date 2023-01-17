package routers

import (
	"encoding/json"
	"net/http"

	"github.com/aidenappl/nu-calendar/db"
	"github.com/aidenappl/nu-calendar/errors"
	"github.com/aidenappl/nu-calendar/query"
)

type HandleEditReferenceEventsReq struct {
	EventReference   *int    `json:"event_reference_id"`
	LocationAddress  *string `json:"location_address"`
	LocationBuilding *string `json:"location_building"`
	LocationRoom     *string `json:"location_room"`
}

func HandleEditReferenceEvents(w http.ResponseWriter, r *http.Request) {
	var body HandleEditReferenceEventsReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		errors.SendError(w, "error decoding json body: "+err.Error(), "", http.StatusBadRequest)
		return
	}

	if body.EventReference == nil {
		errors.SendError(w, "missing or empty field: event_reference_id", "", http.StatusBadRequest)
		return
	}

	// Get the event reference
	refEvent, err := query.GetEventReference(db.DB, query.GetEventReferenceRequest{
		EventReferenceID: *body.EventReference,
	})
	if err != nil {
		errors.SendError(w, "error getting event reference: "+err.Error(), "", http.StatusInternalServerError)
		return
	}

	if refEvent == nil {
		errors.SendError(w, "event reference not found", "", http.StatusBadRequest)
		return
	}

	// Update the event reference
	err = query.UpdateEventReference(db.DB, query.UpdateEventReferenceRequest{
		EventReferenceID: *body.EventReference,
		LocationAddress:  body.LocationAddress,
		LocationBuilding: body.LocationBuilding,
		LocationRoom:     body.LocationRoom,
	})
	if err != nil {
		errors.SendError(w, "error updating event reference: "+err.Error(), "", http.StatusInternalServerError)
		return
	}

	// Respond
	w.WriteHeader(http.StatusOK)
}
