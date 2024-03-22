// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/alecthomas/units"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
)

func printBlocks(blocks []tsdb.BlockReader, writeHeader, humanReadable bool) {
	tw := tabwriter.NewWriter(os.Stdout, 13, 0, 2, ' ', 0)
	defer tw.Flush()

	if writeHeader {
		fmt.Fprintln(tw, "BLOCK ULID\tMIN TIME\tMAX TIME\tDURATION\tNUM SAMPLES\tNUM CHUNKS\tNUM SERIES\tSIZE")
	}

	for _, b := range blocks {
		meta := b.Meta()

		fmt.Fprintf(tw,
			"%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			meta.ULID,
			getFormatedTime(meta.MinTime, humanReadable),
			getFormatedTime(meta.MaxTime, humanReadable),
			time.Duration(meta.MaxTime-meta.MinTime)*time.Millisecond,
			meta.Stats.NumSamples,
			meta.Stats.NumChunks,
			meta.Stats.NumSeries,
			getFormatedBytes(b.Size(), humanReadable),
		)
	}
}

func getFormatedTime(timestamp int64, humanReadable bool) string {
	if humanReadable {
		return time.Unix(timestamp/1000, 0).UTC().String()
	}
	return strconv.FormatInt(timestamp, 10)
}

func getFormatedBytes(bytes int64, humanReadable bool) string {
	if humanReadable {
		return units.Base2Bytes(bytes).String()
	}
	return strconv.FormatInt(bytes, 10)
}

func getSeriesCache() (func(series labels.Labels, ref storage.SeriesRef), func(series labels.Labels) (storage.SeriesRef, bool)) {
	refCache := map[string]storage.SeriesRef{}

	return func(series labels.Labels, ref storage.SeriesRef) {
			stringifiedSeries := fmt.Sprintf("%v", series)
			refCache[stringifiedSeries] = ref
		}, func(series labels.Labels) (storage.SeriesRef, bool) {
			stringifiedSeries := fmt.Sprintf("%v", series)
			foundValue, ok := refCache[stringifiedSeries]
			return foundValue, ok
		}
}
