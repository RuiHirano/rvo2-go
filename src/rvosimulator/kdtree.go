package rvosimulator

import (
	"fmt"
	"math"
)

var (
	MAX_LEAF_SIZE int
)

func init() {
	MAX_LEAF_SIZE = 10
}

type KdTree struct {
	ObstacleTree *ObstacleTreeNode
	AgentTree    []*AgentTreeNode
	Agents       []*Agent
	Obstacles    []*Obstacle
}

type AgentTreeNode struct {
	Begin int
	End   int
	Left  int
	Right int
	MaxX  float64
	MaxY  float64
	MinX  float64
	MinY  float64
}

type ObstacleTreeNode struct {
	Left     *ObstacleTreeNode
	Right    *ObstacleTreeNode
	Obstacle *Obstacle
}

func NewAgentTreeNode() *AgentTreeNode {
	a := &AgentTreeNode{
		Begin: 0,
		End:   0,
		Left:  0,
		Right: 0,
		MaxX:  0,
		MaxY:  0,
		MinX:  0,
		MinY:  0,
	}
	return a
}

func NewObstacleTreeNode() *ObstacleTreeNode {
	o := &ObstacleTreeNode{
		Left:     nil,
		Right:    nil,
		Obstacle: NewObstacle(),
	}
	return o
}

func NewKdTree() *KdTree {
	//agents := make([]*Agent, 0)
	//agentTreeNode := make([]*AgentTreeNode, 0)
	//obstacles := make([]*Obstacle, 0)
	//obstacleTree := make([]*ObstacleTreeNode, 0)
	k := &KdTree{
		ObstacleTree: nil,
	}
	return k
}

// CHECK OK
// ?
func (kt *KdTree) BuildAgentTree() {
	//fmt.Printf("FN: BuildAgentTree %v\n", Sim.Agents)
	// sim ...
	if len(kt.Agents) < len(Sim.Agents) {
		for i := len(kt.Agents); i < len(Sim.Agents); i++ {
			//			fmt.Printf("FN: BuildAgentTree\n ID %v\n Position: %v\n MaxNeighbors: %v\n MaxSpeed; %v\nNeighborDist: %v\n Radius: %v\n TimeHorizon: %v\n TimeHorizonObst: %v\n Velocity: %v\n\n",
			//				Sim.Agents[i].ID, Sim.Agents[i].Position, Sim.Agents[i].MaxNeighbors, Sim.Agents[i].MaxSpeed, Sim.Agents[i].NeighborDist, Sim.Agents[i].Radius, Sim.Agents[i].TimeHorizon, Sim.Agents[i].TimeHorizonObst, Sim.Agents[i].Velocity)

			kt.Agents = append(kt.Agents, Sim.Agents[i])
		}
		// AgentTreeを2*len(kt.Agents)-1で初期化
		kt.AgentTree = make([]*AgentTreeNode, 2*len(kt.Agents)-1)
		for i := 0; i < 2*len(kt.Agents)-1; i++ {
			kt.AgentTree[i] = NewAgentTreeNode()

		}

	}

	if len(kt.Agents) != 0 {
		//		fmt.Printf("FN: To BuildAgentTreeRecursive\n")
		kt.BuildAgentTreeRecursive(0, len(kt.Agents), 0)
	}
}

