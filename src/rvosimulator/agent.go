package rvosimulator

import (
	"math"
)

func init() {
}

type Agent struct {
	ID                int
	Position          *Vector2
	Radius            float64
	TimeHorizon       float64
	TimeHorizonObst   float64
	Velocity          *Vector2
	PrefVelocity      *Vector2
	NewVelocity       *Vector2
	MaxNeighbors      int
	NeighborDist      float64
	MaxSpeed          float64
	ObstacleNeighbors []*ObstacleNeighbor // pair ?
	AgentNeighbors    []*AgentNeighbor    //pair?
	OrcaLines         []*Line
}

type ObstacleNeighbor struct {
	DistSq   float64
	Obstacle *Obstacle
}

type AgentNeighbor struct {
	DistSq float64
	Agent  *Agent
}

// Line :
type Line struct {
	Point     *Vector2
	Direction *Vector2
}

var (
	RVO_EPSILON float64 //A sufficiently small positive number.
)

func init() {
	RVO_EPSILON = 0.00001
}

func NewEmptyAgent() *Agent {
	a := &Agent{
		ID:              0,
		Radius:          float64(0),
		TimeHorizon:     float64(0),
		TimeHorizonObst: float64(0),
		MaxNeighbors:    0,
		NeighborDist:    float64(0),
		MaxSpeed:        float64(0),
	}
	return a
}

func NewAgent(id int, position *Vector2, radius float64, timeHorizon float64, timeHorizonObst float64, velocity *Vector2, newVelocity *Vector2, prefVelocity *Vector2, maxNeighbors int, neighborDist float64, maxSpeed float64, obstacleNeighbors []*ObstacleNeighbor) *Agent {
	a := &Agent{
		ID:                id,
		Position:          position,
		Radius:            radius,
		TimeHorizon:       timeHorizon,
		TimeHorizonObst:   timeHorizonObst,
		Velocity:          velocity,
		PrefVelocity:      prefVelocity,
		NewVelocity:       newVelocity,
		MaxNeighbors:      maxNeighbors,
		NeighborDist:      neighborDist,
		MaxSpeed:          maxSpeed,
		ObstacleNeighbors: obstacleNeighbors,
	}
	return a
}

// CHECK OK
// OK 1
func (a *Agent) ComputeNeighbors() {
	a.ObstacleNeighbors = make([]*ObstacleNeighbor, 0)
	rangeSq := math.Pow(a.TimeHorizonObst*a.MaxSpeed+a.Radius, 2)

	//	fmt.Printf("FN: ComputeNeighbors \n agent check \n ID %v\n Position: %v\n MaxNeighbors: %v\n MaxSpeed; %v\nNeighborDist: %v\n Radius: %v\n TimeHorizon: %v\n TimeHorizonObst: %v\n Velocity: %v\n\n",
	//		a.ID, a.Position, a.MaxNeighbors, a.MaxSpeed, a.NeighborDist, a.Radius, a.TimeHorizon, a.TimeHorizonObst, a.Velocity)

	//	fmt.Printf("FN: ComputeNeighbors \n Obstacle RangeSq %v\n\n",
	//		rangeSq)

	Sim.KdTree.ComputeObstacleNeighbors(a, rangeSq)

	a.AgentNeighbors = make([]*AgentNeighbor, 0)

	if a.MaxNeighbors > 0 {
		rangeSq = math.Pow(a.NeighborDist, 2)
		//		fmt.Printf("FN: ComputeNeighbors \n Agent RangeSq %v\n\n",
		//			rangeSq)
		Sim.KdTree.ComputeAgentNeighbors(a, rangeSq)
	}
}

