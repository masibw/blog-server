package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/crypto/ssh/terminal"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"

	"github.com/masibw/blog-server/usecase"

	"github.com/golang-migrate/migrate/v4"
	"github.com/masibw/blog-server/config"
	"github.com/masibw/blog-server/database"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/masibw/blog-server/domain/dto"
	"gorm.io/driver/mysql"
)

// 管理者用のユーザーを作成するツール
func main() {
	time.Local = time.FixedZone("JST", 9*60*60)

	var mode = flag.String("mode", "create", "specify mode (create,delete) default is create")
	flag.Parse()
	m, err := migrate.New("file://"+os.Getenv("MIGRATION_FILE"), "mysql://"+config.PureDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}

	}

	fmt.Print("mailAddress: ")
	var mailAddress string
	fmt.Scan(&mailAddress)

	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}

	userRepository := database.NewUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepository)

	switch *mode {
	case "create":
		createAdmin(userUC, mailAddress)
	case "delete":
		deleteAdmin(userUC, mailAddress)
	}

}

func createAdmin(userUC *usecase.UserUseCase, mailAddress string) {
	fmt.Print("Password (shorter than 72bytes): ")
	fmt.Println()
	pass, err := ReadPassword()
	password := *(*string)(unsafe.Pointer(&pass))

	if len(password) > 72 {
		log.Fatalf("password too long. shorter than 72bytes :%v", err)
	}
	userDTO := &dto.UserDTO{
		MailAddress: mailAddress,
		Password:    password,
	}
	err = userUC.StoreUser(userDTO)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("admin user created successfully")
}

func deleteAdmin(userUC *usecase.UserUseCase, mailAddress string) {
	err := userUC.DeleteUserByMailAddress(mailAddress)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("admin user deleted successfully")
}

func NewDB() (db *gorm.DB, err error) {

	db, err = gorm.Open(mysql.Open(config.DSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		err = fmt.Errorf("failed to open connection: %w", err)
		return
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()

	if err != nil {
		err = fmt.Errorf("failed to get *sql.DB: %w", err)
		return
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(100)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	if err = sqlDB.Ping(); err != nil {
		err = fmt.Errorf("failed to ping: %w", err)
		return
	}

	return
}

// https://qiita.com/x-color/items/f2b6b0852c1a7484ffff
func ReadPassword() ([]byte, error) {
	// Ctrl+Cのシグナルをキャプチャする
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	defer signal.Stop(signalChan)

	// 現在のターミナルの状態をコピーしておく
	currentState, err := terminal.GetState(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	go func() {
		<-signalChan
		// Ctrl+Cを受信後、ターミナルの状態を先ほどのコピーを用いて元に戻す
		_ = terminal.Restore(int(syscall.Stdin), currentState)
		os.Exit(1)
	}()

	return terminal.ReadPassword(syscall.Stdin)
}