// Position is wrong ,,, node Max is Wrong
// FINISH
func (kt *KdTree) BuildAgentTreeRecursive(begin int, end int, node int) {
	//fmt.Printf("node: %v", node)
	//fmt.Printf("node: %v", kt.AgentTree[node])
	kt.AgentTree[node].Begin = begin
	kt.AgentTree[node].End = end
	kt.AgentTree[node].MinX = kt.Agents[begin].Position.X
	kt.AgentTree[node].MaxX = kt.Agents[begin].Position.X
	kt.AgentTree[node].MinY = kt.Agents[begin].Position.Y
	kt.AgentTree[node].MaxY = kt.Agents[begin].Position.Y

	// i-1番目までのAgentとi番目のAgentのポジションを比較
	for i := begin + 1; i < end; i++ {
		kt.AgentTree[node].MaxX = math.Max(kt.AgentTree[node].MaxX, kt.Agents[i].Position.X)
		kt.AgentTree[node].MinX = math.Min(kt.AgentTree[node].MinX, kt.Agents[i].Position.X)
		kt.AgentTree[node].MaxY = math.Max(kt.AgentTree[node].MaxY, kt.Agents[i].Position.Y)
		kt.AgentTree[node].MinY = math.Min(kt.AgentTree[node].MinY, kt.Agents[i].Position.Y)
		//		fmt.Printf("FN: BuildAgentTreeRecurcive\n Begin %v\n End: %v\n Left: %v\n Right; %v\n MaxX %v\n MaxY: %v\n MinX: %v\n MinY; %v\n\n",
		//			kt.AgentTree[node].Begin, kt.AgentTree[node].End, kt.AgentTree[node].Left, kt.AgentTree[node].Right, kt.AgentTree[node].MaxX, kt.AgentTree[node].MaxY, kt.AgentTree[node].MinX, kt.AgentTree[node].MinY)
	}

	if end-begin > MAX_LEAF_SIZE {
		isVertical := (kt.AgentTree[node].MaxX-kt.AgentTree[node].MinX > kt.AgentTree[node].MaxY-kt.AgentTree[node].MinY)
		left := begin
		right := end
		var splitValue, leftPosition, rightPosition float64

		if isVertical {
			splitValue = 0.5 * (kt.AgentTree[node].MaxX + kt.AgentTree[node].MinX)
			leftPosition = kt.Agents[left].Position.X
			rightPosition = kt.Agents[right-1].Position.X
		} else {
			splitValue = 0.5 * (kt.AgentTree[node].MaxY + kt.AgentTree[node].MinY)
			leftPosition = kt.Agents[left].Position.Y
			rightPosition = kt.Agents[right-1].Position.Y
		}

		for {
			if left < right {
				for {

					if left < right && leftPosition < splitValue {
						left++
					} else {
						break
					}

				}

				for {
					if right > left && rightPosition >= splitValue {
						right--
					} else {
						break
					}
				}

				if left < right {
					// swap array
					t := kt.Agents[left]
					fmt.Printf("FN: right-1 %v\n", right-1)
					kt.Agents[left] = kt.Agents[right-1]
					kt.Agents[right-1] = t

					left++
					right--
				}
			} else {
				break
			}

		}

		if left == begin {
			left++
			right++
		}

		kt.AgentTree[node].Left = node + 1
		kt.AgentTree[node].Right = node + 2*(left-begin)

		kt.BuildAgentTreeRecursive(begin, left, kt.AgentTree[node].Left)
		kt.BuildAgentTreeRecursive(left, end, kt.AgentTree[node].Right)

		//		fmt.Printf("FN: BuildAgentTreeRecurcive\n > MAX_LEAF \n Begin %v\n End: %v\n Left: %v\n Right; %v\n MaxX %v\n MaxY: %v\n MinX: %v\n MinY; %v\n\n",
		//			kt.AgentTree[node].Begin, kt.AgentTree[node].End, kt.AgentTree[node].Left, kt.AgentTree[node].Right, kt.AgentTree[node].MaxX, kt.AgentTree[node].MaxY, kt.AgentTree[node].MinX, kt.AgentTree[node].MinY)
	}
}

// CHECK OK
func (kt *KdTree) BuildObstacleTree() {

	//fmt.Printf("FN: BuildObstacleTree\n")
	kt.DeleteObstacleTree(kt.ObstacleTree)

	// sim.Obstaclesをobstaclesに格納する
	obstacles := make([]*Obstacle, len(Sim.Obstacles))
	//fmt.Printf("FN: BuildObstacleTree\n ID %v\n Point: %v\n\n",
	//	len(obstacles), len(Sim.Obstacles))
	for i := 0; i < len(Sim.Obstacles); i++ {
		obstacles[i] = Sim.Obstacles[i]
		//		fmt.Printf("FN: BuildObstacleTree\n ID %v\n Point: %v\n IsConvex: %v\n UnitDir; %v\n\n",
		//			obstacles[i].ID, obstacles[i].Point, obstacles[i].IsConvex, obstacles[i].UnitDir)
	}

	kt.ObstacleTree = kt.BuildObstacleTreeRecursive(obstacles)
}

