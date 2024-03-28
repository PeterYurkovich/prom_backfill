package main

import "github.com/prometheus/prometheus/model/labels"

/**
											(metric)
										/					\
					(label1="a")			(label1="b")
								|									\
	(label2=randLabel(label1))	(label2=randLabel(label1))

struct node {
	label: "__name__"
	generator: () => string, err
	parent: *node
	children: (*node)[]
}
copy() node, err

add_next_level((*node)[]) err
move to each leaf node, then create a copy of the node at this location


This function type will be used to fill in the blanks in a template
generatorFunction: () => string, err

# Helper functions -  All return () => string, err which can be placed onto a node
staticLabel(labelName)
randomChoice(choiceArray: string[])
randomTemplate(templateString: string (including {}), generatorFunction: () => string)
randomDependantChoice(dependantLabel: string, choiceArray: string[])
randomDependantTemplate(dependantLabel: string, templateString: string (including {}), generatorFunction: () => string)
*/

type node struct {
	label     string
	generator func(*node) (string, error)
	parent    *node
	children  []*node
}

/**
 * This function will be used in the addNextLevel function to replicate the node being inserted
 * to each leaf node. This function clears the children array and parent pointers.
 */
func copy(n *node) (*node, error) {
	if n == nil {
		return nil, nil
	}
	nCopy := *n
	nCopy.children = []*node{}
	nCopy.parent = nil
	return &nCopy, nil
}

/**
 * This function will be used to add the next level of nodes to the tree.
 *
 * @param {*node[]} nodes - The array of nodes to add to the leaf nodes.
 * @param *node - The parent node to iterate down the tree.
 * @returns {error} - An error if there is a problem adding the nodes.
 */
func (parent *node) addNextLevel(newChildren []*node) error {
	if parent == nil {
		return nil
	}
	if len(parent.children) == 0 {
		parent.children = []*node{}
		for _, newChild := range newChildren {
			newChildCopy, err := copy(newChild)
			if err != nil {
				return err
			}
			newChildCopy.parent = parent
			parent.children = append(parent.children, newChildCopy)
		}
		return nil
	}
	for _, child := range parent.children {
		if err := child.addNextLevel(newChildren); err != nil {
			return err
		}
	}
	return nil
}

/**
 * Retrieve a slice of series from the parent node's children.
 *
 * @param {*node} parent - The parent node to retrieve the label and value from.
 * @returns {[]labels.Labels} - A slice of labels containing the label and value.
 */
func (current *node) getSeries(parentValue labels.Labels) ([]labels.Labels, error) {
	var labelAndValues []labels.Labels = []labels.Labels{}
	var labelValueChain labels.Labels

	currentValue, err := current.generator(current)
	if err != nil {
		return nil, err
	}
	currentLabelAndValue := labels.Label{Name: current.label, Value: currentValue}
	if parentValue != nil {
		labelValueChain = append(parentValue, currentLabelAndValue)
	} else {
		labelValueChain = labels.Labels{currentLabelAndValue}
	}

	if len(current.children) == 0 {
		return []labels.Labels{labelValueChain}, nil
	}
	for _, child := range current.children {
		childLabelAndValues, err := child.getSeries(labelValueChain)
		if err != nil {
			return nil, err
		}
		labelAndValues = append(labelAndValues, childLabelAndValues...)
	}
	return labelAndValues, nil
}

func createStaticNodes(label string, values []string) []*node {
	nodes := []*node{}
	for _, value := range values {
		nodes = append(nodes, &node{
			label:     label,
			generator: func(*node) (string, error) { return value, nil },
		})
	}
	return nodes
}
