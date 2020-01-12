import React, { useState } from "react";
import { Scatter } from "react-chartjs-2";
import io from "socket.io-client";
import Slider from "@material-ui/core/Slider";
import { Typography } from "@material-ui/core";

const rvodata = {};

const socket: SocketIOClient.Socket = io();

const calcDataSize = (rvoData: any) => {
    var xmax: number = -Infinity;
    var ymax: number = -Infinity;
    var xmin: number = Infinity;
    var ymin: number = Infinity;

    rvoData.forEach((stepData: any) => {
        stepData.agents.forEach((agent: any) => {
            if (agent.y < ymin) {
                ymin = agent.y;
            }
            if (agent.x < xmin) {
                xmin = agent.x;
            }
            if (agent.y > ymax) {
                ymax = agent.y;
            }
            if (agent.x > xmax) {
                xmax = agent.x;
            }
        });
        stepData.obstacles.forEach((obstacle: any) => {
            obstacle.positions.forEach((position: any) => {
                if (position.y < ymin) {
                    ymin = position.y;
                }
                if (position.x < xmin) {
                    xmin = position.x;
                }
                if (position.y > ymax) {
                    ymax = position.y;
                }
                if (position.x > xmax) {
                    xmax = position.x;
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

const createScatterData = (stepData: any) => {
    var datasets: any = [];

    // Set Agents Dataset
    var agentCoords: any = [];
    stepData.agents.forEach((agent: any) => {
        agentCoords.push({ x: agent.x, y: agent.y });
    });
    datasets.push({
        label: "Agent",
        fill: false,
        backgroundColor: "rgba(75,192,192,0.4)",
        pointBorderColor: "rgba(75,192,192,1)",
        pointBackgroundColor: "#fff",
        pointBorderWidth: 1,
        pointHoverRadius: 5,
        pointHoverBackgroundColor: "rgba(75,192,192,1)",
        pointHoverBorderColor: "rgba(220,220,220,1)",
        pointHoverBorderWidth: 2,
        pointRadius: 1,
        pointHitRadius: 10,
        data: agentCoords
    });

    stepData.obstacles.forEach((obstacle: any) => {
        var obstacleCoords: any = [];
        obstacle.positions.forEach((position: any) => {
            obstacleCoords.push({ x: position.x, y: position.y });
        });
        datasets.push({
            label: "Obstacle" + obstacle.id.toString(),
            fill: false,
            showLine: true,
            backgroundColor: "rgba(75,192,192,0.4)",
            pointBorderColor: "rgba(75,192,192,1)",
            lineTension: 0,
            pointBackgroundColor: "#fff",
            pointBorderWidth: 1,
            pointHoverRadius: 5,
            pointHoverBackgroundColor: "rgba(75,192,192,1)",
            pointHoverBorderColor: "rgba(220,220,220,1)",
            pointHoverBorderWidth: 2,
            pointRadius: 1,
            pointHitRadius: 10,
            data: obstacleCoords
        });
    });
    const data = {
        labels: ["Scatter"],
        datasets: datasets
    };

    return data;
};

interface DataSize {
    xMax: number;
    xMin: number;
    yMax: number;
    yMin: number;
}

const App: React.FC = () => {
    // all rvo data
    const [data, setData] = useState([]);
    // each step rvo data
    const [stepData, setStepData] = useState({});
    const [param, setParam] = useState<any>({});
    const [dataSize, setDataSize] = useState<DataSize>({
        xMax: 10,
        xMin: 0,
        yMax: 10,
        yMin: 0
    });

    socket.on("connect", () => {
        console.log("Socket.IO connected!");
    });
    socket.on("rvo", (strData: string[]) => {
        var rvoData: any = [];
        // dataのparseとsizeを計算
        strData.forEach(value => {
            rvoData.push(JSON.parse(value));
        });
        console.log("rvoData: ", rvoData);
        const size = calcDataSize(rvoData);

        setData(rvoData);
        setStepData(rvoData[0]);
        setDataSize(size);
    });

    socket.on("param", (param: string) => {
        console.log("param: ", param);
        const rvoParam = JSON.parse(param);
        setParam(rvoParam);
    });

    socket.on("disconnect", () => {
        console.log("Socket.IO disconnected!");
    });

    const height = dataSize.yMax - dataSize.yMin;
    const width = dataSize.xMax - dataSize.xMin;

    return (
        <div className="App">
            <h2>RVO2 Simulation Monitor</h2>
            {Object.keys(stepData).length === 0 ? (
                <h2>Loading...</h2>
            ) : (
                <div>
                    <div
                        style={{
                            width: 600,
                            height: (600 * height) / width
                        }}
                    >
                        <Scatter
                            data={createScatterData(stepData)}
                            width={1}
                            height={1}
                            options={{
                                maintainAspectRatio: false,
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
                    <Slider
                        defaultValue={0}
                        aria-labelledby="discrete-slider"
                        valueLabelDisplay="auto"
                        onChange={(object, value) => {
                            if (data.length > value) {
                                if (typeof value === "number") {
                                    setStepData(data[value]);
                                }
                            }
                        }}
                        step={1}
                        marks
                        min={0}
                        max={data.length}
                    />
                    <Typography variant={"body1"}>
                        {"TimeStep: " + param.timeStep}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"NeighborDist: " + param.neighborDist}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"MaxNeighbors: " + param.maxNeighbors}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"TimeHorizon: " + param.timeHorizon}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"TimeHorizonObst: " + param.timeHorizonObst}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"Radius: " + param.radius}
                    </Typography>
                    <Typography variant={"body1"}>
                        {"MaxSpeed: " + param.maxSpeed}
                    </Typography>
                </div>
            )}
        </div>
    );
};

export default App;