// CHECK OK
// FINISH
func (kt *KdTree) BuildObstacleTreeRecursive(obstacles []*Obstacle) *ObstacleTreeNode {
	//	fmt.Printf("FN: BuildObstacleTreeRecursive\n len(obstacles): %v \n\n",
	//		len(obstacles))
	if len(obstacles) == 0 {
		//var node *ObstacleTreeNode
		//		fmt.Printf("FN: BuildObstacleTreeRecursive\n len(obstacles) == 0 \n ObstacleID: %v\n\n",
		//			node)
		return nil
	} else {
		node := NewObstacleTreeNode()
		//fmt.Printf("FN: BuildObstacleTreeRecursive\n")

		optimalSplit := 0
		minLeft := len(obstacles)
		minRight := len(obstacles)

		// Compute optimal split node. 左右の最適な分け方を計算
		for i := 0; i < len(obstacles); i++ {
			leftSize := 0
			rightSize := 0

			var obstacleI1, obstacleI2 *Obstacle
			obstacleI1 = obstacles[i]
			obstacleI2 = obstacleI1.NextObstacle

			for j := 0; j < len(obstacles); j++ {
				if i == j {
					continue
				}
				var obstacleJ1, obstacleJ2 *Obstacle
				obstacleJ1 = obstacles[j]
				obstacleJ2 = obstacleJ1.NextObstacle

				j1LeftOfI := LeftOf(obstacleI1.Point, obstacleI2.Point, obstacleJ1.Point)
				j2LeftOfI := LeftOf(obstacleI1.Point, obstacleI2.Point, obstacleJ2.Point)

				if j1LeftOfI >= -RVO_EPSILON && j2LeftOfI >= -RVO_EPSILON {
					leftSize++
				} else if j1LeftOfI <= RVO_EPSILON && j2LeftOfI <= RVO_EPSILON {
					rightSize++
				} else {
					leftSize++
					rightSize++
				}

				if math.Max(float64(leftSize), float64(rightSize)) > math.Max(float64(minLeft), float64(minRight)) || (math.Max(float64(leftSize), float64(rightSize)) == math.Max(float64(minLeft), float64(minRight)) && math.Min(float64(leftSize), float64(rightSize)) >= math.Min(float64(minLeft), float64(minRight))) {
					break
				}

			}

			// ここを通過する回数がおかしい
			// 二回目のleft, rightがおかしい
			//			fmt.Printf("FN: BuildObstacleTreeRecursive\n i: %v\n leftSize: %v\n rightSize %v \n minLeft:  %v \nminRight:  %v  \n\n",
			//				i, leftSize, rightSize, minLeft, minRight)
			if math.Max(float64(leftSize), float64(rightSize)) < math.Max(float64(minLeft), float64(minRight)) || (math.Max(float64(leftSize), float64(rightSize)) == math.Max(float64(minLeft), float64(minRight)) && math.Min(float64(leftSize), float64(rightSize)) < math.Min(float64(minLeft), float64(minRight))) {
				minLeft = leftSize
				minRight = rightSize
				optimalSplit = i
				//				fmt.Printf("FN: BuildObstacleTreeRecursive optimalSplit &v\n\n", i)
			}

		}

		// Build split node.
		var leftObstacles, rightObstacles []*Obstacle
		leftObstacles = make([]*Obstacle, minLeft)
		rightObstacles = make([]*Obstacle, minRight)

		leftCounter := 0
		rightCounter := 0
		i := optimalSplit //　ここが原因

		var obstacleI1, obstacleI2 *Obstacle
		obstacleI1 = obstacles[i]
		obstacleI2 = obstacleI1.NextObstacle
		fmt.Printf("FN: BuildObstacleTreeRecursive\n obstacleI1: %v\n\n",
			obstacles[0].NextObstacle)

		for j := 0; j < len(obstacles); j++ {
			if i == j {
				continue
			}

			var obstacleJ1, obstacleJ2 *Obstacle
			obstacleJ1 = obstacles[j]
			obstacleJ2 = obstacleJ1.NextObstacle

			var j1LeftOfI, j2LeftOfI float64
			j1LeftOfI = LeftOf(obstacleI1.Point, obstacleI2.Point, obstacleJ1.Point)
			j2LeftOfI = LeftOf(obstacleI1.Point, obstacleI2.Point, obstacleJ2.Point)

			//			fmt.Printf("FN: BuildObstacleTreeRecursive\n j1LeftOfI: %v\n j2LeftOfI: %v \n obI1.Po:  %v \nobI2.Po:  %v \nobJ1.Po:  %v \nobJ2.Po:  %v \n\n",
			//				j1LeftOfI, j2LeftOfI, obstacleI1.Point, obstacleI2.Point, obstacleJ1.Point, obstacleJ2.Point)
			if j1LeftOfI >= -RVO_EPSILON && j2LeftOfI >= -RVO_EPSILON {
				//				fmt.Printf("FN: BuildObstacleTreeRecursive\n if1 \n\n")
				leftObstacles[leftCounter] = obstacles[j]
				leftCounter++
			} else if j1LeftOfI <= RVO_EPSILON && j2LeftOfI <= RVO_EPSILON {
				//				fmt.Printf("FN: BuildObstacleTreeRecursive\n if2 \n\n")
				rightObstacles[rightCounter] = obstacles[j]
				rightCounter++
			} else {
				//				fmt.Printf("FN: BuildObstacleTreeRecursive\n if3 \n\n")
				var t float64
				t = Det(Sub(obstacleI2.Point, obstacleI1.Point), Sub(obstacleJ1.Point, obstacleI1.Point)) / Det(Sub(obstacleI2.Point, obstacleI1.Point), Sub(obstacleJ1.Point, obstacleJ2.Point))

				var splitPoint *Vector2
				splitPoint = Add(obstacleJ1.Point, MulOne(Sub(obstacleJ2.Point, obstacleJ1.Point), t))

				newObstacle := NewObstacle()
				newObstacle.Point = splitPoint
				newObstacle.PrevObstacle = obstacleJ1
				newObstacle.NextObstacle = obstacleJ2
				newObstacle.IsConvex = true
				newObstacle.UnitDir = obstacleJ1.UnitDir
				newObstacle.ID = len(Sim.Obstacles)

				Sim.Obstacles = append(Sim.Obstacles, newObstacle)

				obstacleJ1.NextObstacle = newObstacle
				obstacleJ2.PrevObstacle = newObstacle

				//				fmt.Printf("FN: BuildObstacleTreeRecursive\n leftCounter: %v\n len(leftObstacles) %v \n\n",
				//					leftCounter, len(leftObstacles))
				if j1LeftOfI > 0.0 {
					leftObstacles[leftCounter] = obstacleJ1
					rightObstacles[rightCounter] = newObstacle
					leftCounter++
					rightCounter++
				} else {
					rightObstacles[rightCounter] = obstacleJ1
					leftObstacles[leftCounter] = newObstacle
					leftCounter++
					rightCounter++
				}

			}

		}
		node.Obstacle = obstacleI1
		node.Left = kt.BuildObstacleTreeRecursive(leftObstacles)
		node.Right = kt.BuildObstacleTreeRecursive(rightObstacles)

		//		fmt.Printf("FN: BuildObstacleTreeRecursive\n len(obstacles) != 0 \n ObstacleID: %v\n\n",
		//			node.Obstacle.ID)
		return node
	}
}

