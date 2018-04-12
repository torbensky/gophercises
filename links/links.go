package links

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Find(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	tags := []Link{}
	var linkFn = func(n *html.Node) {
		if n.Data == "a" {
			tags = append(tags, addLink(n))
		}
	}

	// Find all anchors in a depth-first search
	NodeDFS(doc, html.ElementNode, linkFn)

	return tags, nil
}

func addLink(n *html.Node) Link {
	link := Link{}
	// Find the href attribute and get its value
	for _, a := range n.Attr {
		if a.Key == "href" {
			link.Href = a.Val
			break
		}
	}
	// Get the text from this node
	var buffer bytes.Buffer
	var textFn = func(n *html.Node) {
		buffer.WriteString(n.Data)
	}
	NodeDFS(n, html.TextNode, textFn)
	link.Text = buffer.String()

	return link
}

// NodeFn is a function that handles an HTML node
type NodeFn func(*html.Node)

// NodeDFS is a function that performs a depth first search looking for a specific node type.
// The node function is invoked on nodes of a matching type
func NodeDFS(n *html.Node, nt html.NodeType, fn NodeFn) {
	if n.Type == nt {
		fn(n)
	}
	// Depth first search of tree
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		NodeDFS(c, nt, fn)
	}
}
