package helper

import "github.com/fannyhasbi/lab-tools-lending/types"

func TopicInSlice(t types.TopicType, list []types.TopicType) bool {
	for _, b := range list {
		if t == b {
			return true
		}
	}

	return false
}
