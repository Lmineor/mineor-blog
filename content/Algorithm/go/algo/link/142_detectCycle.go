package link

//题意： 给定一个链表，返回链表开始入环的第一个节点。 如果链表无环，则返回 null。
//
//为了表示给定链表中的环，使用整数 pos 来表示链表尾连接到链表中的位置（索引从 0 开始）。 如果 pos 是 -1，则在该链表中没有环。
//
//说明：不允许修改给定的链表。

func detectCycle(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	if head.Next == nil {
		return nil
	}
	slow, fast := head, head
	for fast != nil && fast.Next != nil{
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast{
			break
		}
	}
	if fast == nil || fast.Next == nil{
		return nil
	}
	slow = head
	for slow != fast{
		slow = slow.Next
		fast = fast.Next
	}
	return slow
}
