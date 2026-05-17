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
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("Длина слайса не равна 3")
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка преобразования %v", err)
	}
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка парсинга: %v", err)
	}
	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / float64(mInKm)
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceKm := distance(steps, height)
	durationHours := duration.Hours()
	return distanceKm / durationHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var distanceKm float64
	var speedKmH float64
	var calories float64
	var calErr error
	switch trainingType {
	case "Бег":
		distanceKm = distance(steps, height)
		speedKmH = meanSpeed(steps, height, duration)
		calories, calErr = RunningSpentCalories(steps, weight, height, duration)
		if calErr != nil {
			log.Println(calErr)
			return "", calErr
		}
	case "Ходьба":
		distanceKm = distance(steps, height)
		speedKmH = meanSpeed(steps, height, duration)
		calories, calErr = WalkingSpentCalories(steps, weight, height, duration)
		if calErr != nil {
			log.Println(calErr)
			return "", calErr
		}
	default:
		return "", fmt.Errorf("Неизвестный тип тренировки")
	}
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		trainingType, duration.Hours(), distanceKm, speedKmH, calories,
	)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов должно быть больше нуля, получено: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Вес должен быть больше нуля, получено: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("Рост должен быть больше нуля, получено: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Продолжительность бега должна быть больше нуля, получено: %v", duration)
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / float64(minInH)
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов должно быть больше нуля, получено: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Вес должен быть больше нуля, получено: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("Рост должен быть больше нуля, получено: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Продолжительность ходьбы должна быть больше нуля, получено: %v", duration)
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	baseCalories := (weight * speed * durationInMinutes) / float64(minInH)
	finalCalories := baseCalories * walkingCaloriesCoefficient
	return finalCalories, nil
}
