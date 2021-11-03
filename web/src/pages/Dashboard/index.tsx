import React from 'react';
import { Menu,TrafficInfo,ProxyInfo,ProxyList,ServerPanel } from "../../components";
import "./index.css"

interface ServerStatus {
    Cpu: string
    Mem: string
    TotalUpload: string
    TotalDownload: string
    TotalTrafficUse: string
}

const Dashboard: React.FunctionComponent = (): JSX.Element => {
    let status:ServerStatus = {
        Cpu: "20%",
        Mem: "30%",
        TotalUpload: "10 GB",
        TotalDownload: "10 GB",
        TotalTrafficUse: "20 GB",
    };
    
    return (
        <div className="dashBoard">
            <Menu />
            <div className="content">
                <div className="contentHeader">
                    <div className="systemName">Lrp DashBord</div>
                    <div className="systemInfo">
                        <div className="systemInfoTitle">External IP Address</div>
                        <div>10.0.0.1</div>
                    </div>
                    <div style={{ width: '28%' }}>
                        <div className="systemInfoTitle">Software Version</div>
                        <div>v1.0.0</div>
                    </div>
                </div>
                <div className="contentWidget">
                    <TrafficInfo info="12 MB/S" direction="up"/>
                    <TrafficInfo info="48 MB/S" direction="down"/>
                    <ProxyInfo connNum={30} proxyNum={10} tproxyNum={5} clientNum={5}/>
                </div>
                <ProxyList data={[]}/>
            </div>
            <ServerPanel status={status} clients={[]} />
        </div>
    );
}

export default Dashboard;