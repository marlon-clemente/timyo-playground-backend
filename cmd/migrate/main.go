package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/marlon-clemente/timyo-playground-backend/database/migrations"
	"github.com/marlon-clemente/timyo-playground-backend/packages/config"
)

const migrationsDir = "database/migrations"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// create does not need a DB connection
	if command == "create" {
		if len(os.Args) < 3 {
			log.Fatal("uso: migrate create <nome>")
		}
		runCreate(os.Args[2])
		return
	}

	cfg := config.Load()
	
	if cfg.DatabaseDSN == "" {
		log.Fatal("DATABASE_DSN não está definido")
	}

	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
		fmt.Println("migrações aplicadas com sucesso")

	case "down":
		n := 1
		if len(os.Args) > 2 {
			n, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("argumento inválido: %s", os.Args[2])
			}
		}
		if n <= 0 {
			log.Fatalf("número de steps deve ser positivo, recebido: %d", n)
		}
		if err := m.Steps(-n); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
		fmt.Printf("%d migração(ões) revertida(s)\n", n)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("force requer um número de versão")
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("versão inválida: %s", os.Args[2])
		}
		if err := m.Force(v); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("versão forçada para %d\n", v)

	case "version":
		v, dirty, err := m.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("nenhuma migração aplicada ainda")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("versão: %d, dirty: %v\n", v, dirty)

	default:
		printUsage()
		os.Exit(1)
	}
}

func runCreate(name string) {
	// Validate name: only alphanumeric, underscore, and hyphen
	if !regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(name) {
		log.Fatalf("nome inválido: %q (use apenas [a-z0-9_-])", name)
	}

	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Fatal(err)
	}

	ts := time.Now().Format("20060102150405")
	base := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s", ts, name))

	// Create both up and down files atomically (fail if any exists)
	for _, suffix := range []string{".up.sql", ".down.sql"} {
		path := base + suffix
		// Check if file already exists
		if _, err := os.Stat(path); err == nil {
			log.Fatalf("arquivo já existe: %q", path)
		}
		// Create file with exclusive permissions
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("erro ao criar arquivo: %v", err)
		}
		f.Close()
		fmt.Println("criado:", path)
	}
}


func printUsage() {
	fmt.Println("uso: go run ./cmd/migrate <comando> [args]")
	fmt.Println()
	fmt.Println("comandos:")
	fmt.Println("  up              aplica todas as migrações pendentes")
	fmt.Println("  down [n]        reverte n migrações (padrão: 1)")
	fmt.Println("  force <versão>  força uma versão específica")
	fmt.Println("  version         exibe a versão atual")
	fmt.Println("  create <nome>   cria um novo par de arquivos SQL")
}
