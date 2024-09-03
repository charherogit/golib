package helper

// 维护玩家排行榜，支持插入、删除、查找、更新、排序操作
type Node struct {
	Val   int // 玩家数据
	Left  *Node
	Right *Node
}

func NewNode(val int) *Node {
	return &Node{Val: val}
}

func (n *Node) Insert(val int) {
	if val <= n.Val {
		if n.Left == nil {
			n.Left = NewNode(val)
		} else {
			n.Left.Insert(val)
		}
	} else {
		if n.Right == nil {
			n.Right = NewNode(val)
		} else {
			n.Right.Insert(val)
		}
	}
}

func (n *Node) Search(val int) bool {
	if n == nil {
		return false
	}
	if val < n.Val {
		return n.Left.Search(val)
	} else if val > n.Val {
		return n.Right.Search(val)
	}
	return true
}

func (n *Node) Update(val int) {
	if n == nil {
		return
	}
	if val < n.Val {
		n.Left.Update(val)
	} else if val > n.Val {
		n.Right.Update(val)
	}
	n.Val = val
}

func (n *Node) Delete(val int) *Node {
	if n == nil {
		return nil
	}
	if val < n.Val {
		n.Left = n.Left.Delete(val)
	} else if val > n.Val {
		n.Right = n.Right.Delete(val)
	} else {
		if n.Left == nil {
			return n.Right
		} else if n.Right == nil {
			return n.Left
		}
		min := n.Right.findMin()
		n.Val = min.Val
		n.Right = n.Right.Delete(min.Val)
	}
	return n
}

func (n *Node) findMin() *Node {
	for n.Left != nil {
		n = n.Left
	}
	return n
}

func (n *Node) InOrder() []int {
	if n == nil {
		return nil
	}
	var res []int
	res = append(res, n.Left.InOrder()...)
	res = append(res, n.Val)
	res = append(res, n.Right.InOrder()...)
	return res
}

func (n *Node) PreOrder() []int {
	if n == nil {
		return nil
	}
	var res []int
	res = append(res, n.Val)
	res = append(res, n.Left.PreOrder()...)
	res = append(res, n.Right.PreOrder()...)
	return res
}

func (n *Node) PostOrder() []int {
	if n == nil {
		return nil
	}
	var res []int
	res = append(res, n.Left.PostOrder()...)
	res = append(res, n.Right.PostOrder()...)
	res = append(res, n.Val)
	return res
}

func SortSlice(s []int) []int {
	var root *Node
	for _, v := range s {
		if root == nil {
			root = NewNode(v)
		} else {
			root.Insert(v)
		}
	}
	return root.InOrder()
}

type AVLNode struct {
	Val    int
	Left   *AVLNode
	Right  *AVLNode
	Height int
}

func NewAVLNode(val int) *AVLNode {
	return &AVLNode{Val: val}
}

func (a *AVLNode) Insert(val int) *AVLNode {
	if a == nil {
		return NewAVLNode(val)
	}
	if val < a.Val {
		a.Left = a.Left.Insert(val)
	} else {
		a.Right = a.Right.Insert(val)
	}
	return a.reBalance()
}

func (a *AVLNode) Delete(val int) *AVLNode {
	if a == nil {
		return nil
	}
	if val < a.Val {
		a.Left = a.Left.Delete(val)
	} else if val > a.Val {
		a.Right = a.Right.Delete(val)
	} else {
		if a.Left == nil {
			return a.Right
		} else if a.Right == nil {
			return a.Left
		}
		min := a.Right.findMin()
		a.Val = min.Val
		a.Right = a.Right.Delete(min.Val)
	}
	return a.reBalance()
}

func (a *AVLNode) findMin() *AVLNode {
	for a.Left != nil {
		a = a.Left
	}
	return a
}

func (a *AVLNode) reBalance() *AVLNode {
	a.Height = max(a.Left.Height, a.Right.Height) + 1
	balance := a.Left.Height - a.Right.Height
	if balance > 1 {
		if a.Left.Left.Height < a.Left.Right.Height {
			a.Left = a.Left.leftRotate()
		}
		a = a.rightRotate()
	} else if balance < -1 {
		if a.Right.Left.Height > a.Right.Right.Height {
			a.Right = a.Right.rightRotate()
		}
		a = a.leftRotate()
	}
	return a
}

func (a *AVLNode) leftRotate() *AVLNode {
	root := a.Right
	a.Right = root.Left
	root.Left = a
	a.Height = max(a.Left.Height, a.Right.Height) + 1
	root.Height = max(root.Left.Height, root.Right.Height) + 1
	return root
}

func (a *AVLNode) rightRotate() *AVLNode {
	root := a.Left
	a.Left = root.Right
	root.Right = a
	a.Height = max(a.Left.Height, a.Right.Height) + 1
	root.Height = max(root.Left.Height, root.Right.Height) + 1
	return root
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (a *AVLNode) InOrder() []int {
	if a == nil {
		return nil
	}
	var res []int
	res = append(res, a.Left.InOrder()...)
	res = append(res, a.Val)
	res = append(res, a.Right.InOrder()...)
	return res
}

func (a *AVLNode) PreOrder() []int {
	if a == nil {
		return nil
	}
	var res []int
	res = append(res, a.Val)
	res = append(res, a.Left.PreOrder()...)
	res = append(res, a.Right.PreOrder()...)
	return res
}

func (a *AVLNode) PostOrder() []int {
	if a == nil {
		return nil
	}
	var res []int
	res = append(res, a.Left.PostOrder()...)
	res = append(res, a.Right.PostOrder()...)
	res = append(res, a.Val)
	return res
}

func SortSliceAVL(s []int) []int {
	var root *AVLNode
	for _, v := range s {
		if root == nil {
			root = NewAVLNode(v)
		} else {
			root.Insert(v)
		}
	}
	return root.InOrder()
}
