package filelock

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestFileLockWrites(t *testing.T) {
	var wg sync.WaitGroup
	ctx := NewContext()
	var arr []int

	write := func(val int) {
		defer wg.Done()
		permissions := map[string]bool{
			"a": true,
		}

		defer ctx.WithPermissions(permissions)()
		time.Sleep(10 * time.Millisecond)
		arr = append(arr, val)
	}

	wg.Add(2)

	write(1)
	write(2)

	wg.Wait()

	expected := []int{1, 2}

	if !reflect.DeepEqual(arr, expected) {
		t.Fatalf(`Array is not %+q, got %+q`, expected, arr)
	}
}

// Ensure a writer does not get starved by a reader
func TestFileLockStarvation(t *testing.T) {
	var wg sync.WaitGroup
	
	ctx := NewContext()
	var arr []int

	write := func() {
		defer wg.Done()
		permissions := map[string]bool{
			"a": true,
		}

		defer ctx.WithPermissions(permissions)()
		time.Sleep(10 * time.Millisecond)
		arr = append(arr, 1)
	}

	read := func() {
		defer wg.Done()
		permissions := map[string]bool{
			"a": false,
		}

		defer ctx.WithPermissions(permissions)()
		time.Sleep(10 * time.Millisecond)
		arr = append(arr, 0)
	}

	wg.Add(3)
	read()
	write()
	read()
	wg.Wait()

	expected := []int{0, 1, 0}

	if !reflect.DeepEqual(arr, expected) {
		t.Fatalf(`Array is not %+q, got %+q`, expected, arr)
	}
}
