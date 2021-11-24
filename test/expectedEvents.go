package test

import "fmt"

// This functions goal is to generate the
// expected events for each subscriber-group

func ExpectedEvents(prefixes []string, sufixes []string, groupSizes []int) []string {

	var events []string

	for i, prefix := range prefixes {
		for j := 0; j < groupSizes[i]; j++ {
			e := fmt.Sprintf("%s%d%s", prefix, j, sufixes[i])
			events = append(events, e)
		}
	}

	return events
}
