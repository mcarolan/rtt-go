Feature: Rays

Scenario: Creating and querying a ray
  Given origin ← point(1, 2, 3)
    And direction ← vector(4, 5, 6)
  When r ← ray(origin, direction)
  Then r.origin = origin
    And r.direction = direction

Scenario: Computing a point from a distance
  Given origin ← point(2, 3, 4)
    And direction ← vector(1, 0, 0)
    And r ← ray(origin, direction)
  Then position(r, 0) = point(2, 3, 4)
    And position(r, 1) = point(3, 3, 4)
    And position(r, -1) = point(1, 3, 4)
    And position(r, 2.5) = point(4.5, 3, 4)

Scenario: Translating a ray
  Given r ← ray(point(1, 2, 3), vector(0, 1, 0))
    And m ← translation(3, 4, 5)
  When r2 ← transform(r, m)
  Given origin ← point(4, 6, 8)
    And direction ← vector(0, 1, 0)
  Then r2.origin = origin
    And r2.direction = direction

Scenario: Scaling a ray
  Given r ← ray(point(1, 2, 3), vector(0, 1, 0))
    And m ← scaling(2, 3, 4)
  When r2 ← transform(r, m)
  Given origin ← point(2, 6, 12)
    And direction ← vector(0, 3, 0)
  Then r2.origin = origin
    And r2.direction = direction
