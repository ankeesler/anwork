package sql_test

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/sql"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO: get rid of panics
// TODO: corner cases for setting up db wrong

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
		db, err = sql.Open(logger, "mysql", dsn)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		cleanTestDB(db, dsn)
		db.Close()
	})

	task.RunRepoTests(func() task.Repo {
		return sql.New(logger, db)
	})

	Context("when db is in a weird state", func() {
		BeforeEach(func() {
			_, err := db.Exec("CREATE TABLE tasks (id int)")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error", func() {
			t := task.Task{Name: "task"}
			Expect(sql.New(logger, db).CreateTask(&t)).NotTo(Succeed())
		})
	})
})

func cleanTestDB(db *sql.DB, dsn string) {
	rows, err := db.Query(`SHOW TABLES`)
	Expect(err).NotTo(HaveOccurred())
	defer rows.Close()

	var name string
	for rows.Next() {
		Expect(rows.Scan(&name)).To(Succeed())

		_, err := db.Exec(fmt.Sprintf("DROP TABLE %s", name))
		Expect(err).NotTo(HaveOccurred())
	}
}