// CHECK OK
//FINISH
func (a *Agent) ComputeNewVelocity() {
	a.OrcaLines = make([]*Line, 0)

	invTimeHorizonObst := float64(1) / a.TimeHorizonObst

	/* Create obstacle ORCA lines. */
	for i := 0; i < len(a.ObstacleNeighbors); i++ {

		var obstacle1, obstacle2 *Obstacle
		obstacle1 = a.ObstacleNeighbors[i].Obstacle
		obstacle2 = obstacle1.NextObstacle

		var relativePosition1, relativePosition2 *Vector2
		relativePosition1 = Sub(obstacle1.Point, a.Position)
		relativePosition2 = Sub(obstacle2.Point, a.Position)

		/*
		 * Check if velocity obstacle of obstacle is already taken care of by
		 * previously constructed obstacle ORCA lines.
		 */
		var alreadyCovered bool
		alreadyCovered = false

		for j := 0; j < len(a.OrcaLines); j++ {
			if Det(Sub(MulOne(relativePosition1, invTimeHorizonObst), a.OrcaLines[j].Point), a.OrcaLines[j].Direction)-invTimeHorizonObst*a.Radius >= -RVO_EPSILON && Det(Sub(MulOne(relativePosition2, invTimeHorizonObst), a.OrcaLines[j].Point), a.OrcaLines[j].Direction)-invTimeHorizonObst*a.Radius >= -RVO_EPSILON {
				alreadyCovered = true
				break
			}
		}

		if alreadyCovered {
			continue
		}

		/* Not yet covered. Check for collisions. */

		var distSq1, distSq2, radiusSq float64
		distSq1 = Sqr(relativePosition1)
		distSq2 = Sqr(relativePosition2)

		radiusSq = math.Pow(a.Radius, 2)

		var obstacleVector *Vector2
		obstacleVector = Sub(obstacle2.Point, obstacle1.Point)
		// (-relativePosition1 * obstacleVector) / absSq(obstacleVector)
		var s, distSqLine float64
		s = Mul(Flip(relativePosition1), obstacleVector) / Sqr(obstacleVector)
		// bsSq(-relativePosition1 - s * obstacleVector)
		distSqLine = Sqr(Sub(Flip(relativePosition1), MulOne(obstacleVector, s)))

		var line Line

		if s < 0 && distSq1 <= radiusSq {
			/* Collision with left vertex. Ignore if non-convex. */
			if obstacle1.IsConvex {
				line.Point = NewVector2(0, 0)
				line.Direction = Normalize(NewVector2(-relativePosition1.Y, relativePosition1.X))
				a.OrcaLines = append(a.OrcaLines, &line)
				//				fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(f): %v \n ORCALineDir(f): %v \n\n",
				//					line.Point, line.Direction)
			}

			continue
		} else if s > 1 && distSq2 <= radiusSq {
			/* Collision with right vertex. Ignore if non-convex
			 * or if it will be taken care of by neighoring obstace */
			if obstacle2.IsConvex && Det(relativePosition2, obstacle2.UnitDir) >= 0 {
				line.Point = NewVector2(0, 0)
				line.Direction = Normalize(NewVector2(-relativePosition2.Y, relativePosition2.X))
				a.OrcaLines = append(a.OrcaLines, &line)
				//				fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(g): %v \n ORCALineDir(g): %v \n\n",
				//					line.Point, line.Direction)
			}

			continue
		} else if s >= 0 && s < 1 && distSqLine <= radiusSq {
			/* Collision with obstacle segment. */
			line.Point = NewVector2(0, 0)
			line.Direction = Flip(obstacle1.UnitDir)
			a.OrcaLines = append(a.OrcaLines, &line)
			//			fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(h): %v \n ORCALineDir(h): %v \n\n",
			//				line.Point, line.Direction)
			continue
		}

		/*
		 * No collision.
		 * Compute legs. When obliquely viewed, both legs can come from a single
		 * vertex. Legs extend cut-off line when nonconvex vertex.
		 */

		var leftLegDirection, rightLegDirection *Vector2
		var leg1, leg2 float64

		if s < 0 && distSqLine <= radiusSq {
			/*
			 * Obstacle viewed obliquely so that left vertex
			 * defines velocity obstacle.
			 */
			if !obstacle1.IsConvex {
				/* Ignore obstacle. */
				continue
			}

			obstacle2 = obstacle1

			leg1 = math.Sqrt(distSq1 - radiusSq)
			leftLegDirection = Div(NewVector2(relativePosition1.X*leg1-relativePosition1.Y*a.Radius, relativePosition1.X*a.Radius+relativePosition1.Y*leg1), distSq1)
			rightLegDirection = Div(NewVector2(relativePosition1.X*leg1+relativePosition1.Y*a.Radius, -relativePosition1.X*a.Radius+relativePosition1.Y*leg1), distSq1)
		} else if s > 1 && distSqLine <= radiusSq {
			/*
			 * Obstacle viewed obliquely so that
			 * right vertex defines velocity obstacle.
			 */
			if !obstacle2.IsConvex {
				/* Ignore obstacle. */
				continue
			}

			obstacle1 = obstacle2

			leg2 = math.Sqrt(distSq2 - radiusSq)
			leftLegDirection = Div(NewVector2(relativePosition2.X*leg2-relativePosition2.Y*a.Radius, relativePosition2.X*a.Radius+relativePosition2.Y*leg2), distSq2)
			rightLegDirection = Div(NewVector2(relativePosition2.X*leg2+relativePosition2.Y*a.Radius, -relativePosition2.X*a.Radius+relativePosition2.Y*leg2), distSq2)
		} else {
			/* Usual situation. */
			if obstacle1.IsConvex {
				leg1 = math.Sqrt(distSq1 - radiusSq)
				leftLegDirection = Div(NewVector2(relativePosition1.X*leg1-relativePosition1.Y*a.Radius, relativePosition1.X*a.Radius+relativePosition1.Y*leg1), distSq1)
			} else {
				/* Left vertex non-convex; left leg extends cut-off line. */
				leftLegDirection = Flip(obstacle1.UnitDir)
			}

			if obstacle2.IsConvex {
				leg2 = math.Sqrt(distSq2 - radiusSq)
				rightLegDirection = Div(NewVector2(relativePosition2.X*leg2+relativePosition2.Y*a.Radius, -relativePosition2.X*a.Radius+relativePosition2.Y*leg2), distSq2)
			} else {
				/* Right vertex non-convex; right leg extends cut-off line. */
				rightLegDirection = obstacle1.UnitDir
			}
		}

		/*
		 * Legs can never point into neighboring edge when convex vertex,
		 * take cutoff-line of neighboring edge instead. If velocity projected on
		 * "foreign" leg, no constraint is added.
		 */

		var leftNeighbor *Obstacle
		leftNeighbor = obstacle1.PrevObstacle

		var isLeftLegForeign, isRightLegForeign bool
		isLeftLegForeign = false
		isRightLegForeign = false

		if obstacle1.IsConvex && Det(leftLegDirection, Flip(leftNeighbor.UnitDir)) >= 0 {
			/* Left leg points into obstacle. */
			leftLegDirection = Flip(leftNeighbor.UnitDir)
			isLeftLegForeign = true
		}

		if obstacle2.IsConvex && Det(rightLegDirection, obstacle2.UnitDir) <= 0 {
			/* Right leg points into obstacle. */
			rightLegDirection = obstacle2.UnitDir
			isRightLegForeign = true
		}

		/* Compute cut-off centers. */
		//invTimeHorizonObst * (obstacle1.Point - a.Position)
		leftCutoff := MulOne(Sub(obstacle1.Point, a.Position), invTimeHorizonObst)
		//invTimeHorizonObst * (obstacle2.Point - a.Position)
		rightCutoff := MulOne(Sub(obstacle2.Point, a.Position), invTimeHorizonObst)
		cutoffVec := Sub(rightCutoff, leftCutoff)

		/* Project current velocity on velocity obstacle. */

		/* Check if current velocity is projected on cutoff circles. */
		var t float64
		//t := (obstacle1 == obstacle2 ? 0.5f : ((a.Velocity - leftCutoff) * cutoffVec) / Sqr(cutoffVec));
		if obstacle1 == obstacle2 {
			t = 0.5
		} else {
			t = Mul(Sub(a.Velocity, leftCutoff), cutoffVec) / Sqr(cutoffVec)
		}
		var tLeft, tRight float64
		var unitW *Vector2
		tLeft = Mul(Sub(a.Velocity, leftCutoff), leftLegDirection)
		tRight = Mul(Sub(a.Velocity, rightCutoff), rightLegDirection)

		if (t < 0 && tLeft < 0) || (obstacle1 == obstacle2 && tLeft < 0 && tRight < 0) {
			/* Project on left cut-off circle. */
			unitW = Normalize(Sub(a.Velocity, leftCutoff))

			line.Direction = NewVector2(unitW.Y, -unitW.X)
			line.Point = Add(leftCutoff, MulOne(unitW, a.Radius*invTimeHorizonObst))
			a.OrcaLines = append(a.OrcaLines, &line)
			//			fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(e): %v \n ORCALineDir(e): %v \n\n",
			//				line.Point, line.Direction)
			continue
		} else if t > 1 && tRight < 0 {
			/* Project on right cut-off circle. */
			unitW = Normalize(Sub(a.Velocity, rightCutoff))

			line.Direction = NewVector2(unitW.Y, -unitW.X)
			line.Point = Add(rightCutoff, MulOne(unitW, a.Radius*invTimeHorizonObst))
			a.OrcaLines = append(a.OrcaLines, &line)
			//			fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(d): %v \n ORCALineDir(d): %v \n\n",
			//				line.Point, line.Direction)
			continue
		}

		/*
		 * Project on left leg, right leg, or cut-off line, whichever is closest
		 * to velocity.
		 */
		var distSqCutoff, distSqLeft, distSqRight float64
		if t < 0 || t > 1 || obstacle1 == obstacle2 {
			distSqCutoff = math.Inf(0) // positive infinity
		} else {
			distSqCutoff = Sqr(Sub(a.Velocity, Add(leftCutoff, MulOne(cutoffVec, t))))
		}
		if tLeft < 0 {
			distSqLeft = math.Inf(0) // positive infinity
		} else {
			distSqLeft = Sqr(Sub(a.Velocity, Add(leftCutoff, MulOne(leftLegDirection, tLeft))))
		}
		if tRight < 0 {
			distSqRight = math.Inf(0) // positive infinity
		} else {
			distSqRight = Sqr(Sub(a.Velocity, Add(rightCutoff, MulOne(rightLegDirection, tRight))))
		}

		if distSqCutoff <= distSqLeft && distSqCutoff <= distSqRight {
			/* Project on cut-off line. */
			line.Direction = Flip(obstacle1.UnitDir)
			line.Point = Add(leftCutoff, MulOne(NewVector2(-line.Direction.Y, line.Direction.X), a.Radius*invTimeHorizonObst))
			//line.Point = leftCutoff + a.Radius * invTimeHorizonObst * NewVector2(-line.Direction.Y(), line.Direction.X());
			a.OrcaLines = append(a.OrcaLines, &line)
			//			fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(c): %v \n ORCALineDir(c): %v \n\n",
			//				line.Point, line.Direction)
			continue
		} else if distSqLeft <= distSqRight {
			/* Project on left leg. */
			if isLeftLegForeign {
				continue
			}

			line.Direction = leftLegDirection
			line.Point = Add(leftCutoff, MulOne(NewVector2(-line.Direction.Y, line.Direction.X), a.Radius*invTimeHorizonObst))
			//line.Point = leftCutoff + a.Radius * invTimeHorizonObst * NewVector2(-line.Direction.Y(), line.Direction.X());
			a.OrcaLines = append(a.OrcaLines, &line)
			//			fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(b): %v \n ORCALineDir(b): %v \n\n",
			//				line.Point, line.Direction)
			continue
		} else {
			/* Project on right leg. */
			if isRightLegForeign {
				continue
			}

			line.Direction = Flip(rightLegDirection)
			line.Point = Add(rightCutoff, MulOne(NewVector2(-line.Direction.Y, line.Direction.X), a.Radius*invTimeHorizonObst))
			//line.Point = rightCutoff + a.Radius * invTimeHorizonObst * NewVector2(-line.Direction.Y(), line.Direction.X());
			a.OrcaLines = append(a.OrcaLines, &line)
			//			fmt.Printf("FN: computeNewVelocity \n Obstacle ORCALines \n ORCALinePoint(a): %v \n ORCALineDir(a): %v \n\n",
			//				line.Point, line.Direction)
			continue
		}
	}

	var numObstLines int
	numObstLines = len(a.OrcaLines)

	var invTimeHorizon float64
	invTimeHorizon = float64(1) / a.TimeHorizon

	/* Create agent ORCA lines. */
	for i := 0; i < len(a.AgentNeighbors); i++ {
		var other *Agent
		other = a.AgentNeighbors[i].Agent

		var relativePosition, relativeVelocity *Vector2
		relativePosition = Sub(other.Position, a.Position)
		relativeVelocity = Sub(a.Velocity, other.Velocity)

		var distSq, combinedRadius, combinedRadiusSq float64
		distSq = Sqr(relativePosition)
		combinedRadius = a.Radius + other.Radius
		combinedRadiusSq = math.Pow(combinedRadius, 2)

		var line Line
		var u, w, unitW *Vector2
		var wLengthSq, wLength float64
		if distSq > combinedRadiusSq {
			/* No collision. */
			w = Sub(relativeVelocity, MulOne(relativePosition, invTimeHorizon))
			/* Vector from cutoff center to relative velocity. */

			var dotProduct1, dotProduct2, leg float64
			wLengthSq = Sqr(w)
			dotProduct1 = Mul(w, relativePosition)

			if dotProduct1 < 0 && math.Pow(dotProduct1, 2) > combinedRadiusSq*wLengthSq {
				/* Project on cut-off circle. */
				wLength = math.Sqrt(wLengthSq)
				unitW = Div(w, wLength)

				line.Direction = NewVector2(unitW.Y, -unitW.X)
				u = MulOne(unitW, (combinedRadius*invTimeHorizon - wLength))
			} else {
				/* Project on legs. */
				leg = math.Sqrt(distSq - combinedRadiusSq)

				if Det(relativePosition, w) > 0 {
					/* Project on left leg. */
					line.Direction = Div(NewVector2(relativePosition.X*leg-relativePosition.Y*combinedRadius, relativePosition.X*combinedRadius+relativePosition.Y*leg), distSq)
				} else {
					/* Project on right leg. */
					line.Direction = Flip(Div(NewVector2(relativePosition.X*leg+relativePosition.Y*combinedRadius, -relativePosition.X*combinedRadius+relativePosition.Y*leg), distSq))
				}

				dotProduct2 = Mul(relativeVelocity, line.Direction)

				u = Sub(MulOne(line.Direction, dotProduct2), relativeVelocity)
			}
		} else {
			/* Collision. Project on cut-off circle of time timeStep. */
			var invTimeStep float64
			invTimeStep = float64(1) / Sim.TimeStep

			/* Vector from cutoff center to relative velocity. */
			w = Sub(relativeVelocity, MulOne(relativePosition, invTimeStep))

			wLength = Abs(w)
			unitW = Div(w, wLength)

			line.Direction = NewVector2(unitW.Y, -unitW.X)
			u = MulOne(unitW, combinedRadius*invTimeStep-wLength)
		}

		line.Point = Add(a.Velocity, MulOne(u, 0.5))
		a.OrcaLines = append(a.OrcaLines, &line)
		//		fmt.Printf("FN: computeNewVelocity \n Agent ORCALines \n ORCALinePoint(a): %v \n ORCALineDir(a): %v \n\n",
		//			line.Point, line.Direction)
	}
	//fmt.Printf("orca %v", a.OrcaLines)

	lineFail := a.LinearProgram2(a.OrcaLines, a.MaxSpeed, a.PrefVelocity, false)
	//	fmt.Printf("FN: computeNewVelocity \n lineFail: %v \n len(ORCALines): %v \n newVelocity: %v \n\n",
	//		lineFail, len(a.OrcaLines), a.NewVelocity)
	if lineFail < len(a.OrcaLines) {
		a.LinearProgram3(a.OrcaLines, numObstLines, lineFail, a.MaxSpeed)
	}

}

