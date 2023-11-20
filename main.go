package main

import (
	"log"
	"golang.org/x/sys/unix"
	"os"
)

// print debug statements
var dl *log.Logger = log.New(os.Stdout, "[DEBUG]: ", log.Lshortfile)
// print error debug statements
var el *log.Logger = log.New(os.Stderr, "[ERROR]: ", log.Lshortfile)

func main() {
	// create an epoll instance
	efd, err__create_epoll := unix.EpollCreate(1)
	if err__create_epoll != nil {
		el.Fatalf("failed to create an epoll instance.\n")
	}

	dl.Printf("created new epoll instance.\n")
	// create an EpollEvent including the fd that we want to observe
	// and the events we are interested in
	epoll_event := unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd: 0, // fd to stdin
	}

	// add the desired fd to the `interest` list
	err__add_event := unix.EpollCtl(efd, unix.EPOLL_CTL_ADD, int(epoll_event.Fd), &epoll_event)
	if err__add_event != nil {
		el.Fatalf("failed to add new Epoll Event\n")
	}

	dl.Printf("added new epoll event.\n")

	// create an `event store` where epoll saves all events 
	max_events := 1
	event_store := make([]unix.EpollEvent, max_events)

	dl.Printf("waiting to for input on stdin ...\n")
	// wait for events
	n_events, err__wait_events := unix.EpollWait(efd, event_store, -1) 
	if err__wait_events != nil {
		el.Fatalf("failed when waiting for events.\n")
	}

	dl.Printf("got %d events.\n", n_events)
	read_fd := event_store[0].Fd

	buffer := make([]byte, 20)
	n_bytes, err__read_bytes := unix.Read(int(read_fd), buffer)
	if err__read_bytes != nil {
		el.Fatalf("failed to read data from fd.\n")
	}

	dl.Printf("read %d bytes from fd\n", n_bytes)
	dl.Printf("Data:\n%s\n", string(buffer))
}

