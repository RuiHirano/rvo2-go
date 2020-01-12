import React, { useState } from "react";
import { Scatter } from "react-chartjs-2";
import io from "socket.io-client";
import Slider from "@material-ui/core/Slider";

const rvodata = {};

const socket: SocketIOClient.Socket = io();

const calcDataSize = (rvoData: any) => {
    var xmax: number = -Infinity;
    var ymax: number = -Infinity;
    var xmin: number = Infinity;
    var ymin: number = Infinity;

    rvoData.forEach((stepData: any) => {
        stepData.agents.forEach((agent: any) => {
            if (agent.lat < ymin) {
                ymin = agent.lat;
            }
            if (agent.lon < xmin) {
                xmin = agent.lon;
            }
            if (agent.lat > ymax) {
                ymax = agent.lat;
            }
            if (agent.lon > xmax) {
                xmax = agent.lon;
            }
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
    console.log("stepData: %v\n", stepData);
    var coords: any = [];
    stepData.agents.forEach((agent: any) => {
        coords.push({ x: agent.lon, y: agent.lat });
    });
    const data = {
        labels: ["Scatter"],
        datasets: [
            {
                label: "Monitor",
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
                data: coords
            }
        ]
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
    const [dataSize, setDataSize] = useState<DataSize>({
        xMax: 10,
        xMin: 0,
        yMax: 10,
        yMin: 0
    });
    console.log("test3");

    socket.on("connect", () => {
        console.log("Socket.IO connected!");
        socket.emit("some:event", "test");
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

    socket.on("disconnect", () => {
        console.log("Socket.IO disconnected!");
    });

    console.log(
        "size2: ",
        parseInt(
            (
                (200 * (dataSize.yMax - dataSize.yMin)) /
                (dataSize.xMax - dataSize.xMin)
            ).toString()
        )
    );

    const height = dataSize.yMax - dataSize.yMin;
    const width = dataSize.xMax - dataSize.xMin;
    console.log("hw: ", height, width);

    return (
        <div className="App">
            <h2>RVO2 Simulation Monitor</h2>
            {Object.keys(stepData).length === 0 ? (
                <h2>Loading...</h2>
            ) : (
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
                                            min: dataSize.yMin,
                                            max: dataSize.yMax
                                        }
                                    }
                                ],
                                xAxes: [
                                    {
                                        ticks: {
                                            beginAtZero: true,
                                            min: dataSize.xMin,
                                            max: dataSize.xMax
                                        }
                                    }
                                ]
                            }
                        }}
                    />
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
                </div>
            )}
        </div>
    );
};

export default App;
