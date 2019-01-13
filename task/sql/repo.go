// Package sql contains a task.Repo implementation which uses a SQL database
// to store data.
package sql

import (
	"context"
	stdlibsql "database/sql"
	"fmt"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/ankeesler/anwork/task"
)

type repo struct {
	logger lager.Logger

	db *DB

	tablesCreated bool
}

// New returns a task.Repo that stores task.Task's in an SQL database.
func New(logger lager.Logger, db *DB) task.Repo {
	return &repo{logger: logger, db: db}
}

func (r *repo) CreateTask(task *task.Task) error {
	logger := r.logger.Session("create-task")
	logger.Debug("begin", lager.Data{"task": task})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	q := `INSERT INTO tasks (name, start_date, priority, state) VALUES (?, ?, ?, ?)`
	stmt, err := r.db.Prepare(ctx, logger, q)
	if err != nil {
		logger.Error("prepare", err)
		return err
	}
	defer stmt.Close(logger)

	result, err := stmt.Exec(
		ctx,
		logger,
		task.Name,
		task.StartDate,
		task.Priority,
		task.State,
	)
	if err != nil {
		logger.Error("exec", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("result-last-insert-id", err)
		return err
	}
	task.ID = int(id)

	return nil
}

func (r *repo) Tasks() ([]*task.Task, error) {
	logger := r.logger.Session("tasks")
	logger.Debug("begin")
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return nil, err
	}

	ctx, cancel := makeCtx()
	defer cancel()
	rows, err := r.db.Query(ctx, logger, "SELECT * FROM tasks")
	if err != nil {
		logger.Error("get-tasks", err)
		return nil, err
	}

	tasks := make([]*task.Task, 0)
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				logger.Error("rows-next", err)
				return nil, err
			} else {
				break
			}
		}

		task := new(task.Task)
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.StartDate,
			&task.Priority,
			&task.State,
		); err != nil {
			logger.Error("rows-scan", err)
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *repo) FindTaskByID(id int) (*task.Task, error) {
	logger := r.logger.Session("find-by-id")
	logger.Debug("begin", lager.Data{"id": id})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return nil, err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	q := fmt.Sprintf(`SELECT * FROM tasks WHERE id = %d`, id)
	row := r.db.QueryRow(ctx, logger, q)

	task := new(task.Task)
	if err := row.Scan(
		&task.ID,
		&task.Name,
		&task.StartDate,
		&task.Priority,
		&task.State,
	); err != nil {
		if err == stdlibsql.ErrNoRows {
			return nil, nil
		} else {
			logger.Error("scan-row", err)
			return nil, err
		}
	}

	return task, nil
}

func (r *repo) FindTaskByName(name string) (*task.Task, error) {
	logger := r.logger.Session("find-task-by-name")
	logger.Debug("begin", lager.Data{"name": name})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return nil, err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	q := fmt.Sprintf(`SELECT * FROM tasks WHERE name = '%s'`, name)
	row := r.db.QueryRow(ctx, logger, q)

	task := new(task.Task)
	if err := row.Scan(
		&task.ID,
		&task.Name,
		&task.StartDate,
		&task.Priority,
		&task.State,
	); err != nil {
		if err == stdlibsql.ErrNoRows {
			return nil, nil
		} else {
			logger.Error("scan-row", err)
			return nil, err
		}
	}

	return task, nil
}

func (r *repo) UpdateTask(task *task.Task) error {
	logger := r.logger.Session("update-task")
	logger.Debug("begin", lager.Data{"task": task})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	if exists, err := r.taskExists(ctx, logger, task.ID); err != nil {
		logger.Error("task-exists", err)
		return err
	} else if !exists {
		return fmt.Errorf("unknown task with id %d", task.ID)
	}

	q := fmt.Sprintf(`
UPDATE tasks
SET name = '%s', start_date = %d, priority = %d, state = '%s'
WHERE id = %d`,
		task.Name, task.StartDate, task.Priority, task.State, task.ID)
	_, err := r.db.Exec(ctx, logger, q)
	if err != nil {
		logger.Error("exec", err)
		return err
	}

	return nil
}

func (r *repo) DeleteTask(task *task.Task) error {
	logger := r.logger.Session("delete-task")
	logger.Debug("begin", lager.Data{"task": task})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	q = fmt.Sprintf(`DELETE FROM tasks WHERE id = %d`, task.ID)
	_, err = r.db.Exec(ctx, logger, q)
	if err != nil {
		logger.Error("exec", err)
		return err
	}

	return nil
}

