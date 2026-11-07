package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const AdminsFile = "admins.txt"

func LoadAdminVKUserIDs(path string) (map[int64]struct{}, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[int64]struct{}{}, nil
		}
		return nil, err
	}
	defer file.Close()

	ids := make(map[int64]struct{})
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		id, err := strconv.ParseInt(line, 10, 64)
		if err != nil || id < 1 {
			return nil, fmt.Errorf("%s:%d: некорректный vk id", path, lineNo)
		}
		ids[id] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}
