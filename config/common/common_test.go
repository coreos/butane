package common

import (
	"reflect"
	"testing"

	"github.com/coreos/vcontext/tree"
)

func TestCamel(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{},
		{
			"foo",
			"foo",
		},
		{
			"snake_case",
			"snakeCase",
		},
		{
			"long_snake_case",
			"longSnakeCase",
		},
		{
			"camelAlready",
			"camelAlready",
		},
	}

	for i, test := range tests {
		if camel(test.in) != test.out {
			t.Errorf("#%d: expected %q got %q", i, test.out, camel(test.in))
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		in  tree.Node
		out tree.Node
	}{
		{},
		{
			tree.Leaf{
				Marker: tree.MarkerFromIndices(1, 2),
			},
			tree.Leaf{
				Marker: tree.MarkerFromIndices(1, 2),
			},
		},
		{
			tree.MapNode{
				Marker: tree.MarkerFromIndices(1, 2),
				Children: map[string]tree.Node{
					"foo_bar": tree.Leaf{
						tree.MarkerFromIndices(3, 4),
					},
				},
				Keys: map[string]tree.Leaf{
					"foo_bar": tree.Leaf{
						tree.MarkerFromIndices(3, 4),
					},
				},
			},
			tree.MapNode{
				Marker: tree.MarkerFromIndices(1, 2),
				Children: map[string]tree.Node{
					"fooBar": tree.Leaf{
						tree.MarkerFromIndices(3, 4),
					},
				},
				Keys: map[string]tree.Leaf{
					"fooBar": tree.Leaf{
						tree.MarkerFromIndices(3, 4),
					},
				},
			},
		},
		{
			tree.SliceNode{
				Marker: tree.MarkerFromIndices(5, 6),
				Children: []tree.Node{
					tree.MapNode{
						Marker: tree.MarkerFromIndices(1, 2),
						Children: map[string]tree.Node{
							"foo_bar": tree.Leaf{
								tree.MarkerFromIndices(3, 4),
							},
						},
						Keys: map[string]tree.Leaf{
							"foo_bar": tree.Leaf{
								tree.MarkerFromIndices(3, 4),
							},
						},
					},
				},
			},
			tree.SliceNode{
				Marker: tree.MarkerFromIndices(5, 6),
				Children: []tree.Node{
					tree.MapNode{
						Marker: tree.MarkerFromIndices(1, 2),
						Children: map[string]tree.Node{
							"fooBar": tree.Leaf{
								tree.MarkerFromIndices(3, 4),
							},
						},
						Keys: map[string]tree.Leaf{
							"fooBar": tree.Leaf{
								tree.MarkerFromIndices(3, 4),
							},
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		actual := ToCamelCase(test.in)
		if !reflect.DeepEqual(actual, test.out) {
			t.Errorf("#%d: expected %+v got %+v", i, test.out, actual)
		}
	}
}
