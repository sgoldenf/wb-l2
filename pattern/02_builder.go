package pattern

import "math"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
	Строитель - порождающий паттерн проектирования, который инкапсулирует создание объекта и позволяет разделить его на различные этапы.

	Применимость:
	- Когда процесс создания нового объекта не должен зависеть от того, из каких частей он состоит и как они связаны
	- Когда необходимо обеспечить получение различных вариаций объекта в процессе его создания
	- В системе могут существовать сложные объекты, создание которых за одну операцию затруднительно или невозможно.
	Требуется поэтапное построение объектов с контролем результатов выполнения каждого этапа

	Плюсы и минусы:
	+ позволяет изменять внутреннее представление продукта
	+ изолирует код, реализующий конструирование и представление
	+ дает более тонкий контроль над процессом конструирования.
	- алгоритм создания сложного объекта не должен зависеть от того, из каких частей состоит объект и как они стыкуются между собой
	- процесс конструирования должен обеспечивать различные представления конструируемого объекта
	- ConcreteBuilder и создаваемый им объект жестко связаны

	Примеры использования на практике:
	Составление sql запросов, юнит-тесты
*/

// TransformMatrix - тип матриц трансформаций для афинных преобразований.
type TransformMatrix [4][4]float64

// Builder задает абстрактный интерфейс для создания объекта Product.
type Builder interface {
	buildTransformMatrix(x, y, z float64)
	getMatrix() TransformMatrix
}

// Director конструирует объект, вызывая методы строителя.
type Director struct {
	builder Builder
}

// Construct - метод для клиента.
func (d *Director) Construct(x, y, z float64) TransformMatrix {
	d.builder.buildTransformMatrix(x, y, z)
	return d.builder.getMatrix()
}

// RotationMatrixBuilder - конкретный строитель для матриц поворота.
type RotationMatrixBuilder struct {
	matrix TransformMatrix
}

func (b *RotationMatrixBuilder) buildTransformMatrix(x, y, z float64) {
	b.matrix[0][0] = math.Cos(y) * math.Cos(z)
	b.matrix[0][1] = -math.Sin(z)*math.Cos(x) + math.Cos(z)*math.Sin(y)*math.Sin(x)
	b.matrix[0][2] = math.Sin(z)*math.Sin(x) + math.Cos(z)*math.Sin(y)*math.Cos(x)
	b.matrix[1][0] = math.Sin(z) * math.Cos(y)
	b.matrix[1][1] = math.Cos(z)*math.Cos(x) + math.Sin(z)*math.Sin(y)*math.Sin(x)
	b.matrix[1][2] = -math.Cos(z)*math.Sin(x) + math.Sin(z)*math.Sin(y)*math.Cos(x)
	b.matrix[2][0] = -math.Sin(y)
	b.matrix[2][1] = math.Cos(y) * math.Sin(x)
	b.matrix[2][2] = math.Cos(x) * math.Cos(y)
	b.matrix[3][3] = 1
}

func (b *RotationMatrixBuilder) getMatrix() TransformMatrix {
	return b.matrix
}

// TranslationMatrixBuilder - конкретный строитель для матриц переноса.
type TranslationMatrixBuilder struct {
	matrix TransformMatrix
}

func (b *TranslationMatrixBuilder) buildTransformMatrix(x, y, z float64) {
	translateValues := [3]float64{x, y, z}
	for i := 0; i < 3; i++ {
		b.matrix[i][i] = 1
		b.matrix[i][3] = translateValues[i]
	}
	b.matrix[3][3] = 1
}

func (b *TranslationMatrixBuilder) getMatrix() TransformMatrix {
	return b.matrix
}

// ScaleMatrixBuilder - конкретный строитель для матриц масштабирования.
type ScaleMatrixBuilder struct {
	matrix TransformMatrix
}

func (b *ScaleMatrixBuilder) buildTransformMatrix(x, y, z float64) {
	scaleValues := [4]float64{x, y, z, 1}
	for i := 0; i < 4; i++ {
		b.matrix[i][i] = scaleValues[i]
	}
}

func (b *ScaleMatrixBuilder) getMatrix() TransformMatrix {
	return b.matrix
}
