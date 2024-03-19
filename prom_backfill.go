package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

func main() {
	err := os.Mkdir("tsdb-test", 0700)
	noErr(err)

	db, err := tsdb.Open("tsdb-test", nil, nil, tsdb.DefaultOptions(), nil)
	noErr(err)

	app := db.Appender(context.Background())
	seriesBuilder := labels.NewScratchBuilder(2)
	seriesBuilder.Add("doo", "da")
	seriesBuilder.Add("foo", "bar")
	series := seriesBuilder.Labels()

	ref, err := app.Append(0, series, time.Now().Unix(), 100)
	noErr(err)

	for i := 0.0; i < 100; i++ {
		_, err = app.Append(ref, series, time.Now().Unix()+1000*int64(i), 100+i)
		noErr(err)
	}

	err = app.Commit()
	noErr(err)

	querier, err := db.Querier(math.MinInt64, math.MaxInt64)
	noErr(err)
	ss := querier.Select(context.Background(), false, nil, labels.MustNewMatcher(labels.MatchEqual, "foo", "bar"))

	for ss.Next() {
		series := ss.At()
		fmt.Println("series:", series.Labels().String())

		it := series.Iterator(nil)
		for it.Next() != chunkenc.ValNone {
			_, v := it.At()
			fmt.Println("sample", v)
		}

		fmt.Println("it.Err():", it.Err())
	}
	fmt.Println("ss.Err():", ss.Err())
	ws := ss.Warnings()
	if len(ws) > 0 {
		fmt.Println("warnings:", ws)
	}
	err = querier.Close()
	noErr(err)

	err = db.Close()
	noErr(err)

	err = os.RemoveAll("tsdb-test")
	noErr(err)

	fake := randomCNRBT()
	fmt.Printf("%+v\n", fake)

	fakeTwo := randomCNRBT()
	isEqual := cmp.Equal(fake, fakeTwo)
	fmt.Println(isEqual)

}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

type container_network_receive_bytes_total struct {
	Container, Endpoint, Id, Instance, Interface, Job, Metrics_Path, Name, Namespace, Node, Pod, Prometheus, Service string
	Value                                                                                                            float32
}

func randomCNRBT() container_network_receive_bytes_total {
	cnrbt := container_network_receive_bytes_total{
		"Container",
		"Endpoint",
		"Id",
		"Instance",
		"Interface",
		"Job",
		"Metrics_Path",
		"Name",
		"Namespace",
		"Node",
		"Pod",
		"Prometheus",
		"Service",
		123.42,
	}
	return cnrbt
}
