import React from 'react';
import { BsArrowBarUp, BsArrowBarDown } from "react-icons/bs";
import "./index.css"

interface IProps {
    info: string;
    direction: string
}

const TrafficInfo: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    return (
        <div className="trafficInfo">
            <div className="trafficIcon">
                {props.direction === "up" ? <BsArrowBarUp /> : <BsArrowBarDown />}
            </div>
            <div className="trafficContent">
                <div className="trafficContentTiTle">{props.direction === "up" ? "Upload" : "Download"}</div>
                <div>{props.info}</div>
            </div>
        </div>
    );
}

export default TrafficInfo;