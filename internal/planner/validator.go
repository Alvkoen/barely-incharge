package planner

import (
	"fmt"
	"slices"
	"time"
)

func ValidateBlocks(blocks []TimeBlock, busyBlocks []TimeBlock) error {
	// Create map for O(1) busy block lookup instead of O(n) search
	busySet := make(map[string]bool, len(busyBlocks))
	for _, b := range busyBlocks {
		busySet[blockKey(b)] = true
	}

	allBlocks := make([]TimeBlock, 0, len(blocks)+len(busyBlocks))
	allBlocks = append(allBlocks, blocks...)
	allBlocks = append(allBlocks, busyBlocks...)

	slices.SortFunc(allBlocks, func(a, b TimeBlock) int {
		return a.Start.Compare(b.Start)
	})

	for i := 0; i < len(allBlocks)-1; i++ {
		if blocksOverlap(allBlocks[i], allBlocks[i+1]) {
			block1, block2 := allBlocks[i], allBlocks[i+1]

			isBusy1 := busySet[blockKey(block1)]
			isBusy2 := busySet[blockKey(block2)]

			if isBusy1 || isBusy2 {
				plannedBlock, busyBlock := block1, block2
				if isBusy1 {
					plannedBlock, busyBlock = block2, block1
				}
				return fmt.Errorf("block '%s' (%s-%s) overlaps with busy time '%s' (%s-%s)",
					plannedBlock.Title, formatTime(plannedBlock.Start), formatTime(plannedBlock.End),
					busyBlock.Title, formatTime(busyBlock.Start), formatTime(busyBlock.End))
			}

			return fmt.Errorf("blocks overlap: '%s' (%s-%s) and '%s' (%s-%s)",
				block1.Title, formatTime(block1.Start), formatTime(block1.End),
				block2.Title, formatTime(block2.Start), formatTime(block2.End))
		}
	}

	return nil
}

func blockKey(b TimeBlock) string {
	return b.Title + "|" + b.Start.String() + "|" + b.End.String()
}

func formatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

func blocksOverlap(a, b TimeBlock) bool {
	return a.Start.Before(b.End) && b.Start.Before(a.End)
}
