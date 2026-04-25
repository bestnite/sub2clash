package common

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// BuiltSub 保存最终输出所需的完整 YAML 树。
//
// 这里刻意不再保存整份 typed 配置副本：
// - root 是整个转换流程的最终产物
// - 所有常规输出都直接从 root 序列化
// - nodeList 模式也从 root 中提取 proxies，而不是依赖额外状态
type BuiltSub struct {
	root *yaml.Node
}

// MarshalYAML 让 BuiltSub 在输出时直接复用 patch 后的 YAML 树，
// 从而避免再次经过 struct round-trip 丢失未知字段。
func (b *BuiltSub) MarshalYAML() (any, error) {
	if b == nil || b.root == nil {
		return nil, nil
	}
	if b.root.Kind == yaml.DocumentNode {
		if len(b.root.Content) == 0 {
			return nil, nil
		}
		return b.root.Content[0], nil
	}
	return b.root, nil
}

// MarshalNodeListYAML 从最终 YAML 树中提取 proxies 节点，构造 nodeList 模式输出。
// 这样 nodeList 也直接复用最终 root，而不是依赖额外的 typed struct 副本。
func (b *BuiltSub) MarshalNodeListYAML() ([]byte, error) {
	if b == nil || b.root == nil {
		return yaml.Marshal(&yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"})
	}

	proxiesNode, err := GetYAMLPath(b.root, "proxies")
	if err != nil {
		return nil, err
	}

	root := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	if proxiesNode != nil && !isNullYAMLNode(proxiesNode) {
		setMappingValue(root, "proxies", cloneYAMLNode(proxiesNode))
	}

	return yaml.Marshal(root)
}

// ParseYAMLDocument 把原始 YAML 解析成 DocumentNode，
// 并确保根内容最终是一个可写入的 mapping 节点。
func ParseYAMLDocument(data []byte) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if _, err := rootMappingNode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// HasYAMLPath 判断某个点路径是否存在。
// 这里仅关心“是否找到节点”，不关心节点具体类型。
func HasYAMLPath(doc *yaml.Node, path string) bool {
	current, err := GetYAMLPath(doc, path)
	return err == nil && current != nil
}

// GetYAMLPath 按 a.b.c 这种点路径向下查找节点。
// 当前实现只支持 mapping 之间的逐层下钻，不处理数组索引路径。
func GetYAMLPath(doc *yaml.Node, path string) (*yaml.Node, error) {
	segments := splitYAMLPath(path)
	if len(segments) == 0 {
		return nil, fmt.Errorf("yaml path is empty")
	}

	current, err := rootMappingNode(doc)
	if err != nil {
		return nil, err
	}

	for _, segment := range segments {
		next := findMappingValue(current, segment)
		if next == nil {
			return nil, nil
		}
		current = next
	}

	return current, nil
}

// SetYAMLPath 按点路径写入一个值；不存在的中间层会自动补成 mapping。
// 例如 a.b.c=1 会在缺失时依次创建 a 和 b 两层对象节点。
func SetYAMLPath(doc *yaml.Node, path string, value any) error {
	segments := splitYAMLPath(path)
	if len(segments) == 0 {
		return fmt.Errorf("yaml path is empty")
	}

	current, err := rootMappingNode(doc)
	if err != nil {
		return err
	}

	for idx, segment := range segments[:len(segments)-1] {
		next := findMappingValue(current, segment)
		if next == nil {
			next = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			setMappingValue(current, segment, next)
		}
		if next.Kind != yaml.MappingNode {
			return fmt.Errorf("yaml path %q segment %q is not a mapping", path, strings.Join(segments[:idx+1], "."))
		}
		current = next
	}

	encoded, err := encodeYAMLNode(value)
	if err != nil {
		return err
	}

	setMappingValue(current, segments[len(segments)-1], encoded)
	return nil
}

// EnsureYAMLSequencePath 确保某个路径最终是 sequence（YAML 数组）节点。
// 不存在时会自动创建，已存在但类型不匹配时返回错误。
func EnsureYAMLSequencePath(doc *yaml.Node, path string) (*yaml.Node, error) {
	return ensureYAMLPathKind(doc, path, yaml.SequenceNode, "!!seq")
}

// EnsureYAMLMappingPath 确保某个路径最终是 mapping（YAML 对象）节点。
func EnsureYAMLMappingPath(doc *yaml.Node, path string) (*yaml.Node, error) {
	return ensureYAMLPathKind(doc, path, yaml.MappingNode, "!!map")
}

// SetYAMLMappingField 在一个 mapping 节点里设置单个字段。
// 它等价于“在当前对象上写 key: value”。
func SetYAMLMappingField(node *yaml.Node, key string, value any) error {
	if node == nil || node.Kind != yaml.MappingNode {
		return fmt.Errorf("yaml node is not a mapping")
	}

	encoded, err := encodeYAMLNode(value)
	if err != nil {
		return err
	}

	setMappingValue(node, key, encoded)
	return nil
}

// AppendYAMLSequenceValue 向 sequence 节点末尾追加一个元素。
func AppendYAMLSequenceValue(node *yaml.Node, value any) error {
	if node == nil || node.Kind != yaml.SequenceNode {
		return fmt.Errorf("yaml node is not a sequence")
	}

	encoded, err := encodeYAMLNode(value)
	if err != nil {
		return err
	}

	node.Content = append(node.Content, encoded)
	return nil
}

// FindYAMLSequenceMappingByStringField 在 YAML 数组中查找一个对象元素，
// 要求该对象存在指定字段且字段值等于目标字符串。
//
// 例如在 proxy-groups 里按 name 查找：
//   - name: 节点选择
//     type: select
func FindYAMLSequenceMappingByStringField(node *yaml.Node, field string, value string) *yaml.Node {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}

	for _, item := range node.Content {
		if item == nil || item.Kind != yaml.MappingNode {
			continue
		}
		fieldNode := findMappingValue(item, field)
		if fieldNode == nil || fieldNode.Kind != yaml.ScalarNode {
			continue
		}
		if fieldNode.Value == value {
			return item
		}
	}

	return nil
}

