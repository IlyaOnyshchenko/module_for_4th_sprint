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

// Функция parseTraining() преобразует входные данные в виде строки в информацию о количестве шагов, виде активности, её продолжительности
// Также внутри функции обрабатываются возможные ошибки
func parseTraining(data string) (int, string, time.Duration, error) {
	// Входная строка разбивается по разделителю "," и части становятся элементами слайса dataInput
	dataInput := strings.Split(data, ",")
	// Так как формат вводимой строки "6000,Ходьба,1h00m" (разделённый двумя запятыми), то проверяем длину получившегося слайса
	// Возвращаем возможную ошибку в виде сообщения "недопустимый формат данных"
	if len(dataInput) != 3 {
		return 0, "", 0, fmt.Errorf("недопустимый формат данных")
	}
	// Находим количество шагов steps путём преобразования строки в число. Количество шагов - первый элемент dataInput
	// Также не забываем обработать ошибку преобразования
	steps, err := strconv.Atoi(dataInput[0])
	if err != nil {
		return 0, "", 0, err
	}
	// Значение шагов не может быть отрицательным или равным нулю
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("нулевые шаги")
	}
	// Вид активности - переменная типа string. Элемент слайса dataInput (второй) по типу совпадает со string. Присваиваем значение переменной activType
	activType := dataInput[1]
	// Находим продолжительность прогулки walkDur (третий элемент dataInput). Парсим строку в переменную time.Duration, обрабатывая возможную ошибку
	walkDur, err := time.ParseDuration(dataInput[2])
	if err != nil {
		return 0, "", 0, err
	}
	// Проверяем, чтобы продолжительность была положительна и отлична от нуля
	if walkDur <= 0 {
		return 0, "", 0, fmt.Errorf("нулевая продолжительность")
	}
	// Возвращаем из функции полученные значения и отсутствие ошибок
	return steps, activType, walkDur, nil
}

// Функция distance является внутренней вспомогательной функцией для расчета пройденного расстояния в километрах
// При известноим количестве шагов steps и росте height рассчитывается длина шага stepLen с использованием константы stepLengthCoefficient
func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	// Так как длина шага и их количество выражены в метрах, используем константу mInKm для перевода значения в километры
	return (float64(steps) * stepLen) / mInKm
}

// Функция meanSpeed является внутренней вспомогательной функцией для расчета средней скорости
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Здесь обработаем ошибку ввода продолжительности активности, так как согласно ТЗ нам это необходимо сделать.
	// Личное мнение автора кода - нижняя проверка излишняя, так как не равенство нулю и не отрицательное значение
	// проверяется ниже в теле функции TrainingInfo() при вызове функции parseTraining(). В алгоритме этой функции аргументы для вызова функций meanSpeed() и distance()
	// проверены.
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	// Величина скорости нам необходима выраженная в км/ч. Для этого воспользуемся методом hour из пакета time
	return dist / (float64(duration) / float64(time.Hour))
}

// Функция TrainingInfo() выводит данные о типе тренировки, её длительности, преодоленной дистанции, скорости и сожжённых калориях в виде форматированного вывода
func TrainingInfo(data string, weight, height float64) (string, error) {
	// Находим количество шагов steps и продолжительность прогулки walkDur с помощью вспомогательной функции parseTraining()
	// Узнаем вид активности и присваиваем её переменной activType, обрабатываем ошибку при вызове функции err. Согласно ТЗ, их нужно вывести в лог
	steps, activType, walkDur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// Получив входные данные для расчета пройденного пути distance и средней скорости meanSpeed,
	// вызываем вспомогательные функции и присваиваем названным переменным значения функций
	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, walkDur)
	// конвертируем полученную переменную walkDur типа time.Duration в количество часов типа float64
	convWalkDur := walkDur.Hours()
	// Необходимо вывести результат в зависимости от вида выполняемой физической активности activType
	switch activType {
	// В случае бега рассчитываем сожжённые калории с помощью функции RunningSpentCalories() и присваиваем значение переменной runningSpentCalories
	case "Бег":
		runningSpentCalories, err := RunningSpentCalories(steps, weight, height, walkDur)
		// Выводим лог ошибки в случае возникновения
		if err != nil {
			log.Println(err)
		}
		// Итоговый результат присваиваем переменной res в виде форматированного вывода и возвращаем её значение из функции и значение отсутствия ошибки nil
		res := fmt.Sprintf("Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", convWalkDur, distance, meanSpeed, runningSpentCalories)
		return res, nil
	// В случае ходьбы рассчитываем сожжённые калории с помощью функции WalkingSpentCalories() и присваиваем значение переменной walkingSpentCalories
	case "Ходьба":
		walkingSpentCalories, err := WalkingSpentCalories(steps, weight, height, walkDur)
		// Выводим лог ошибки в случае возникновения
		if err != nil {
			log.Println(err)
		}
		// Итоговый результат присваиваем переменной res в виде форматированного вывода и возвращаем её значение из функции и значение отсутствия ошибки nil
		res := fmt.Sprintf("Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", convWalkDur, distance, meanSpeed, walkingSpentCalories)
		return res, nil
	default:
		// Если при вводе было введено ни "Бег", ни "Ходьба", выводим ошибку "неизвестный тип тренировки"
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

// Функция RunningSpentCalories() подсчитывает количество сожжённых калорий по известному количеству шагов steps, росту height и весу weight,
// а также продолжительности бега duration
// Является экспортируемой, но в других пакетах не встречается
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Так как является экспортируемой, необходимо проверить корректность аргументов функции:
	// Количество шагов steps не должно быть нулевым или отрицательным
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов отрицательным быть не может")
	}
	// Вес weight и рост heght не должны быть нулевыми или отрицательными. Если хотя бы одна из величин отрицательна, это уже ошибка
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("рост и вес не могут быть отрицательными или нулевыми")
	}
	// Продолжительность активности не должна быть нулевой или отрицательной
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность бега - величина не отрицательная")
	}
	// Для расчета калорий нам нужно значение средней скорости meanSpeed
	meanSpeed := meanSpeed(steps, height, duration)
	// Продолжительность активности необходимо перевести в минуты. Присваиваем переменной durationInMinutes
	durationInMinutes := float64(duration) / float64(time.Minute)
	// Итоговое значение функции выводим по формуле с использованием константы minInH. Также возвращаем значение отсутствия ошибки
	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

// Функция RunningSpentCalories() подсчитывает количество сожжённых калорий по известному количеству шагов steps, росту height и весу weight,
// а также продолжительности ходьбы duration
// Является экспортируемой, встречается в пакете daysteps при расчете дневной активности с помощью DayActionInfo
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Необходимость проверки входных аргументов функции та же, что и у RunningSpentCalories()
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов отрицательным быть не может")
	}
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("рост и вес не могут быть отрицательными или нулевыми")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность бега - величина не отрицательная")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := float64(duration) / float64(time.Minute)
	// Итоговое значение функции выводим по формуле с поправкой на коэффициент для расчета калорий при ходьбе walkingCaloriesCoefficient
	// Также возвращаем значение отсутствия ошибки
	return ((weight * meanSpeed * durationInMinutes) / minInH) * walkingCaloriesCoefficient, nil
}
