package xutil

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	EmptyError           = errors.New("queue is empty")
	IndexOutOfRangeError = errors.New("index out of range")
	ItemNotFoundError    = errors.New("item not found")
)

func NewQueue(n int) *Queue {
	if n <= 0 {
		n = 1024
	}
	return &Queue{
		step: n,
		list: make([]interface{}, n+1),
	}
}

type Queue struct {
	step   int
	debug  bool
	logger *logrus.Logger

	list []interface{}
	head int
	tail int
	mu   sync.RWMutex
}

// 获取队列容量
func (q *Queue) cap() int {
	return len(q.list)
}

// 获取队列长度
func (q *Queue) len() int {
	switch {
	case q.tail > q.head:
		return q.tail - q.head
	case q.tail < q.head:
		return q.tail + q.cap() - q.head
	default:
		return 0
	}
}

// 获取指定下标的前一个下标
func (q *Queue) previous(i int) int {
	return (i - 1 + q.cap()) % q.cap()
}

// 获取指定下标的后一个下标
func (q *Queue) next(i int) int {
	return (i + 1) % q.cap()
}

// 获取头下标的前一个下标
func (q *Queue) headPrevious() int {
	return q.previous(q.head)
}

// 获取头下标的后一个下标
func (q *Queue) headNext() int {
	return q.next(q.head)
}

// 获取尾下标的前一个下标
func (q *Queue) tailPrevious() int {
	return q.previous(q.tail)
}

// 获取尾下标的后一个下标
func (q *Queue) tailNext() int {
	return q.next(q.tail)
}

// 判断队列是否为空
func (q *Queue) isEmpty() bool {
	return q.head == q.tail
}

// 判断队列是否已满
func (q *Queue) isFull() bool {
	return q.tailNext() == q.head
}

// 获取指定下标的元素
func (q *Queue) get(i int) (interface{}, error) {
	if q.isEmpty() {
		return nil, EmptyError
	}
	if !q.isValidIndex(i) {
		return nil, IndexOutOfRangeError
	}
	return q.list[i], nil
}

// 判断下标是否合法
func (q *Queue) isValidIndex(i int) bool {
	if q.isEmpty() {
		return false
	}
	if q.head < q.tail {
		return i >= q.head && i < q.tail
	}
	if q.head > q.tail {
		return i >= q.head || i < q.tail
	}
	return false
}

// 从队列尾向队列中添加一个元素
func (q *Queue) push(item interface{}) {
	// 扩容
	q.grow()

	// 添加元素
	q.list[q.tail] = item
	q.tail = q.tailNext()

	q.printDebug("push")
}

// 队列数组扩容
func (q *Queue) grow() {
	if !q.isFull() {
		return
	}
	nl := make([]interface{}, q.cap()+q.step)
	q.migrate(nl)
}

// 队列数据迁移至新队列数组
func (q *Queue) migrate(dst []interface{}) {
	q.printDebug("migrate before")

	i, j := q.head, 0
	for ; i != q.tail; i, j = q.next(i), j+1 {
		dst[j] = q.list[i]
	}

	q.list = dst
	q.head = 0
	q.tail = j

	q.printDebug("migrate after")
}

// 从队列头弹出一个元素
func (q *Queue) pop() (interface{}, bool) {
	// 判断是否为空
	if q.isEmpty() {
		q.printDebug("pop nothing")
		return nil, false
	}

	// 缩容
	q.reduce()

	// 弹出元素
	item := q.list[q.head]
	q.list[q.head] = nil
	q.head = q.headNext()

	q.printDebug("pop")

	return item, true
}

// 缩容
func (q *Queue) reduce() {
	if q.cap()-q.len() < 6*q.step {
		return
	}
	nl := make([]interface{}, q.cap()-3*q.step)
	q.migrate(nl)
}

// 输出 debug 信息
func (q *Queue) printDebug(tag string) {
	if !q.debug {
		return
	}
	if q.logger == nil {
		fmt.Printf("[queue] %s => %s.\n", tag, q.status())
		return
	}
	q.logger.Debugf("[queue] %s => %s.", tag, q.status())
}

// 获取队列状态
func (q *Queue) status() string {
	return fmt.Sprintf("head: %-4d, tail: %-4d, len: %-4d, cap: %-4d", q.head, q.tail, q.len(), q.cap())
}

// 获取队列头元素
func (q *Queue) getHeadItem() (interface{}, error) {
	return q.get(q.head)
}

// 获取队列尾元素
func (q *Queue) getTailItem() (interface{}, error) {
	return q.get(q.tailPrevious())
}

// ************************* export *************************

// EnableDebug 开启 debug 模式，该模式下会输出队列的操作日志
func (q *Queue) EnableDebug() *Queue {
	q.debug = true
	return q
}

// SetLogger 设置自定义日志
func (q *Queue) SetLogger(logger *logrus.Logger) *Queue {
	q.logger = logger
	return q
}

// Get 获取指定下标的元素
func (q *Queue) Get(i int) (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.get(i)
}

// GetHeadItem 获取队列头元素
func (q *Queue) GetHeadItem() (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.getHeadItem()
}

// GetTailItem 获取队列尾元素
func (q *Queue) GetTailItem() (interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.getTailItem()
}

// Status 获取队列状态
func (q *Queue) Status() string {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.status()
}

// IsEmpty 判断队列是否为空
func (q *Queue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.isEmpty()
}

// Cap 获取队列容量
func (q *Queue) Cap() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.cap()
}

// Len 获取队列长度
func (q *Queue) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.len()
}

// Push 从队列尾向队列中添加一个元素
func (q *Queue) Push(b interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.push(b)
}

// Pop 从队列头弹出一个元素
func (q *Queue) Pop() (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.pop()
}

