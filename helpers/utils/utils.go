package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog/log"

	emailNormalizer "github.com/dimuska139/go-email-normalizer"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const DEFAULT_PHONENUMBER_COUNTRY = "NG"

type DbConfig struct {
	User     string
	Password string
	Database string
}

type NormalizePhonenumberOpts struct {
	Phonenumber string
	CountryCode string
}

func NormalizePhonenumber(opts NormalizePhonenumberOpts) (string, error) {
	if opts.CountryCode == "" {
		opts.CountryCode = DEFAULT_PHONENUMBER_COUNTRY
	}

	n, err := phonenumbers.Parse(opts.Phonenumber, opts.CountryCode)
	if err != nil {
		return "", err
	}

	return phonenumbers.Format(n, phonenumbers.E164), nil
}

// TODO: Test NormalizeEmail
func NormalizeEmail(email string) string {
	n := emailNormalizer.NewNormalizer()
	return n.Normalize(strings.ToLower(email))
}

func DropMigrations(dbUrl string, source source.Driver) error {
	m, err := migrate.NewWithSourceInstance("iofs", source, dbUrl)

	if err != nil {
		return err
	}

	if err := m.Drop(); err != nil {
		return err
	}

	return nil
}

func RunMigrations(dbUrl string, source source.Driver) error {
	m, err := migrate.NewWithSourceInstance("iofs", source, dbUrl)

	if err != nil {
		return err
	}

	log.Info().Msg("running migrations")

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func WaitForDb(dbUrl string) error {
	db, err := sqlx.Open("postgres", dbUrl)

	if err != nil {
		return err
	}

	defer db.Close()

	log.Info().Msg("waiting for database")

	for timeout := time.After(time.Second * 30); ; {
		select {
		case <-timeout:
			return errors.New("database connection timed out")
		default:
			if err := db.Ping(); err != nil {
				log.Info().Msg("waiting...")
				time.Sleep(time.Millisecond * 500)
				continue
			}
			log.Info().Msg("database connected")
			return nil
		}
	}
}

func PullImage(ctx context.Context, cli *client.Client, imageName string) error {
	log.Info().Msgf("pulling docker image: %s", imageName)
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}

func StopContainer(ctx context.Context, cli *client.Client, id string) error {
	timeOut := int(30 * time.Second)
	stopOptions := container.StopOptions{
		Timeout: &timeOut,
	}
	return cli.ContainerStop(ctx, id, stopOptions)
}

func GetFreePort() (int, error) {
	tcpAddress, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err == nil {
		if listener, err := net.ListenTCP("tcp", tcpAddress); err == nil {
			defer listener.Close()
			return listener.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return 0, err
}

func CreateDbContainer(ctx context.Context, cli *client.Client, config DbConfig) (string, string, error) {
	imageName := "postgres:14"

	if err := PullImage(ctx, cli, imageName); err != nil {
		return "", "", err
	}

	port, err := GetFreePort()
	if err != nil {
		return "", "", err
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Env: []string{
				"POSTGRES_USER=" + config.User,
				"POSTGRES_PASSWORD=" + config.Password,
				"POSTGRES_DB=" + config.Database,
				"listen_addresses = '*'",
			},
			ExposedPorts: nat.PortSet{
				"5432/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			AutoRemove: true,
			PortBindings: nat.PortMap{
				"5432/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprint(port),
					},
				},
			},
		},
		nil, nil, "",
	)
	if err != nil {
		return "", "", err
	}

	dbHost := "0.0.0.0"

	if linuxRunningOnDockerDesktop(ctx, cli) {
		dbHost = "gateway.docker.internal"
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", err
	}

	dbUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		dbHost,
		port,
		config.Database,
	)

	return dbUrl, resp.ID, nil
}

func linuxRunningOnDockerDesktop(ctx context.Context, dockerClient *client.Client) bool {
	info, _ := dockerClient.Info(ctx)
	return info.Name == "docker-desktop" && runtime.GOOS == "linux"
}
