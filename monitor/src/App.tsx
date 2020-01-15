import React, { useState, useEffect } from "react";
import { Scatter } from "react-chartjs-2";
import io from "socket.io-client";
import Slider from "@material-ui/core/Slider";
import { Typography, Button } from "@material-ui/core";

interface Coord {
    X: number;
    Y: number;
}

interface ORCALine {
    Point: Coord;
    Direction: Coord;
}

interface Agent {
    ID: number;
    Position: Coord;
    Velocity: Coord;
    ORCALines: ORCALine[];
    IsShowLine: boolean;
}

interface Obstacle {
    ID: number;
    Positions: Coord[];
}

interface StepData {
    Index: number;
    Obstacles: Obstacle[];
    Agents: Agent[];
}

interface Param {
    TimeStep: number;
    NeighborDist: number;
    MaxNeighbors: number;
    TimeHorizon: number;
    TimeHorizonObst: number;
    Radius: number;
    MaxSpeed: number;
}

interface DataSize {
    xMax: number;
    xMin: number;
    yMax: number;
    yMin: number;
}

const socket: SocketIOClient.Socket = io();

const mulVec = (vec1: Coord, vec2: Coord): number => {
    return vec1.X * vec2.X + vec1.Y * vec2.Y;
};

const calcDataSize = (data: StepData[]) => {
    var xmax: number = -Infinity;
    var ymax: number = -Infinity;
    var xmin: number = Infinity;
    var ymin: number = Infinity;

    data.forEach((stepData: StepData) => {
        stepData.Agents.forEach((agent: Agent) => {
            if (agent.Position.Y < ymin) {
                ymin = agent.Position.Y;
            }
            if (agent.Position.X < xmin) {
                xmin = agent.Position.X;
            }
            if (agent.Position.Y > ymax) {
                ymax = agent.Position.Y;
            }
            if (agent.Position.X > xmax) {
                xmax = agent.Position.X;
            }
        });
        stepData.Obstacles.forEach((obstacle: Obstacle) => {
            obstacle.Positions.forEach((position: Coord) => {
                if (position.Y < ymin) {
                    ymin = position.Y;
                }
                if (position.X < xmin) {
                    xmin = position.X;
                }
                if (position.Y > ymax) {
                    ymax = position.Y;
                }
                if (position.X > xmax) {
                    xmax = position.X;
                }
            });
        });
    });

    const size: DataSize = {
        xMax: xmax,
        xMin: xmin,
        yMax: ymax,
        yMin: ymin
    };
    console.log("size: ", size);

    return size;
};

const createScatterData = (
    stepData: StepData,
    dataSize: DataSize,
    param: Param,
    isFill: boolean
) => {
    var datasets: any = [];

    // Set Agents Dataset
    stepData.Agents.forEach((agent: Agent) => {
        const agentCoord = [{ x: agent.Position.X, y: agent.Position.Y }];
        // set ORCA Lines
        if (agent.IsShowLine) {
            const epsiron: number = 0.000001;
            agent.ORCALines.forEach((line: ORCALine, index: number) => {
                var orcaCoords: any = [];
                // y =  line.Direction.Y /line.Direction.X * (x - line.Point.X) + line.Point.Y;
                // epsiron: to show line, to add a few trand

                orcaCoords.push({
                    x: dataSize.xMin,
                    y:
                        (line.Direction.Y / (line.Direction.X + epsiron)) *
                            (dataSize.xMin -
                                agent.Position.X -
                                param.TimeStep * line.Point.X) +
                        agent.Position.Y +
                        param.TimeStep * line.Point.Y
                });
                orcaCoords.push({
                    x: dataSize.xMax,
                    y:
                        (line.Direction.Y / (line.Direction.X + epsiron)) *
                            (dataSize.xMax -
                                agent.Position.X -
                                param.TimeStep * line.Point.X) +
                        agent.Position.Y +
                        param.TimeStep * line.Point.Y
                });
                let fill: boolean | string = false;
                let fillColor: string = "rgba(255,255,255,0.5)";
                let backgroundColor: string = "rgba(255,255,255,0.5)";
                if (isFill) {
                    if (mulVec(line.Direction, agent.Velocity) > 0) {
                        fill = "start";
                        fillColor = "rgba(0,255,0,0.5)";
                        backgroundColor = "rgba(0,255,0,0.5)";
                    } else if (mulVec(line.Direction, agent.Velocity) < 0) {
                        fill = "end";
                        fillColor = "rgba(255,0,0,0.5)";
                        backgroundColor = "rgba(255,0,0,0.5)";
                    }
                }
                datasets.push({
                    label:
                        "OrcaLine" +
                        agent.ID.toString() +
                        "-" +
                        index.toString(),
                    fill: fill,
                    fillColor: fillColor,
                    backgroundColor: backgroundColor,
                    showLine: true,
                    borderColor: "rgba(255,0,0,1)",
                    borderWidth: 1,
                    lineTension: 0,
                    data: orcaCoords
                });
            });
            datasets.push({
                id: agent.ID,
                label: "Agent" + agent.ID.toString(),
                fill: false,
                backgroundColor: "rgba(75,192,192,1)",
                pointBorderColor: "rgba(255,0,0,1)",
                pointBackgroundColor: "rgba(255,0,0,1)",
                pointBorderWidth: 1,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: "rgba(75,192,192,1)",
                pointHoverBorderColor: "rgba(220,220,220,1)",
                pointHoverBorderWidth: 2,
                pointRadius: 5,
                pointHitRadius: 10,
                data: agentCoord
            });
        } else {
            datasets.push({
                id: agent.ID,
                label: "Agent" + agent.ID.toString(),
                fill: false,
                backgroundColor: "rgba(0,0,255,1)",
                pointBorderColor: "rgba(0,0,255,1)",
                pointBackgroundColor: "rgba(0,0,255,1)",
                pointBorderWidth: 1,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: "rgba(75,192,192,1)",
                pointHoverBorderColor: "rgba(220,220,220,1)",
                pointHoverBorderWidth: 2,
                pointRadius: 2,
                pointHitRadius: 10,
                data: agentCoord
            });
        }
    });

    // Set Obstacles Dataset
    stepData.Obstacles.forEach((obstacle: Obstacle) => {
        var obstacleCoords: any = [];
        obstacle.Positions.forEach((position: Coord) => {
            obstacleCoords.push({ x: position.X, y: position.Y });
        });
        datasets.push({
            label: "Obstacle" + obstacle.ID.toString(),
            fill: false,
            showLine: true,
            borderColor: "rgba(0,0,0,1)",
            pointHoverRadius: 0,
            pointBorderWidth: 0,
            pointRadius: 0,
            lineTension: 0,
            data: obstacleCoords
        });
    });
    const data = {
        labels: ["Scatter"],
        datasets: datasets
    };

    return data;
};

