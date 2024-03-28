package main

/**
 * container, endpoint, id, instance, interface, job, metrics_path, name, namespace, node, pod, prometheus, service, __name__
 */

func create_cnrbt_tree() *node {
	root_node := &node{
		label:     "__name__",
		generator: func(self *node) (string, error) { return "container_network_receive_bytes_total", nil },
		parent:    nil,
		children:  []*node{},
	}

	root_node.addNextLevel([]*node{
		{
			label:     "container",
			generator: func(self *node) (string, error) { return "POD	", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "endpoint",
			generator: func(self *node) (string, error) { return "https-metrics	", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label: "id",
			generator: func(self *node) (string, error) {
				return "/kubepods.slice/kubepods-burstable.slice/kubepods-burstable-podaaaaaaaaaaaaaaaaaaaaaaa.slice/aaaa-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil
			},
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "instance",
			generator: func(self *node) (string, error) { return "00.0.000.000:00000", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "interface",
			generator: func(self *node) (string, error) { return "eth0", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "job",
			generator: func(self *node) (string, error) { return "kubelet", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "metrics_path",
			generator: func(self *node) (string, error) { return "/metrics/cadvisor", nil },
		},
	})

	root_node.addNextLevel(
		createStaticNodes("pod", []string{
			"alertmanager-main-0",
			"alertmanager-main-1",
		}))

	root_node.addNextLevel([]*node{
		{
			label:     "namespace",
			generator: func(self *node) (string, error) { return "openshift-monitoring", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label: "name",
			generator: func(self *node) (string, error) {
				return "k8s_POD_alertmanager-main-0_openshift-monitoring_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa_a", nil
			},
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "namespace",
			generator: func(self *node) (string, error) { return "openshift-monitoring", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "node",
			generator: func(self *node) (string, error) { return "ip-00-0-000-000.ec2.internal", nil },
		},
		{
			label:     "node",
			generator: func(self *node) (string, error) { return "ip-00-0-000-001.ec2.internal", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "prometheus",
			generator: func(self *node) (string, error) { return "openshift-monitoring/k8s", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "service",
			generator: func(self *node) (string, error) { return "kubelet", nil },
		},
	})

	return root_node
}
