package list

import (
	"testing"
)

func TestMainer(t *testing.T) {

	// testing head_insert
	/*
		headPtr := &ListNode{
			Data: 1,
			Next: nil,
		}
		headNode := headPtr.HeadInsert(10)
		headNode.Traverse()
	*/

	// testing tail_insert
	/*
		headPtr := &ListNode{
			Data: 1,
			Next: nil,
		}
		headPtr.TailInsert(10)
		headPtr.Traverse()
	*/

	// testing find_value
	/*
		headPtr := &ListNode{
			Data: 1,
			Next: &ListNode{
				Data: 10,
				Next: nil,
			},
		}
		t.Log(headPtr.FindValue(10).Data)
	*/

	// testing delete node

	headPtr := &ListNode{
		Data: 1,
		Next: &ListNode{
			Data: 10,
			Next: nil,
		},
	}
	headPtr.Traverse()
	headPtr.Delete(headPtr)
	headPtr.Traverse()

	// testing reverse
	/*
		headPtr := &ListNode{
			Data: 1,
			Next: &ListNode{
				Data: 2,
				Next: &ListNode{Data: 3},
			},
		}
		head.Traverse()

		reverse := head.Reverse()
		reverse.Traverse()
	*/
}
