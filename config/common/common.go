package common

import (
	"bytes"
	"encoding/json"
	"strings"

	vyaml "github.com/coreos/vcontext/yaml"
	"github.com/coreos/vcontext/tree"
	"gopkg.in/yaml.v3"
)

type TranslateOptions struct {
	Pretty bool
	Strict bool
}

type Common struct {
	Version string `yaml:"version"`
	Variant string `yaml:"variant"`
}

// Misc helpers
func Unmarshal(data []byte, to interface{}, strict bool) (tree.Node, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(strict)
	if err := dec.Decode(to); err != nil {
		return nil, err
	}
	return vyaml.UnmarshalToContext(data)
}

func Marshal(from interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(from, "", "  ")
	}
	return json.Marshal(from)
}

func camel(in string) string {
	words := strings.Split(in, "_")
	for i, word := range words[1:] {
		words[i+1] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func ToCamelCase(t tree.Node) tree.Node {
	switch n := t.(type) {
	case tree.MapNode:
		m := tree.MapNode{
			Children:  make(map[string]tree.Node, len(n.Children)),
			Keys:      make(map[string]tree.Leaf, len(n.Keys)),
			Marker:    n.Marker,
		}
		for k, v := range n.Children {
			m.Children[camel(k)] = ToCamelCase(v)
		}
		for k, v := range n.Keys {
			m.Keys[camel(k)] = v
		}
		return m
	case tree.SliceNode:
		s := tree.SliceNode{
			Children: make([]tree.Node, 0, len(n.Children)),
			Marker:   n.Marker,
		}
		for _, v := range n.Children {
			s.Children = append(s.Children, ToCamelCase(v))
		}
		return s
	default: // leaf
		return t
	}
}
