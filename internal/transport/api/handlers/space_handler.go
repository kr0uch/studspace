package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
)

type SpaceService interface {
	ListActiveSpaces(ctx context.Context) ([]entities.Space, error)
	GetSpaceAvailability(ctx context.Context, in dto.GetSpaceAvailability) (*dto.SpaceAvailability, error)
}

type SpaceController struct {
	spaceService SpaceService
}

func NewSpaceController(spaceService SpaceService) *SpaceController {
	return &SpaceController{spaceService: spaceService}
}

func (c *SpaceController) ListSpacesHandler(w http.ResponseWriter, r *http.Request) error {
	spaces, err := c.spaceService.ListActiveSpaces(r.Context())
	if err != nil {
		return err
	}

	out := make([]api.Space, 0, len(spaces))
	for _, space := range spaces {
		out = append(out, api.Space{Space: space})
	}
	return json.NewEncoder(w).Encode(out)
}

func (c *SpaceController) GetAvailabilityHandler(w http.ResponseWriter, r *http.Request) error {
	spaceID, err := parseSpaceID(r.PathValue("id"))
	if err != nil {
		return err
	}

	date, err := parseDateField(r.URL.Query().Get("date"), "date")
	if err != nil {
		return err
	}

	availability, err := c.spaceService.GetSpaceAvailability(r.Context(), dto.GetSpaceAvailability{
		SpaceID: spaceID,
		Date:    date,
	})
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(api.NewSpaceAvailability(*availability))
}

func parseSpaceID(raw string) (int16, error) {
	id, err := strconv.ParseInt(raw, 10, 16)
	if err != nil || id < 1 {
		return 0, errors.Validation(map[string]string{"id": "некорректный идентификатор"})
	}
	return int16(id), nil
}

func spaceName(spaces []entities.Space, spaceID int16) string {
	for _, space := range spaces {
		if space.SpaceID == spaceID {
			return space.Name
		}
	}
	return ""
}
