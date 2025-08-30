package daysteps

import (
	"fmt"
	"log"
	"module_for_4th_sprint/internal/spentcalories"
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
	// Создаём слайс dataInput из двух элементов типа string, создавая элементы из строки по разделителю ",".
	// Также проверяем, что длина равна двум, чтобы в дальнейшем коде не выйти за границы слайса
	dataInput := strings.Split(data, ",")
	if len(dataInput) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}
	// Первый элемент слайса - шаги. Преобразуем строку числа шагов в число, обрабатывая возможную ошибку, присваиваем значение переменной steps, ошибку - err
	steps, err := strconv.Atoi(dataInput[0])
	if err != nil {
		return 0, 0, err
	}
	// Проверяем, чтобы число шагов было положительное и отличное от нуля
	if steps <= 0 {
		return 0, 0, fmt.Errorf("отрицательные шаги")
	}
	// Второй элемент слайса - продолжительность прогулки. Парсим строку в переменную walkDur, обрабатывая возможную ошибку err
	walkDur, err := time.ParseDuration(dataInput[1])
	if err != nil {
		return 0, 0, err
	}
	// Проверяем, чтобы продолжительность была положительна и отлична от нуля
	if walkDur <= 0 {
		return 0, 0, fmt.Errorf("отрицательная продолжительность")
	}
	// Возвращаем из функции полученные значения и отсутствие ошибок
	return steps, walkDur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем значения количества шагов steps, продолжительности прогулки walkDur и выводим на экран возможные ошибки err.
	// Возвращаем пустую строку, согласно ТЗ.
	steps, walkDur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	// Считаем количество пройденного пути dist в км, используя константы stepLength и mInKm
	dist := (float64(steps) * stepLength) / mInKm
	// Получаем значение количества сожженных калорий и присваеиваем её переменной calories через функцию WalkingSpentCalories из пакета spentcalories
	// Также обрабатываем возможные ошибки и выводим их на экран
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkDur)
	if err != nil {
		log.Println(err)
		return ""
	}
	// Итоговый результат присваиваем переменной res в виде форматированного вывода и возвращаем её значение из функции
	res := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, dist, calories)
	return res
}
