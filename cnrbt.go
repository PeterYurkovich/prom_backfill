package main

import (
	"math/rand"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

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

func randomizeCnrbtValue(exemplar *container_network_receive_bytes_total) func() *container_network_receive_bytes_total {

	return func() *container_network_receive_bytes_total {
		exemplar.Value = rand.Float32()
		return exemplar
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
