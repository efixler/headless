package browser

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMustSetMaxTabsForTabs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewChrome(ctx, Headless(true))
	_, err := b.AcquireTab()
	if err == nil {
		t.Fatal("Expected error")
	}
	if !errors.Is(err, ErrMaxTabsNotSet) {
		t.Fatalf("Expected ErrMaxTabsNotSet, got %v", err)
	}
}

func TestErrorWhenTooManyTabs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewChrome(ctx, Headless(true), MaxTabs(1), TabAcquireTimeout(1*time.Second))
	_, err := b.AcquireTab()
	if err != nil {
		t.Fatalf("Unexpected error acquiring tab: %v", err)
	}
	_, err = b.AcquireTab()
	if err == nil {
		t.Fatal("Expected error")
	}
	t.Logf("Got error: %v", err)
	if !errors.Is(err, ErrMaxTabs) {
		t.Fatalf("Expected ErrMaxTabs, got %v", err)
	}
}

func TestAcquireTab(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewChrome(ctx, Headless(true), MaxTabs(1), TabAcquireTimeout(10*time.Second))
	tab, err := b.AcquireTab()
	if err != nil {
		t.Fatalf("Unexpected error acquiring tab: %v", err)
	}
	go func() {
		time.Sleep(3 * time.Second)
		tab.Get("https://www.mozilla.org/en-US/contact/", nil)
		t.Log("First tab done")
	}()
	_, err = b.AcquireTab()
	if err != nil {
		t.Fatalf("Unexpected error acquiring tab: %v", err)
	}
}
