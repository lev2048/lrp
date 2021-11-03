import React from 'react';
import { BsCpuFill } from "react-icons/bs";
import { MdOutlineDataSaverOff, MdOutlineSelectAll } from "react-icons/md";
import "./index.css";

interface ServerStatus {
    cpu: string
    mem: string
    totalUpload: string
    totalDownload: string
    totalTrafficUse: string
}

interface ClientInfo {
    id: string
    ip: string
    mark: string
    online: boolean
    proxyInfos: any
}

interface IProps {
    status: ServerStatus
    clients: ClientInfo[]
}

const ServerPanel: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    let clients: JSX.Element[] = props.clients?.map((v, k) => (
        <div className="clientItem" key={k}>
            <div className={`clientStatus ${v.online ? "online" : "offline"}`}></div>
            <div className="clientInfo">
                <div>mk: {v.mark && v.id.substr(0, 4)}</div>
                <div>ip: {v.ip}</div>
                <div>id: {v.id}</div>1
            </div>
            <div className="clientProxy">
                <div>Proxy</div>
                <div>{v.proxyInfos.length}</div>
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
                    <div>{props.status.cpu}</div>
                </div>
            </div>
            <div className="infoItem">
                <div className="infoIcon">
                    <MdOutlineSelectAll />
                </div>
                <div className="infoText">
                    <div className="infoTitle">Mem</div>
                    <div>{props.status.mem}</div>
                </div>
            </div>
            <div className="infoItem">
                <div className="infoIcon">
                    <MdOutlineDataSaverOff />
                </div>
                <div className="infoText">
                    <div className="infoTitle">DataUse</div>
                        <div>{props.status.totalTrafficUse} [ {props.status.totalUpload} / {props.status.totalDownload} ]</div>
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