// splitYAMLPath 把 a.b.c 这种点路径拆成 [a b c]。
// 空片段会被忽略，避免出现连续点号时产生无意义路径段。
func splitYAMLPath(path string) []string {
	parts := strings.Split(path, ".")
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		segments = append(segments, part)
	}
	return segments
}

// ensureYAMLPathKind 是 EnsureYAMLSequencePath / EnsureYAMLMappingPath 的底层实现。
// 它会：
// 1. 逐层确保中间节点存在且都是 mapping
// 2. 确保最后一个节点存在，且类型符合预期
func ensureYAMLPathKind(doc *yaml.Node, path string, kind yaml.Kind, tag string) (*yaml.Node, error) {
	segments := splitYAMLPath(path)
	if len(segments) == 0 {
		return nil, fmt.Errorf("yaml path is empty")
	}

	current, err := rootMappingNode(doc)
	if err != nil {
		return nil, err
	}

	// 跳过最后一个元素在后面处理
	for idx, segment := range segments[:len(segments)-1] {
		next := findMappingValue(current, segment)
		if next == nil || isNullYAMLNode(next) {
			next = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			setMappingValue(current, segment, next)
		}
		if next.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("yaml path %q segment %q is not a mapping", path, strings.Join(segments[:idx+1], "."))
		}
		current = next
	}

	lastSegment := segments[len(segments)-1]
	node := findMappingValue(current, lastSegment)
	if node == nil || isNullYAMLNode(node) {
		node = &yaml.Node{Kind: kind, Tag: tag}
		setMappingValue(current, lastSegment, node)
	}
	if node.Kind != kind {
		return nil, fmt.Errorf("yaml path %q is not a %s", path, yamlKindName(kind))
	}

	return node, nil
}

// rootMappingNode 统一把“文档根”整理成一个可操作的 mapping 节点。
//
// yaml.v3 通常把整份 YAML 包在 DocumentNode 下，真正的内容位于 Content[0]。
// 当前项目的 patch 逻辑都假定最外层是 key-value 结构，因此这里会：
// 1. 处理空文档
// 2. 取出 DocumentNode 的实际根内容
// 3. 确保该根内容是 mapping
func rootMappingNode(doc *yaml.Node) (*yaml.Node, error) {
	if doc == nil {
		return nil, fmt.Errorf("yaml document is nil")
	}

	root := doc
	if doc.Kind == 0 {
		doc.Kind = yaml.DocumentNode
	}
	if doc.Kind == yaml.DocumentNode {
		if len(doc.Content) == 0 {
			doc.Content = append(doc.Content, &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"})
		}
		root = doc.Content[0]
	}

	if root.Kind == 0 {
		root.Kind = yaml.MappingNode
		root.Tag = "!!map"
	}
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("yaml root must be a mapping node")
	}

	return root, nil
}

// isNullYAMLNode 判断一个节点是否为空/未初始化/null。
// 这让我们在“路径不存在”和“路径存在但值为 null”时都能按缺失处理。
func isNullYAMLNode(node *yaml.Node) bool {
	if node == nil {
		return true
	}
	if node.Kind == 0 {
		return true
	}
	return node.Kind == yaml.ScalarNode && node.Tag == "!!null"
}

// yamlKindName 仅用于生成更可读的错误信息。
func yamlKindName(kind yaml.Kind) string {
	switch kind {
	case yaml.MappingNode:
		return "mapping"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.DocumentNode:
		return "document"
	default:
		return "node"
	}
}

// findMappingValue 在 mapping 节点中按 key 查找对应的 value 节点。
//
// 需要注意：yaml.v3 的 MappingNode.Content 不是 map，而是交替存储：
// [key1, value1, key2, value2, ...]
// 所以这里每次 idx += 2，依次跳过一个完整的 key-value 对。
func findMappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for idx := 0; idx+1 < len(node.Content); idx += 2 {
		if node.Content[idx].Value == key {
			return node.Content[idx+1]
		}
	}
	return nil
}

// setMappingValue 在 mapping 节点中设置 key 对应的 value。
// 如果 key 已存在，就原位替换；否则在末尾追加一组新的 key-value。
func setMappingValue(node *yaml.Node, key string, value *yaml.Node) {
	for idx := 0; idx+1 < len(node.Content); idx += 2 {
		if node.Content[idx].Value == key {
			node.Content[idx+1] = value
			return
		}
	}

	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		value,
	)
}

// encodeYAMLNode 把普通 Go 值编码成 *yaml.Node，方便统一塞回 YAML 树。
// 如果 Encode 产生的是 DocumentNode，这里会自动取出它的实际内容节点。
func encodeYAMLNode(value any) (*yaml.Node, error) {
	var node yaml.Node
	if err := node.Encode(value); err != nil {
		return nil, err
	}
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}, nil
		}
		return node.Content[0], nil
	}
	return &node, nil
}

// cloneYAMLNode 深拷贝一个节点树，避免把同一个子树同时挂到多个输出根下。
func cloneYAMLNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}

	clone := *node
	if len(node.Content) != 0 {
		clone.Content = make([]*yaml.Node, len(node.Content))
		for i := range node.Content {
			clone.Content[i] = cloneYAMLNode(node.Content[i])
		}
	}
	return &clone
}
