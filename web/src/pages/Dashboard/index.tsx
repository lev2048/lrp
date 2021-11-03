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

interface Proxy {
    id: string;
    cid: string;
    info: string;
    mark: string;
    isTemp: boolean;
    status: number;
    connNum: number;
}

interface Client {
    id: string;
    ip: string;
    mark: string;
    online: boolean;
    proxyInfos: Proxy[];
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
    clientInfos: Client[];
}

const getStatus = (url: string) => axios.get(url).then((res: any) => {
    return res.data.info as ServerStatus
})
const getServerInfo = (url: string) => axios.get(url).then((res: any) => res.data.info as ServerInfo)

const Dashboard: React.FunctionComponent = (): JSX.Element => {
    const status = useSWR<ServerStatus, any>('http://10.0.0.216/v1/status', getStatus, { refreshInterval: 1000 });
    const serverInfo = useSWR<ServerInfo, any>('http://10.0.0.216/v1/dashbord', getServerInfo, { refreshInterval: 5000 });
    if (!status.data || !serverInfo.data) return <div>loading</div>
    return (
        <div className="dashBoard">
            <Menu />
            <div className="content">
                <div className="contentHeader">
                    <div className="systemName">Lrp DashBord</div>
                    <div className="systemInfo">
                        <div className="systemInfoTitle">External IP Address</div>
                        <div>{serverInfo.data.externalIp}</div>
                    </div>
                    <div style={{ width: '28%' }}>
                        <div className="systemInfoTitle">Software Version</div>
                        <div>{serverInfo.data.version}</div>
                    </div>
                </div>
                <div className="contentWidget">
                    <TrafficInfo info={status.data.upload} direction="up" />
                    <TrafficInfo info={status.data.download} direction="down" />
                    <ProxyInfo connNum={serverInfo.data.connTotal} proxyNum={serverInfo.data.proxyNum} tproxyNum={serverInfo.data.tProxyNum} clientNum={serverInfo.data.clientNum} />
                </div>
                <ProxyList data={serverInfo.data.clientInfos} />
            </div>
            <ServerPanel status={status.data} clients={serverInfo.data.clientInfos} />
        </div>
    );
}

export default Dashboard;