import React from 'react';
import axios from 'axios';
import useSWR from 'swr';
import { Menu, TrafficInfo, ProxyInfo, ProxyList, ServerPanel } from "../../components";
import "./index.css"

interface ServerStatus {
    cpu: string
    mem: string
    upload: string
    download: string
    totalUpload: string
    totalDownload: string
    totalTrafficUse: string
}

interface ProxyInfo {
    id: string;
    info: string;
    mark: string;
    isTemp: boolean;
    status: number;
    connNum: number;
}

interface ClientInfo {
    id: string;
    ip: string;
    mark: string;
    online: boolean;
    proxyInfos: ProxyInfo[];
}

interface ServerInfo {
    version: string;
    protocol: string;
    proxyNum: number;
    tProxyNum: number;
    connTotal: number;
    clientNum: number;
    externalIp: string;
    serverPort: string;
    clientInfos: ClientInfo[];
}

const getStatus = (url: string) => axios.get(url).then((res: any) => res.data.info as ServerStatus)
const getServerInfo = (url: string) => axios.get(url).then((res: any) => res.data.info as ServerInfo)

const Dashboard: React.FunctionComponent = (): JSX.Element => {
    const [status]: [status: ServerStatus,error: any] = useSWR('/v1/status', getStatus, { refreshInterval: 1000 }) as any;
    const [serverInfo]: [serverInfo: ServerInfo,error: any] = useSWR('/v1/dashbord', getServerInfo, { refreshInterval: 5000 }) as any;

    return (
        <div className="dashBoard">
            <Menu />
            <div className="content">
                <div className="contentHeader">
                    <div className="systemName">Lrp DashBord</div>
                    <div className="systemInfo">
                        <div className="systemInfoTitle">External IP Address</div>
                        <div>{serverInfo.externalIp}</div>
                    </div>
                    <div style={{ width: '28%' }}>
                        <div className="systemInfoTitle">Software Version</div>
                        <div>{serverInfo.version}</div>
                    </div>
                </div>
                <div className="contentWidget">
                    <TrafficInfo info={status.upload} direction="up" />
                    <TrafficInfo info={status.download} direction="down" />
                    <ProxyInfo connNum={30} proxyNum={10} tproxyNum={5} clientNum={5} />
                </div>
                <ProxyList data={serverInfo.clientInfos} />
            </div>
            <ServerPanel status={status} clients={serverInfo.clientInfos} />
        </div>
    );
}

export default Dashboard;