//OK
func (kt *KdTree) ComputeAgentNeighbors(agent *Agent, rangeSq float64) {
	kt.QueryAgentTreeRecursive(agent, rangeSq, 0)
}

//OK
func (kt *KdTree) ComputeObstacleNeighbors(agent *Agent, rangeSq float64) {
	kt.QueryObstacleTreeRecursive(agent, rangeSq, kt.ObstacleTree)
}

//OK
func (kt *KdTree) DeleteObstacleTree(node *ObstacleTreeNode) {
	if node != nil {
		kt.DeleteObstacleTree(node.Left)
		kt.DeleteObstacleTree(node.Right)
		//delete node
	}
}

// CHECK OK
// FINISH
func (kt *KdTree) QueryAgentTreeRecursive(agent *Agent, rangeSq float64, node int) {
	if kt.AgentTree[node].End-kt.AgentTree[node].Begin <= MAX_LEAF_SIZE {
		for i := kt.AgentTree[node].Begin; i < kt.AgentTree[node].End; i++ {
			agent.InsertAgentNeighbor(kt.Agents[i], rangeSq)
		}
	} else {
		var distSqLeft, distSqRight float64
		distSqLeft = math.Pow(math.Max(0, kt.AgentTree[kt.AgentTree[node].Left].MinX-agent.Position.X), 2) + math.Pow(math.Max(0, agent.Position.X-kt.AgentTree[kt.AgentTree[node].Left].MaxX), 2) + math.Pow(math.Max(0, kt.AgentTree[kt.AgentTree[node].Left].MinY-agent.Position.Y), 2) + math.Pow(math.Max(0, agent.Position.Y-kt.AgentTree[kt.AgentTree[node].Left].MaxY), 2)

		distSqRight = math.Pow(math.Max(0, kt.AgentTree[kt.AgentTree[node].Right].MinX-agent.Position.X), 2) + math.Pow(math.Max(0, agent.Position.X-kt.AgentTree[kt.AgentTree[node].Right].MaxX), 2) + math.Pow(math.Max(0, kt.AgentTree[kt.AgentTree[node].Right].MinY-agent.Position.Y), 2) + math.Pow(math.Max(0, agent.Position.Y-kt.AgentTree[kt.AgentTree[node].Right].MaxY), 2)

		//		fmt.Printf("FN: QueryAgentTreeRecursive \n distSqLeft: %v\n\n",
		//			distSqLeft)

		//		fmt.Printf("FN: QueryAgentTreeRecursive \n distSqRight: %v\n\n",
		//			distSqRight)
		if distSqLeft < distSqRight {
			if distSqLeft < rangeSq {
				kt.QueryAgentTreeRecursive(agent, rangeSq, kt.AgentTree[node].Left)

				if distSqRight < rangeSq {
					kt.QueryAgentTreeRecursive(agent, rangeSq, kt.AgentTree[node].Right)
				}
			}
		} else {
			if distSqRight < rangeSq {
				kt.QueryAgentTreeRecursive(agent, rangeSq, kt.AgentTree[node].Right)

				if distSqLeft < rangeSq {
					kt.QueryAgentTreeRecursive(agent, rangeSq, kt.AgentTree[node].Left)
				}
			}
		}

	}

}