func (r *repo) CreateEvent(event *task.Event) error {
	logger := r.logger.Session("create-event")
	logger.Debug("begin", lager.Data{"event": event})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	q := `INSERT INTO events (title, date, type, task_id) VALUES (?, ?, ?, ?)`
	stmt, err := r.db.Prepare(ctx, logger, q)
	if err != nil {
		logger.Error("prepare", err)
		return err
	}
	defer stmt.Close(logger)

	result, err := stmt.Exec(
		ctx,
		logger,
		event.Title,
		event.Date,
		event.Type,
		event.TaskID,
	)
	if err != nil {
		logger.Error("exec", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("last-insert-id", err)
		return err
	}
	event.ID = int(id)

	return nil
}

func (r *repo) Events() ([]*task.Event, error) {
	logger := r.logger.Session("events")
	logger.Debug("begin")
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return nil, err
	}

	ctx, cancel := makeCtx()
	defer cancel()
	rows, err := r.db.Query(ctx, logger, "SELECT * FROM events")
	if err != nil {
		logger.Error("get-events", err)
		return nil, err
	}

	events := make([]*task.Event, 0)
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				logger.Error("rows-next", err)
				return nil, err
			} else {
				break
			}
		}

		event := new(task.Event)
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.Type,
			&event.TaskID,
		); err != nil {
			logger.Error("rows-scan", err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *repo) FindEventByID(id int) (*task.Event, error) {
	logger := r.logger.Session("find-event-by-id")
	logger.Debug("begin", lager.Data{"id": id})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return nil, err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	q := fmt.Sprintf(`SELECT * FROM events WHERE id = %d`, id)
	row := r.db.QueryRow(ctx, logger, q)

	event := new(task.Event)
	if err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Date,
		&event.Type,
		&event.TaskID,
	); err != nil {
		if err == stdlibsql.ErrNoRows {
			return nil, nil
		} else {
			logger.Error("scan", err)
			return nil, err
		}
	}

	return event, nil
}

func (r *repo) DeleteEvent(event *task.Event) error {
	logger := r.logger.Session("delete-event")
	logger.Debug("begin", lager.Data{"event": event})
	defer logger.Debug("end")

	if err := r.ensureTablesExist(logger); err != nil {
		logger.Error("ensure-tables", err)
		return err
	}

	ctx, cancel := makeCtx()
	defer cancel()

	q := fmt.Sprintf(`DELETE FROM events WHERE id = %d`, event.ID)
	_, err := r.db.Exec(ctx, logger, q)
	if err != nil {
		logger.Error("exec", err)
		return err
	}

	return nil
}

func (r *repo) ensureTablesExist(logger lager.Logger) error {
	if r.tablesCreated {
		return nil
	}

	if exists, err := r.tablesExist(logger); err != nil {
		return err
	} else if !exists {
		ctx, cancel := makeCtx()
		defer cancel()

		q := `
CREATE TABLE tasks (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  start_date bigint NOT NULL,
  priority int NOT NULL,
  state varchar(16) NOT NULL
)
`
		_, err = r.db.Exec(ctx, logger, q)
		if err != nil {
			r.logger.Error("create-tasks-table", err)
			return err
		}

		q = `
CREATE TABLE events (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title varchar(255) NOT NULL,
  date bigint NOT NULL,
  type int NOT NULL,
  task_id int NOT NULL
)
`
		_, err = r.db.Exec(ctx, logger, q)
		if err != nil {
			r.logger.Error("create-events-table", err)
			return err
		}
	}

	r.tablesCreated = true

	return nil
}

func (r *repo) tablesExist(logger lager.Logger) (bool, error) {
	ctx, cancel := makeCtx()
	defer cancel()

	q := `SHOW TABLES`
	rows, err := r.db.Query(ctx, logger, q)

	if err != nil {
		return false, err
	} else if !rows.Next() {
		return false, nil
	} else {
		return true, nil
	}
}

func (r *repo) taskExists(
	ctx context.Context,
	logger lager.Logger,
	id int,
) (bool, error) {
	q := fmt.Sprintf("SELECT id FROM tasks WHERE id = %d", id)
	if err := r.db.QueryRow(ctx, logger, q).Scan(&id); err == stdlibsql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func makeCtx() (context.Context, func()) {
	return context.WithTimeout(context.Background(), time.Second*3)
}
