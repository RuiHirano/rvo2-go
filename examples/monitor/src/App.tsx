import React from "react";
import { Scatter } from "react-chartjs-2";
import * as io from "socket.io-client";

const rvodata = {};

const data = {
    labels: ["Scatter"],
    datasets: [
        {
            label: "My First dataset",
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
            data: [
                { x: 65, y: 75 },
                { x: 59, y: 49 },
                { x: 80, y: 90 },
                { x: 81, y: 29 },
                { x: 56, y: 36 },
                { x: 55, y: 25 },
                { x: 40, y: 18 }
            ]
        }
    ]
};

const App: React.FC = () => {
    const getData = () => {
        console.log("get data");
    };

    const socket = io();
    socket.on("connect", () => {
        console.log("Socket.IO connected!");
    });
    socket.on("rvo", getData());
    socket.on("disconnect", () => {
        console.log("Socket.IO disconnected!");
    });

    return (
        <div className="App">
            <h2>Line Example</h2>
            <Scatter data={data} />
        </div>
    );
};

export default App;
