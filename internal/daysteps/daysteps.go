package daysteps

import (
	"fmt"
	"log"
	"main/internal/spentcalories"
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
		return 0, 0, fmt.Errorf("неверный формат данных")
	}
	steps, err := strconv.Atoi(dataInput[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("отрицательные шаги")
	}
	walkDur, err := time.ParseDuration(dataInput[1])
	if err != nil {
		return 0, 0, err
	}
	if walkDur <= 0 {
		return 0, 0, fmt.Errorf("отрицательная продолжительность")
	}
	return steps, walkDur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, walkDur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	dist := (float64(steps) * stepLength) / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkDur)
	if err != nil {
		return ""
	}
	res := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, dist, calories)
	return res
}