// !
// FINISH
func (a *Agent) InsertAgentNeighbor(agent *Agent, rangeSq float64) {

	if a != agent { // ?
		distSq := Sqr(Sub(a.Position, agent.Position))

		// 2Agent間の距離が半径よりも近かった場合
		if distSq < rangeSq {
			if len(a.AgentNeighbors) < a.MaxNeighbors {
				a.AgentNeighbors = append(a.AgentNeighbors, &AgentNeighbor{DistSq: distSq, Agent: agent})
			}

			i := len(a.AgentNeighbors) - 1

			for {
				if i != 0 && distSq < a.AgentNeighbors[i-1].DistSq {
					a.AgentNeighbors[i] = a.AgentNeighbors[i-1]
					i--
				} else {
					break
				}
			}

			a.AgentNeighbors[i] = &AgentNeighbor{
				DistSq: distSq,
				Agent:  agent,
			}

			if len(a.AgentNeighbors) == a.MaxNeighbors {
				rangeSq = a.AgentNeighbors[len(a.AgentNeighbors)-1].DistSq
			}
		}
	}
}

// FINISH
func (a *Agent) InsertObstacleNeighbor(obstacle *Obstacle, rangeSq float64) {
	nextObstacle := obstacle.NextObstacle

	distSq := DistSqPointLineSegment(obstacle.Point, nextObstacle.Point, a.Position)

	if distSq < rangeSq {
		a.ObstacleNeighbors = append(a.ObstacleNeighbors, &ObstacleNeighbor{DistSq: distSq, Obstacle: obstacle})

		i := len(a.ObstacleNeighbors) - 1

		for {
			if i != 0 && distSq < a.ObstacleNeighbors[i-1].DistSq {
				a.ObstacleNeighbors[i] = a.ObstacleNeighbors[i-1]
				i--
			} else {
				break
			}
		}

		a.ObstacleNeighbors[i] = &ObstacleNeighbor{
			DistSq:   distSq,
			Obstacle: obstacle,
		}
	}
}

