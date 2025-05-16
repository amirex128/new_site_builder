# Memory Loader

A Go scheduler library built on top of [gocron](https://github.com/go-co-op/gocron) to provide easy job scheduling with a functional options pattern.

## Features

- Simple, fluent API for scheduling jobs
- Multiple scheduling options: Duration, CronJob, Daily, Weekly, Monthly, and OneTime
- Flexible handler support (accepts any function type)
- Functional options pattern for easy configuration
- Builder pattern for easy job registration
- Preconfigured job scheduling patterns

## Installation

```go
go get github.com/go-co-op/gocron/v2
```

## Usage

### Basic Setup

```go
import (
    "time"
    "github.com/yourorg/memoryloader"
)

// Create a logger
logger := yourLoggerImplementation

// Initialize scheduler with a logger
scheduler := memoryloader.NewScheduler(logger)

// Add jobs using fluent API
scheduler.Duration(
    memoryloader.DurationManager{
        Handler: func() {
            // Your job logic here
            fmt.Println("Duration job executed")
        },
        Duration: time.Second * 30,
    },
).Start()
```

### Using Functional Options Pattern

The library supports a functional options pattern similar to the example you provided:

```go
import (
    "time"
    "github.com/yourorg/memoryloader"
)

// Create a logger
logger := yourLoggerImplementation

// Register scheduler with options
scheduler, err := memoryloader.RegisterScheduler(
    memoryloader.WithLogger(logger),
    memoryloader.WithDurationJobs(
        memoryloader.DurationManager{
            Handler: func() {
                fmt.Println("Running every 5 minutes")
            },
            Duration: time.Minute * 5,
        },
    ),
    memoryloader.WithCronJobs(
        memoryloader.CronJobManager{
            Handler: func() {
                fmt.Println("Running at midnight")
            },
            Crontab: "0 0 * * *", // Run at midnight
        },
    ),
)

if err != nil {
    log.Fatalf("Failed to register scheduler: %v", err)
}

// Start the scheduler
scheduler.Start()
```

### Using Multiple Job Types

```go
scheduler := memoryloader.NewScheduler(logger)

// Add multiple job types using fluent API
scheduler.Duration(
    memoryloader.DurationManager{
        Handler: func() {
            fmt.Println("Running every 5 minutes")
        },
        Duration: time.Minute * 5,
    },
).CronJob(
    memoryloader.CronJobManager{
        Handler: func() {
            fmt.Println("Running at midnight")
        },
        Crontab: "0 0 * * *", // Run at midnight
    },
).Daily(
    memoryloader.DailyManager{
        Handler: func() {
            fmt.Println("Running daily at 8 AM and 5 PM")
        },
        Interval: 1,
        AtTimes: []time.Time{
            time.Date(0, 0, 0, 8, 0, 0, 0, time.Local),
            time.Date(0, 0, 0, 17, 0, 0, 0, time.Local),
        },
    },
).Start()
```

## Flexible Handler Support

The library now supports any type of handler function:

```go
// Simple function with no parameters
scheduler.Duration(
    memoryloader.DurationManager{
        Handler: func() {
            fmt.Println("Simple job")
        },
        Duration: time.Minute,
    },
)

// Function with parameters
scheduler.CronJob(
    memoryloader.CronJobManager{
        Handler: func(name string, count int) {
            fmt.Printf("Job %s ran %d times\n", name, count)
        },
        Crontab: "*/5 * * * *",
    },
)

// Method from a struct
type MyService struct {
    Name string
}

func (s *MyService) DoWork() {
    fmt.Printf("Service %s is working\n", s.Name)
}

myService := &MyService{Name: "BackgroundService"}
scheduler.Daily(
    memoryloader.DailyManager{
        Handler: myService.DoWork,
        Interval: 1,
        AtTimes: []time.Time{time.Date(0, 0, 0, 12, 0, 0, 0, time.Local)},
    },
)
```

## Job Types

### Duration Jobs

Runs jobs at fixed time intervals.

```go
scheduler.Duration(
    memoryloader.DurationManager{
        Handler: func() {
            // Run every 10 seconds
        },
        Duration: time.Second * 10,
    },
)
```

### Random Duration Jobs

Runs jobs at random intervals between min and max durations.

```go
scheduler.RandomDuration(
    memoryloader.RandomDurationManager{
        Handler: func() {
            // Run at random intervals between 1 and 5 minutes
        },
        Min: time.Minute,
        Max: time.Minute * 5,
    },
)
```

### Cron Jobs

Runs jobs according to cron expressions.

```go
scheduler.CronJob(
    memoryloader.CronJobManager{
        Handler: func() {
            // Run at 15 minutes past every hour
        },
        Crontab: "15 * * * *",
    },
)
```

### Daily Jobs

Runs jobs daily at specific times.

```go
scheduler.Daily(
    memoryloader.DailyManager{
        Handler: func() {
            // Run every 2 days at 9 AM and 6 PM
        },
        Interval: 2, // Every 2 days
        AtTimes: []time.Time{
            time.Date(0, 0, 0, 9, 0, 0, 0, time.Local),
            time.Date(0, 0, 0, 18, 0, 0, 0, time.Local),
        },
    },
)
```

### Weekly Jobs

Runs jobs on specific days of the week.

```go
scheduler.Weekly(
    memoryloader.WeeklyManager{
        Handler: func() {
            // Run every Monday and Friday at 10 AM
        },
        Interval: 1, // Every week
        Weekdays: []time.Weekday{time.Monday, time.Friday},
        AtTimes: []time.Time{time.Date(0, 0, 0, 10, 0, 0, 0, time.Local)},
    },
)
```

### Monthly Jobs

Runs jobs on specific days of the month.

```go
scheduler.Monthly(
    memoryloader.MonthlyManager{
        Handler: func() {
            // Run on the 1st and 15th of every month at noon
        },
        Interval: 1, // Every month
        DaysOfMonth: []int{1, 15},
        AtTimes: []time.Time{time.Date(0, 0, 0, 12, 0, 0, 0, time.Local)},
    },
)
```

### One-Time Jobs

Runs jobs once at a specific time.

```go
scheduler.OneTime(
    memoryloader.OneTimeManager{
        Handler: func() {
            // Run once at the specified time
        },
        StartAt: time.Date(2023, 12, 31, 23, 59, 59, 0, time.Local), // New Year's Eve
    },
)
```

## Job Management

### Starting the Scheduler

```go
// Register your jobs, then start the scheduler
scheduler.Start()
```

### Shutting Down

```go
// Gracefully shutdown when your application exits
err := scheduler.Shutdown()
if err != nil {
    // Handle error
}
```

### Job Options

Additional options can be applied to jobs:

```go
// First add a job
scheduler.Duration(
    memoryloader.DurationManager{
        Handler: func() {
            fmt.Println("Cleanup job")
        },
        Duration: time.Hour,
    },
)

// Get the job info
jobs := scheduler.JobInfo()
if len(jobs) > 0 {
    // Apply options to the first job
    scheduler.WithJobOptions(jobs[0], 
        gocron.WithTags("background", "cleanup"),
        gocron.WithName("CleanupJob"),
    )
}
```

### Removing Jobs

```go
// Get the job info
jobs := scheduler.JobInfo()
if len(jobs) > 0 {
    // Remove a job by its ID
    scheduler.RemoveJob(jobs[0])
}

// Remove jobs by tags
scheduler.RemoveByTags("cleanup", "background")
```

## Manager Struct Reference

### DurationManager

```go
type DurationManager struct {
    Handler  interface{} // Any function type
    Duration time.Duration // Time between executions
}
```

### RandomDurationManager

```go
type RandomDurationManager struct {
    Handler interface{} // Any function type
    Min     time.Duration // Minimum time between executions
    Max     time.Duration // Maximum time between executions
}
```

### CronJobManager

```go
type CronJobManager struct {
    Handler interface{} // Any function type
    Crontab string // Cron expression (e.g. "0 * * * *")
}
```

### DailyManager

```go
type DailyManager struct {
    Handler  interface{} // Any function type
    Interval int // Days between runs (1 = every day)
    AtTimes  []time.Time // Times to run each day
}
```

### WeeklyManager

```go
type WeeklyManager struct {
    Handler  interface{} // Any function type
    Interval int // Weeks between runs (1 = every week)
    Weekdays []time.Weekday // Days of the week to run
    AtTimes  []time.Time // Times to run on those days
}
```

### MonthlyManager

```go
type MonthlyManager struct {
    Handler     interface{} // Any function type
    Interval    int // Months between runs (1 = every month)
    DaysOfMonth []int // Days of the month to run (1-31)
    AtTimes     []time.Time // Times to run on those days
}
```

### OneTimeManager

```go
type OneTimeManager struct {
    Handler interface{} // Any function type
    StartAt time.Time // When to run the job (use zero time for immediate execution)
}
```

## Functional Options

```go
// Register a scheduler with various options
scheduler, err := memoryloader.RegisterScheduler(
    // Configure the logger
    memoryloader.WithLogger(logger),
    
    // Add duration jobs
    memoryloader.WithDurationJobs(
        memoryloader.DurationManager{
            Handler: func() { 
                fmt.Println("Duration job") 
            },
            Duration: time.Minute,
        },
    ),
    
    // Add cron jobs
    memoryloader.WithCronJobs(
        memoryloader.CronJobManager{
            Handler: func() { 
                fmt.Println("Cron job") 
            },
            Crontab: "0 * * * *",
        },
    ),
    
    // Add daily jobs
    memoryloader.WithDailyJobs(
        memoryloader.DailyManager{
            Handler: func() { 
                fmt.Println("Daily job") 
            },
            Interval: 1,
            AtTimes: []time.Time{time.Date(0, 0, 0, 9, 0, 0, 0, time.Local)},
        },
    ),
)

if err != nil {
    log.Fatalf("Failed to register scheduler: %v", err)
}

scheduler.Start()
```

## License

MIT
