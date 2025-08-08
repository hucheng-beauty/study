package model

// ResourcePackCodeAttributeTree is the alias to mark the root node of code_attribute_tree.
type ResourcePackCodeAttributeTree Node

type Node struct {
    Name                 string               `json:"name"`
    Children             []*Node              `json:"children"`
    Details              map[string]any       `json:"details"`
    PatternMatchingRules PatternMatchingRules `json:"pattern_matching_rules"`
}

type PatternMatchingRules struct {
    IsNameIgnored bool `json:"is_name_ignored"`
    IsNameMust    bool `json:"is_name_must"`
    IsChildMust   bool `json:"is_child_must"`
}