// FINISH
func (a *Agent) Update() {
	//	fmt.Printf("FN: Update  \n NewVelocity: %v \n Position: %v\n TimeStep %v\n MulOne %v\n AddSum %v\n\n",
	//		a.NewVelocity, a.Position, Sim.TimeStep, MulOne(a.NewVelocity, Sim.TimeStep), AddSum(a.Position, MulOne(a.Velocity, Sim.TimeStep)))
	a.Velocity = a.NewVelocity
	a.Position = Add(a.Position, MulOne(a.Velocity, Sim.TimeStep))
}

// FINISH
func (a *Agent) LinearProgram1(lines []*Line, lineNo int, radius float64, optVelocity *Vector2, directionOpt bool) bool {
	var dotProduct, discriminant float64
	dotProduct = Mul(lines[lineNo].Point, lines[lineNo].Direction)
	discriminant = math.Pow(dotProduct, 2) + math.Pow(radius, 2) - Sqr(lines[lineNo].Point)

	if discriminant < 0 {
		/* Max speed circle fully invalidates line lineNo. */
		return false
	}

	var sqrtDiscriminant, tLeft, tRight float64
	sqrtDiscriminant = math.Sqrt(discriminant)
	tLeft = -dotProduct - sqrtDiscriminant
	tRight = -dotProduct + sqrtDiscriminant

	for i := 0; i < lineNo; i++ {
		var denominator, numerator float64
		denominator = Det(lines[lineNo].Direction, lines[i].Direction)
		numerator = Det(lines[i].Direction, Sub(lines[lineNo].Point, lines[i].Point))

		if math.Abs(denominator) <= RVO_EPSILON {
			/* Lines lineNo and i are (almost) parallel. */
			if numerator < 0 {
				//				fmt.Printf("FN: LinearProgram1  false1 \n")
				return false
			} else {
				continue
			}
		}

		var t float64
		t = numerator / denominator

		if denominator >= 0 {
			/* Line i bounds line lineNo on the right. */
			tRight = math.Min(tRight, t)
		} else {
			/* Line i bounds line lineNo on the left. */
			tLeft = math.Max(tLeft, t)
		}

		if tLeft > tRight {
			//			fmt.Printf("FN: LinearProgram1  false2 \n")
			return false
		}
	}

	if directionOpt {
		/* Optimize direction. */
		if Mul(optVelocity, lines[lineNo].Direction) > 0 {
			/* Take right extreme. */
			a.NewVelocity = Add(lines[lineNo].Point, MulOne(lines[lineNo].Direction, tRight))

			//			fmt.Printf("FN: LinearProgram1  \n a.NewVelocity1: %v \n\n",
			//				a.NewVelocity)
		} else {
			/* Take left extreme. */
			a.NewVelocity = Add(lines[lineNo].Point, MulOne(lines[lineNo].Direction, tLeft))

			//			fmt.Printf("FN: LinearProgram1  \n a.NewVelocity2: %v \n\n",
			//				a.NewVelocity)
		}
	} else {
		/* Optimize closest point. */
		t := Mul(lines[lineNo].Direction, Sub(optVelocity, lines[lineNo].Point))

		if t < tLeft {
			a.NewVelocity = Add(lines[lineNo].Point, MulOne(lines[lineNo].Direction, tLeft))

			//			fmt.Printf("FN: LinearProgram1  \n a.NewVelocity3: %v \n\n",
			//				a.NewVelocity)
		} else if t > tRight {
			a.NewVelocity = Add(lines[lineNo].Point, MulOne(lines[lineNo].Direction, tRight))

			//			fmt.Printf("FN: LinearProgram1  \n a.NewVelocity4: %v \n\n",
			//				a.NewVelocity)
		} else {
			a.NewVelocity = Add(lines[lineNo].Point, MulOne(lines[lineNo].Direction, t))

			//			fmt.Printf("FN: LinearProgram1  \n a.NewVelocity5: %v \n a.NewVelocity: %v\n\n",
			//				a.NewVelocity, a.NewVelocity)
		}
	}

	return true
}

