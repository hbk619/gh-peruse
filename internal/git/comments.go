package git

import (
	"slices"
	"time"
)

func SortCommentsInPlace(commentList []Comment) {
	slices.SortFunc(commentList, func(i, j Comment) int {
		return time.Time.Compare(i.CreatedAt, j.CreatedAt)
	})
}
