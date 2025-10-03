package parse

import (
	"strconv"
	"strings"

	e "github.com/KirillLich/todoapi/internal/errors"
)

func ParseId(path string, previousKey string) (int, error) {
	arr := strings.Split(path, "/")
	if arr[len(arr)-1] == previousKey {
		return 0, e.ErrInvalidId
	}
	var sId string
	for i := 0; i < len(arr); i++ {
		if arr[i] == previousKey {
			i++
			sId = arr[i]
			break
		}
	}
	id, err := strconv.Atoi(sId)
	if err == strconv.ErrSyntax {
		return 0, e.ErrInvalidId
	}
	return id, err
}