// FINISH
func (a *Agent) LinearProgram2(lines []*Line, radius float64, optVelocity *Vector2, directionOpt bool) int {
	//	fmt.Printf("FN: LinearProgram2\n")
	if directionOpt {
		/*
		 * Optimize direction. Note that the optimization velocity is of unit
		 * length in this case.
		 */
		a.NewVelocity = MulOne(optVelocity, radius)

		//		fmt.Printf("FN: LinearProgram2  \n a.NewVelocity1: %v \n\n",
		//			a.NewVelocity)
	} else if Sqr(optVelocity) > math.Pow(radius, 2) {
		/* Optimize closest point and outside circle. */
		a.NewVelocity = MulOne(Normalize(optVelocity), radius)

		//		fmt.Printf("FN: LinearProgram2  \n a.NewVelocity2: %v \n\n",
		//			a.NewVelocity)
	} else {
		/* Optimize closest point and inside circle. */
		a.NewVelocity = optVelocity

		//		fmt.Printf("FN: LinearProgram2  \n a.NewVelocity3: %v \n\n",
		//			a.NewVelocity)
	}

	for i := 0; i < len(lines); i++ {
		//		fmt.Printf("FN: LinearProgram2 i %v \n\n",
		//			i)
		if Det(lines[i].Direction, Sub(lines[i].Point, a.NewVelocity)) > 0 {
			/* Result does not satisfy constraint i. Compute new optimal a.NewVelocity. */
			var tempResult *Vector2
			tempResult = a.NewVelocity

			//Cause
			//			fmt.Printf("FN: LinearProgram2  \n toLP1\n len(lines): %v \n radius: %v\n optVelocity: %v \n directonOpt: %v \n a.NewVelocity: %v \n\n", len(lines), radius, optVelocity, directionOpt, a.NewVelocity)
			if a.LinearProgram1(lines, i, radius, optVelocity, directionOpt) == false {
				a.NewVelocity = tempResult

				//				fmt.Printf("FN: LinearProgram2  \n a.NewVelocity4: %v \n\n",
				//					a.NewVelocity)
				return i
			}
		}
	}
	return len(lines)
}

