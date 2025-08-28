package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataInput := strings.Split(data, ",")
	if len(dataInput) != 2 {
		return 0, 0, fmt.Errorf("Неверный формат данных")
	}
	steps, err := strconv.Atoi(dataInput[0])
	if err != nil {
		return 0, 0, err
	}
	walkDur, err := time.ParseDuration(dataInput[1])
	if err != nil {
		return 0, 0, err
	}
	return steps, walkDur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, walkDur, err := parsePackage(data)
	if err != nil || steps == 0 {
		return ""
	}
	dist := (float64(steps) * stepLength) / mInKm
	calories := WalkingSpentCalories(steps, weight, height, walkDur)
	res := fmt.Sprintf("Количество шагов: %d.\n", steps, "Дистанция составила %d км.\n", dist, "Вы сожгли %d ккал.", calories)
	return res
}