const setParamType = (anyParam: any): Param => {
    var param: Param = {
        TimeStep: anyParam.timeStep,
        NeighborDist: anyParam.neighborDist,
        MaxNeighbors: anyParam.maxNeighbors,
        TimeHorizon: anyParam.timeHorizon,
        TimeHorizonObst: anyParam.timeHorizonObst,
        Radius: anyParam.radius,
        MaxSpeed: anyParam.maxSpeed
    };
    return param;
};

const setRVODataType = (anyData: any): StepData[] => {
    let data: StepData[] = [];
    anyData.forEach((element: any, index: number) => {
        let obstacles: Obstacle[] = [];
        element.obstacles.forEach((obstacle: any) => {
            let positions: Coord[] = [];
            obstacle.positions.forEach((position: any) => {
                const pos: Coord = {
                    X: position.x,
                    Y: position.y
                };
                positions.push(pos);
            });
            const obst: Obstacle = {
                ID: obstacle.id,
                Positions: positions
            };
            obstacles.push(obst);
        });

        let agents: Agent[] = [];
        element.agents.forEach((agent: any) => {
            let orcaLines: ORCALine[] = [];
            agent.orcaLines.forEach((line: any) => {
                const point: Coord = {
                    X: line.point.x,
                    Y: line.point.y
                };
                const direction: Coord = {
                    X: line.direction.x,
                    Y: line.direction.y
                };
                const orcaLine: ORCALine = {
                    Point: point,
                    Direction: direction
                };
                orcaLines.push(orcaLine);
            });
            const position: Coord = {
                X: agent.position.x,
                Y: agent.position.y
            };
            const velocity: Coord = {
                X: agent.velocity.x,
                Y: agent.velocity.y
            };
            const ag: Agent = {
                ID: agent.id,
                Position: position,
                Velocity: velocity,
                ORCALines: orcaLines,
                IsShowLine: false
            };
            agents.push(ag);
        });

        const stepData: StepData = {
            Obstacles: obstacles,
            Agents: agents,
            Index: index
        };
        data.push(stepData);
    });
    return data;
};

const defaultParam: Param = {
    TimeStep: 0,
    NeighborDist: 0,
    MaxNeighbors: 0,
    TimeHorizon: 0,
    TimeHorizonObst: 0,
    Radius: 0,
    MaxSpeed: 0
};

