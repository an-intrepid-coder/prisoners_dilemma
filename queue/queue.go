/* NOTE: In order to be able to provide "*OrNil"-style functions, Get() 
   and Consume() both return pointers to the desired contents and must
   be dereferenced.  */

package queue

type Queue struct {
    capacity int
    size int
    contents []int
    front int
    back int
}

func (q *Queue) init(n int) {
    q.capacity = n + 1 
    q.size = 0
    q.contents = make([]int, q.capacity, q.capacity) 
    q.front = 0
    q.back = 0
}

func MakeQueue(n int) Queue {
    q := Queue{}
    q.init(n)
    return q
}

func (q *Queue) Empty() bool { 
    return q.front == q.back
}

func (q *Queue) Full() bool { 
    i := (q.front + 1) % q.capacity
    return i == q.back
}

func (q* Queue) Insert(k int) { 
    if !q.Full() {
        q.contents[q.front] = k
        q.front = (q.front + 1) % q.capacity
        q.size++
    } 
}

func (q* Queue) GetOrNil() *int {
    if !q.Empty() {
        i := (q.front - 1) % q.capacity
        return &q.contents[i]
    }
    return nil 
}

func (q* Queue) Del() {
    if !q.Empty() {
        q.back = (q.back + 1) % q.capacity
        q.size--
    } 
}

func (q* Queue) ConsumeOrNil() *int {
    if !q.Empty() {
        e := q.GetOrNil()
        q.Del()
        return e
    }
    return nil 
}

func (q* Queue) Contents() []int {
    r := []int{}
    for i := q.back; i != q.front; i = (i + 1) % q.capacity {
        r = append(r, q.contents[i])
    }
    return r
}

func (q* Queue) Size() int {
    return q.size
}

func (q* Queue) Capacity() int {
    return q.capacity
}

