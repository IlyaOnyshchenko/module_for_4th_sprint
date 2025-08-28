package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	dataInput := strings.Split(data, ",")
	if len(dataInput) != 3 {
		return 0, "", 0, fmt.Errorf("Недопустимый формат данных")
	}
	steps, err := strconv.Atoi(dataInput[0])
	if err != nil {
		return 0, "", 0, err
	}
	activType := dataInput[1]
	walkDur, err := time.ParseDuration(dataInput[2])
	if err != nil {
		return 0, "", 0, err
	}
	return steps, activType, walkDur, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	return (float64(steps) * stepLen) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / (float64(duration) / float64(time.Hour))
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activType, walkDur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, walkDur)
	switch activType {
	case "Бег":
		runningSpentCalories, err := RunningSpentCalories(steps, weight, height, walkDur)
		if err != nil {
			log.Println(err)
		}
		res := fmt.Sprintf("Тип тренировки: Бег\n", "Длительность: %d ч.\n", walkDur, "Дистанция: %d км.\n", distance, "Скорость: %d км/ч", meanSpeed, "Сожгли калорий: %d", runningSpentCalories)
		return res, nil
	case "Ходьба":
		walkingSpentCalories, err := WalkingSpentCalories(steps, weight, height, walkDur)
		if err != nil {
			log.Println(err)
		}
		res := fmt.Sprintf("Тип тренировки: Бег\n", "Длительность: %d ч.\n", walkDur, "Дистанция: %d км.\n", distance, "Скорость: %d км/ч", meanSpeed, "Сожгли калорий: %d", walkingSpentCalories)
		return res, nil
	default:
		fmt.Println("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 {
		return 0, fmt.Errorf("Количество шагов отрицательным быть не может")
	}
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("Рост и вес не могут быть отрицательными или нулевыми")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Продолжительность бега - величина не отрицательная")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := float64(duration) / float64(time.Minute)
	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 {
		return 0, fmt.Errorf("Количество шагов отрицательным быть не может")
	}
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("Рост и вес не могут быть отрицательными или нулевыми")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Продолжительность бега - величина не отрицательная")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := float64(duration) / float64(time.Minute)
	return ((weight * meanSpeed * durationInMinutes) / minInH) * walkingCaloriesCoefficient, nil
}
