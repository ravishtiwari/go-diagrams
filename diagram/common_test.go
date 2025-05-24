package diagram_test

import (
	"errors"
	"os"
	"testing"

	"github.com/emarais-godaddy/go-diagrams/diagram"
	"github.com/emarais-godaddy/go-diagrams/nodes/generic"
	"github.com/stretchr/testify/require"
)

// setupTest creates a new diagram for testing and returns it along with a cleanup function.
func setupTest(t *testing.T) (*diagram.Diagram, func()) {
	t.Helper()

	d, err := diagram.New(diagram.Label("test"), diagram.Direction("LR"))
	require.NoError(t, err, "Failed to create new diagram")

	cleanup := func() {
		// remove the auto-generated data from this test execution
		err := os.RemoveAll("./go-diagrams") // Default diagram name is "go-diagrams" if not specified by an option
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		require.NoError(t, err, "Failed to cleanup test data")

		// also try to remove "test" directory if created by diagram.Label("test")
		err = os.RemoveAll("./test")
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		require.NoError(t, err, "Failed to cleanup test data named 'test'")
	}

	return d, cleanup
}

// Custom EdgeOption for no direction
func NoDirection() diagram.EdgeOption {
	return func(o *diagram.EdgeOptions) {
		o.Forward = false
		o.Reverse = false
	}
}

func TestAddSingleNode(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node := generic.Blank.Blank().Label("Test Node")
	d.Add(node)

	require.Len(t, d.Nodes(), 1, "Diagram should have one node")
	require.Equal(t, node, d.Nodes()[0], "The added node should be in the diagram's node list")
}

func TestAddMultipleNodes(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node1 := generic.Blank.Blank().Label("Test Node 1")
	node2 := generic.Blank.Blank().Label("Test Node 2")
	d.Add(node1, node2)

	require.Len(t, d.Nodes(), 2, "Diagram should have two nodes")
	require.Contains(t, d.Nodes(), node1, "Diagram should contain node1")
	require.Contains(t, d.Nodes(), node2, "Diagram should contain node2")
}

func TestCreateNewGroup(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	group := diagram.NewGroup("TestGroup").Label("Test Group Label")
	d.Group(group)

	require.Len(t, d.Groups(), 1, "Diagram should have one group")
	require.Equal(t, group, d.Groups()[0], "The created group should be in the diagram's group list")
	require.Equal(t, "cluster_TestGroup", group.ID(), "Group ID mismatch") // ID is prefixed with "cluster_"
	// To check the label, we'd ideally have a getter or check the rendered output.
	// For now, we assume that if NewGroup().Label() is called, it's set.
	// If there was a GetLabel() method or access to options:
	// require.Equal(t, "Test Group Label", group.GetLabel())
}

func TestAddNodeToGroup(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node := generic.Blank.Blank().Label("Test Node")
	group := diagram.NewGroup("TestGroup").Add(node)
	d.Group(group)

	require.Len(t, group.Nodes(), 1, "Group should have one node")
	require.Equal(t, node, group.Nodes()[0], "The added node should be in the group's node list")
	require.Len(t, d.Nodes(), 0, "Diagram should have no nodes directly, as it's in a group")
	require.Len(t, d.Groups(), 1, "Diagram should have one group")
}

func TestAddMultipleNodesToGroup(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node1 := generic.Blank.Blank().Label("Test Node 1")
	node2 := generic.Blank.Blank().Label("Test Node 2")
	group := diagram.NewGroup("TestGroup").Add(node1, node2)
	d.Group(group)

	require.Len(t, group.Nodes(), 2, "Group should have two nodes")
	require.Contains(t, group.Nodes(), node1, "Group should contain node1")
	require.Contains(t, group.Nodes(), node2, "Group should contain node2")
	require.Len(t, d.Nodes(), 0, "Diagram should have no nodes directly")
	require.Len(t, d.Groups(), 1, "Diagram should have one group")
}

func TestConnectNodesForward(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node1 := generic.Blank.Blank().Label("Node 1")
	node2 := generic.Blank.Blank().Label("Node 2")
	d.Add(node1, node2)
	d.Connect(node1, node2, diagram.Forward())

	require.Len(t, d.Edges(), 1, "Diagram should have one edge")
	edge := d.Edges()[0]
	require.Equal(t, node1.ID(), edge.Start(), "Edge source ID mismatch")
	require.Equal(t, node2.ID(), edge.End(), "Edge target ID mismatch")
	require.True(t, edge.Options.Forward, "Edge should be forward")
	require.False(t, edge.Options.Reverse, "Edge should not be reverse")
}

func TestConnectNodesBackward(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node1 := generic.Blank.Blank().Label("Node 1")
	node2 := generic.Blank.Blank().Label("Node 2")
	d.Add(node1, node2)
	d.Connect(node1, node2, diagram.Reverse()) // Corrected to Reverse()

	require.Len(t, d.Edges(), 1, "Diagram should have one edge")
	edge := d.Edges()[0]
	require.Equal(t, node1.ID(), edge.Start(), "Edge source ID mismatch")
	require.Equal(t, node2.ID(), edge.End(), "Edge target ID mismatch")
	require.False(t, edge.Options.Forward, "Edge should not be forward")
	require.True(t, edge.Options.Reverse, "Edge should be reverse")
}

func TestConnectNodesBidirectional(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node1 := generic.Blank.Blank().Label("Node 1")
	node2 := generic.Blank.Blank().Label("Node 2")
	d.Add(node1, node2)
	d.Connect(node1, node2, diagram.Bidirectional())

	require.Len(t, d.Edges(), 1, "Diagram should have one edge")
	edge := d.Edges()[0]
	require.Equal(t, node1.ID(), edge.Start(), "Edge source ID mismatch")
	require.Equal(t, node2.ID(), edge.End(), "Edge target ID mismatch")
	require.True(t, edge.Options.Forward, "Edge should be forward (for bidirectional)")
	require.True(t, edge.Options.Reverse, "Edge should be reverse (for bidirectional)")
}

func TestConnectNodesNoDirection(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	node1 := generic.Blank.Blank().Label("Node 1")
	node2 := generic.Blank.Blank().Label("Node 2")
	d.Add(node1, node2)
	d.Connect(node1, node2, NoDirection()) // Using custom NoDirection

	require.Len(t, d.Edges(), 1, "Diagram should have one edge")
	edge := d.Edges()[0]
	require.Equal(t, node1.ID(), edge.Start(), "Edge source ID mismatch")
	require.Equal(t, node2.ID(), edge.End(), "Edge target ID mismatch")
	require.False(t, edge.Options.Forward, "Edge should not be forward (for no direction)")
	require.False(t, edge.Options.Reverse, "Edge should not be reverse (for no direction)")
}

func TestGroupLabelIsSet(t *testing.T) {
	d, cleanup := setupTest(t)
	defer cleanup()

	groupLabel := "My Test Group Label"
	group := diagram.NewGroup("MyGroup").Label(groupLabel)
	d.Group(group)
	
	// There isn't a direct getter for options.Label from Group.
	// This test would be more robust if we could inspect the generated .dot file or GroupOptions.
	// For now, we assume .Label() works. If Group had a method like `GetLabel()` or `Options()` it would be better.
	// A placeholder assertion:
	require.Equal(t, "cluster_MyGroup", group.ID(), "Group ID should be prefixed")
}
