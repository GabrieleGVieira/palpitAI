package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/config"
	"github.com/gabrielevieira/palpitai/backend/internal/data/importer"
	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
	"github.com/gabrielevieira/palpitai/backend/internal/data/normalizer"
	"github.com/gabrielevieira/palpitai/backend/internal/data/repository"
	"github.com/gabrielevieira/palpitai/backend/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	importType := flag.String("type", "", "import type: matches, goalscorers, fifa-ranking")
	filePath := flag.String("file", "", "CSV file path")
	flag.Parse()

	if *importType == "" || *filePath == "" {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/import-history --type=matches --file=./data/results.csv")
		os.Exit(2)
	}

	_ = godotenv.Load()
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.Migrate(ctx, db); err != nil {
		logger.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	teamsRepo := repository.NewTeamsRepository(db)
	matchesRepo := repository.NewHistoricalMatchesRepository(db)
	goalscorersRepo := repository.NewHistoricalGoalscorersRepository(db)
	rankingsRepo := repository.NewFifaRankingsRepository(db)
	importLogsRepo := repository.NewImportLogsRepository(db)
	teamNormalizer := normalizer.NewTeamNameNormalizer(teamsRepo)

	logID, err := importLogsRepo.Start(ctx, *importType, *filePath)
	if err != nil {
		logger.Error("import log start failed", "error", err)
		os.Exit(1)
	}

	var result importer.ImportResult
	switch *importType {
	case "matches":
		result, err = importer.NewHistoricalMatchesImporter(teamNormalizer, matchesRepo).ImportFile(ctx, *filePath)
	case "goalscorers":
		result, err = importer.NewHistoricalGoalscorersImporter(teamNormalizer, matchesRepo, goalscorersRepo).ImportFile(ctx, *filePath)
	case "fifa-ranking":
		result, err = importer.NewFifaRankingImporter(teamNormalizer, rankingsRepo).ImportFile(ctx, *filePath)
	default:
		err = fmt.Errorf("unsupported import type: %s", *importType)
	}

	status := "completed"
	var errorMessage *string
	if err != nil {
		status = "failed"
		message := err.Error()
		errorMessage = &message
	}
	if finishErr := importLogsRepo.Finish(ctx, logID, status, result.ImportSummary, errorMessage); finishErr != nil {
		logger.Error("import log finish failed", "error", finishErr)
	}
	if err != nil {
		logger.Error("import failed", "type", *importType, "file", *filePath, "error", err)
		os.Exit(1)
	}

	printReport(*importType, *filePath, status, result.ImportSummary)
}

func printReport(importType string, filePath string, status string, summary models.ImportSummary) {
	fmt.Println("Import finished")
	fmt.Printf("type: %s\n", importType)
	fmt.Printf("file: %s\n", filePath)
	fmt.Printf("status: %s\n", status)
	fmt.Printf("processed: %d\n", summary.ProcessedCount)
	fmt.Printf("inserted: %d\n", summary.InsertedCount)
	fmt.Printf("skipped: %d\n", summary.SkippedCount)
	fmt.Printf("errors: %d\n", summary.ErrorCount)
}
