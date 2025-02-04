package healthchecks

import (
	"sync"
	"time"

	"github.com/kingstonduy/go-core/metadata"
)

func NewHealthChecker(opts ...HealthCheckOption) HealthChecker {
	options := HealthCheckOptions{
		Name:        metadata.DefaultServiceName,
		Version:     metadata.DefaultServiceVersion,
		Description: metadata.DefaultServiceDescription,
	}

	for _, opt := range opts {
		opt(&options)
	}

	app := &HealthCheckerApplication{
		livenessCheckers:  make(map[string]HealthCheckHandler),
		readinessCheckers: make(map[string]HealthCheckHandler),
		Options:           options,
	}

	return app
}

func (app *HealthCheckerApplication) runChecks(checks map[string]HealthCheckHandler) ApplicationHealthDetailed {
	var (
		start     = time.Now()
		wg        sync.WaitGroup
		checklist = make(chan Integration, len(checks))
		result    = ApplicationHealthDetailed{
			Name:         app.Options.Name,
			Version:      app.Options.Version,
			Description:  app.Options.Description,
			Status:       true,
			Date:         start.Format(time.RFC3339),
			Duration:     0,
			Integrations: []Integration{},
		}
	)

	wg.Add(len(checks))
	for name, handler := range checks {
		go func(name string, handler HealthCheckHandler) {
			checklist <- handler.Check(name)
			wg.Done()
		}(name, handler)
	}

	go func() {
		wg.Wait()
		close(checklist)
		result.Duration = time.Since(start).Nanoseconds()
	}()

	for chk := range checklist {
		if !chk.Status {
			result.Status = false
		}
		result.Integrations = append(result.Integrations, chk)
	}

	return result
}

// AddLivenessCheck implements HealthChecker.
func (app *HealthCheckerApplication) AddLivenessCheck(name string, check HealthCheckHandler) {
	app.checksMutex.Lock()
	defer app.checksMutex.Unlock()
	app.livenessCheckers[name] = check
}

// AddReadinessCheck implements HealthChecker.
func (app *HealthCheckerApplication) AddReadinessCheck(name string, check HealthCheckHandler) {
	app.checksMutex.Lock()
	defer app.checksMutex.Unlock()
	app.readinessCheckers[name] = check
}

// LivenessCheck implements HealthChecker.
func (app *HealthCheckerApplication) LivenessCheck() ApplicationHealthDetailed {
	return app.runChecks(app.livenessCheckers)
}

// RedinessCheck implements HealthChecker.
func (app *HealthCheckerApplication) RedinessCheck() ApplicationHealthDetailed {
	return app.runChecks(app.readinessCheckers)
}
