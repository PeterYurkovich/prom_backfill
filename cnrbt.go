package main

/**
 * container, endpoint, id, instance, interface, job, metrics_path, name, namespace, node, pod, prometheus, service, __name__
 */

func create_cnrbt_tree() *node {
	root_node := &node{
		label:     "__name__",
		generator: func() (string, error) { return "container_network_receive_bytes_total", nil },
		parent:    nil,
		children:  []*node{},
	}

	root_node.addNextLevel([]*node{
		{
			label:     "container",
			generator: func() (string, error) { return "container1", nil },
		},
		{
			label:     "container",
			generator: func() (string, error) { return "container2", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "endpoint",
			generator: func() (string, error) { return "endpoint1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "id",
			generator: func() (string, error) { return "id1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "instance",
			generator: func() (string, error) { return "instance1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "interface",
			generator: func() (string, error) { return "interface1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "job",
			generator: func() (string, error) { return "job1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "metrics_path",
			generator: func() (string, error) { return "metrics_path1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "name",
			generator: func() (string, error) { return "name1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "namespace",
			generator: func() (string, error) { return "namespace1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "node",
			generator: func() (string, error) { return "node1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "pod",
			generator: func() (string, error) { return "pod1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "prometheus",
			generator: func() (string, error) { return "prometheus1", nil },
		},
	})

	root_node.addNextLevel([]*node{
		{
			label:     "service",
			generator: func() (string, error) { return "service1", nil },
		},
		{
			label:     "service",
			generator: func() (string, error) { return "service2", nil },
		},
	})

	return root_node
}
