package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

func mergeMapsNode(targetNode, sourceNode *yaml.Node) {
	mergedKeys := make(map[string]bool)
	for i, targetContent := range targetNode.Content {
		if i%2 == 0 {
			targetKey := targetContent.Value
			for j, sourceContent := range sourceNode.Content {
				if j%2 == 0 {
					sourceKey := sourceContent.Value
					if targetKey == sourceKey {
						if targetNode.Content[i+1].Kind == yaml.MappingNode && sourceNode.Content[j+1].Kind == yaml.MappingNode {
							mergeMapsNode(targetNode.Content[i+1], sourceNode.Content[j+1])
						} else {
							targetNode.Content[i+1] = sourceNode.Content[j+1]
						}
						mergedKeys[sourceKey] = true
						break
					}
				}
			}
		}
	}

	for i := 0; i < len(sourceNode.Content); i += 2 {
		sourceKey := sourceNode.Content[i].Value
		if !mergedKeys[sourceKey] {
			targetNode.Content = append(targetNode.Content, sourceNode.Content[i], sourceNode.Content[i+1])
		}
	}
}

func MergeYAMLFromString(path string, sourceYAMLString string) error {
	targetYAMLContent, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var targetNode yaml.Node
	if len(targetYAMLContent) == 0 {
		targetNode = yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{Kind: yaml.MappingNode}}}
	} else {
		err = yaml.Unmarshal(targetYAMLContent, &targetNode)
		if err != nil {
			return err
		}
	}

	var sourceNode yaml.Node
	err = yaml.Unmarshal([]byte(sourceYAMLString), &sourceNode)
	if err != nil {
		return err
	}

	if targetNode.Kind == yaml.DocumentNode && sourceNode.Kind == yaml.DocumentNode {
		if len(targetNode.Content) > 0 && len(sourceNode.Content) > 0 {
			if targetNode.Content[0].Kind == yaml.MappingNode && sourceNode.Content[0].Kind == yaml.MappingNode {
				mergeMapsNode(targetNode.Content[0], sourceNode.Content[0])
			}
		}
	}

	mergedYAMLContent, err := yaml.Marshal(&targetNode)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, mergedYAMLContent, 0644)
	return err
}
