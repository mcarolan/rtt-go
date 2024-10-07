package ray

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"rtt/matrix"
	"rtt/shared"
	"rtt/sharedtest"
	"rtt/transformations"
	"rtt/tuple"
	"rtt/tupletest"
	"testing"

	"github.com/cucumber/godog"
)

func aRotation(ctx context.Context, variable, over string, value float64) (context.Context, error) {
	if over == "x" {
		return context.WithValue(ctx, sharedtest.Variables{Name: variable}, transformations.RotationX(math.Pi/value)), nil
	} else if over == "y" {
		return context.WithValue(ctx, sharedtest.Variables{Name: variable}, transformations.RotationY(math.Pi/value)), nil
	} else if over == "z" {
		return context.WithValue(ctx, sharedtest.Variables{Name: variable}, transformations.RotationZ(math.Pi/value)), nil
	} else {
		return ctx, fmt.Errorf("Unknown component %s", over)
	}
}

func aMatrixMul(ctx context.Context, variable, m1Var, m2Var string) (context.Context, error) {
	m1 := ctx.Value(sharedtest.Variables{Name: m1Var}).(*matrix.Matrix)
	m2 := ctx.Value(sharedtest.Variables{Name: m2Var}).(*matrix.Matrix)

	result := m1.Multiply(m2)

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, result), nil
}

func aRayFromVariables(ctx context.Context, variable, originVariable, directionVariable string) (context.Context, error) {
	origin := ctx.Value(sharedtest.Variables{Name: originVariable}).(*tuple.Tuple)
	direction := ctx.Value(sharedtest.Variables{Name: directionVariable}).(*tuple.Tuple)

	ray := NewRay(*origin, *direction)

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, ray), nil
}

func aPointLightFromVariables(ctx context.Context, variable, positionVariable, intensityVariable string) (context.Context, error) {
	position := ctx.Value(sharedtest.Variables{Name: positionVariable}).(*tuple.Tuple)
	intensity := ctx.Value(sharedtest.Variables{Name: intensityVariable}).(*tuple.Tuple)

	pointLight := NewPointLight(*position, *intensity)

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, pointLight), nil
}

func aMaterial(ctx context.Context, variable, positionVariable, intensityVariable string) (context.Context, error) {
	material := NewMaterial()

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, material), nil
}

func aRayFromValues(ctx context.Context, variable string, originX, originY, originZ, directionX, directionY, directionZ float64) (context.Context, error) {
	origin := tuple.Point(originX, originY, originZ)
	direction := tuple.Vector(directionX, directionY, directionZ)

	ray := NewRay(*origin, *direction)

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, ray), nil
}

func aSphere(ctx context.Context, variable string) (context.Context, error) {
	sphere := NewSphere()
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, sphere), nil
}

func aIntersect(ctx context.Context, variable, sphereVariable, rayVariable string) (context.Context, error) {
	sphere := ctx.Value(sharedtest.Variables{Name: sphereVariable}).(*Sphere)
	ray := ctx.Value(sharedtest.Variables{Name: rayVariable}).(*Ray)

	result := sphere.Intersect(ray)

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, result), nil
}

func aIntersection(ctx context.Context, variable string, t float64, sphereVariable string) (context.Context, error) {
	sphere := ctx.Value(sharedtest.Variables{Name: sphereVariable}).(*Sphere)

	result := sphere.Intersection(t)

	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, result), nil
}

func aTranslationMatrix(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	matrix := transformations.Translation(x, y, z)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, matrix), nil
}

func aScalingMatrix(ctx context.Context, variable string, x, y, z float64) (context.Context, error) {
	matrix := transformations.Scaling(x, y, z)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, matrix), nil
}

func aTransform(ctx context.Context, variable, rayVariable, matrixVariable string) (context.Context, error) {
	matrix := ctx.Value(sharedtest.Variables{Name: matrixVariable}).(*matrix.Matrix)
	ray := ctx.Value(sharedtest.Variables{Name: rayVariable}).(*Ray)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, ray.Transform(matrix)), nil
}

func aIntersections2(ctx context.Context, variable, i1Variable, i2Variable string) (context.Context, error) {
	i1 := ctx.Value(sharedtest.Variables{Name: i1Variable}).(*Intersection)
	i2 := ctx.Value(sharedtest.Variables{Name: i2Variable}).(*Intersection)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, []Intersection{*i1, *i2}), nil
}