const App: React.FC = () => {
    const [loading, setLoading] = useState<boolean>(true);
    // all rvo data
    const [data, setData] = useState<StepData[]>([]);
    const [isFill, setIsFill] = useState<boolean>(false);
    const [isScale, setIsScale] = useState<boolean>(false);
    // each step rvo data
    //const [stepData, setStepData] = useState<StepData>(defaultStepData);
    const [showIndex, setShowIndex] = useState<number>(0);
    const [param, setParam] = useState<Param>(defaultParam);
    const [dataSize, setDataSize] = useState<DataSize>({
        xMax: 10,
        xMin: 0,
        yMax: 10,
        yMin: 0
    });

    useEffect(() => {}, [data]);

    socket.on("connect", () => {
        console.log("Socket.IO connected!");
    });
    socket.on("rvo", (strData: string[]) => {
        console.log("strData: ", strData);
        var anyData: any = [];
        // dataのparseとsizeを計算
        strData.forEach(value => {
            anyData.push(JSON.parse(value));
        });
        console.log("anyData: ", anyData);
        const data: StepData[] = setRVODataType(anyData);
        const size: DataSize = calcDataSize(data);
        console.log("data: ", data);

        setData(data);
        //setStepData(data[0]);
        setDataSize(size);
        setLoading(false);
    });

    socket.on("param", (strParam: string) => {
        console.log("param: ", strParam);
        const rvoParam = JSON.parse(strParam);
        const param: Param = setParamType(rvoParam);
        setParam(param);
    });

    socket.on("disconnect", () => {
        console.log("Socket.IO disconnected!");
    });
    const height = dataSize.yMax - dataSize.yMin;
    const width = dataSize.xMax - dataSize.xMin;

    return (
        <div className="App">
            <h2>RVO2 Simulation Monitor</h2>
            {loading ? (
                <h2>Loading...</h2>
            ) : (
                <div>
                    <div
                        style={
                            isScale
                                ? {
                                      width: (600 * width) / height,
                                      height: 600
                                  }
                                : { width: 600, height: 600 }
                        }
                    >
                        <Scatter
                            data={createScatterData(
                                data[showIndex],
                                dataSize,
                                param,
                                isFill
                            )}
                            width={1}
                            height={1}
                            options={{
								legend: {
									display: false
								},
                                maintainAspectRatio: false,
                                events: ["click"],
                                onClick: function(e: any, el: any) {
                                    if (!el || el.length === 0) return;
                                    const index: number =
                                        el[0]._chart.data.datasets[
                                            el[0]._datasetIndex
                                        ].id;
                                    // data内の全てのステップ変更
                                    var data2 = [...data];
                                    data2.forEach((sd: StepData) => {
                                        sd.Agents[index].IsShowLine = !sd
                                            .Agents[index].IsShowLine;
                                    });
                                    //setStepData(sd);
                                    setData(data2);
                                },
                                scales: {
                                    yAxes: [
                                        {
                                            ticks: {
                                                beginAtZero: true,
                                                min: dataSize.yMin - height / 6,
                                                max: dataSize.yMax + height / 6
                                            }
                                        }
                                    ],
                                    xAxes: [
                                        {
                                            ticks: {
                                                beginAtZero: true,
                                                min: dataSize.xMin - width / 6,
                                                max: dataSize.xMax + width / 6
                                            }
                                        }
                                    ]
                                }
                            }}
                        />
                    </div>
                    <div
                        style={
                            isScale
                                ? {
                                      width: (600 * width) / height
                                  }
                                : { width: 600 }
                        }
                    >
                        <Slider
                            defaultValue={0}
                            aria-labelledby="discrete-slider"
                            valueLabelDisplay="auto"
                            onChange={(object, index) => {
                                if (data.length > index) {
                                    if (typeof index === "number") {
                                        //setStepData(data[index]);
                                        setShowIndex(index);
                                    }
                                }
                            }}
                            step={1}
                            marks
                            min={0}
                            max={data.length}
                        />
                    </div>
                    <Button
                        variant={"contained"}
                        onClick={() => setIsFill(!isFill)}
                    >
                        {"Fill: " + isFill}
                    </Button>
                    <Button
                        variant={"contained"}
                        onClick={() => setIsScale(!isScale)}
                    >
                        {"Auto Scale: " + isScale}
                    </Button>
                    <Typography variant={"body1"}>
                        {"AgentNum: " + data[showIndex].Agents.length}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"ObstacleNum: " + data[showIndex].Obstacles.length}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"TimeStep: " + param.TimeStep}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"NeighborDist: " + param.NeighborDist}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"MaxNeighbors: " + param.MaxNeighbors}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"TimeHorizon: " + param.TimeHorizon}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"TimeHorizonObst: " + param.TimeHorizonObst}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"Radius: " + param.Radius}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"MaxSpeed: " + param.MaxSpeed}
                    </Typography>
                </div>
            )}
        </div>
    );
};

export default App;
