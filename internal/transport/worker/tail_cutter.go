package worker

import (
	"context"
	"sync"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

type TailCutter struct {
	mutex sync.RWMutex
	wg    sync.WaitGroup
	label string
	// context
	ctx    context.Context
	cancel context.CancelFunc
	// timer event duration
	timerDuration time.Duration
	timer         *time.Timer
	tailQueue     chan string
	// get tail use case
	tailListUC usecase.TailGetUseCase[string]
	// cut tail use case
	tailCutUC usecase.TailCutUseCase[string]
	// config
	tailCut         bool
	tailCutDuration time.Duration
	// logger
	log logger.Logger
}

func NewTailCutter(
	label string,
	parentCtx context.Context,
	timerDuration time.Duration,
	tailCut bool,
	tailCutDuration time.Duration,
	tailListUC usecase.TailGetUseCase[string],
	tailCutUC usecase.TailCutUseCase[string],
	logger logger.Logger,
) *TailCutter {
	ctx, cancel := context.WithCancel(parentCtx)
	return &TailCutter{
		label:           label,
		ctx:             ctx,
		cancel:          cancel,
		tailCut:         tailCut,
		tailCutDuration: tailCutDuration,
		timerDuration:   timerDuration,
		tailListUC:      tailListUC,
		tailCutUC:       tailCutUC,
		log:             logger.GetLogger("tail cut worker"),
	}
}

func (tc *TailCutter) Start(workerCount int) error {
	tc.log.Debugf("starting %s tail cutter worker ", tc.label)
	defer tc.log.Debugf("started %s tail cutter worker ", tc.label)
	// timer
	if tc.timer == nil {
		tc.timer = time.NewTimer(tc.timerDuration)
	} else {
		tc.timer.Reset(tc.timerDuration)
	}
	// queue
	tc.tailQueue = make(chan string, 128)
	// launch timer event listener
	tc.wg.Add(1)
	go tc.timerEventListener()
	// create workers
	for i := 0; i < workerCount; i++ {
		tc.wg.Add(1)
		go tc.tailCutWorker(i)
	}

	return nil
}

func (tc *TailCutter) Stop() {
	// timer
	if tc.timer != nil {
		tc.timer.Stop()
	}
	// cancel ctx
	if tc.cancel != nil {
		tc.cancel()
	}
	// close queue
	if tc.tailQueue != nil {
		close(tc.tailQueue)
	}

	// waiting for stop
	tc.wg.Wait()
}

func (tc *TailCutter) SetUp(timerDuration time.Duration, tailCut bool, tailCutDuration time.Duration) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	tc.timerDuration = timerDuration
	tc.tailCut = tailCut
	tc.tailCutDuration = tailCutDuration
}

func (tc *TailCutter) timerEventListener() {
	tc.log.Debugf("start %s tail cutter timer event listener", tc.label)
	defer tc.log.Debugf("finish %s tail cutter timer event listener", tc.label)
	defer tc.wg.Done()

	for {
		select {
		case <-tc.ctx.Done():
			tc.log.Debugf("stop %s tail cutter timer event listener", tc.label)
			return
		case <-tc.timer.C:
			tc.timeEventProcess(time.Now().Add(-tc.timerDuration))
			tc.timer.Reset(tc.timerDuration)
		}
	}
}

func (tc *TailCutter) timeEventProcess(tailLabel time.Time) {
	tc.log.Debugf("start %s tail cutter timer event processor", tc.label)
	defer tc.log.Debugf("finish %s tail cutter timer event processor", tc.label)

	// checks
	queueLength := len(tc.tailQueue)
	tc.log.Debugf("%s tail cutter queue length %v", tc.label, queueLength)
	if queueLength > 0 {
		tc.log.Debugf("%s tail cutter still processing from preior event, exit", tc.label)

		return
	}

	// processing
	ids, err := tc.tailListUC.GetTail(tc.ctx, tailLabel)
	if err != nil {
		tc.log.Errorf("%s tail cutter get tail list error: %v", tc.label, err)

		return
	}

	// put data into queue
	for _, id := range ids {
		select {
		case tc.tailQueue <- id:
		case <-tc.ctx.Done():
			tc.log.Debugf("stop %s tail cutter timer event processor", tc.label)
			return
		}
	}
}

func (tc *TailCutter) tailCutWorker(index int) {
	tc.log.Debugf("start %s tail cutter worker %v", tc.label, index)
	defer tc.log.Debugf("finish %s tail cutter worker %v", tc.label, index)
	defer tc.wg.Done()

	for {
		select {
		case <-tc.ctx.Done():
			tc.log.Debugf("stop %s tail cutter worker %v", tc.label, index)
			return
		case id, ok := <-tc.tailQueue:
			if !ok {
				tc.log.Debugf("stop %s tail cutter worker %v", tc.label, index)
				return
			}
			err := tc.tailCutUC.Cut(tc.ctx, id)
			if err != nil {
				tc.log.Errorf("%s tail cutter worker %v error: %v", tc.label, index, err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}
