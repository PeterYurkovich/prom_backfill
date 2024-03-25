package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	tsdb_errors "github.com/prometheus/prometheus/tsdb/errors"
)

func main() {
	start_time_days := flag.Int64("start", 14, "Number of days in the past to start backfilling (default 14 days)")
	end_time_days := flag.Int64("end", 0, "Number of days in the past to end backfilling (default 0 days [now])")
	cnrbt := flag.Bool("cnrbt", false, "Use the container_network_receive_bytes_total tree")
	flag.Parse()

	if *start_time_days < 0 || *end_time_days < 0 {
		fmt.Println("Start Time and End Time must be positive")
		os.Exit(1)
	}
	if *start_time_days < *end_time_days {
		fmt.Println("End Time must be greater than Start Time")
		os.Exit(1)
	}

	flag.Parse()

	allSeries := []labels.Labels{}
	if *cnrbt {
		crnbt_tree := create_cnrbt_tree()
		cnrbtSeries, err := crnbt_tree.getSeries(nil)
		noErr(err)
		allSeries = append(allSeries, cnrbtSeries...)
	}

	if len(allSeries) == 0 {
		fmt.Println("No series selected")
		os.Exit(1)
	}

	err := os.Mkdir("tsdb", 0700)
	defer os.RemoveAll("tsdb")
	noErr(err)

	createBlocks("tsdb", false, allSeries, *start_time_days, *end_time_days)

	currDir, err := os.Getwd()
	noErr(err)
	err = Tar(currDir+"/tsdb", "tsdb")
	noErr(err)
	err = Gzip(currDir+"/tsdb/tsdb.tar", "")
	noErr(err)
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getCompatibleBlockDuration(maxBlockDuration int64) int64 {
	blockDuration := tsdb.DefaultBlockDuration
	if maxBlockDuration > tsdb.DefaultBlockDuration {
		ranges := tsdb.ExponentialBlockRanges(tsdb.DefaultBlockDuration, 10, 3)
		idx := len(ranges) - 1 // Use largest range if user asked for something enormous.
		for i, v := range ranges {
			if v > maxBlockDuration {
				idx = i - 1
				break
			}
		}
		blockDuration = ranges[idx]
	}
	return blockDuration
}

func createBlocks(outputDir string, quiet bool, series []labels.Labels, startDay, endDay int64) (returnErr error) {
	mint := time.Now().UnixMilli() - startDay*24*time.Hour.Milliseconds()
	maxt := time.Now().UnixMilli() - endDay*24*time.Hour.Milliseconds()
	maxSamplesInAppender := 5000
	blockDuration := getCompatibleBlockDuration(2 * time.Hour.Milliseconds())
	mint = blockDuration * (mint / blockDuration)

	db, err := tsdb.OpenDBReadOnly(outputDir, nil)
	if err != nil {
		return err
	}
	defer func() {
		returnErr = tsdb_errors.NewMulti(returnErr, db.Close()).Err()
	}()

	var wroteHeader = false
	setSeriesCache, getSeriesCache := getSeriesCache()

	for t := mint; t <= maxt; t += blockDuration {
		tsUpper := t + blockDuration

		err := func() error {
			w, err := tsdb.NewBlockWriter(log.NewNopLogger(), outputDir, 2*blockDuration)
			if err != nil {
				return fmt.Errorf("block writer: %w", err)
			}
			defer func() {
				err = tsdb_errors.NewMulti(err, w.Close()).Err()
			}()

			ctx := context.Background()
			app := w.Appender(ctx)
			samplesCount := 0
			for i := t; i < tsUpper; i += 30 * time.Second.Milliseconds() {
				for _, series := range series {
					ref, ok := getSeriesCache(series)
					if !ok {
						newRef, err := app.Append(0, series, i, rand.Float64())
						noErr(err)
						setSeriesCache(series, newRef)
						continue
					} else {
						_, err = app.Append(ref, series, i, rand.Float64())
						noErr(err)
					}

					samplesCount++
					if samplesCount < maxSamplesInAppender {
						continue
					}

					// If we arrive here, the samples count is greater than the maxSamplesInAppender.
					// Therefore the old appender is committed and a new one is created.
					// This prevents keeping too many samples lined up in an appender and thus in RAM.
					if err := app.Commit(); err != nil {
						return fmt.Errorf("commit: %w", err)
					}

					app = w.Appender(ctx)
					samplesCount = 0
				}
			}

			if err := app.Commit(); err != nil {
				return fmt.Errorf("commit: %w", err)
			}

			block, err := w.Flush(ctx)
			switch {
			case err == nil:
				if quiet {
					break
				}
				blocks, err := db.Blocks()
				if err != nil {
					return fmt.Errorf("get blocks: %w", err)
				}
				for _, b := range blocks {
					if b.Meta().ULID == block {
						printBlocks([]tsdb.BlockReader{b}, !wroteHeader, false)
						wroteHeader = true
						break
					}
				}
			case errors.Is(err, tsdb.ErrNoSeriesAppended):
			default:
				return fmt.Errorf("flush: %w", err)
			}

			return nil
		}()
		if err != nil {
			return fmt.Errorf("process blocks: %w", err)
		}
	}
	return nil
}
