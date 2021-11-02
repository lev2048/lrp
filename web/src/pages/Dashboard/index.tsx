import React, { useEffect, useState } from 'react';
import { BsFillMenuAppFill, BsArrowBarDown, BsArrowBarUp, BsCpuFill,BsPlusCircleDotted } from "react-icons/bs";
import { MdSettingsPower, MdOutlineDataSaverOff, MdOutlineSelectAll } from "react-icons/md";
import { BiPlanet,BiTerminal,BiTrash } from "react-icons/bi";
import {Menu} from "../../components";
import "./index.css"

const Dashboard: React.FunctionComponent = (): JSX.Element => {
    return <div className="dashBoard">
        <Menu/>
        <div className="content">
            <div className="contentTitle">
                <div className="systemName">Lrp DashBord</div>
                <div style={{ width: '28%' }}>
                    <div className="contentTitleName">External IP Address</div>
                    <div>10.0.0.1</div>
                </div>
                <div style={{ width: '28%' }}>
                    <div className="contentTitleName">Software Version</div>
                    <div>v1.0.0</div>
                </div>
            </div>
            <div className="contentWidget">
                <div className="traffic">
                    <div className="trafficIcon">
                        <BsArrowBarDown />
                    </div>
                    <div className="trafficInfo">
                        <div className="trafficTitle">Download</div>
                        <div>120 MB/s</div>
                    </div>
                </div>
                <div className="traffic" style={{ marginLeft: '1%' }}>
                    <div className="trafficIcon">
                        <BsArrowBarUp />
                    </div>
                    <div className="trafficInfo">
                        <div className="trafficTitle">Upload</div>
                        <div>64 MB/s</div>
                    </div>
                </div>
                <div className="proxyInfo">
                    <div className="proxyItem">
                        <div className="proxyItemTitle">Client</div>
                        <div>4</div>
                    </div>
                    <div className="splitLine"></div>
                    <div className="proxyItem">
                        <div className="proxyItemTitle">Proxy</div>
                        <div>4/2</div>
                    </div>
                    <div className="splitLine"></div>
                    <div className="proxyItem">
                        <div className="proxyItemTitle">Conn</div>
                        <div>12</div>
                    </div>
                </div>
            </div>
            <div className="contentProxyList">
                <div className="proxyListHeader">
                    <div className="InfoTitle">ProxyList</div>
                    <div className="proxyAdd"><BsPlusCircleDotted/></div>
                </div>
                <div className="proxyListContent">
                    <div className="proxyListItem">
                        <div className="proxyIcon"><BiPlanet/></div>
                        <div style={{ width: '21%' }}>c5tuk9b765s2jfnbtfog</div>
                        <div style={{ width: '30%' }}>10.0.0.220:8801 = 10.0.0.1:80</div>
                        <div style={{ width: '10%' }}>30 conn</div>
                        <div className="proxyStatus">Running</div>
                        <div className="proxyDel"><BiTrash/></div>
                    </div>
                </div>
                <div className="proxyListContent">
                    <div className="proxyListItem">
                        <div className="proxyIcon"><BiTerminal/></div>
                        <div style={{ width: '21%' }}>c5tuk9b765s2jfnbtfog</div>
                        <div style={{ width: '30%' }}>10.0.0.220:8801 = 10.0.0.1:80</div>
                        <div style={{ width: '10%' }}>30 conn</div>
                        <div className="proxyStatus">Running</div>
                        <div className="proxyDel"><BiTrash/></div>
                    </div>
                </div>
            </div>
        </div>
        <div className="serverInfo">
            <div className="InfoTitle">Information</div>
            <div className="infoItem">
                <div className="infoIcon">
                    <BsCpuFill />
                </div>
                <div className="infoText">
                    <div className="infoTitle">Cpu</div>
                    <div>23%</div>
                </div>
            </div>
            <div className="infoItem">
                <div className="infoIcon">
                    <MdOutlineSelectAll />
                </div>
                <div className="infoText">
                    <div className="infoTitle">Mem</div>
                    <div>23%</div>
                </div>
            </div>
            <div className="infoItem">
                <div className="infoIcon">
                    <MdOutlineDataSaverOff />
                </div>
                <div className="infoText">
                    <div className="infoTitle">DataUse</div>
                    <div>
                        <div>100 GB [ 10 GB / 90 GB ]</div>
                    </div>
                </div>
            </div>
            <div className="InfoTitle" style={{ marginTop: '35px' }}>ClientList</div>
            <div className="clientList">
                <div className="clientItem">
                    <div className="clientStatus online"></div>
                    <div className="clientInfo">
                        <div>mark: home</div>
                        <div>ip: 10.0.0.1</div>
                        <div>id: c5tuk9b765s2jfnbtfog</div>
                    </div>
                    <div className="clientProxy">
                        <div>Proxy</div>
                        <div>1</div>
                    </div>
                </div>
                <div className="clientItem">
                    <div className="clientStatus offline"></div>
                    <div className="clientInfo">
                        <div>mark: company</div>
                        <div>ip: 192.168.3.122</div>
                        <div>id: c5tuk9b765s2jfnbtfag</div>
                    </div>
                    <div className="clientProxy">
                        <div>Proxy</div>
                        <div>4</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
}

export default Dashboard;