package sql_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/sql"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO: improve performance by combining setup queries?

var _ = Describe("SQL Repo", func() {
	var (
		logger lager.Logger

		dsn string
		db  *sql.DB
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("sql")

		var ok bool
		dsn, ok = os.LookupEnv("ANWORK_TEST_SQL_DSN")
		Expect(ok).To(BeTrue(), "ANWORK_TEST_SQL_DSN env var must be set!")

		var err error
		db, err = sql.Open("mysql", dsn)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		logger = logger.Session("after-each")
		cleanTestDB(logger, db, dsn)
		db.Close(logger)
	})

	task.RunRepoTests(func() task.Repo {
		return sql.New(logger, db)
	})

	Context("when db is in a weird state", func() {
		BeforeEach(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			_, err := db.Exec(
				ctx,
				logger.Session("before-each"),
				"CREATE TABLE tasks (id int)",
			)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error", func() {
			t := task.Task{Name: "task"}
			Expect(sql.New(logger, db).CreateTask(&t)).NotTo(Succeed())
		})
	})

	Context("when we try to delete a task before all its events", func() {
		var repo task.Repo

		BeforeEach(func() {
			repo = sql.New(logger, db)

			t := task.Task{Name: "task"}
			Expect(repo.CreateTask(&t)).To(Succeed())
			for i := 0; i < 3; i++ {
				e := task.Event{Title: fmt.Sprintf("event-%d", i), TaskID: t.ID}
				Expect(repo.CreateEvent(&e)).To(Succeed())
			}

			Expect(repo.Tasks()).To(HaveLen(1))
			Expect(repo.Events()).To(HaveLen(3))

			Expect(repo.DeleteTask(&t)).To(Succeed())
		})

		It("deletes all of the events too", func() {
			Expect(repo.Events()).To(HaveLen(0))
		})
	})

	Context("benchmarking", func() {
		Measure("CRUD'ing 10 tasks with one repo", func(b Benchmarker) {
			repo := sql.New(logger, db)
			runtime := b.Time("runtime", func() {
				for i := 0; i < 10; i++ {
					name := fmt.Sprintf("task-%d", i)

					// C
					t := task.Task{Name: name}
					Expect(repo.CreateTask(&t)).To(Succeed())

					// R
					_, err := repo.Tasks()
					Expect(err).NotTo(HaveOccurred())

					// U
					t.Priority = 99
					Expect(repo.UpdateTask(&t)).To(Succeed())

					// D
					Expect(repo.DeleteTask(&t)).To(Succeed())
				}
			})
			Expect(runtime.Seconds()).To(BeNumerically("<", .5))
		}, 5)
	})
})

func cleanTestDB(logger lager.Logger, db *sql.DB, dsn string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := db.Query(ctx, logger, `SHOW TABLES`)
	Expect(err).NotTo(HaveOccurred())
	defer rows.Close()

	var name string
	for rows.Next() {

		Expect(rows.Scan(&name)).To(Succeed())

		_, err := db.Exec(ctx, logger, fmt.Sprintf("DROP TABLE %s", name))
		Expect(err).NotTo(HaveOccurred())
	}

	Expect(rows.Err()).NotTo(HaveOccurred())
}
