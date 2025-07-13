package runner

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"likha/config"

	"likha/generator/factory"
	"likha/generator/types"
	"likha/output"
	output_types "likha/output/types"
	"likha/progress"

	tea "github.com/charmbracelet/bubbletea"
)

// Runner coordinates the data generation process.
type Runner struct {
	config       *config.Config
	count        int64
	generators   map[string]types.Generator
	fieldOrder   []string
	writer       output_types.Writer
	prog         *tea.Program
	progressChan chan progress.ProgressMsg // Channel to send progress updates to the Bubble Tea model
}

// Job represents a single row generation task.
type Job struct {
	Index int64
}

// Result holds the data for a generated row.
type Result struct {
	Index int64
	Data  map[string]interface{}
	Err   error
}

// NewRunner creates and initializes a new Runner.
func NewRunner(cfg *config.Config, count int64) (*Runner, error) {
	gens := make(map[string]types.Generator)
	fieldOrder := make([]string, len(cfg.Fields))

	// First pass: create all non-foreignkey generators to ensure dependencies are available.
	for i, f := range cfg.Fields {
		fieldOrder[i] = f.Name
		if f.Generator.Type != "foreignkey" {
			g, err := factory.NewGenerator(f.Generator, gens)
			if err != nil {
				return nil, fmt.Errorf("error creating generator for field '%s': %w", f.Name, err)
			}
			gens[f.Name] = g
		}
	}
	// Second pass: create foreignkey generators which may depend on others.
	for _, f := range cfg.Fields {
		if f.Generator.Type == "foreignkey" {
			g, err := factory.NewGenerator(f.Generator, gens)
			if err != nil {
				return nil, fmt.Errorf("error creating foreignkey generator for field '%s': %w", f.Name, err)
			}
			gens[f.Name] = g
		}
	}

	// Create the output file.
	file, err := os.Create(cfg.Output.File)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file '%s': %w", cfg.Output.File, err)
	}

	// Create the appropriate writer.
	writer, err := output.NewWriter(&cfg.Output, file)
	if err != nil {
		file.Close() // Clean up the file if writer creation fails.
		return nil, fmt.Errorf("failed to create writer: %w", err)
	}

	// Initialize the progress bar model and get its update channel.
	progressModel, progressChan := progress.NewModel(count)
	p := tea.NewProgram(progressModel)

	return &Runner{
		config:       cfg,
		count:        count,
		generators:   gens,
		fieldOrder:   fieldOrder,
		writer:       writer,
		prog:         p,
		progressChan: progressChan, // Store the channel to send updates
	}, nil
}

// Run starts the generation process using a worker pool and shows a progress bar.
func (r *Runner) Run() error {
	// Ensure the writer is closed and the file handle is released on exit.
	if c, ok := r.writer.(interface{ Close() error }); ok {
		defer c.Close()
	}

	// Write the header row for formats that support it (e.g., CSV).
	if err := r.writer.WriteHeader(r.fieldOrder); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Set up a worker pool to parallelize generation.
	numWorkers := runtime.NumCPU()
	jobs := make(chan Job, numWorkers)
	results := make(chan Result, numWorkers)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go r.worker(&wg, jobs, results)
	}

	// Start the actual work in a goroutine
	go func() {
		defer func() {
			// Wait for all workers to finish
			wg.Wait()
			close(results) // Close the results channel after all workers are done
			// Close the progress channel to signal the progress model that no more updates are coming
			close(r.progressChan)
		}()

		// Goroutine to feed jobs to the workers.
		go func() {
			for i := int64(0); i < r.count; i++ {
				jobs <- Job{Index: i}
			}
			close(jobs) // Close the jobs channel after all jobs are sent
		}()

		// Process results from the workers.
		var processedCount int64
		var hasError bool // Flag to indicate if an error has occurred

		for processedCount < r.count {
			if hasError {
				// If an error has occurred, break from processing further results.
				// The tea.Program will quit via the ErrorMsg handler.
				break
			}

			select {
			case result, ok := <-results:
				if !ok {
					// Results channel closed unexpectedly before all records were processed.
					if !hasError {
						r.progressChan <- progress.ProgressMsg{
							Current: processedCount,
							Total:   r.count,
							Done:    false,
							Error:   fmt.Errorf("results channel closed unexpectedly before all records processed (%d/%d)", processedCount, r.count),
						}
					}
					hasError = true
					break
				}

				if result.Err != nil {
					// Send error to progress bar via the channel
					r.progressChan <- progress.ProgressMsg{
						Current: processedCount,
						Total:   r.count,
						Done:    false,
						Error:   result.Err,
					}
					hasError = true
					break
				} else {
					if err := r.writer.WriteRow(result.Data); err != nil {
						// Send error to progress bar and stop via the channel
						r.progressChan <- progress.ProgressMsg{
							Current: processedCount,
							Total:   r.count,
							Done:    false,
							Error:   fmt.Errorf("failed to write row %d: %w", result.Index, err),
						}
						hasError = true
						break
					}
				}
				processedCount++
				// Send progress update to the progress bar via the channel
				r.progressChan <- progress.ProgressMsg{
					Current: processedCount,
					Total:   r.count,
					Done:    processedCount == r.count,
					Error:   nil,
				}
			}
		}
	}()

	// Start the Bubble Tea progress bar in the main thread (this blocks until quit)
	// The program will quit when 100% progress is reached or an ErrorMsg is sent.
	_, err := r.prog.Run()
	return err
}

// worker is the function run by each goroutine in the pool.
// It receives jobs, generates data, and sends results back.
func (r *Runner) worker(wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for job := range jobs {
		rowData := make(map[string]interface{})
		// We must generate fields in the order specified in the config
		// to ensure dependencies like foreign keys are met.
		for _, fieldName := range r.fieldOrder {
			gen, ok := r.generators[fieldName]
			if !ok {
				results <- Result{Index: job.Index, Err: fmt.Errorf("internal error: generator for field '%s' not found", fieldName)}
				goto nextJob // Use goto to break out of nested loops on error
			}
			val, err := gen.Generate(rowData)
			if err != nil {
				results <- Result{Index: job.Index, Err: fmt.Errorf("field '%s': %w", fieldName, err)}
				goto nextJob // Use goto to break out of nested loops on error
			}
			rowData[fieldName] = val
		}
		results <- Result{Index: job.Index, Data: rowData}
	nextJob: // Label for goto
	}
}
