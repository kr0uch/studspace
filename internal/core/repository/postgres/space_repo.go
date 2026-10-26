package postgres

import (
	"context"

	"github.com/kr0uch/studspace/internal/core/models/dao"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	dbpg "github.com/kr0uch/studspace/pkg/database/postgres"
	"github.com/kr0uch/studspace/pkg/errors"
)

type SpaceRepository struct {
	postgresDB *dbpg.PostgresDB
}

func NewSpaceRepository(postgresDB *dbpg.PostgresDB) *SpaceRepository {
	return &SpaceRepository{postgresDB: postgresDB}
}

type spaceAmenityRow struct {
	SpaceID   int16  `db:"space_id"`
	AmenityID int16  `db:"amenity_id"`
	Name      string `db:"name"`
	Icon      string `db:"icon"`
}

func (r *SpaceRepository) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	var rows []dao.Space
	if err := r.postgresDB.SelectContext(ctx, &rows, qSpaceListActive); err != nil {
		return nil, wrapErr(err)
	}

	var amenityRows []spaceAmenityRow
	if err := r.postgresDB.SelectContext(ctx, &amenityRows, qSpaceListAmenities); err != nil {
		return nil, wrapErr(err)
	}

	amenitiesBySpace := make(map[int16][]dao.Amenity, len(rows))
	for _, row := range amenityRows {
		amenitiesBySpace[row.SpaceID] = append(amenitiesBySpace[row.SpaceID], dao.Amenity{
			AmenityID: row.AmenityID,
			Name:      row.Name,
			Icon:      row.Icon,
		})
	}

	result := make([]entities.Space, 0, len(rows))
	for _, row := range rows {
		space, err := row.ToEntity(amenitiesBySpace[row.SpaceID])
		if err != nil {
			return nil, errors.ErrInternal.Wrap(err)
		}
		result = append(result, space)
	}
	return result, nil
}

func (r *SpaceRepository) GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error) {
	var row dao.Space
	if err := r.postgresDB.GetContext(ctx, &row, qSpaceGetByID, spaceID); err != nil {
		return nil, wrapErr(err)
	}

	var amenityRows []spaceAmenityRow
	if err := r.postgresDB.SelectContext(ctx, &amenityRows, qSpaceListAmenitiesBySpace, spaceID); err != nil {
		return nil, wrapErr(err)
	}

	amenities := make([]dao.Amenity, 0, len(amenityRows))
	for _, amenity := range amenityRows {
		amenities = append(amenities, dao.Amenity{
			AmenityID: amenity.AmenityID,
			Name:      amenity.Name,
			Icon:      amenity.Icon,
		})
	}

	space, err := row.ToEntity(amenities)
	if err != nil {
		return nil, errors.ErrInternal.Wrap(err)
	}
	return &space, nil
}
