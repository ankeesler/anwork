// Package sql contains a task.Repo implementation which uses a SQL database
// to store data.
package sql

import (
	stdlibsql "database/sql"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/ankeesler/anwork/task"
)

type repo struct {
	logger lager.Logger

	db *DB
}

// New returns a task.Repo that stores task.Task's in an SQL database.
func New(logger lager.Logger, db *DB) task.Repo {
	return &repo{logger: logger, db: db}
}

func (r *repo) CreateTask(task *task.Task) error {
	r.logger.Debug("create-task-begin", lager.Data{"task": task})
	defer r.logger.Debug("create-task-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return err
	}

	stmt, err := r.db.Prepare(`INSERT INTO tasks (name, start_date, priority, state) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(task.Name, task.StartDate, task.Priority, task.State)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = int(id)

	return nil
}

func (r *repo) Tasks() ([]*task.Task, error) {
	r.logger.Debug("tasks-begin")
	defer r.logger.Debug("tasks-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return nil, err
	}

	rows, err := r.db.Query("SELECT * FROM tasks")
	if err != nil {
		r.logger.Error("get-tasks", err)
		return nil, err
	}

	tasks := make([]*task.Task, 0)
	for rows.Next() {
		if rows.Err() != nil {
			r.logger.Error("rows-next", err)
			return nil, err
		}

		task := new(task.Task)
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.StartDate,
			&task.Priority,
			&task.State,
		); err != nil {
			r.logger.Error("rows-scan", err)
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *repo) FindTaskByID(id int) (*task.Task, error) {
	r.logger.Debug("find-task-by-id-begin", lager.Data{"id": id})
	defer r.logger.Debug("find-task-by-id-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return nil, err
	}

	q := fmt.Sprintf(`SELECT * FROM tasks WHERE id = %d`, id)
	row := r.db.QueryRow(q)

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
			return nil, err
		}
	}

	return task, nil
}

func (r *repo) FindTaskByName(name string) (*task.Task, error) {
	r.logger.Debug("find-task-by-name-begin", lager.Data{"name": name})
	defer r.logger.Debug("find-task-by-name-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return nil, err
	}

	q := fmt.Sprintf(`SELECT * FROM tasks WHERE name = '%s'`, name)
	row := r.db.QueryRow(q)

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
			return nil, err
		}
	}

	return task, nil
}

func (r *repo) UpdateTask(task *task.Task) error {
	r.logger.Debug("update-task-begin", lager.Data{"task": task})
	defer r.logger.Debug("update-task-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return err
	}

	q := fmt.Sprintf(`
UPDATE tasks
SET name = '%s', start_date = %d, priority = %d, state = '%s'
WHERE id = %d`,
		task.Name, task.StartDate, task.Priority, task.State, task.ID)
	result, err := r.db.Exec(q)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return fmt.Errorf("unknown task with id %d", task.ID)
	}

	return nil
}

func (r *repo) DeleteTask(task *task.Task) error {
	r.logger.Debug("delete-task-begin", lager.Data{"task": task})
	defer r.logger.Debug("delete-task-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return err
	}

	q := fmt.Sprintf(`DELETE FROM tasks WHERE id = %d`, task.ID)
	result, err := r.db.Exec(q)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return fmt.Errorf("unknown task with id %d", task.ID)
	}

	return nil
}

func (r *repo) CreateEvent(event *task.Event) error {
	r.logger.Debug("create-event-begin", lager.Data{"event": event})
	defer r.logger.Debug("create-event-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return err
	}

	stmt, err := r.db.Prepare(`INSERT INTO events (title, date, type, task_id) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(event.Title, event.Date, event.Type, event.TaskID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = int(id)

	return nil
}

func (r *repo) Events() ([]*task.Event, error) {
	r.logger.Debug("events-begin")
	defer r.logger.Debug("events-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return nil, err
	}

	rows, err := r.db.Query("SELECT * FROM events")
	if err != nil {
		r.logger.Error("get-tasks", err)
		return nil, err
	}

	events := make([]*task.Event, 0)
	for rows.Next() {
		if rows.Err() != nil {
			r.logger.Error("rows-next", err)
			return nil, err
		}

		event := new(task.Event)
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.Type,
			&event.TaskID,
		); err != nil {
			r.logger.Error("rows-scan", err)
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *repo) FindEventByID(id int) (*task.Event, error) {
	r.logger.Debug("find-event-by-id-begin", lager.Data{"id": id})
	defer r.logger.Debug("find-event-by-id-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return nil, err
	}

	q := fmt.Sprintf(`SELECT * FROM events WHERE id = %d`, id)
	row := r.db.QueryRow(q)

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
			return nil, err
		}
	}

	return event, nil
}

func (r *repo) DeleteEvent(event *task.Event) error {
	r.logger.Debug("delete-event-begin", lager.Data{"event": event})
	defer r.logger.Debug("delete-event-end")

	if err := r.ensureTablesExist(); err != nil {
		r.logger.Error("ensure-tables", err)
		return err
	}

	q := fmt.Sprintf(`DELETE FROM events WHERE id = %d`, event.ID)
	result, err := r.db.Exec(q)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return fmt.Errorf("unknown event with id %d", event.ID)
	}

	return nil
}

func (r *repo) ensureTablesExist() error {
	if exists, err := r.tablesExist(); err != nil {
		return err
	} else if !exists {
		q := `
CREATE TABLE tasks (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  start_date bigint NOT NULL,
  priority int NOT NULL,
  state varchar(16) NOT NULL
)
`
		_, err = r.db.Query(q)
		if err != nil {
			return err
		}

		// TODO: can this be combined with the above?
		q = `
CREATE TABLE events (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title varchar(255) NOT NULL,
  date bigint NOT NULL,
  type int NOT NULL,
  task_id int NOT NULL,
  FOREIGN KEY (task_id) REFERENCES tasks(id)
)
`
		_, err = r.db.Query(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repo) tablesExist() (bool, error) {
	q := `SHOW TABLES`
	rows, err := r.db.Query(q)

	if err != nil {
		return false, err
	} else if !rows.Next() {
		return false, nil
	} else {
		return true, nil
	}
}