// CHECK OK
// FINISH
func (kt *KdTree) QueryObstacleTreeRecursive(agent *Agent, rangeSq float64, node *ObstacleTreeNode) {
	//	fmt.Printf("QueryObstacleTreeRecursive %v\n", node.Obstacle)
	if node == nil {
		return
	} else {
		var obstacle1, obstacle2 *Obstacle
		obstacle1 = node.Obstacle
		obstacle2 = obstacle1.NextObstacle

		var agentLeftOfLine, distSqLine float64
		agentLeftOfLine = LeftOf(obstacle1.Point, obstacle2.Point, agent.Position)

		var tNode *ObstacleTreeNode
		if agentLeftOfLine >= 0 {
			tNode = node.Left
		} else {
			tNode = node.Right
		}
		kt.QueryObstacleTreeRecursive(agent, rangeSq, tNode)

		distSqLine = math.Pow(agentLeftOfLine, 2) / Sqr(Sub(obstacle2.Point, obstacle1.Point))

		//		fmt.Printf("FN: QueryObstacleTreeRecursive \n distSqLine: %v\n\n",
		//			distSqLine)

		if distSqLine < rangeSq {
			if agentLeftOfLine < 0 {
				/*
				 * Try obstacle at this node only if agent is on Right side of
				 * obstacle (and can see obstacle).
				 */
				agent.InsertObstacleNeighbor(node.Obstacle, rangeSq)
			}

			/* Try other side of line. */
			var t2Node *ObstacleTreeNode
			if agentLeftOfLine >= 0 {
				t2Node = node.Right
			} else {
				t2Node = node.Left
			}
			kt.QueryObstacleTreeRecursive(agent, rangeSq, t2Node)

		}
	}
}

