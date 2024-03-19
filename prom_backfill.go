package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
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

	faker := randomCnrbt()

	fake := faker()
	fmt.Printf("%+v\n", fake)

	fakeTwo := faker()
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

func randomCnrbt() func() *container_network_receive_bytes_total {
	cnrbtOptions := struct {
		container, endpoint, id, instance, interface_, job, metricsPath, name, namespace, node, pod, prometheus, service []string
	}{
		container:   []string{"container1", "container2", "container3"},
		endpoint:    []string{"endpoint1", "endpoint2", "endpoint3"},
		id:          []string{"id1", "id2", "id3"},
		instance:    []string{"instance1", "instance2", "instance3"},
		interface_:  []string{"interface1", "interface2", "interface3"},
		job:         []string{"job1", "job2", "job3"},
		metricsPath: []string{"metricsPath1", "metricsPath2", "metricsPath3"},
		name:        []string{"name1", "name2", "name3"},
		namespace:   []string{"namespace1", "namespace2", "namespace3"},
		node:        []string{"node1", "node2", "node3"},
		pod:         []string{"pod1", "pod2", "pod3"},
		prometheus:  []string{"prometheus1", "prometheus2", "prometheus3"},
		service:     []string{"service1", "service2", "service3"},
	}

	randOption := func(slice []string) string {
		return slice[rand.Intn(len(slice))]
	}

	return func() *container_network_receive_bytes_total {
		return &container_network_receive_bytes_total{
			randOption(cnrbtOptions.container),
			randOption(cnrbtOptions.endpoint),
			randOption(cnrbtOptions.id),
			randOption(cnrbtOptions.instance),
			randOption(cnrbtOptions.interface_),
			randOption(cnrbtOptions.job),
			randOption(cnrbtOptions.metricsPath),
			randOption(cnrbtOptions.name),
			randOption(cnrbtOptions.namespace),
			randOption(cnrbtOptions.node),
			randOption(cnrbtOptions.pod),
			randOption(cnrbtOptions.prometheus),
			randOption(cnrbtOptions.service),
			rand.Float32(),
		}
	}
}
