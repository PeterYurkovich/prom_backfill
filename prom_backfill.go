package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
)

func main() {
	err := os.Mkdir("tsdb-test", 0700)
	noErr(err)

	db, err := tsdb.Open("tsdb-test", nil, nil, tsdb.DefaultOptions(), nil)
	noErr(err)

	app := db.Appender(context.Background())

	cnrbtCache, setCnrbtRef := getSeriesCache()

	startTime := time.Now().Unix() - 1000*60*60*10 // 10 hours ago
	endTime := time.Now().Unix() - 1000*60*60*9    // 9 hours ago
	randomCnrbtGenerator := randomCnrbt()

	for i := startTime; i < endTime; i += 1000 * 30 {
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
		}
	}

	err = app.Commit()
	noErr(err)

	// querier, err := db.Querier(math.MinInt64, math.MaxInt64)
	// noErr(err)
	// ss := querier.Select(context.Background(), false, nil, labels.MustNewMatcher(labels.MatchEqual, "container", "container1"))

	// for ss.Next() {
	// series := ss.At()
	// fmt.Println("series:", series.Labels().String())

	// it := series.Iterator(nil)
	// for it.Next() != chunkenc.ValNone {
	// _, v := it.At()
	// fmt.Println("sample", v)
	// }

	// fmt.Println("it.Err():", it.Err())
	// }
	// fmt.Println("ss.Err():", ss.Err())
	// ws := ss.Warnings()
	// if len(ws) > 0 {
	// fmt.Println("warnings:", ws)
	// }
	// err = querier.Close()
	// noErr(err)

	err = db.Close()
	noErr(err)

	err = os.RemoveAll("tsdb-test")
	noErr(err)
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

func getSeriesCache() (func(cnrbt *container_network_receive_bytes_total) (labels.Labels, storage.SeriesRef), func(cnrbt *container_network_receive_bytes_total, ref storage.SeriesRef)) {
	seriesCache := map[container_network_receive_bytes_total]labels.Labels{}
	refCache := map[container_network_receive_bytes_total]storage.SeriesRef{}

	return func(cnrbt *container_network_receive_bytes_total) (labels.Labels, storage.SeriesRef) {
			cachedSeries, ok := seriesCache[*cnrbt]
			if ok {
				cachedRef, ok := refCache[*cnrbt]
				if ok {
					return cachedSeries, cachedRef
				}
				return cachedSeries, 0
			}

			seriesBuilder := labels.NewScratchBuilder(13)

			seriesBuilder.Add("container", cnrbt.Container)
			seriesBuilder.Add("endpoint", cnrbt.Endpoint)
			seriesBuilder.Add("id", cnrbt.Id)
			seriesBuilder.Add("instance", cnrbt.Instance)
			seriesBuilder.Add("interface", cnrbt.Interface)
			seriesBuilder.Add("job", cnrbt.Job)
			seriesBuilder.Add("metrics_path", cnrbt.Metrics_Path)
			seriesBuilder.Add("name", cnrbt.Name)
			seriesBuilder.Add("namespace", cnrbt.Namespace)
			seriesBuilder.Add("node", cnrbt.Node)
			seriesBuilder.Add("pod", cnrbt.Pod)
			seriesBuilder.Add("prometheus", cnrbt.Prometheus)
			seriesBuilder.Add("service", cnrbt.Service)

			series := seriesBuilder.Labels()
			seriesCache[*cnrbt] = series

			return series, 0
		}, func(cnrbt *container_network_receive_bytes_total, ref storage.SeriesRef) {
			refCache[*cnrbt] = ref
		}
}