// FINISH
func (kt *KdTree) QueryVisibility(q1 *Vector2, q2 *Vector2, radius float64) bool {
	result := kt.QueryVisibilityRecursive(q1, q2, radius, kt.ObstacleTree)
	return result
}

// FINISH
func (kt *KdTree) QueryVisibilityRecursive(q1 *Vector2, q2 *Vector2, radius float64, node *ObstacleTreeNode) bool {
	if node == nil {
		return true
	} else {
		var obstacle1, obstacle2 *Obstacle
		obstacle1 = node.Obstacle
		obstacle2 = obstacle1.NextObstacle

		var q1LeftOfI, q2LeftOfI, invLengthI float64
		q1LeftOfI = LeftOf(obstacle1.Point, obstacle2.Point, q1)
		q2LeftOfI = LeftOf(obstacle1.Point, obstacle2.Point, q2)
		invLengthI = float64(1.0) / Sqr(Sub(obstacle2.Point, obstacle1.Point))

		if q1LeftOfI >= 0 && q2LeftOfI >= 0 {
			return kt.QueryVisibilityRecursive(q1, q2, radius, node.Left) && ((math.Pow(q1LeftOfI, 2)*invLengthI >= math.Pow(radius, 2) && math.Pow(q2LeftOfI, 2)*invLengthI >= math.Pow(radius, 2)) || kt.QueryVisibilityRecursive(q1, q2, radius, node.Right))
		} else if q1LeftOfI <= 0 && q2LeftOfI <= 0 {
			return kt.QueryVisibilityRecursive(q1, q2, radius, node.Right) && ((math.Pow(q1LeftOfI, 2)*invLengthI >= math.Pow(radius, 2) && math.Pow(q2LeftOfI, 2)*invLengthI >= math.Pow(radius, 2)) || kt.QueryVisibilityRecursive(q1, q2, radius, node.Left))
		} else if q1LeftOfI >= 0 && q2LeftOfI <= 0 {
			/* One can see through obstacle from Left to Right. */
			return kt.QueryVisibilityRecursive(q1, q2, radius, node.Left) && kt.QueryVisibilityRecursive(q1, q2, radius, node.Right)
		} else {
			var point1LeftOfQ, point2LeftOfQ, invLengthQ float64
			point1LeftOfQ = LeftOf(q1, q2, obstacle1.Point)
			point2LeftOfQ = LeftOf(q1, q2, obstacle2.Point)
			invLengthQ = float64(1.0) / Sqr(Sub(q2, q1))

			return (point1LeftOfQ*point2LeftOfQ >= 0 && math.Pow(point1LeftOfQ, 2)*invLengthQ > math.Pow(radius, 2) && math.Pow(point2LeftOfQ, 2)*invLengthQ > math.Pow(radius, 2) && kt.QueryVisibilityRecursive(q1, q2, radius, node.Left) && kt.QueryVisibilityRecursive(q1, q2, radius, node.Right))
		}
	}
}