// SlideN 从队列尾向队列中添加一个元素，如果添加完元素队列长度大于 n，则截取（只保留）队列后 n 个元素。类似滑动窗口
func (q *Queue) SlideN(b interface{}, n int) (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 添加元素
	q.push(b)

	// 判断是否可以滑动
	if q.len() <= n {
		q.printDebug("slide deny")
		return nil, false
	}

	// 将队列长度控制在 n
	var ok bool
	var item interface{}
	for q.len() > n {
		item, ok = q.pop()
	}

	q.printDebug("slide allow")

	return item, ok
}

// SlideFunc 需要删除的元素返回 true，否则返回 false。
// 从队列头开始直到第一个不需要删除的元素出现，前面的元素全部删除
type SlideFunc func(item interface{}) bool

// Slide 从队列尾部添加一个元素，并滑动窗口
func (q *Queue) Slide(b interface{}, f SlideFunc) (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 添加元素
	q.push(b)

	// 将队列控制在指定条件内
	var ok bool
	var item interface{}
	for f(q.list[q.head]) {
		item, ok = q.pop()
	}

	q.printDebug("slide allow")

	return item, ok
}

// Walk 遍历队列
// reverse false：从头到尾遍历，true：从尾到头遍历
func (q *Queue) Walk(f func(item interface{}), reverse bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if reverse {
		for i := q.tailPrevious(); i != q.headPrevious(); i = q.previous(i) {
			f(q.list[i])
		}
	} else {
		for i := q.head; i != q.tail; i = q.next(i) {
			f(q.list[i])
		}
	}
}

// FindFunc 是符合条件的元素返回 true，否则返回 false
type FindFunc func(item interface{}) bool

// Find 遍历队列，返回第一个符合条件的元素
// reverse false：从头到尾遍历，true：从尾到头遍历
func (q *Queue) Find(f FindFunc, reverse bool) (int, interface{}, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.isEmpty() {
		return 0, nil, EmptyError
	}

	if reverse {
		for i := q.tailPrevious(); i != q.headPrevious(); i = q.previous(i) {
			if item, _ := q.get(i); f(item) {
				return i, item, nil
			}
		}
	} else {
		for i := q.head; i != q.tail; i = q.next(i) {
			if item, _ := q.get(i); f(item) {
				return i, item, nil
			}
		}
	}

	return 0, nil, ItemNotFoundError
}

// FindAll 遍历队列，返回全部符合条件的元素
func (q *Queue) FindAll(f FindFunc) []interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.isEmpty() {
		return []interface{}{}
	}

	all := make([]interface{}, 0)
	for i := q.head; i != q.tail; i = q.next(i) {
		if item, _ := q.get(i); f(item) {
			all = append(all, item)
		}
	}

	return all
}

// GetTerminalItemsN 获取队列前/后 n 个 item
func (q *Queue) GetTerminalItemsN(n int, reverse bool) []interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if n > q.len() {
		n = q.len()
	}

	items := make([]interface{}, 0, n)
	if reverse {
		for i, j := 0, q.tailPrevious(); i < n && j != q.headPrevious(); i, j = i+1, q.previous(j) {
			item, _ := q.get(j)
			items = append(items, item)
		}
	} else {
		for i, j := 0, q.head; i < n && j != q.tail; i, j = i+1, q.next(j) {
			item, _ := q.get(j)
			items = append(items, item)
		}
	}

	return items
}

// GetTerminalItems 获取队列前/后多个符合条件的 item，遇到第一个不符合条件的 item 停止遍历
func (q *Queue) GetTerminalItems(f FindFunc, reverse bool) []interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	items := make([]interface{}, 0)
	if reverse {
		for i := q.tailPrevious(); i != q.headPrevious(); i = q.previous(i) {
			item, _ := q.get(i)
			if !f(item) {
				break
			}
			items = append(items, item)
		}
	} else {
		for i := q.head; i != q.tail; i = q.next(i) {
			item, _ := q.get(i)
			if !f(item) {
				break
			}
			items = append(items, item)
		}
	}

	return items
}

// Window
// reverse false：从头到尾遍历，true：从尾到头遍历
// 返回结果包含 start item，不包含 stop item
func (q *Queue) Window(start FindFunc, stop FindFunc, reverse bool) ([]interface{}, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.isEmpty() {
		return []interface{}{}, false
	}

	isStart := false
	isStop := false
	window := make([]interface{}, 0)
	if reverse {
		for i := q.tailPrevious(); i != q.headPrevious(); i = q.previous(i) {
			item, _ := q.get(i)
			if !isStart && start(item) {
				isStart = true
			}
			if !isStop && stop(item) {
				isStop = true
			}
			if isStop {
				break
			}
			if isStart && !isStop {
				window = append(window, item)
			}
		}
	} else {
		for i := q.head; i != q.tail; i = q.next(i) {
			item, _ := q.get(i)
			if !isStart && start(item) {
				isStart = true
			}
			if isStart && !isStop {
				window = append(window, item)
			}
			if !isStop && stop(item) {
				isStop = true
			}
			if isStop {
				break
			}
		}
	}

	return window, isStart && isStop
}

// Reset 重置队列
func (q *Queue) Reset(data []interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	initCap := (len(data)/q.step+1)*q.step + 1

	list := make([]interface{}, initCap)
	for i, item := range data {
		list[i] = item
	}

	q.list = list
	q.head = 0
	q.tail = len(data)
}

// CopyItems 复制队列列表
func (q *Queue) CopyItems() []interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	cpy := make([]interface{}, 0, q.len())
	for i := q.head; i != q.tail; i = q.next(i) {
		cpy = append(cpy, q.list[i])
	}

	return cpy
}
