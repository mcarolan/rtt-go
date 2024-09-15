package ray

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"rtt/matrix"
	"rtt/shared"
	"rtt/transformations"
	"rtt/tuple"
	"testing"

	"github.com/cucumber/godog"
)

type variables struct{ name string }

var tupleVariableName = `([a-z]+[0-9]*)`

func aPoint(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := tuple.Point(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aVector(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	p := tuple.Vector(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, p), nil
}

func aRayFromVariables(ctx context.Context, variable, originVariable, directionVariable string) (context.Context, error) {
	origin := ctx.Value(variables{name: originVariable}).(*tuple.Tuple)
	direction := ctx.Value(variables{name: directionVariable}).(*tuple.Tuple)

	ray := NewRay(*origin, *direction)

	return context.WithValue(ctx, variables{name: variable}, ray), nil
}

func aRayFromValues(ctx context.Context, variable string, originX, originY, originZ, directionX, directionY, directionZ float64) (context.Context, error) {
	origin := tuple.Point(originX, originY, originZ)
	direction := tuple.Vector(directionX, directionY, directionZ)

	ray := NewRay(*origin, *direction)

	return context.WithValue(ctx, variables{name: variable}, ray), nil
}

func aSphere(ctx context.Context, variable string) (context.Context, error) {
	sphere := NewSphere()
	return context.WithValue(ctx, variables{name: variable}, sphere), nil
}

func aIntersect(ctx context.Context, variable, sphereVariable, rayVariable string) (context.Context, error) {
	sphere := ctx.Value(variables{name: sphereVariable}).(*Sphere)
	ray := ctx.Value(variables{name: rayVariable}).(*Ray)

	result, err := sphere.Intersect(ray)

	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, variables{name: variable}, result), nil
}

func aIntersection(ctx context.Context, variable string, t float64, sphereVariable string) (context.Context, error) {
	sphere := ctx.Value(variables{name: sphereVariable}).(*Sphere)

	result := sphere.Intersection(t)

	return context.WithValue(ctx, variables{name: variable}, result), nil
}

func aTranslationMatrix(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	matrix := transformations.Translation(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, matrix), nil
}

func aScalingMatrix(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	matrix := transformations.Scaling(x, y, z)
	return context.WithValue(ctx, variables{name: variable}, matrix), nil
}

func aTransform(ctx context.Context, variable, rayVariable, matrixVariable string) (context.Context, error) {
	matrix := ctx.Value(variables{name: matrixVariable}).(*matrix.Matrix)
	ray := ctx.Value(variables{name: rayVariable}).(*Ray)
	return context.WithValue(ctx, variables{name: variable}, ray.Transform(matrix)), nil
}

func aIntersections2(ctx context.Context, variable, i1Variable, i2Variable string) (context.Context, error) {
	i1 := ctx.Value(variables{name: i1Variable}).(*Intersection)
	i2 := ctx.Value(variables{name: i2Variable}).(*Intersection)
	return context.WithValue(ctx, variables{name: variable}, []Intersection{*i1, *i2}), nil
}

func aIntersections4(ctx context.Context, variable, i1Variable, i2Variable, i3Variable, i4Variable string) (context.Context, error) {
	i1 := ctx.Value(variables{name: i1Variable}).(*Intersection)
	i2 := ctx.Value(variables{name: i2Variable}).(*Intersection)
	i3 := ctx.Value(variables{name: i3Variable}).(*Intersection)
	i4 := ctx.Value(variables{name: i4Variable}).(*Intersection)
	return context.WithValue(ctx, variables{name: variable}, []Intersection{*i1, *i2, *i3, *i4}), nil
}

func aHit(ctx context.Context, variable, intersectionsVariable string) (context.Context, error) {
	intersections := ctx.Value(variables{name: intersectionsVariable}).([]Intersection)
	result := Hit(intersections)
	return context.WithValue(ctx, variables{name: variable}, result), nil
}

func setTransform(ctx context.Context, sphereVariable, matrixVariable string) (context.Context, error) {
	sphere := ctx.Value(variables{name: sphereVariable}).(*Sphere)
	matrix := ctx.Value(variables{name: matrixVariable}).(*matrix.Matrix)
	sphere.Transformation = *matrix
	return ctx, nil
}

func assertRayComponent(ctx context.Context, rayVariable, component, tupleVariable string) (context.Context, error) {
	ray := ctx.Value(variables{name: rayVariable}).(*Ray)
	t := ctx.Value(variables{name: tupleVariable}).(*tuple.Tuple)

	var expected tuple.Tuple

	if component == "origin" {
		expected = ray.Origin
	} else if component == "direction" {
		expected = ray.Direction
	} else {
		return ctx, fmt.Errorf("unknown component %s", component)
	}

	if !tuple.CompareTuple(t, &expected) {
		return ctx, fmt.Errorf("Error %+v != %+v!", t, expected)
	}

	return ctx, nil
}

func assertIntersectionsT(ctx context.Context, intersectionVariable string, index int, t float64) (context.Context, error) {
	intersections := ctx.Value(variables{name: intersectionVariable}).([]Intersection)
	intersection := intersections[index]

	if !shared.CompareFloat(intersection.T, t) {
		return ctx, fmt.Errorf("Error %+v != %+v!", intersection.T, t)
	}

	return ctx, nil
}

func assertIntersectionsObject(ctx context.Context, intersectionVariable string, index int, objectVariable string) (context.Context, error) {
	intersections := ctx.Value(variables{name: intersectionVariable}).([]Intersection)
	intersection := intersections[index]

	object := ctx.Value(variables{name: objectVariable}).(*Sphere)

	if intersection.Object != object.Id {
		return ctx, fmt.Errorf("Error %d != %d!", intersection.Object, object.Id)
	}

	return ctx, nil
}

func assertSphereTransform(ctx context.Context, sphereVariable, matrixVariable string) (context.Context, error) {
	sphere := ctx.Value(variables{name: sphereVariable}).(*Sphere)
	m := matrix.Identity

	if matrixVariable != "id" {
		m = ctx.Value(variables{name: matrixVariable}).(*matrix.Matrix)
	}

	if !sphere.Transformation.Equals(m) {
		return ctx, fmt.Errorf("Error %+v != %+v!", sphere.Transformation, m)
	}

	return ctx, nil
}

func assertIntersectionT(ctx context.Context, intersectionVariable string, t float64) (context.Context, error) {
	intersection := ctx.Value(variables{name: intersectionVariable}).(*Intersection)

	if !shared.CompareFloat(intersection.T, t) {
		return ctx, fmt.Errorf("Error %+v != %+v!", intersection.T, t)
	}

	return ctx, nil
}

func assertIntersectionObject(ctx context.Context, intersectionVariable, objectVariable string) (context.Context, error) {
	intersection := ctx.Value(variables{name: intersectionVariable}).(*Intersection)
	object := ctx.Value(variables{name: objectVariable}).(*Sphere)

	if intersection.Object != object.Id {
		return ctx, fmt.Errorf("Error %d != %d!", intersection.Object, object.Id)
	}

	return ctx, nil
}

func assertArrayCount(ctx context.Context, variable string, expected int) (context.Context, error) {
	intersections := ctx.Value(variables{name: variable})

	value := reflect.ValueOf(intersections)

	if value.Kind() == reflect.Slice {
		if value.Len() != expected {
			return ctx, fmt.Errorf("Error count %d not %d!", value.Len(), expected)
		}
	} else {
		return ctx, errors.New("Not a slice")
	}

	return ctx, nil
}

func assertArrayComponent(ctx context.Context, variable string, i int, expected float64) (context.Context, error) {
	intersections := ctx.Value(variables{name: variable}).([]float64)
	value := intersections[i]

	if !shared.CompareFloat(value, expected) {
		return ctx, fmt.Errorf("Error %f != %f!", expected, value)
	}

	return ctx, nil
}

func assertIntersectionEquals(ctx context.Context, i1Variable, i2Variable string) (context.Context, error) {
	i1 := ctx.Value(variables{name: i1Variable}).(*Intersection)
	i2 := ctx.Value(variables{name: i2Variable}).(*Intersection)

	if !shared.CompareFloat(i1.T, i2.T) || i1.Object != i2.Object {
		return ctx, fmt.Errorf("Error %+v != %+v!", i1, i2)
	}

	return ctx, nil
}

func assertIntersectionNothing(ctx context.Context, variable string) (context.Context, error) {
	i := ctx.Value(variables{name: variable}).(*Intersection)

	if i != nil {
		return ctx, fmt.Errorf("Error %+v is not nothing!", i)
	}

	return ctx, nil
}

func assertRayPosition(ctx context.Context, rayVariable string, t, x, y, z float64) (context.Context, error) {
	ray := ctx.Value(variables{name: rayVariable}).(*Ray)

	expected := tuple.Point(x, y, z)
	actual := ray.Position(t)

	if !tuple.CompareTuple(actual, expected) {
		return ctx, fmt.Errorf("Error %+v != %+v!", actual, expected)
	}

	return ctx, nil
}

func constructors(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^(.+) ← point\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aPoint)

	regex = fmt.Sprintf(`^(.+) ← vector\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aVector)

	regex = fmt.Sprintf(`^(.+) ← ray\(%s, %s\)$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, aRayFromVariables)

	regex = fmt.Sprintf(`^(.+) ← ray\(point\(%s, %s, %s\), vector\(%s, %s, %s\)\)$`, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aRayFromValues)

	ctx.Step(`^(.+) ← sphere\(\)$`, aSphere)

	regex = fmt.Sprintf(`^(.+) ← intersect\(%s, %s\)$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, aIntersect)

	regex = fmt.Sprintf(`^(.+) ← translation\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aTranslationMatrix)

	regex = fmt.Sprintf(`^(.+) ← scaling\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, aScalingMatrix)

	regex = fmt.Sprintf(`^(.+) ← transform\(%s, %s\)$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, aTransform)

	regex = fmt.Sprintf(`^(.+) ← intersection\(%s, %s\)$`, shared.Decimal, tupleVariableName)
	ctx.Step(regex, aIntersection)

	regex = fmt.Sprintf(`^(.+) ← intersections\(%s, %s\)$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, aIntersections2)

	regex = fmt.Sprintf(`^(.+) ← intersections\(%s, %s, %s, %s\)$`, tupleVariableName, tupleVariableName, tupleVariableName, tupleVariableName)
	ctx.Step(regex, aIntersections4)

	regex = fmt.Sprintf(`^(.+) ← hit\(%s\)$`, tupleVariableName)
	ctx.Step(regex, aHit)
}

func assertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^%s.(origin|direction) = %s$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, assertRayComponent)
	regex = fmt.Sprintf(`^position\((.+), %s\) = point\(%s, %s, %s\)$`, shared.Decimal, shared.Decimal, shared.Decimal, shared.Decimal)
	ctx.Step(regex, assertRayPosition)
	regex = fmt.Sprintf(`^%s.count = %s$`, tupleVariableName, shared.PosInt)
	ctx.Step(regex, assertArrayCount)
	regex = fmt.Sprintf(`^%s\[%s\] = %s$`, tupleVariableName, shared.PosInt, shared.Decimal)
	ctx.Step(regex, assertArrayComponent)
	regex = fmt.Sprintf(`^%s.t = %s$`, tupleVariableName, shared.Decimal)
	ctx.Step(regex, assertIntersectionT)
	regex = fmt.Sprintf(`^%s.object = %s$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, assertIntersectionObject)
	regex = fmt.Sprintf(`^%s\[%s\].t = %s$`, tupleVariableName, shared.PosInt, shared.Decimal)
	ctx.Step(regex, assertIntersectionsT)
	regex = fmt.Sprintf(`^%s\[%s\].object = %s$`, tupleVariableName, shared.PosInt, tupleVariableName)
	ctx.Step(regex, assertIntersectionsObject)
	regex = fmt.Sprintf(`^%s = %s$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, assertIntersectionEquals)
	regex = fmt.Sprintf(`^%s is nothing$`, tupleVariableName)
	ctx.Step(regex, assertIntersectionNothing)
	regex = fmt.Sprintf(`^%s.transform = %s$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, assertSphereTransform)
}

func setters(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^set_transform\(%s, %s\)$`, tupleVariableName, tupleVariableName)
	ctx.Step(regex, setTransform)
}

func initializeScenario(ctx *godog.ScenarioContext) {
	constructors(ctx)
	assertions(ctx)
	setters(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: initializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/rays.feature", "features/spheres.feature", "features/intersections.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero exit status")
	}
}
