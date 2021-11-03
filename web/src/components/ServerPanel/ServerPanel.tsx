import React from 'react';
import { BsCpuFill } from "react-icons/bs";
import { MdOutlineDataSaverOff, MdOutlineSelectAll } from "react-icons/md";
import "./index.css";

interface ServerStatus {
    Cpu: string
    Mem: string
    TotalUpload: string
    TotalDownload: string
    TotalTrafficUse: string
}

interface ClientInfo {
    Id: string
    Ip: string
    Mark: string
    IsOnline: boolean
    ProxyNum: number
}

interface IProps {
    status: ServerStatus
    clients: ClientInfo[]
}

const ServerPanel: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    let clients: JSX.Element[] = props.clients?.map((v, k) => (
        <div className="clientItem" key={k}>
            <div className={`clientStatus ${v.IsOnline ? "online" : "offline"}`}></div>
            <div className="clientInfo">
                <div>mk: {v.Mark && v.Id.substr(0, 4)}</div>
                <div>ip: {v.Ip}</div>
                <div>id: {v.Id}</div>1
            </div>
            <div className="clientProxy">
                <div>Proxy</div>
                <div>{v.ProxyNum}</div>
            </div>
        </div>
    ))
    return (
        <div className="serverInfo">
            <div className="InfoTitle">Information</div>
            <div className="infoItem">
                <div className="infoIcon">
                    <BsCpuFill />
                </div>
                <div className="infoText">
                    <div className="infoTitle">Cpu</div>
                    <div>{props.status.Cpu}</div>
                </div>
            </div>
            <div className="infoItem">
                <div className="infoIcon">
                    <MdOutlineSelectAll />
                </div>
                <div className="infoText">
                    <div className="infoTitle">Mem</div>
                    <div>{props.status.Mem}</div>
                </div>
            </div>
            <div className="infoItem">
                <div className="infoIcon">
                    <MdOutlineDataSaverOff />
                </div>
                <div className="infoText">
                    <div className="infoTitle">DataUse</div>
                        <div>{props.status.TotalTrafficUse} [ {props.status.TotalUpload} / {props.status.TotalDownload} ]</div>
                </div>
            </div>
            <div className="InfoTitle">ClientList</div>
            <div className="clientList">
                {clients && <div className="clientEmpty"><img src="/images/noClient.png" alt=""></img></div>}
            </div>
        </div>
    );
}

export default ServerPanel;