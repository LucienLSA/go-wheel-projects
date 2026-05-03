package lru

import "fmt"

type Node struct {
	value int
	key   int
	prev  *Node
	next  *Node
}

func NewNode(key int, value int) *Node {
	return &Node{
		key:   key,
		value: value,
	}
}

type LRUCache struct {
	capacity    int
	dummy       *Node
	key_to_node map[int]*Node
}

func Constructor(capacity int) LRUCache {
	dummy := NewNode(0, 0)
	dummy.next = dummy
	dummy.prev = dummy
	return LRUCache{
		capacity:    capacity,
		dummy:       dummy,
		key_to_node: make(map[int]*Node),
	}
}

func (this *LRUCache) remove(x *Node) {
	x.prev.next = x.next
	x.next.prev = x.prev
}

func (this *LRUCache) get_node(key int) *Node {
	if node, ok := this.key_to_node[key]; ok {
		this.remove(node)
		this.put_front(node)
		return node
	}
	return nil
}

func (this *LRUCache) put_front(x *Node) {
	x.prev = this.dummy
	x.next = this.dummy.next
	x.prev.next = x
	x.next.prev = x
}

func (this *LRUCache) Get(key int) int {
	node := this.get_node(key)
	if node != nil {
		return node.value
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	node := this.get_node(key)
	if node != nil {
		node.value = value
		return
	}
	node = NewNode(key, value)
	this.key_to_node[key] = node
	this.put_front(node)
	if len(this.key_to_node) > this.capacity {
		back_node := this.dummy.prev
		delete(this.key_to_node, back_node.key)
		this.remove(back_node)
	}
}

func main() {
	capacity := 3
	obj := Constructor(capacity)
	obj.Put(1, 1)
	obj.Put(2, 2)
	param_1 := obj.Get(1) // 预期输出 1
	fmt.Println(param_1)

	// 额外测试：满容量淘汰
	obj.Put(3, 3)
	obj.Put(4, 4)           // 淘汰 2
	fmt.Println(obj.Get(2)) // 输出 -1
}