func aIntersections4(ctx context.Context, variable, i1Variable, i2Variable, i3Variable, i4Variable string) (context.Context, error) {
	i1 := ctx.Value(sharedtest.Variables{Name: i1Variable}).(*Intersection)
	i2 := ctx.Value(sharedtest.Variables{Name: i2Variable}).(*Intersection)
	i3 := ctx.Value(sharedtest.Variables{Name: i3Variable}).(*Intersection)
	i4 := ctx.Value(sharedtest.Variables{Name: i4Variable}).(*Intersection)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, []Intersection{*i1, *i2, *i3, *i4}), nil
}

func aHit(ctx context.Context, variable, intersectionsVariable string) (context.Context, error) {
	intersections := ctx.Value(sharedtest.Variables{Name: intersectionsVariable}).([]Intersection)
	result := Hit(intersections)
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, result), nil
}

func aNormalAt(ctx context.Context, variable, sphereVariable, xStr, yStr, zStr string) (context.Context, error) {
	sphere := ctx.Value(sharedtest.Variables{Name: sphereVariable}).(*Sphere)
	x, y, z, err := sharedtest.ParseXYZ(xStr, yStr, zStr)

	if err != nil {
		return ctx, err
	}

	result := sphere.NormalAt(*tuple.Point(x, y, z))
	return context.WithValue(ctx, sharedtest.Variables{Name: variable}, result), nil
}

