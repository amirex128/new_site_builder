package sfmemoryloader

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

// DurationManager handles jobs that run at fixed duration intervals
type DurationManager struct {
	Handler  interface{}
	Duration time.Duration
}

// RandomDurationManager handles jobs with random intervals between min and max durations
type RandomDurationManager struct {
	Handler interface{}
	Min     time.Duration
	Max     time.Duration
}

// CronJobManager handles jobs using cron expression patterns
type CronJobManager struct {
	Handler interface{}
	Crontab string
}

// DailyManager handles jobs that run daily at specific times
type DailyManager struct {
	Handler  interface{}
	Interval int         // Number of days between runs (1 = every day)
	AtTimes  []time.Time // Specific times to run each day
}

// WeeklyManager handles jobs that run on specific days of the week
type WeeklyManager struct {
	Handler  interface{}
	Interval int            // Number of weeks between runs (1 = every week)
	Weekdays []time.Weekday // Days of the week to run
	AtTimes  []time.Time    // Specific times to run on those days
}

// MonthlyManager handles jobs that run on specific days of the month
type MonthlyManager struct {
	Handler     interface{}
	Interval    int         // Number of months between runs (1 = every month)
	DaysOfMonth []int       // Days of the month to run (1-31)
	AtTimes     []time.Time // Specific times to run on those days
}

// OneTimeManager handles jobs that run only once at a specific time
type OneTimeManager struct {
	Handler interface{}
	StartAt time.Time // When to run the job
}

// Scheduler is the main scheduler struct
type Scheduler struct {
	scheduler gocron.Scheduler
	logger    Logger
	jobs      []gocron.Job
}

// NewScheduler creates a new scheduler with the provided logger
func NewScheduler(logger Logger) *Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.ErrorWithCategory(
			Category.System.General,
			SubCategory.Operation.Initialization,
			"failed to create scheduler",
			map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
		)
	}
	return &Scheduler{
		scheduler: scheduler,
		logger:    logger,
		jobs:      make([]gocron.Job, 0),
	}
}

func (s *Scheduler) Duration(jobs ...DurationManager) *Scheduler {
	for _, job := range jobs {
		task := gocron.NewTask(job.Handler)
		j, err := s.scheduler.NewJob(
			gocron.DurationJob(
				job.Duration,
			),
			task,
		)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"loader error",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
		} else {
			s.jobs = append(s.jobs, j)
		}
	}
	return s
}

func (s *Scheduler) RandomDuration(jobs ...RandomDurationManager) *Scheduler {
	for _, job := range jobs {
		task := gocron.NewTask(job.Handler)
		j, err := s.scheduler.NewJob(
			gocron.DurationRandomJob(
				job.Min,
				job.Max,
			),
			task,
		)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"loader error",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
		} else {
			s.jobs = append(s.jobs, j)
		}
	}
	return s
}

func (s *Scheduler) CronJob(jobs ...CronJobManager) *Scheduler {
	for _, job := range jobs {
		task := gocron.NewTask(job.Handler)
		j, err := s.scheduler.NewJob(
			gocron.CronJob(
				job.Crontab,
				false,
			),
			task,
		)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"loader error",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
		} else {
			s.jobs = append(s.jobs, j)
		}
	}
	return s
}

func (s *Scheduler) Daily(jobs ...DailyManager) *Scheduler {
	for _, job := range jobs {
		var atTimes gocron.AtTimes
		if len(job.AtTimes) > 0 {
			atTimes = createAtTimes(job.AtTimes)
		}

		interval := job.Interval
		if interval <= 0 {
			interval = 1 // Default to every day
		}

		task := gocron.NewTask(job.Handler)
		j, err := s.scheduler.NewJob(
			gocron.DailyJob(
				uint(interval),
				atTimes,
			),
			task,
		)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"loader error",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
		} else {
			s.jobs = append(s.jobs, j)
		}
	}
	return s
}

func (s *Scheduler) Weekly(jobs ...WeeklyManager) *Scheduler {
	for _, job := range jobs {
		var atTimes gocron.AtTimes
		if len(job.AtTimes) > 0 {
			atTimes = createAtTimes(job.AtTimes)
		}

		var weekdays gocron.Weekdays
		if len(job.Weekdays) > 0 {
			weekdays = createWeekdays(job.Weekdays)
		} else {
			weekdays = gocron.NewWeekdays(time.Monday)
		}

		interval := job.Interval
		if interval <= 0 {
			interval = 1 // Default to every week
		}

		task := gocron.NewTask(job.Handler)
		j, err := s.scheduler.NewJob(
			gocron.WeeklyJob(
				uint(interval),
				weekdays,
				atTimes,
			),
			task,
		)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"loader error",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
		} else {
			s.jobs = append(s.jobs, j)
		}
	}
	return s
}

func (s *Scheduler) Monthly(jobs ...MonthlyManager) *Scheduler {
	for _, job := range jobs {
		var atTimes gocron.AtTimes
		if len(job.AtTimes) > 0 {
			atTimes = createAtTimes(job.AtTimes)
		}

		var daysOfMonth gocron.DaysOfTheMonth
		if len(job.DaysOfMonth) > 0 {
			firstDay := job.DaysOfMonth[0]
			restDays := job.DaysOfMonth[1:]
			daysOfMonth = gocron.NewDaysOfTheMonth(firstDay, restDays...)
		} else {
			daysOfMonth = gocron.NewDaysOfTheMonth(1)
		}

		interval := job.Interval
		if interval <= 0 {
			interval = 1 // Default to every month
		}

		task := gocron.NewTask(job.Handler)
		j, err := s.scheduler.NewJob(
			gocron.MonthlyJob(
				uint(interval),
				daysOfMonth,
				atTimes,
			),
			task,
		)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"loader error",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
		} else {
			s.jobs = append(s.jobs, j)
		}
	}
	return s
}

