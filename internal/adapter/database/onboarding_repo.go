package database

import (
	"context"
	"time"
	"errors"
	
	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/core/erro"

	go_core_observ "github.com/eliezerraj/go-core/observability"
	go_core_pg "github.com/eliezerraj/go-core/database/pg"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var tracerProvider go_core_observ.TracerProvider
var childLogger = log.With().Str("adapter", "database").Logger()

type WorkerRepository struct {
	DatabasePGServer *go_core_pg.DatabasePGServer
}

func NewWorkerRepository(databasePGServer *go_core_pg.DatabasePGServer) *WorkerRepository{
	childLogger.Debug().Msg("NewWorkerRepository")

	return &WorkerRepository{
		DatabasePGServer: databasePGServer,
	}
}

func (w WorkerRepository) AddPerson(ctx context.Context, tx pgx.Tx, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("AddPerson")

	span := tracerProvider.Span(ctx, "database.AddPerson")
	defer span.End()

	query := `INSERT INTO person (	person_id, 
									name,
									created_at,
									tenant_id) 
									VALUES($1, $2, $3, $4) RETURNING id`

	onboarding.Person.CreatedAt = time.Now()

	row := tx.QueryRow(ctx, query,  onboarding.Person.PersonID,  
									onboarding.Person.Name,
									onboarding.Person.CreatedAt,
									onboarding.Person.TenantID)

	var id int
	
	if err := row.Scan(&id); err != nil {
		return nil, errors.New(err.Error())
	}

	onboarding.Person.ID = id

	return onboarding, nil
}

func (w WorkerRepository) GetPerson(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("GetPerson")
	
	span := tracerProvider.Span(ctx, "database.GetPerson")
	defer span.End()

	conn, err := w.DatabasePGServer.Acquire(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer w.DatabasePGServer.Release(conn)

	res_person := model.Person{}
	res_onboarding := model.Onboarding{Person: &res_person}

	query := `SELECT id,
					person_id,	 
					name,
					created_at,
					updated_at 
				FROM public.person 
				WHERE person_id =$1`

	rows, err := conn.Query(ctx, query, onboarding.Person.PersonID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan( 	&res_person.ID,
							&res_person.PersonID, 
							&res_person.Name, 
							&res_person.CreatedAt,
							&res_person.UpdatedAt,
						)
		if err != nil {
			return nil, errors.New(err.Error())
        }
		return &res_onboarding, nil
	}
	
	return nil, erro.ErrNotFound
}

func (w WorkerRepository) UpdatePerson(ctx context.Context, tx pgx.Tx, onboarding *model.Onboarding) (int64, error){
	childLogger.Debug().Msg("UpdatePerson")

	span := tracerProvider.Span(ctx, "database.UpdatePerson")
	defer span.End()

	t_updateAt := time.Now()
	onboarding.Person.UpdatedAt = &t_updateAt

	query := `Update public.person
				set name = $2, 
					updated_at = $3
				where person_id = $1`

	row, err := tx.Exec(ctx, query, onboarding.Person.PersonID,  
									onboarding.Person.Name,
									onboarding.Person.UpdatedAt)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	if int(row.RowsAffected()) == 0 {
		return 0, erro.ErrUpdateRows
	}
	childLogger.Debug().Int("rowsAffected : ",int(row.RowsAffected())).Msg("")
	
	return row.RowsAffected(), nil
}

func (w WorkerRepository) ListPerson(ctx context.Context, onboarding *model.Onboarding) (*[]model.Onboarding, error){
	childLogger.Debug().Msg("ListPerson")
	
	span := tracerProvider.Span(ctx, "database.ListPerson")
	defer span.End()

	conn, err := w.DatabasePGServer.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("error acquire")
		return nil, errors.New(err.Error())
	}
	defer w.DatabasePGServer.Release(conn)

	res_onboarding_list := []model.Onboarding{}

	query := `SELECT id,
					person_id, 
					name,
					created_at,
					updated_at 
				FROM public.person
				WHERE person_id >= $1 
				ORDER BY person_id asc`

	rows, err := conn.Query(ctx, query, onboarding.Person.PersonID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		res_person := model.Person{}
		res_onboarding := model.Onboarding{Person: &res_person}

		err := rows.Scan( 	&res_person.ID,
							&res_person.PersonID, 
							&res_person.Name, 
							&res_person.CreatedAt,
							&res_person.UpdatedAt,
						)
		if err != nil {
			return nil, errors.New(err.Error())
        }
		res_onboarding_list = append(res_onboarding_list, res_onboarding)
	}
	
	return &res_onboarding_list, nil
}