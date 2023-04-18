package mappings

type MappingNode struct {
	TargetType Type
	Children   []*MappingNode
	Source     SourceMapping
}

type SourceMapping struct {
	Source string
	Mapped bool
}

func Find(base string, mappingTree []*MappingNode, isTarget func(node *MappingNode, path string) bool) (*MappingNode, string) {
	var result *MappingNode
	var resultPath string
	for _, mapNode := range mappingTree {
		mapNode.Inspect(base, func(fullPath string, node *MappingNode) bool {
			if isTarget(node, fullPath) {
				result = node
				resultPath = fullPath
				return false
			}

			return true
		})

		if result != nil {
			return result, resultPath
		}
	}

	return result, resultPath
}

func (m *MappingNode) GetNode() *MappingNode {
	return m
}

type InspectionFunc func(fullPath string, node *MappingNode) bool

func (m *MappingNode) Inspect(base string, inspect InspectionFunc) {
	var inspectRec func(base string, m *MappingNode, inspect InspectionFunc) bool
	inspectRec = func(base string, m *MappingNode, inspect InspectionFunc) bool {
		for i := range m.Children {
			if !inspect(base+m.TargetType.ArgumentName+"."+m.Children[i].TargetType.ArgumentName, m.Children[i]) {
				return false
			}
		}

		for i := range m.Children {
			if !inspectRec(base+m.TargetType.ArgumentName+".", m.Children[i], inspect) {
				return false
			}
		}

		return true
	}
	if !inspect(base+m.TargetType.ArgumentName, m) {
		return
	}

	inspectRec(base, m, inspect)
}
