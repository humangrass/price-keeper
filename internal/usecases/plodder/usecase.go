package plodder

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/humangrass/price-keeper/domain/models"
	"github.com/humangrass/price-keeper/domain/repository"
	"github.com/humangrass/price-keeper/pgk/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type UseCase struct {
	pairRepo repository.PairsRepository
	logger   *logger.Logger

	activePairs     map[uuid.UUID]models.Pair
	pairsMutex      sync.RWMutex
	refreshInterval time.Duration
	cron            *cron.Cron
}

func NewPlodderUseCase(
	pairRepo repository.PairsRepository,
	logger *logger.Logger,
	refreshInterval time.Duration,
) *UseCase {
	return &UseCase{
		pairRepo:        pairRepo,
		logger:          logger,
		activePairs:     make(map[uuid.UUID]models.Pair),
		cron:            cron.New(cron.WithSeconds()),
		refreshInterval: refreshInterval,
	}
}

func (uc *UseCase) Run(ctx context.Context) {
	var wg sync.WaitGroup

	uc.cron.Start()

	wg.Add(1)
	go uc.taskUpdater(ctx, &wg)

	<-ctx.Done()

	uc.cron.Stop()

	wg.Wait()
}

func (uc *UseCase) taskUpdater(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(uc.refreshInterval)
	defer ticker.Stop()

	if err := uc.updateTasks(ctx); err != nil {
		uc.logger.Sugar().Error("failed to update tasks", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			uc.logger.Sugar().Info("Stopping task updater")
			return
		case <-ticker.C:
			if err := uc.updateTasks(ctx); err != nil {
				uc.logger.Sugar().Error("failed to update tasks", zap.Error(err))
			}
		}
	}
}

func (uc *UseCase) updateTasks(ctx context.Context) error {
	pairs, err := uc.pairRepo.GetActivePairs(ctx)
	if err != nil {
		return err
	}

	newActivePairs := make(map[uuid.UUID]models.Pair)
	for _, pair := range pairs {
		newActivePairs[pair.UUID] = pair
	}

	uc.pairsMutex.Lock()
	defer uc.pairsMutex.Unlock()

	for uuid, pair := range newActivePairs {
		if _, exists := uc.activePairs[uuid]; !exists {
			uc.addTaskToCron(pair)
		}
	}

	for uuid, pair := range uc.activePairs {
		if _, exists := newActivePairs[uuid]; !exists {
			uc.removeTaskFromCron(pair)
		}
	}
	uc.activePairs = newActivePairs

	return nil
}

func (uc *UseCase) addTaskToCron(pair models.Pair) {
	schedule := fmt.Sprintf("@every %s", time.Duration(pair.Timeframe).String())

	_, err := uc.cron.AddFunc(schedule, func() {
		if err := uc.processPair(context.Background(), pair); err != nil {
			uc.logger.Sugar().Error("failed to process pair",
				zap.String("pair", pair.UUID.String()),
				zap.Error(err))
		}
	})

	if err != nil {
		uc.logger.Sugar().Error("failed to add task to cron",
			zap.String("pair", pair.UUID.String()),
			zap.Error(err))
		return
	}

	uc.logger.Sugar().Debug("Added new task to cron",
		zap.String("pair", pair.UUID.String()),
		zap.String("schedule", schedule))
}

func (uc *UseCase) removeTaskFromCron(removedPair models.Pair) {
	// TODO: приходится пересоздавать крон множество раз, плохое решение
	// Вот бы переиспользовать старые сущности
	// entries := uc.cron.Entries()
	uc.cron.Stop()
	uc.cron = cron.New(cron.WithSeconds())

	// for _, entry := range entries {
	// 	fmt.Println(entry)
	// }

	for _, pair := range uc.activePairs {
		if removedPair == pair {
			continue
		}
		uc.addTaskToCron(pair)
	}

	uc.cron.Start()
}

func (uc *UseCase) processPair(ctx context.Context, pair models.Pair) error {
	uc.logger.Sugar().Debug("Process pair", zap.String("pair", pair.String()))

	return nil
}
