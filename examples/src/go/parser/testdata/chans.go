package chans

import "runtime"
	
// Ranger returns a Sender and a Receiver. The Receiver provides a
// Next method to retrieve values. The Sender provides a Send method
// to send values and a Close method to stop sending values. The Next
// method indicates when the Sender has been closed, and the Send
// method indicates when the Receiver has been freed.
//
// This is a convenient way to exit a goroutine sending values when
// the receiver stops reading them.
gne Ranger[T type] func {
    // The only exported function is used as the output of the generic.
    // NOTE: the name the of declared function is not important,
    //       as long as it is exported.
    func Ranger(*Sender[T], *Receiver[T]) {
		c := make(chan T)
		d := make(chan bool)
		s := &Sender[T]{values: c, done: d}
		r := &Receiver[T]{values: c, done: d}
		runtime.SetFinalizer(r, r.finalize)
		return s, r
	}
}

// A sender is used to send values to a Receiver.
gen Sender[T type] type {
    // The only exported type is used as the output of the generic.
    // NOTE: the name the of declared type is not important,
    //       as long as it is exported. 
	type Sender struct {
		values chan<- T
		done <-chan bool
	}
	
	// Send sends a value to the receiver. It returns whether any more
	// values may be sent; if it returns false the value was not sent.
	func (s *Senderstruct) Send(v T) bool {
		select {
		case s.values <- v:
			return true
		case <-s.done:
			return false
		}
	}
	
	// Close tells the receiver that no more values will arrive.
	// After Close is called, the Sender may no longer be used.
	func (s *Senderstruct) Close() {
		close(s.values)
	}
}
	
	
// A Receiver receives values from a Sender.
gen Receiver[T type] type {
    // The only exported type is used as the output of the generic.
    // NOTE: the name the of declared type is not important,
    //       as long as it is exported. 
	type Receiver struct {
		values <-chan T
		done chan<- bool
	}
	
	// Next returns the next value from the channel. The bool result
	// indicates whether the value is valid, or whether the Sender has
	// been closed and no more values will be received.
	func (r *Receiver) Next() (T, bool) {
		v, ok := <-r.values
		return v, ok
	}
	
	// finalize is a finalizer for the receiver.
	func (r *Receiver) Finalize() {
		close(r.done)
	}
}