// ! ?? projLine make size
//FINISH
func (a *Agent) LinearProgram3(lines []*Line, numObstLines int, beginLine int, radius float64) {
	//	fmt.Printf("FN: LinearProgram3  \n lineSize: %v\n beginLine: %v \n\n", len(lines), beginLine)
	var distance float64
	distance = 0.0
	for i := beginLine; i < len(lines); i++ {
		//		fmt.Printf("FN: LinearProgram3  \n lineDirection: %v\n linePoint: %v\n a.NewVelocity: %v \n\n", lines[i].Direction, lines[i].Point, a.NewVelocity)
		if Det(lines[i].Direction, Sub(lines[i].Point, a.NewVelocity)) > distance {
			/* Result does not satisfy constraint of line i. */
			// ? std::vector<Line> projLines(lines.begin(), lines.begin() + static_cast<ptrdiff_t>(numObstLines));
			var projLines []*Line
			projLines = make([]*Line, 0) // ?

			for j := numObstLines; j < i; j++ {
				var line Line
				var determinant float64
				determinant = Det(lines[i].Direction, lines[j].Direction)

				if math.Abs(determinant) <= RVO_EPSILON {
					/* Line i and line j are parallel. */
					if Mul(lines[i].Direction, lines[j].Direction) > 0 {
						/* Line i and line j point in the same direction. */
						continue
					} else {
						/* Line i and line j point in opposite direction. */
						line.Point = MulOne(Add(lines[i].Point, lines[j].Point), 0.5)
					}
				} else {
					// line.Point = lines[i].Point + (det(lines[j].Direction, lines[i].Point - lines[j].Point) / determinant) * lines[i].Direction;
					line.Point = Add(lines[i].Point, MulOne(lines[i].Direction, (Det(lines[j].Direction, Sub(lines[i].Point, lines[j].Point))/determinant)))
				}

				line.Direction = Normalize(Sub(lines[j].Direction, lines[i].Direction))

				projLines = append(projLines, &line)
			}

			var tempResult *Vector2
			tempResult = a.NewVelocity

			if a.LinearProgram2(projLines, radius, NewVector2(-lines[i].Direction.Y, lines[i].Direction.X), true) < len(projLines) {
				/* This should in principle not happen.  The a.NewVelocity is by definition
				 * already in the feasible region of this linear program. If it fails,
				 * it is due to small floating point error, and the current a.NewVelocity is
				 * kept.
				 */

				a.NewVelocity = tempResult

				//				fmt.Printf("FN: LinearProgram3  \n a.NewVelocity1: %v \n\n",
				//					a.NewVelocity)
			}

			distance = Det(lines[i].Direction, Sub(lines[i].Point, a.NewVelocity))
		}
	}
}
