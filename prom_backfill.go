package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/prometheus/tsdb"
	tsdb_errors "github.com/prometheus/prometheus/tsdb/errors"
)

func main() {
	err := os.Mkdir(".tsdb_test", 0700)
	noErr(err)

	createBlocks(".tsdb_test", false)
	// err = os.RemoveAll("tsdb_test")
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

func createBlocks(outputDir string, quiet bool) (returnErr error) {
	mint := time.Now().Unix() - 6*24*time.Hour.Milliseconds() // 6 days go
	maxt := time.Now().Unix() - 5*24*time.Hour.Milliseconds() // 5 days ago
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
			cnrbtCache, setCnrbtRef := getSeriesCache()
			randomCnrbtGenerator := randomCnrbt()
			for i := t; i < tsUpper; i += 30 * time.Second.Milliseconds() {
				for j := 0; j < 100; j++ {
					cnrbt := randomCnrbtGenerator()
					labels, cachedRef := cnrbtCache(cnrbt)
					if cachedRef == 0 {
						newRef, err := app.Append(0, labels, i, 100)
						noErr(err)
						setCnrbtRef(cnrbt, newRef)
					} else {
						_, err = app.Append(cachedRef, labels, i, 100)
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
