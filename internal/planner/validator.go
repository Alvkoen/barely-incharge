package planner

import (
	"fmt"
	"slices"
)

func ValidateBlocks(blocks []TimeBlock, busyBlocks []TimeBlock) error {
	allBlocks := make([]TimeBlock, 0, len(blocks)+len(busyBlocks))
	allBlocks = append(allBlocks, blocks...)
	allBlocks = append(allBlocks, busyBlocks...)

	slices.SortFunc(allBlocks, func(a, b TimeBlock) int {
		return a.Start.Compare(b.Start)
	})

	for i := 0; i < len(allBlocks)-1; i++ {
		if blocksOverlap(allBlocks[i], allBlocks[i+1]) {
			block1 := allBlocks[i]
			block2 := allBlocks[i+1]

			isBusy1 := isInList(block1, busyBlocks)
			isBusy2 := isInList(block2, busyBlocks)

			if isBusy1 || isBusy2 {
				busyBlock := block1
				plannedBlock := block2
				if !isBusy2 {
					busyBlock = block2
					plannedBlock = block1
				}
				return fmt.Errorf("block '%s' (%s-%s) overlaps with busy time '%s' (%s-%s)",
					plannedBlock.Title, plannedBlock.Start.Format("15:04"), plannedBlock.End.Format("15:04"),
					busyBlock.Title, busyBlock.Start.Format("15:04"), busyBlock.End.Format("15:04"))
			}

			return fmt.Errorf("blocks overlap: '%s' (%s-%s) and '%s' (%s-%s)",
				block1.Title, block1.Start.Format("15:04"), block1.End.Format("15:04"),
				block2.Title, block2.Start.Format("15:04"), block2.End.Format("15:04"))
		}
	}

	return nil
}

func isInList(block TimeBlock, list []TimeBlock) bool {
	return slices.ContainsFunc(list, func(b TimeBlock) bool {
		return b.Title == block.Title && b.Start.Equal(block.Start) && b.End.Equal(block.End)
	})
}

func blocksOverlap(a, b TimeBlock) bool {
	return a.Start.Before(b.End) && b.Start.Before(a.End)
}