func setTransform(ctx context.Context, sphereVariable, matrixVariable string) (context.Context, error) {
	sphere := ctx.Value(sharedtest.Variables{Name: sphereVariable}).(*Sphere)
	matrix := ctx.Value(sharedtest.Variables{Name: matrixVariable}).(*matrix.Matrix)
	err := sphere.SetTransform(matrix)

	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func assertRayComponent(ctx context.Context, rayVariable, component, tupleVariable string) (context.Context, error) {
	ray := ctx.Value(sharedtest.Variables{Name: rayVariable}).(*Ray)
	t := ctx.Value(sharedtest.Variables{Name: tupleVariable}).(*tuple.Tuple)

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

func assertPointLightComponent(ctx context.Context, pointLightVariable, component, tupleVariable string) (context.Context, error) {
	pointLight := ctx.Value(sharedtest.Variables{Name: pointLightVariable}).(*PointLight)
	t := ctx.Value(sharedtest.Variables{Name: tupleVariable}).(*tuple.Tuple)

	var expected tuple.Tuple

	if component == "position" {
		expected = pointLight.Position
	} else if component == "intensity" {
		expected = pointLight.Intensity
	} else {
		return ctx, fmt.Errorf("unknown component %s", component)
	}

	if !tuple.CompareTuple(t, &expected) {
		return ctx, fmt.Errorf("Error %+v != %+v!", t, expected)
	}

	return ctx, nil
}

func assertMaterialComponent(ctx context.Context, pointLightVariable, component, tupleVariable string) (context.Context, error) {
	pointLight := ctx.Value(sharedtest.Variables{Name: pointLightVariable}).(*PointLight)
	t := ctx.Value(sharedtest.Variables{Name: tupleVariable}).(*tuple.Tuple)

	var expected tuple.Tuple

	if component == "position" {
		expected = pointLight.Position
	} else if component == "intensity" {
		expected = pointLight.Intensity
	} else {
		return ctx, fmt.Errorf("unknown component %s", component)
	}

	if !tuple.CompareTuple(t, &expected) {
		return ctx, fmt.Errorf("Error %+v != %+v!", t, expected)
	}

	return ctx, nil
}

func assertIntersectionsT(ctx context.Context, intersectionVariable string, index int, t float64) (context.Context, error) {
	intersections := ctx.Value(sharedtest.Variables{Name: intersectionVariable}).([]Intersection)
	intersection := intersections[index]

	if !shared.CompareFloat(intersection.T, t) {
		return ctx, fmt.Errorf("Error %+v != %+v!", intersection.T, t)
	}

	return ctx, nil
}

func assertIntersectionsObject(ctx context.Context, intersectionVariable string, index int, objectVariable string) (context.Context, error) {
	intersections := ctx.Value(sharedtest.Variables{Name: intersectionVariable}).([]Intersection)
	intersection := intersections[index]

	object := ctx.Value(sharedtest.Variables{Name: objectVariable}).(*Sphere)

	if intersection.Object != object.Id {
		return ctx, fmt.Errorf("Error %d != %d!", intersection.Object, object.Id)
	}

	return ctx, nil
}

func assertSphereTransform(ctx context.Context, sphereVariable, matrixVariable string) (context.Context, error) {
	sphere := ctx.Value(sharedtest.Variables{Name: sphereVariable}).(*Sphere)
	m := matrix.Identity

	if matrixVariable != "id" {
		m = ctx.Value(sharedtest.Variables{Name: matrixVariable}).(*matrix.Matrix)
	}

	if !sphere.transformation.Equals(m) {
		return ctx, fmt.Errorf("Error %+v != %+v!", sphere.transformation, m)
	}

	return ctx, nil
}

func assertEqualsVector(ctx context.Context, tupleVariable, xStr, yStr, zStr string) (context.Context, error) {
	actual := ctx.Value(sharedtest.Variables{Name: tupleVariable}).(*tuple.Tuple)
	x, y, z, err := sharedtest.ParseXYZ(xStr, yStr, zStr)

	if err != nil {
		return ctx, err
	}

	expected := tuple.Vector(x, y, z)

	if !tuple.CompareTuple(expected, actual) {
		return ctx, fmt.Errorf("Error %+v != %+v!", expected, actual)
	}

	return ctx, nil
}

func assertIntersectionT(ctx context.Context, intersectionVariable string, t float64) (context.Context, error) {
	intersection := ctx.Value(sharedtest.Variables{Name: intersectionVariable}).(*Intersection)

	if !shared.CompareFloat(intersection.T, t) {
		return ctx, fmt.Errorf("Error %+v != %+v!", intersection.T, t)
	}

	return ctx, nil
}

func assertIntersectionObject(ctx context.Context, intersectionVariable, objectVariable string) (context.Context, error) {
	intersection := ctx.Value(sharedtest.Variables{Name: intersectionVariable}).(*Intersection)
	object := ctx.Value(sharedtest.Variables{Name: objectVariable}).(*Sphere)

	if intersection.Object != object.Id {
		return ctx, fmt.Errorf("Error %d != %d!", intersection.Object, object.Id)
	}

	return ctx, nil
}

func assertArrayCount(ctx context.Context, variable string, expected int) (context.Context, error) {
	intersections := ctx.Value(sharedtest.Variables{Name: variable})

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
	intersections := ctx.Value(sharedtest.Variables{Name: variable}).([]float64)
	value := intersections[i]

	if !shared.CompareFloat(value, expected) {
		return ctx, fmt.Errorf("Error %f != %f!", expected, value)
	}

	return ctx, nil
}

func assertIntersectionEquals(ctx context.Context, i1Variable, i2Variable string) (context.Context, error) {
	i1 := ctx.Value(sharedtest.Variables{Name: i1Variable}).(*Intersection)
	i2 := ctx.Value(sharedtest.Variables{Name: i2Variable}).(*Intersection)

	if !shared.CompareFloat(i1.T, i2.T) || i1.Object != i2.Object {
		return ctx, fmt.Errorf("Error %+v != %+v!", i1, i2)
	}

	return ctx, nil
}

func assertIntersectionNothing(ctx context.Context, variable string) (context.Context, error) {
	i := ctx.Value(sharedtest.Variables{Name: variable}).(*Intersection)

	if i != nil {
		return ctx, fmt.Errorf("Error %+v is not nothing!", i)
	}

	return ctx, nil
}

func assertRayPosition(ctx context.Context, rayVariable string, t, x, y, z float64) (context.Context, error) {
	ray := ctx.Value(sharedtest.Variables{Name: rayVariable}).(*Ray)

	expected := tuple.Point(x, y, z)
	actual := ray.Position(t)

	if !tuple.CompareTuple(actual, expected) {
		return ctx, fmt.Errorf("Error %+v != %+v!", actual, expected)
	}

	return ctx, nil
}

func constructors(ctx *godog.ScenarioContext) {
	tupletest.AddConstructPoint(ctx)
	tupletest.AddConstructVector(ctx)

	regex := fmt.Sprintf(`^(.+) ← point_light\(%s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aPointLightFromVariables)

	regex = fmt.Sprintf(`^(.+) ← material\(\)$`)
	ctx.Step(regex, aMaterial)

	regex = fmt.Sprintf(`^(.+) ← ray\(%s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aRayFromVariables)

	regex = fmt.Sprintf(`^(.+) ← ray\(point\(%s, %s, %s\), vector\(%s, %s, %s\)\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	ctx.Step(regex, aRayFromValues)

	ctx.Step(`^(.+) ← sphere\(\)$`, aSphere)

	regex = fmt.Sprintf(`^(.+) ← intersect\(%s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aIntersect)

	regex = fmt.Sprintf(`^(.+) ← translation\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	ctx.Step(regex, aTranslationMatrix)

	regex = fmt.Sprintf(`^(.+) ← scaling\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	ctx.Step(regex, aScalingMatrix)

	regex = fmt.Sprintf(`^(.+) ← transform\(%s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aTransform)

	regex = fmt.Sprintf(`^(.+) ← intersection\(%s, %s\)$`, sharedtest.Decimal, sharedtest.TupleVariableName)
	ctx.Step(regex, aIntersection)

	regex = fmt.Sprintf(`^(.+) ← intersections\(%s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aIntersections2)

	regex = fmt.Sprintf(`^(.+) ← intersections\(%s, %s, %s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aIntersections4)

	regex = fmt.Sprintf(`^(.+) ← hit\(%s\)$`, sharedtest.TupleVariableName)
	ctx.Step(regex, aHit)

	regex = fmt.Sprintf(`^(.+) ← normal_at\(%s, point\(%s, %s, %s\)\)$`, sharedtest.TupleVariableName, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	ctx.Step(regex, aNormalAt)

	regex = `^(.+) ← rotation_(.)\(π\/(\d+)\)$`
	ctx.Step(regex, aRotation)

	regex = fmt.Sprintf(`^(.+) ← %s \* %s$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, aMatrixMul)

	tupletest.AddConstructColor(ctx)
}

func assertions(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^%s.(origin|direction) = %s$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, assertRayComponent)
	regex = fmt.Sprintf(`^position\((.+), %s\) = point\(%s, %s, %s\)$`, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal, sharedtest.Decimal)
	ctx.Step(regex, assertRayPosition)
	regex = fmt.Sprintf(`^%s.count = %s$`, sharedtest.TupleVariableName, sharedtest.PosInt)
	ctx.Step(regex, assertArrayCount)
	regex = fmt.Sprintf(`^%s\[%s\] = %s$`, sharedtest.TupleVariableName, sharedtest.PosInt, sharedtest.Decimal)
	ctx.Step(regex, assertArrayComponent)
	regex = fmt.Sprintf(`^%s.t = %s$`, sharedtest.TupleVariableName, sharedtest.Decimal)
	ctx.Step(regex, assertIntersectionT)
	regex = fmt.Sprintf(`^%s.object = %s$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, assertIntersectionObject)
	regex = fmt.Sprintf(`^%s\[%s\].t = %s$`, sharedtest.TupleVariableName, sharedtest.PosInt, sharedtest.Decimal)
	ctx.Step(regex, assertIntersectionsT)
	regex = fmt.Sprintf(`^%s\[%s\].object = %s$`, sharedtest.TupleVariableName, sharedtest.PosInt, sharedtest.TupleVariableName)
	ctx.Step(regex, assertIntersectionsObject)
	regex = fmt.Sprintf(`^%s = %s$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, assertIntersectionEquals)
	regex = fmt.Sprintf(`^%s is nothing$`, sharedtest.TupleVariableName)
	ctx.Step(regex, assertIntersectionNothing)
	regex = fmt.Sprintf(`^%s.transform = %s$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, assertSphereTransform)

	regex = fmt.Sprintf(`^%s.(position|intensity) = %s$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, assertPointLightComponent)

	regex = fmt.Sprintf(`^%s.(ambient|color|diffuse|shininess|specular) = %s$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
	ctx.Step(regex, assertMaterialComponent)

	tupletest.AddCompareNormalize(ctx)
	tupletest.AddCompareVector(ctx)
}

func setters(ctx *godog.ScenarioContext) {
	regex := fmt.Sprintf(`^set_transform\(%s, %s\)$`, sharedtest.TupleVariableName, sharedtest.TupleVariableName)
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
			Paths:    []string{"features/rays.feature", "features/spheres.feature", "features/intersections.feature", "features/lights.feature", "features/materials.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero exit status")
	}
}
