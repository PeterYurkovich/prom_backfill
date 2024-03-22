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
		{
			label:     "endpoint",
			generator: func() (string, error) { return "endpoint2", nil },
		},
	})

	return root_node
}
