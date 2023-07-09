package util

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/onsi/ginkgo/v2"
)

// MatchSnapshot takes in any type and tries to match its stringer with the stored one
// if no previous snapshots exist it will create them.
// This function is meant to be used with ginkgo.
func MatchSnapshot(args ...interface{}) {
	t := ginkgo.GinkgoT()
	snaps.MatchSnapshot(t, args)
}