func (s *Scheduler) OneTime(jobs ...OneTimeManager) *Scheduler {
	for _, job := range jobs {
		var startAtOption gocron.OneTimeJobStartAtOption

		if job.StartAt.IsZero() {
			startAtOption = gocron.OneTimeJobStartImmediately()
		} else {
			startAtOption = gocron.OneTimeJobStartDateTime(job.StartAt)
		}

		task := gocron.NewTask(job.Handler)
		j, err := s.scheduler.NewJob(
			gocron.OneTimeJob(
				startAtOption,
			),
			task,
		)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"loader error",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
		} else {
			s.jobs = append(s.jobs, j)
		}
	}
	return s
}

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}

// Helper functions to create gocron types
func createAtTimes(times []time.Time) gocron.AtTimes {
	var atTimes []gocron.AtTime
	for _, t := range times {
		atTime := gocron.NewAtTime(uint(t.Hour()), uint(t.Minute()), uint(t.Second()))
		atTimes = append(atTimes, atTime)
	}

	if len(atTimes) > 0 {
		first := atTimes[0]
		rest := atTimes[1:]
		return gocron.NewAtTimes(first, rest...)
	}

	return gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))
}

func createWeekdays(days []time.Weekday) gocron.Weekdays {
	if len(days) > 0 {
		first := days[0]
		rest := days[1:]
		return gocron.NewWeekdays(first, rest...)
	}

	return gocron.NewWeekdays(time.Monday)
}

// WithJobOptions adds options to a job
func (s *Scheduler) WithJobOptions(jobID gocron.Job, options ...gocron.JobOption) error {
	var lastError error
	for _, option := range options {
		_, err := s.scheduler.Update(jobID.ID(), nil, nil, option)
		if err != nil {
			s.logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Status.Error,
				"error updating job",
				map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
			)
			lastError = err
		}
	}
	return lastError
}

// JobInfo returns information about all currently scheduled jobs
func (s *Scheduler) JobInfo() []gocron.Job {
	return s.scheduler.Jobs()
}

// RemoveJob removes a job by its ID
func (s *Scheduler) RemoveJob(jobID gocron.Job) {
	err := s.scheduler.RemoveJob(jobID.ID())
	if err != nil {
		s.logger.ErrorWithCategory(
			Category.System.General,
			SubCategory.Status.Error,
			"error removing job",
			map[string]interface{}{ExtraKey.Error.ErrorMessage: err.Error()},
		)
	}
}

// RemoveByTags removes all jobs with the given tags
func (s *Scheduler) RemoveByTags(tags ...string) {
	s.scheduler.RemoveByTags(tags...)
}

// Configuration option type
type SchedulerOption func(*Scheduler)

// RegisterScheduler registers jobs with various configuration options
func RegisterScheduler(options ...SchedulerOption) (*Scheduler, error) {
	// Create a new scheduler with a default logger
	scheduler := &Scheduler{
		jobs: make([]gocron.Job, 0),
	}

	// Apply all options
	for _, option := range options {
		option(scheduler)
	}

	// If no logger was provided, return an error
	if scheduler.logger == nil {
		return nil, ErrLoggerRequired
	}

	// If no scheduler was created, create one
	if scheduler.scheduler == nil {
		s, err := gocron.NewScheduler()
		if err != nil {
			return nil, err
		}
		scheduler.scheduler = s
	}

	return scheduler, nil
}

// WithLogger configures the logger for the scheduler
func WithLogger(logger Logger) SchedulerOption {
	return func(s *Scheduler) {
		s.logger = logger
	}
}

// WithDurationJobs adds duration jobs to the scheduler
func WithDurationJobs(jobs ...DurationManager) SchedulerOption {
	return func(s *Scheduler) {
		s.Duration(jobs...)
	}
}

// WithRandomDurationJobs adds random duration jobs to the scheduler
func WithRandomDurationJobs(jobs ...RandomDurationManager) SchedulerOption {
	return func(s *Scheduler) {
		s.RandomDuration(jobs...)
	}
}

// WithCronJobs adds cron jobs to the scheduler
func WithCronJobs(jobs ...CronJobManager) SchedulerOption {
	return func(s *Scheduler) {
		s.CronJob(jobs...)
	}
}

// WithDailyJobs adds daily jobs to the scheduler
func WithDailyJobs(jobs ...DailyManager) SchedulerOption {
	return func(s *Scheduler) {
		s.Daily(jobs...)
	}
}

// WithWeeklyJobs adds weekly jobs to the scheduler
func WithWeeklyJobs(jobs ...WeeklyManager) SchedulerOption {
	return func(s *Scheduler) {
		s.Weekly(jobs...)
	}
}

// WithMonthlyJobs adds monthly jobs to the scheduler
func WithMonthlyJobs(jobs ...MonthlyManager) SchedulerOption {
	return func(s *Scheduler) {
		s.Monthly(jobs...)
	}
}

// WithOneTimeJobs adds one-time jobs to the scheduler
func WithOneTimeJobs(jobs ...OneTimeManager) SchedulerOption {
	return func(s *Scheduler) {
		s.OneTime(jobs...)
	}
}

// Error definitions
var (
	ErrLoggerRequired = errorString("logger is required")
)

// errorString is a simple implementation of error
type errorString string

func (e errorString) Error() string {
	return string(e)
}
